package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

type Server struct {
	store *Store
}

type contextKey string

const (
	requestIDContextKey   contextKey = "request_id"
	sessionInfoContextKey contextKey = "session_info"
	sqliteSchemaVersion   string     = "2026-03-v1"
)

var requestIDSeq uint64

func NewServer(store *Store) *Server {
	return &Server{store: store}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", s.health)

	// Auth API
	mux.HandleFunc("POST /api/merchant/auth/login", s.login)

	// Dashboard
	mux.HandleFunc("GET /api/merchant/dashboard", s.dashboardData)

	// Shop APIs
	mux.HandleFunc("GET /api/merchant/orders", s.listOrders)
	mux.HandleFunc("POST /api/merchant/orders", s.createOrder)
	mux.HandleFunc("GET /api/merchant/products", s.listProducts)
	mux.HandleFunc("POST /api/merchant/products", s.createProduct)
	mux.HandleFunc("GET /api/merchant/customers", s.listCustomers)
	mux.HandleFunc("POST /api/merchant/customers", s.createCustomer)

	// Clinic APIs
	mux.HandleFunc("GET /api/merchant/appointments", s.listAppointments)
	mux.HandleFunc("POST /api/merchant/appointments", s.createAppointment)
	mux.HandleFunc("GET /api/merchant/visits", s.listVisits)
	mux.HandleFunc("POST /api/merchant/visits", s.createVisit)
	mux.HandleFunc("GET /api/merchant/prescriptions", s.listPrescriptions)
	mux.HandleFunc("POST /api/merchant/prescriptions", s.createPrescription)
	mux.HandleFunc("GET /api/merchant/followups", s.listFollowUps)
	mux.HandleFunc("POST /api/merchant/followups", s.createFollowUp)

	// Shopify placeholders
	mux.HandleFunc("POST /api/merchant/shopify/oauth/start", s.shopifyOAuthStart)
	mux.HandleFunc("POST /api/merchant/shopify/webhooks", s.shopifyWebhook)

	return s.loggingMiddleware(corsMiddleware(mux))
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	resp := map[string]any{
		"status": "ok",
		"service": map[string]any{
			"name":     "petwell-merchant-backend",
			"time_utc": time.Now().UTC().Format(time.RFC3339),
		},
		"store": map[string]any{
			"type": "memory",
		},
		"db": map[string]any{
			"enabled":        false,
			"reachable":      true,
			"schema_version": "",
		},
	}

	statusCode := http.StatusOK
	if s.store.db != nil {
		resp["store"] = map[string]any{
			"type": "sqlite",
		}
		dbPayload := map[string]any{
			"enabled":        true,
			"reachable":      true,
			"schema_version": sqliteSchemaVersion,
		}
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := s.store.db.PingContext(ctx); err != nil {
			dbPayload["reachable"] = false
			dbPayload["error"] = err.Error()
			resp["status"] = "degraded"
			statusCode = http.StatusServiceUnavailable
		}
		resp["db"] = dbPayload
	}

	jsonResponse(w, statusCode, resp)
}

func (s *Server) dashboardData(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = session.MerchantType
	}
	orders := s.store.ListOrders(session.TenantID)
	products := s.store.ListProducts(session.TenantID)
	appointments := s.store.ListAppointments(session.TenantID)
	visits := s.store.ListVisits(session.TenantID)
	followups := s.store.ListFollowUps(session.TenantID)

	if mode == "clinic" {
		jsonResponse(w, http.StatusOK, map[string]any{
			"mode": mode,
			"kpis": []map[string]any{
				{"label": "今日预约", "value": len(appointments)},
				{"label": "就诊中", "value": len(visits)},
				{"label": "回访待办", "value": len(followups)},
				{"label": "已登记客户", "value": len(s.store.ListCustomers(session.TenantID))},
			},
		})
		return
	}

	lowStock := 0
	for _, p := range products {
		if p.Stock <= 10 {
			lowStock++
		}
	}
	totalAmount := 0.0
	for _, o := range orders {
		totalAmount += o.AmountHKD
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"mode": mode,
		"kpis": []map[string]any{
			{"label": "订单总数", "value": len(orders)},
			{"label": "商品总数", "value": len(products)},
			{"label": "库存预警", "value": lowStock},
			{"label": "累计GMV(HKD)", "value": totalAmount},
		},
	})
}

