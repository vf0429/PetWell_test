package internal

import (
	"fmt"
	"sync"
	"time"
)

type Store struct {
	mu            sync.RWMutex
	orders        map[string]Order
	products      map[string]Product
	customers     map[string]Customer
	appointments  map[string]Appointment
	visits        map[string]Visit
	prescriptions map[string]Prescription
	followUps     map[string]FollowUp
}

func NewStore() *Store {
	return &Store{
		orders:        make(map[string]Order),
		products:      make(map[string]Product),
		customers:     make(map[string]Customer),
		appointments:  make(map[string]Appointment),
		visits:        make(map[string]Visit),
		prescriptions: make(map[string]Prescription),
		followUps:     make(map[string]FollowUp),
	}
}

func newID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func (s *Store) CreateAppointment(in Appointment) Appointment {
	s.mu.Lock()
	defer s.mu.Unlock()
	in.ID = newID("appt")
	if in.Status == "" {
		in.Status = "pending"
	}
	s.appointments[in.ID] = in
	return in
}

func (s *Store) CreateOrder(in Order) Order {
	s.mu.Lock()
	defer s.mu.Unlock()
	in.ID = newID("ord")
	if in.Status == "" {
		in.Status = "pending"
	}
	if in.CreatedAt.IsZero() {
		in.CreatedAt = time.Now()
	}
	s.orders[in.ID] = in
	return in
}

func (s *Store) ListOrders() []Order {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Order, 0, len(s.orders))
	for _, v := range s.orders {
		out = append(out, v)
	}
	return out
}

func (s *Store) CreateProduct(in Product) Product {
	s.mu.Lock()
	defer s.mu.Unlock()
	in.ID = newID("prd")
	if in.UpdatedAt.IsZero() {
		in.UpdatedAt = time.Now()
	}
	s.products[in.ID] = in
	return in
}

func (s *Store) ListProducts() []Product {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Product, 0, len(s.products))
	for _, v := range s.products {
		out = append(out, v)
	}
	return out
}

func (s *Store) CreateCustomer(in Customer) Customer {
	s.mu.Lock()
	defer s.mu.Unlock()
	in.ID = newID("cus")
	if in.CreatedAt.IsZero() {
		in.CreatedAt = time.Now()
	}
	if in.LastActive.IsZero() {
		in.LastActive = in.CreatedAt
	}
	s.customers[in.ID] = in
	return in
}

func (s *Store) ListCustomers() []Customer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Customer, 0, len(s.customers))
	for _, v := range s.customers {
		out = append(out, v)
	}
	return out
}

func (s *Store) ListAppointments() []Appointment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Appointment, 0, len(s.appointments))
	for _, v := range s.appointments {
		out = append(out, v)
	}
	return out
}

func (s *Store) CreateVisit(in Visit) Visit {
	s.mu.Lock()
	defer s.mu.Unlock()
	in.ID = newID("visit")
	if in.Status == "" {
		in.Status = "in_progress"
	}
	if in.CheckInAt.IsZero() {
		in.CheckInAt = time.Now()
	}
	s.visits[in.ID] = in
	return in
}

func (s *Store) ListVisits() []Visit {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Visit, 0, len(s.visits))
	for _, v := range s.visits {
		out = append(out, v)
	}
	return out
}

func (s *Store) CreatePrescription(in Prescription) Prescription {
	s.mu.Lock()
	defer s.mu.Unlock()
	in.ID = newID("rx")
	s.prescriptions[in.ID] = in
	return in
}

func (s *Store) ListPrescriptions() []Prescription {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Prescription, 0, len(s.prescriptions))
	for _, v := range s.prescriptions {
		out = append(out, v)
	}
	return out
}

func (s *Store) CreateFollowUp(in FollowUp) FollowUp {
	s.mu.Lock()
	defer s.mu.Unlock()
	in.ID = newID("fu")
	if in.Status == "" {
		in.Status = "pending"
	}
	s.followUps[in.ID] = in
	return in
}

func (s *Store) ListFollowUps() []FollowUp {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]FollowUp, 0, len(s.followUps))
	for _, v := range s.followUps {
		out = append(out, v)
	}
	return out
}

func (s *Store) SeedDemoData() {
	if len(s.orders) > 0 || len(s.appointments) > 0 {
		return
	}

	s.CreateCustomer(Customer{
		TenantID: "tenant_demo",
		Name:     "Amy Chan",
		Email:    "amy@example.com",
		Phone:    "+85290000001",
	})
	s.CreateCustomer(Customer{
		TenantID: "tenant_demo",
		Name:     "Ken Wong",
		Email:    "ken@example.com",
		Phone:    "+85290000002",
	})

	s.CreateProduct(Product{
		TenantID: "tenant_demo",
		Name:     "Omega-3 Supplement",
		SKU:      "PW-OMG3-01",
		Stock:    18,
		PriceHKD: 198,
	})
	s.CreateProduct(Product{
		TenantID: "tenant_demo",
		Name:     "Urinary Cat Food",
		SKU:      "PW-UR-08",
		Stock:    5,
		PriceHKD: 320,
	})

	s.CreateOrder(Order{
		TenantID:   "tenant_demo",
		CustomerID: "seed_customer_1",
		Status:     "paid",
		AmountHKD:  398,
	})
	s.CreateOrder(Order{
		TenantID:   "tenant_demo",
		CustomerID: "seed_customer_2",
		Status:     "pending",
		AmountHKD:  722,
	})

	s.CreateAppointment(Appointment{
		TenantID:       "tenant_demo",
		PetID:          "pet_001",
		DoctorID:       "doc_001",
		ScheduledAt:    time.Now().Add(2 * time.Hour),
		Status:         "confirmed",
		ChiefComplaint: "skin itching",
	})
	s.CreateVisit(Visit{
		TenantID: "tenant_demo",
		PetID:    "pet_008",
		DoctorID: "doc_002",
		Status:   "in_progress",
	})
}
