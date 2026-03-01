package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	store *Store
}

func NewServer(store *Store) *Server {
	return &Server{store: store}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, http.StatusOK, map[string]string{"status": "ok"})
	})

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

	return loggingMiddleware(corsMiddleware(mux))
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
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
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
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
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
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
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
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
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
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
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
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if !s.store.IsVisitInTenant(in.VisitID, session.TenantID) {
		jsonResponse(w, http.StatusForbidden, map[string]string{"error": "visit does not belong to current tenant"})
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
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if !s.store.IsVisitInTenant(in.VisitID, session.TenantID) {
		jsonResponse(w, http.StatusForbidden, map[string]string{"error": "visit does not belong to current tenant"})
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if req.ShopDomain == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "shop_domain is required"})
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{
		"message": "OAuth start placeholder created",
		"next":    "build Shopify authorization URL and redirect merchant",
	})
}

func (s *Server) shopifyWebhook(w http.ResponseWriter, r *http.Request) {
	var event ShopifyWebhookEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	jsonResponse(w, http.StatusAccepted, map[string]string{
		"message": "webhook accepted",
		"topic":   event.Topic,
	})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var req AuthLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	resp, err := s.store.Login(req)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "invalid credentials") {
			jsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
			return
		}
		jsonResponse(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	jsonResponse(w, http.StatusOK, resp)
}

func (s *Server) requireSession(w http.ResponseWriter, r *http.Request) (SessionInfo, bool) {
	sessionID := strings.TrimSpace(r.Header.Get("X-Session-ID"))
	if sessionID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "missing session"})
		return SessionInfo{}, false
	}
	session, err := s.store.GetSession(sessionID)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "invalid session"})
		return SessionInfo{}, false
	}
	return session, true
}

func jsonResponse(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Shopify-Hmac-SHA256, X-Session-ID")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