func (s *Server) createOrder(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	var in Order
	if !decodeJSONBody(w, r, &in) {
		return
	}
	in.TenantID = session.TenantID
	created := s.store.CreateOrder(in)
	jsonResponse(w, http.StatusCreated, created)
}

func (s *Server) listOrders(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	jsonResponse(w, http.StatusOK, s.store.ListOrders(session.TenantID))
}

func (s *Server) createProduct(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	var in Product
	if !decodeJSONBody(w, r, &in) {
		return
	}
	in.TenantID = session.TenantID
	created := s.store.CreateProduct(in)
	jsonResponse(w, http.StatusCreated, created)
}

func (s *Server) listProducts(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	jsonResponse(w, http.StatusOK, s.store.ListProducts(session.TenantID))
}

func (s *Server) createCustomer(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	var in Customer
	if !decodeJSONBody(w, r, &in) {
		return
	}
	in.TenantID = session.TenantID
	created := s.store.CreateCustomer(in)
	jsonResponse(w, http.StatusCreated, created)
}

func (s *Server) listCustomers(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	jsonResponse(w, http.StatusOK, s.store.ListCustomers(session.TenantID))
}

func (s *Server) createAppointment(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	var in Appointment
	if !decodeJSONBody(w, r, &in) {
		return
	}
	in.TenantID = session.TenantID
	created := s.store.CreateAppointment(in)
	jsonResponse(w, http.StatusCreated, created)
}

func (s *Server) listAppointments(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	jsonResponse(w, http.StatusOK, s.store.ListAppointments(session.TenantID))
}

func (s *Server) createVisit(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	var in Visit
	if !decodeJSONBody(w, r, &in) {
		return
	}
	in.TenantID = session.TenantID
	created := s.store.CreateVisit(in)
	jsonResponse(w, http.StatusCreated, created)
}

func (s *Server) listVisits(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	jsonResponse(w, http.StatusOK, s.store.ListVisits(session.TenantID))
}

func (s *Server) createPrescription(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	var in Prescription
	if !decodeJSONBody(w, r, &in) {
		return
	}
	if !s.store.IsVisitInTenant(in.VisitID, session.TenantID) {
		writeAPIError(w, r, http.StatusForbidden, "forbidden_visit_tenant_mismatch", "visit does not belong to current tenant")
		return
	}
	created := s.store.CreatePrescription(in)
	jsonResponse(w, http.StatusCreated, created)
}

func (s *Server) listPrescriptions(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	jsonResponse(w, http.StatusOK, s.store.ListPrescriptions(session.TenantID))
}

func (s *Server) createFollowUp(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	var in FollowUp
	if !decodeJSONBody(w, r, &in) {
		return
	}
	if !s.store.IsVisitInTenant(in.VisitID, session.TenantID) {
		writeAPIError(w, r, http.StatusForbidden, "forbidden_visit_tenant_mismatch", "visit does not belong to current tenant")
		return
	}
	created := s.store.CreateFollowUp(in)
	jsonResponse(w, http.StatusCreated, created)
}

func (s *Server) listFollowUps(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	jsonResponse(w, http.StatusOK, s.store.ListFollowUps(session.TenantID))
}

func (s *Server) shopifyOAuthStart(w http.ResponseWriter, r *http.Request) {
	var req ShopifyOAuthRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	if req.ShopDomain == "" {
		writeAPIError(w, r, http.StatusBadRequest, "invalid_request_shop_domain_required", "shop_domain is required")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{
		"message": "OAuth start placeholder created",
		"next":    "build Shopify authorization URL and redirect merchant",
	})
}

