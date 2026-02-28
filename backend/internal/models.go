package internal

import "time"

type Appointment struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	PetID          string    `json:"pet_id"`
	DoctorID       string    `json:"doctor_id"`
	ScheduledAt    time.Time `json:"scheduled_at"`
	Status         string    `json:"status"`
	ChiefComplaint string    `json:"chief_complaint"`
}

type Visit struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	AppointmentID string    `json:"appointment_id"`
	PetID         string    `json:"pet_id"`
	DoctorID      string    `json:"doctor_id"`
	Status        string    `json:"status"`
	CheckInAt     time.Time `json:"checkin_at"`
	CheckOutAt    time.Time `json:"checkout_at,omitempty"`
}

type Prescription struct {
	ID           string `json:"id"`
	VisitID      string `json:"visit_id"`
	MedicineName string `json:"medicine_name"`
	Dosage       string `json:"dosage"`
	Frequency    string `json:"frequency"`
	DurationDays int    `json:"duration_days"`
	Notes        string `json:"notes"`
}

type FollowUp struct {
	ID          string    `json:"id"`
	VisitID     string    `json:"visit_id"`
	DueAt       time.Time `json:"due_at"`
	Channel     string    `json:"channel"`
	Status      string    `json:"status"`
	ResultNotes string    `json:"result_notes"`
}

type ShopifyOAuthRequest struct {
	ShopDomain string `json:"shop_domain"`
}

type ShopifyWebhookEvent struct {
	Topic   string                 `json:"topic"`
	Payload map[string]interface{} `json:"payload"`
}

type Order struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	CustomerID string    `json:"customer_id"`
	Status     string    `json:"status"`
	AmountHKD  float64   `json:"amount_hkd"`
	CreatedAt  time.Time `json:"created_at"`
}

type Product struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	SKU       string    `json:"sku"`
	Stock     int       `json:"stock"`
	PriceHKD  float64   `json:"price_hkd"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Customer struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	CreatedAt  time.Time `json:"created_at"`
	LastActive time.Time `json:"last_active"`
}
