package internal

import (
	"encoding/json"
	"log"
	"net/http"
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
	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "shop"
	}
	orders := s.store.ListOrders()
	products := s.store.ListProducts()
	appointments := s.store.ListAppointments()
	visits := s.store.ListVisits()
	followups := s.store.ListFollowUps()

	if mode == "clinic" {
		jsonResponse(w, http.StatusOK, map[string]any{
			"mode": mode,
			"kpis": []map[string]any{
				{"label": "今日预约", "value": len(appointments)},
				{"label": "就诊中", "value": len(visits)},
				{"label": "回访待办", "value": len(followups)},
				{"label": "已登记客户", "value": len(s.store.ListCustomers())},
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
	var in Order
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	created := s.store.CreateOrder(in)
	jsonResponse(w, http.StatusCreated, created)
}

func (s *Server) listOrders(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, s.store.ListOrders())
}

func (s *Server) createProduct(w http.ResponseWriter, r *http.Request) {
	var in Product
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	created := s.store.CreateProduct(in)
	jsonResponse(w, http.StatusCreated, created)
}

func (s *Server) listProducts(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, s.store.ListProducts())
}

func (s *Server) createCustomer(w http.ResponseWriter, r *http.Request) {
	var in Customer
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	created := s.store.CreateCustomer(in)
	jsonResponse(w, http.StatusCreated, created)
}

func (s *Server) listCustomers(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, s.store.ListCustomers())
}

func (s *Server) createAppointment(w http.ResponseWriter, r *http.Request) {
	var in Appointment
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	created := s.store.CreateAppointment(in)
	jsonResponse(w, http.StatusCreated, created)
}

func (s *Server) listAppointments(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, s.store.ListAppointments())
}

func (s *Server) createVisit(w http.ResponseWriter, r *http.Request) {
	var in Visit
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	created := s.store.CreateVisit(in)
	jsonResponse(w, http.StatusCreated, created)
}

func (s *Server) listVisits(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, s.store.ListVisits())
}

func (s *Server) createPrescription(w http.ResponseWriter, r *http.Request) {
	var in Prescription
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	created := s.store.CreatePrescription(in)
	jsonResponse(w, http.StatusCreated, created)
}

func (s *Server) listPrescriptions(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, s.store.ListPrescriptions())
}

func (s *Server) createFollowUp(w http.ResponseWriter, r *http.Request) {
	var in FollowUp
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	created := s.store.CreateFollowUp(in)
	jsonResponse(w, http.StatusCreated, created)
}

func (s *Server) listFollowUps(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, s.store.ListFollowUps())
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

func jsonResponse(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Shopify-Hmac-SHA256")
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
