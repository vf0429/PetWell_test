package internal

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type staffUser struct {
	ID           string
	TenantID     string
	MerchantType string
	Role         string
	Name         string
	Email        string
	Phone        string
	AuthProvider string
	PasswordHash string
}

type Store struct {
	mu            sync.RWMutex
	db            *sql.DB
	orders        map[string]Order
	products      map[string]Product
	customers     map[string]Customer
	appointments  map[string]Appointment
	visits        map[string]Visit
	prescriptions map[string]Prescription
	followUps     map[string]FollowUp
	users         map[string]staffUser
	sessions      map[string]SessionInfo
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
		users:         make(map[string]staffUser),
		sessions:      make(map[string]SessionInfo),
	}
}

func NewSQLiteStore(dbPath string) (*Store, error) {
	if dbPath == "" {
		dbPath = "./db/petwell_merchant.db"
	}
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	store := NewStore()
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)", filepath.ToSlash(dbPath))
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	if err := initSQLiteSchema(db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("init schema: %w", err)
	}

	store.db = db
	return store, nil
}

func (s *Store) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func initSQLiteSchema(db *sql.DB) error {
	const schema = `
CREATE TABLE IF NOT EXISTS orders (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  customer_id TEXT NOT NULL,
  status TEXT NOT NULL,
  amount_hkd REAL NOT NULL,
  created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  name TEXT NOT NULL,
  sku TEXT NOT NULL,
  stock INTEGER NOT NULL,
  price_hkd REAL NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS customers (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  phone TEXT NOT NULL,
  created_at TEXT NOT NULL,
  last_active TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS appointments (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  pet_id TEXT NOT NULL,
  doctor_id TEXT NOT NULL,
  scheduled_at TEXT NOT NULL,
  status TEXT NOT NULL,
  chief_complaint TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS visits (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  appointment_id TEXT NOT NULL,
  pet_id TEXT NOT NULL,
  doctor_id TEXT NOT NULL,
  status TEXT NOT NULL,
  checkin_at TEXT NOT NULL,
  checkout_at TEXT
);

CREATE TABLE IF NOT EXISTS prescriptions (
  id TEXT PRIMARY KEY,
  visit_id TEXT NOT NULL,
  medicine_name TEXT NOT NULL,
  dosage TEXT NOT NULL,
  frequency TEXT NOT NULL,
  duration_days INTEGER NOT NULL,
  notes TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS followups (
  id TEXT PRIMARY KEY,
  visit_id TEXT NOT NULL,
  due_at TEXT,
  channel TEXT NOT NULL,
  status TEXT NOT NULL,
  result_notes TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS staff_users (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  merchant_type TEXT NOT NULL,
  role TEXT NOT NULL,
  name TEXT NOT NULL,
  email TEXT UNIQUE,
  phone TEXT UNIQUE,
  auth_provider TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  tenant_id TEXT NOT NULL,
  merchant_type TEXT NOT NULL,
  role TEXT NOT NULL,
  display_name TEXT NOT NULL,
  method TEXT NOT NULL,
  created_at TEXT NOT NULL,
  expires_at TEXT NOT NULL
);
`
	_, err := db.Exec(schema)
	return err
}

func newID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func timeToText(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339Nano)
}

func timeFromText(v string) time.Time {
	if v == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339Nano, v)
	if err != nil {
		return time.Time{}
	}
	return t
}

func normalizePhone(countryCode, phone string) string {
	cc := strings.TrimSpace(countryCode)
	p := strings.TrimSpace(phone)
	if p == "" {
		return ""
	}
	if strings.HasPrefix(p, "+") {
		return p
	}
	return cc + p
}

