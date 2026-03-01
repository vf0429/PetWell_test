#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8090}"
LOGIN_EMAIL="${LOGIN_EMAIL:-shop_admin@petwell.com}"
LOGIN_PASSWORD="${LOGIN_PASSWORD:-Shop123456}"

if ! command -v jq >/dev/null 2>&1; then
  echo "jq is required. Please install jq first."
  exit 1
fi

echo "== Health =="
curl -sS "$BASE_URL/health" | jq .

echo "== Login =="
LOGIN_RESP=$(curl -sS -X POST "$BASE_URL/api/merchant/auth/login" \
  -H 'Content-Type: application/json' \
  -d "{\"method\":\"email\",\"email\":\"$LOGIN_EMAIL\",\"password\":\"$LOGIN_PASSWORD\"}")
echo "$LOGIN_RESP" | jq .
SESSION_ID=$(echo "$LOGIN_RESP" | jq -r '.session_id')
if [[ -z "$SESSION_ID" || "$SESSION_ID" == "null" ]]; then
  echo "Login failed: session_id is empty"
  exit 1
fi
AUTH_HEADER="X-Session-ID: $SESSION_ID"

echo "== Create customer =="
CUSTOMER=$(curl -sS -X POST "$BASE_URL/api/merchant/customers" \
  -H "$AUTH_HEADER" \
  -H 'Content-Type: application/json' \
  -d '{
    "name":"Smoke Owner",
    "email":"owner-smoke@test.com",
    "phone":"+85291231234"
  }')
echo "$CUSTOMER" | jq .
CUSTOMER_ID=$(echo "$CUSTOMER" | jq -r '.id')

echo "== Create product =="
PRODUCT=$(curl -sS -X POST "$BASE_URL/api/merchant/products" \
  -H "$AUTH_HEADER" \
  -H 'Content-Type: application/json' \
  -d '{
    "name":"Digestive Treat",
    "sku":"PW-DIG-01",
    "stock":22,
    "price_hkd":188
  }')
echo "$PRODUCT" | jq .

echo "== Create order =="
ORDER=$(curl -sS -X POST "$BASE_URL/api/merchant/orders" \
  -H "$AUTH_HEADER" \
  -H 'Content-Type: application/json' \
  -d "{
    \"customer_id\":\"$CUSTOMER_ID\",
    \"status\":\"paid\",
    \"amount_hkd\":188
  }")
echo "$ORDER" | jq .

echo "== Create appointment =="
APPT=$(curl -sS -X POST "$BASE_URL/api/merchant/appointments" \
  -H "$AUTH_HEADER" \
  -H 'Content-Type: application/json' \
  -d '{
    "pet_id":"pet_123",
    "doctor_id":"doc_1",
    "scheduled_at":"2026-03-01T10:00:00Z",
    "chief_complaint":"vomiting"
  }')
echo "$APPT" | jq .

echo "== Create visit =="
VISIT=$(curl -sS -X POST "$BASE_URL/api/merchant/visits" \
  -H "$AUTH_HEADER" \
  -H 'Content-Type: application/json' \
  -d '{
    "pet_id":"pet_123",
    "doctor_id":"doc_1",
    "appointment_id":"manual_link"
  }')
echo "$VISIT" | jq .
VISIT_ID=$(echo "$VISIT" | jq -r '.id')

echo "== Create prescription =="
curl -sS -X POST "$BASE_URL/api/merchant/prescriptions" \
  -H "$AUTH_HEADER" \
  -H 'Content-Type: application/json' \
  -d "{\"visit_id\":\"$VISIT_ID\",\"medicine_name\":\"Antibiotic\",\"dosage\":\"5mg\",\"frequency\":\"BID\",\"duration_days\":5}" | jq .

echo "== Create follow-up =="
curl -sS -X POST "$BASE_URL/api/merchant/followups" \
  -H "$AUTH_HEADER" \
  -H 'Content-Type: application/json' \
  -d "{\"visit_id\":\"$VISIT_ID\",\"due_at\":\"2026-03-03T09:00:00Z\",\"channel\":\"phone\"}" | jq .

echo "== Shopify placeholders =="
curl -sS -X POST "$BASE_URL/api/merchant/shopify/oauth/start" \
  -H 'Content-Type: application/json' \
  -d '{"shop_domain":"demo-store.myshopify.com"}' | jq .
curl -sS -X POST "$BASE_URL/api/merchant/shopify/webhooks" \
  -H 'Content-Type: application/json' \
  -d '{"topic":"orders/create","payload":{"id":"1001"}}' | jq .

echo "== List endpoints (tenant-filtered) =="
curl -sS "$BASE_URL/api/merchant/dashboard?mode=shop" -H "$AUTH_HEADER" | jq .
curl -sS "$BASE_URL/api/merchant/dashboard?mode=clinic" -H "$AUTH_HEADER" | jq .
curl -sS "$BASE_URL/api/merchant/orders" -H "$AUTH_HEADER" | jq .
curl -sS "$BASE_URL/api/merchant/products" -H "$AUTH_HEADER" | jq .
curl -sS "$BASE_URL/api/merchant/customers" -H "$AUTH_HEADER" | jq .
curl -sS "$BASE_URL/api/merchant/appointments" -H "$AUTH_HEADER" | jq .
curl -sS "$BASE_URL/api/merchant/visits" -H "$AUTH_HEADER" | jq .
curl -sS "$BASE_URL/api/merchant/prescriptions" -H "$AUTH_HEADER" | jq .
curl -sS "$BASE_URL/api/merchant/followups" -H "$AUTH_HEADER" | jq .

echo "Smoke tests done."