func (s *Server) shopifyWebhook(w http.ResponseWriter, r *http.Request) {
	var event ShopifyWebhookEvent
	if !decodeJSONBody(w, r, &event) {
		return
	}
	jsonResponse(w, http.StatusAccepted, map[string]string{
		"message": "webhook accepted",
		"topic":   event.Topic,
	})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var req AuthLoginRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	resp, err := s.store.Login(req)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "invalid credentials") {
			writeAPIError(w, r, http.StatusUnauthorized, "invalid_credentials", "invalid credentials")
			return
		}
		log.Printf("request_id=%s login_error=%v", requestIDFromContext(r.Context()), err)
		writeAPIError(w, r, http.StatusInternalServerError, "internal_error", "internal server error")
		return
	}
	jsonResponse(w, http.StatusOK, resp)
}

func (s *Server) requireSession(w http.ResponseWriter, r *http.Request) (SessionInfo, bool) {
	sessionID := strings.TrimSpace(r.Header.Get("X-Session-ID"))
	if sessionID == "" {
		writeAPIError(w, r, http.StatusUnauthorized, "missing_session", "missing session")
		return SessionInfo{}, false
	}
	if cached, ok := sessionFromContext(r.Context()); ok && cached.SessionID == sessionID {
		return cached, true
	}
	session, err := s.store.GetSession(sessionID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "expired") {
			writeAPIError(w, r, http.StatusUnauthorized, "session_expired", "session expired")
			return SessionInfo{}, false
		}
		writeAPIError(w, r, http.StatusUnauthorized, "invalid_session", "invalid session")
		return SessionInfo{}, false
	}
	return session, true
}

type apiErrorResponse struct {
	Error     apiErrorDetail `json:"error"`
	RequestID string         `json:"request_id"`
}

type apiErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		writeAPIError(w, r, http.StatusBadRequest, "invalid_json", "invalid JSON")
		return false
	}
	var extra json.RawMessage
	if err := decoder.Decode(&extra); err != io.EOF {
		writeAPIError(w, r, http.StatusBadRequest, "invalid_json", "invalid JSON")
		return false
	}
	return true
}

func writeAPIError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	jsonResponse(w, status, apiErrorResponse{
		Error: apiErrorDetail{
			Code:    code,
			Message: message,
		},
		RequestID: requestIDFromContext(r.Context()),
	})
}

func requestIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(requestIDContextKey).(string)
	return strings.TrimSpace(v)
}

func sessionFromContext(ctx context.Context) (SessionInfo, bool) {
	v, ok := ctx.Value(sessionInfoContextKey).(SessionInfo)
	return v, ok
}

func newRequestID() string {
	seq := atomic.AddUint64(&requestIDSeq, 1)
	return fmt.Sprintf("req_%d_%06d", time.Now().UnixMilli(), seq%1_000_000)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(b)
}

func logFieldOrDash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "-"
	}
	return v
}

func jsonResponse(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Shopify-Hmac-SHA256, X-Session-ID, X-Request-ID")
		w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := strings.TrimSpace(r.Header.Get("X-Request-ID"))
		if requestID == "" {
			requestID = newRequestID()
		}

		ctx := context.WithValue(r.Context(), requestIDContextKey, requestID)

		var tenantID, userID string
		sessionID := strings.TrimSpace(r.Header.Get("X-Session-ID"))
		if sessionID != "" {
			if session, err := s.store.GetSession(sessionID); err == nil {
				ctx = context.WithValue(ctx, sessionInfoContextKey, session)
				tenantID = session.TenantID
				userID = session.UserID
			}
		}

		r = r.WithContext(ctx)
		w.Header().Set("X-Request-ID", requestID)

		recorder := &statusRecorder{ResponseWriter: w}
		startedAt := time.Now()
		next.ServeHTTP(recorder, r)

		status := recorder.status
		if status == 0 {
			status = http.StatusOK
		}
		log.Printf(
			"request_id=%s method=%s path=%s status=%d duration_ms=%d tenant_id=%s user_id=%s",
			requestID,
			r.Method,
			r.URL.Path,
			status,
			time.Since(startedAt).Milliseconds(),
			logFieldOrDash(tenantID),
			logFieldOrDash(userID),
		)
	})
}