func hashPassword(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (s *Store) Login(req AuthLoginRequest) (AuthLoginResponse, error) {
	method := strings.ToLower(strings.TrimSpace(req.Method))
	passwordHash := hashPassword(strings.TrimSpace(req.Password))

	var u staffUser
	var found bool

	if s.db != nil {
		var err error
		u, found, err = s.loginFromSQLite(method, req, passwordHash)
		if err != nil {
			return AuthLoginResponse{}, err
		}
	} else {
		u, found = s.loginFromMemory(method, req, passwordHash)
	}

	if !found {
		return AuthLoginResponse{}, errors.New("invalid credentials")
	}

	session := SessionInfo{
		SessionID:    newID("sess"),
		UserID:       u.ID,
		TenantID:     u.TenantID,
		MerchantType: u.MerchantType,
		Role:         u.Role,
		DisplayName:  u.Name,
		Method:       method,
	}

	if s.db != nil {
		now := time.Now()
		_, err := s.db.Exec(
			`INSERT INTO sessions (id, user_id, tenant_id, merchant_type, role, display_name, method, created_at, expires_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			session.SessionID, session.UserID, session.TenantID, session.MerchantType, session.Role, session.DisplayName, session.Method,
			timeToText(now), timeToText(now.Add(24*time.Hour)),
		)
		if err != nil {
			return AuthLoginResponse{}, fmt.Errorf("create session: %w", err)
		}
	} else {
		s.mu.Lock()
		s.sessions[session.SessionID] = session
		s.mu.Unlock()
	}

	return AuthLoginResponse{
		SessionID: session.SessionID,
		User: AuthSessionUser{
			ID:           session.UserID,
			TenantID:     session.TenantID,
			MerchantType: session.MerchantType,
			Role:         session.Role,
			DisplayName:  session.DisplayName,
			Method:       session.Method,
		},
	}, nil
}

func (s *Store) loginFromSQLite(method string, req AuthLoginRequest, passwordHash string) (staffUser, bool, error) {
	var (
		q    string
		args []any
	)

	switch method {
	case "email":
		email := strings.TrimSpace(req.Email)
		if email == "" || strings.TrimSpace(req.Password) == "" {
			return staffUser{}, false, nil
		}
		q = `SELECT id, tenant_id, merchant_type, role, name, email, phone, auth_provider, password_hash
			 FROM staff_users WHERE email = ? LIMIT 1`
		args = []any{email}
	case "phone":
		phone := normalizePhone(req.CountryCode, req.Phone)
		if phone == "" || strings.TrimSpace(req.Password) == "" {
			return staffUser{}, false, nil
		}
		q = `SELECT id, tenant_id, merchant_type, role, name, email, phone, auth_provider, password_hash
			 FROM staff_users WHERE phone = ? LIMIT 1`
		args = []any{phone}
	case "google":
		email := strings.TrimSpace(req.Email)
		if email != "" {
			q = `SELECT id, tenant_id, merchant_type, role, name, email, phone, auth_provider, password_hash
				 FROM staff_users WHERE email = ? AND auth_provider = 'google' LIMIT 1`
			args = []any{email}
		} else {
			q = `SELECT id, tenant_id, merchant_type, role, name, email, phone, auth_provider, password_hash
				 FROM staff_users WHERE auth_provider = 'google' LIMIT 1`
		}
	default:
		return staffUser{}, false, nil
	}

	row := s.db.QueryRow(q, args...)
	var u staffUser
	if err := row.Scan(&u.ID, &u.TenantID, &u.MerchantType, &u.Role, &u.Name, &u.Email, &u.Phone, &u.AuthProvider, &u.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return staffUser{}, false, nil
		}
		return staffUser{}, false, err
	}

	if method != "google" && u.PasswordHash != passwordHash {
		return staffUser{}, false, nil
	}
	if method == "google" && u.AuthProvider != "google" {
		return staffUser{}, false, nil
	}
	return u, true, nil
}

func (s *Store) loginFromMemory(method string, req AuthLoginRequest, passwordHash string) (staffUser, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, u := range s.users {
		switch method {
		case "email":
			if strings.EqualFold(u.Email, strings.TrimSpace(req.Email)) && u.PasswordHash == passwordHash {
				return u, true
			}
		case "phone":
			if u.Phone == normalizePhone(req.CountryCode, req.Phone) && u.PasswordHash == passwordHash {
				return u, true
			}
		case "google":
			if u.AuthProvider == "google" {
				if strings.TrimSpace(req.Email) == "" || strings.EqualFold(u.Email, strings.TrimSpace(req.Email)) {
					return u, true
				}
			}
		}
	}
	return staffUser{}, false
}

func (s *Store) GetSession(sessionID string) (SessionInfo, error) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return SessionInfo{}, errors.New("missing session id")
	}

	if s.db != nil {
		row := s.db.QueryRow(
			`SELECT id, user_id, tenant_id, merchant_type, role, display_name, method, expires_at
			 FROM sessions WHERE id = ? LIMIT 1`,
			sessionID,
		)
		var (
			out       SessionInfo
			expiresAt string
		)
		if err := row.Scan(&out.SessionID, &out.UserID, &out.TenantID, &out.MerchantType, &out.Role, &out.DisplayName, &out.Method, &expiresAt); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return SessionInfo{}, errors.New("invalid session")
			}
			return SessionInfo{}, err
		}
		if exp := timeFromText(expiresAt); !exp.IsZero() && time.Now().After(exp) {
			return SessionInfo{}, errors.New("session expired")
		}
		return out, nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	out, ok := s.sessions[sessionID]
	if !ok {
		return SessionInfo{}, errors.New("invalid session")
	}
	return out, nil
}

func (s *Store) IsVisitInTenant(visitID, tenantID string) bool {
	if visitID == "" || tenantID == "" {
		return false
	}
	if s.db != nil {
		row := s.db.QueryRow(`SELECT COUNT(1) FROM visits WHERE id = ? AND tenant_id = ?`, visitID, tenantID)
		var c int
		if err := row.Scan(&c); err != nil {
			return false
		}
		return c > 0
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.visits[visitID]
	return ok && v.TenantID == tenantID
}

func (s *Store) CreateAppointment(in Appointment) Appointment {
	s.mu.Lock()
	defer s.mu.Unlock()
	in.ID = newID("appt")
	if in.Status == "" {
		in.Status = "pending"
	}
	if in.ScheduledAt.IsZero() {
		in.ScheduledAt = time.Now().Add(1 * time.Hour)
	}

	if s.db != nil {
		_, err := s.db.Exec(
			`INSERT INTO appointments (id, tenant_id, pet_id, doctor_id, scheduled_at, status, chief_complaint)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			in.ID, in.TenantID, in.PetID, in.DoctorID, timeToText(in.ScheduledAt), in.Status, in.ChiefComplaint,
		)
		if err == nil {
			return in
		}
		log.Printf("sqlite CreateAppointment failed, fallback to memory: %v", err)
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

	if s.db != nil {
		_, err := s.db.Exec(
			`INSERT INTO orders (id, tenant_id, customer_id, status, amount_hkd, created_at)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			in.ID, in.TenantID, in.CustomerID, in.Status, in.AmountHKD, timeToText(in.CreatedAt),
		)
		if err == nil {
			return in
		}
		log.Printf("sqlite CreateOrder failed, fallback to memory: %v", err)
	}

	s.orders[in.ID] = in
	return in
}

func (s *Store) ListOrders(tenantID string) []Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.db != nil {
		rows, err := s.db.Query(
			`SELECT id, tenant_id, customer_id, status, amount_hkd, created_at
			 FROM orders WHERE tenant_id = ? ORDER BY created_at DESC`,
			tenantID,
		)
		if err == nil {
			defer rows.Close()
			out := make([]Order, 0)
			for rows.Next() {
				var o Order
				var createdAt string
				if scanErr := rows.Scan(&o.ID, &o.TenantID, &o.CustomerID, &o.Status, &o.AmountHKD, &createdAt); scanErr != nil {
					continue
				}
				o.CreatedAt = timeFromText(createdAt)
				out = append(out, o)
			}
			if rows.Err() == nil {
				return out
			}
		} else {
			log.Printf("sqlite ListOrders failed, fallback to memory: %v", err)
		}
	}

	out := make([]Order, 0, len(s.orders))
	for _, v := range s.orders {
		if v.TenantID == tenantID || tenantID == "" {
			out = append(out, v)
		}
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

	if s.db != nil {
		_, err := s.db.Exec(
			`INSERT INTO products (id, tenant_id, name, sku, stock, price_hkd, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			in.ID, in.TenantID, in.Name, in.SKU, in.Stock, in.PriceHKD, timeToText(in.UpdatedAt),
		)
		if err == nil {
			return in
		}
		log.Printf("sqlite CreateProduct failed, fallback to memory: %v", err)
	}

	s.products[in.ID] = in
	return in
}

func (s *Store) ListProducts(tenantID string) []Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.db != nil {
		rows, err := s.db.Query(
			`SELECT id, tenant_id, name, sku, stock, price_hkd, updated_at
			 FROM products WHERE tenant_id = ? ORDER BY updated_at DESC`,
			tenantID,
		)
		if err == nil {
			defer rows.Close()
			out := make([]Product, 0)
			for rows.Next() {
				var p Product
				var updatedAt string
				if scanErr := rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.SKU, &p.Stock, &p.PriceHKD, &updatedAt); scanErr != nil {
					continue
				}
				p.UpdatedAt = timeFromText(updatedAt)
				out = append(out, p)
			}
			if rows.Err() == nil {
				return out
			}
		} else {
			log.Printf("sqlite ListProducts failed, fallback to memory: %v", err)
		}
	}

	out := make([]Product, 0, len(s.products))
	for _, v := range s.products {
		if v.TenantID == tenantID || tenantID == "" {
			out = append(out, v)
		}
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

	if s.db != nil {
		_, err := s.db.Exec(
			`INSERT INTO customers (id, tenant_id, name, email, phone, created_at, last_active)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			in.ID, in.TenantID, in.Name, in.Email, in.Phone, timeToText(in.CreatedAt), timeToText(in.LastActive),
		)
		if err == nil {
			return in
		}
		log.Printf("sqlite CreateCustomer failed, fallback to memory: %v", err)
	}

	s.customers[in.ID] = in
	return in
}

func (s *Store) ListCustomers(tenantID string) []Customer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.db != nil {
		rows, err := s.db.Query(
			`SELECT id, tenant_id, name, email, phone, created_at, last_active
			 FROM customers WHERE tenant_id = ? ORDER BY created_at DESC`,
			tenantID,
		)
		if err == nil {
			defer rows.Close()
			out := make([]Customer, 0)
			for rows.Next() {
				var c Customer
				var createdAt, lastActive string
				if scanErr := rows.Scan(&c.ID, &c.TenantID, &c.Name, &c.Email, &c.Phone, &createdAt, &lastActive); scanErr != nil {
					continue
				}
				c.CreatedAt = timeFromText(createdAt)
				c.LastActive = timeFromText(lastActive)
				out = append(out, c)
			}
			if rows.Err() == nil {
				return out
			}
		} else {
			log.Printf("sqlite ListCustomers failed, fallback to memory: %v", err)
		}
	}

	out := make([]Customer, 0, len(s.customers))
	for _, v := range s.customers {
		if v.TenantID == tenantID || tenantID == "" {
			out = append(out, v)
		}
	}
	return out
}

func (s *Store) ListAppointments(tenantID string) []Appointment {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.db != nil {
		rows, err := s.db.Query(
			`SELECT id, tenant_id, pet_id, doctor_id, scheduled_at, status, chief_complaint
			 FROM appointments WHERE tenant_id = ? ORDER BY scheduled_at DESC`,
			tenantID,
		)
		if err == nil {
			defer rows.Close()
			out := make([]Appointment, 0)
			for rows.Next() {
				var a Appointment
				var scheduledAt string
				if scanErr := rows.Scan(&a.ID, &a.TenantID, &a.PetID, &a.DoctorID, &scheduledAt, &a.Status, &a.ChiefComplaint); scanErr != nil {
					continue
				}
				a.ScheduledAt = timeFromText(scheduledAt)
				out = append(out, a)
			}
			if rows.Err() == nil {
				return out
			}
		} else {
			log.Printf("sqlite ListAppointments failed, fallback to memory: %v", err)
		}
	}

	out := make([]Appointment, 0, len(s.appointments))
	for _, v := range s.appointments {
		if v.TenantID == tenantID || tenantID == "" {
			out = append(out, v)
		}
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

	if s.db != nil {
		_, err := s.db.Exec(
			`INSERT INTO visits (id, tenant_id, appointment_id, pet_id, doctor_id, status, checkin_at, checkout_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			in.ID, in.TenantID, in.AppointmentID, in.PetID, in.DoctorID, in.Status, timeToText(in.CheckInAt), timeToText(in.CheckOutAt),
		)
		if err == nil {
			return in
		}
		log.Printf("sqlite CreateVisit failed, fallback to memory: %v", err)
	}

	s.visits[in.ID] = in
	return in
}

func (s *Store) ListVisits(tenantID string) []Visit {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.db != nil {
		rows, err := s.db.Query(
			`SELECT id, tenant_id, appointment_id, pet_id, doctor_id, status, checkin_at, checkout_at
			 FROM visits WHERE tenant_id = ? ORDER BY checkin_at DESC`,
			tenantID,
		)
		if err == nil {
			defer rows.Close()
			out := make([]Visit, 0)
			for rows.Next() {
				var v Visit
				var checkInAt, checkOutAt sql.NullString
				if scanErr := rows.Scan(&v.ID, &v.TenantID, &v.AppointmentID, &v.PetID, &v.DoctorID, &v.Status, &checkInAt, &checkOutAt); scanErr != nil {
					continue
				}
				if checkInAt.Valid {
					v.CheckInAt = timeFromText(checkInAt.String)
				}
				if checkOutAt.Valid {
					v.CheckOutAt = timeFromText(checkOutAt.String)
				}
				out = append(out, v)
			}
			if rows.Err() == nil {
				return out
			}
		} else {
			log.Printf("sqlite ListVisits failed, fallback to memory: %v", err)
		}
	}

	out := make([]Visit, 0, len(s.visits))
	for _, v := range s.visits {
		if v.TenantID == tenantID || tenantID == "" {
			out = append(out, v)
		}
	}
	return out
}

func (s *Store) CreatePrescription(in Prescription) Prescription {
	s.mu.Lock()
	defer s.mu.Unlock()
	in.ID = newID("rx")

	if s.db != nil {
		_, err := s.db.Exec(
			`INSERT INTO prescriptions (id, visit_id, medicine_name, dosage, frequency, duration_days, notes)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			in.ID, in.VisitID, in.MedicineName, in.Dosage, in.Frequency, in.DurationDays, in.Notes,
		)
		if err == nil {
			return in
		}
		log.Printf("sqlite CreatePrescription failed, fallback to memory: %v", err)
	}

	s.prescriptions[in.ID] = in
	return in
}

func (s *Store) ListPrescriptions(tenantID string) []Prescription {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.db != nil {
		rows, err := s.db.Query(
			`SELECT p.id, p.visit_id, p.medicine_name, p.dosage, p.frequency, p.duration_days, p.notes
			 FROM prescriptions p
			 JOIN visits v ON p.visit_id = v.id
			 WHERE v.tenant_id = ?
			 ORDER BY p.id DESC`,
			tenantID,
		)
		if err == nil {
			defer rows.Close()
			out := make([]Prescription, 0)
			for rows.Next() {
				var p Prescription
				if scanErr := rows.Scan(&p.ID, &p.VisitID, &p.MedicineName, &p.Dosage, &p.Frequency, &p.DurationDays, &p.Notes); scanErr != nil {
					continue
				}
				out = append(out, p)
			}
			if rows.Err() == nil {
				return out
			}
		} else {
			log.Printf("sqlite ListPrescriptions failed, fallback to memory: %v", err)
		}
	}

	out := make([]Prescription, 0, len(s.prescriptions))
	for _, p := range s.prescriptions {
		visit, ok := s.visits[p.VisitID]
		if ok && (visit.TenantID == tenantID || tenantID == "") {
			out = append(out, p)
		}
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

	if s.db != nil {
		_, err := s.db.Exec(
			`INSERT INTO followups (id, visit_id, due_at, channel, status, result_notes)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			in.ID, in.VisitID, timeToText(in.DueAt), in.Channel, in.Status, in.ResultNotes,
		)
		if err == nil {
			return in
		}
		log.Printf("sqlite CreateFollowUp failed, fallback to memory: %v", err)
	}

	s.followUps[in.ID] = in
	return in
}

func (s *Store) ListFollowUps(tenantID string) []FollowUp {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.db != nil {
		rows, err := s.db.Query(
			`SELECT f.id, f.visit_id, f.due_at, f.channel, f.status, f.result_notes
			 FROM followups f
			 JOIN visits v ON f.visit_id = v.id
			 WHERE v.tenant_id = ?
			 ORDER BY f.id DESC`,
			tenantID,
		)
		if err == nil {
			defer rows.Close()
			out := make([]FollowUp, 0)
			for rows.Next() {
				var f FollowUp
				var dueAt sql.NullString
				if scanErr := rows.Scan(&f.ID, &f.VisitID, &dueAt, &f.Channel, &f.Status, &f.ResultNotes); scanErr != nil {
					continue
				}
				if dueAt.Valid {
					f.DueAt = timeFromText(dueAt.String)
				}
				out = append(out, f)
			}
			if rows.Err() == nil {
				return out
			}
		} else {
			log.Printf("sqlite ListFollowUps failed, fallback to memory: %v", err)
		}
	}

	out := make([]FollowUp, 0, len(s.followUps))
	for _, f := range s.followUps {
		visit, ok := s.visits[f.VisitID]
		if ok && (visit.TenantID == tenantID || tenantID == "") {
			out = append(out, f)
		}
	}
	return out
}

func (s *Store) SeedDemoData() {
	if len(s.ListOrders("tenant_shop_demo")) > 0 || len(s.ListAppointments("tenant_clinic_demo")) > 0 {
		return
	}

	s.seedDemoUsers()

	s.CreateCustomer(Customer{
		TenantID: "tenant_shop_demo",
		Name:     "Amy Chan",
		Email:    "amy@example.com",
		Phone:    "+85290000001",
	})
	s.CreateCustomer(Customer{
		TenantID: "tenant_shop_demo",
		Name:     "Ken Wong",
		Email:    "ken@example.com",
		Phone:    "+85290000002",
	})

	s.CreateProduct(Product{
		TenantID: "tenant_shop_demo",
		Name:     "Omega-3 Supplement",
		SKU:      "PW-OMG3-01",
		Stock:    18,
		PriceHKD: 198,
	})
	s.CreateProduct(Product{
		TenantID: "tenant_shop_demo",
		Name:     "Urinary Cat Food",
		SKU:      "PW-UR-08",
		Stock:    5,
		PriceHKD: 320,
	})

	s.CreateOrder(Order{
		TenantID:   "tenant_shop_demo",
		CustomerID: "seed_customer_1",
		Status:     "paid",
		AmountHKD:  398,
	})
	s.CreateOrder(Order{
		TenantID:   "tenant_shop_demo",
		CustomerID: "seed_customer_2",
		Status:     "pending",
		AmountHKD:  722,
	})

	s.CreateAppointment(Appointment{
		TenantID:       "tenant_clinic_demo",
		PetID:          "pet_001",
		DoctorID:       "doc_001",
		ScheduledAt:    time.Now().Add(2 * time.Hour),
		Status:         "confirmed",
		ChiefComplaint: "skin itching",
	})
	v := s.CreateVisit(Visit{
		TenantID: "tenant_clinic_demo",
		PetID:    "pet_008",
		DoctorID: "doc_002",
		Status:   "in_progress",
	})
	s.CreatePrescription(Prescription{
		VisitID:      v.ID,
		MedicineName: "Probiotic",
		Dosage:       "1 capsule",
		Frequency:    "BID",
		DurationDays: 7,
		Notes:        "After meal",
	})
}

func (s *Store) seedDemoUsers() {
	users := []staffUser{
		{
			ID:           "usr_shop_admin",
			TenantID:     "tenant_shop_demo",
			MerchantType: "shop",
			Role:         "admin",
			Name:         "Shop Admin",
			Email:        "shop_admin@petwell.com",
			Phone:        "+85290000011",
			AuthProvider: "password",
			PasswordHash: hashPassword("Shop123456"),
		},
		{
			ID:           "usr_clinic_admin",
			TenantID:     "tenant_clinic_demo",
			MerchantType: "clinic",
			Role:         "admin",
			Name:         "Clinic Admin",
			Email:        "clinic_admin@petwell.com",
			Phone:        "+85290000022",
			AuthProvider: "password",
			PasswordHash: hashPassword("Clinic123456"),
		},
		{
			ID:           "usr_google_shop",
			TenantID:     "tenant_shop_demo",
			MerchantType: "shop",
			Role:         "admin",
			Name:         "Google Shop Admin",
			Email:        "google_shop@petwell.com",
			Phone:        "+85290000033",
			AuthProvider: "google",
			PasswordHash: "",
		},
	}

	if s.db != nil {
		for _, u := range users {
			_, _ = s.db.Exec(
				`INSERT OR IGNORE INTO staff_users (id, tenant_id, merchant_type, role, name, email, phone, auth_provider, password_hash, created_at)
				 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				u.ID, u.TenantID, u.MerchantType, u.Role, u.Name, u.Email, u.Phone, u.AuthProvider, u.PasswordHash, timeToText(time.Now()),
			)
		}
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range users {
		s.users[u.ID] = u
	}
}
