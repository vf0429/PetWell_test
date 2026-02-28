#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8090}"

echo "== Health =="
curl -sS "$BASE_URL/health" | jq .

echo "== Create appointment =="
APPT=$(curl -sS -X POST "$BASE_URL/api/merchant/appointments" \
  -H 'Content-Type: application/json' \
  -d '{
    "tenant_id":"tenant_demo",
    "pet_id":"pet_123",
    "doctor_id":"doc_1",
    "scheduled_at":"2026-03-01T10:00:00Z",
    "chief_complaint":"vomiting"
  }')
echo "$APPT" | jq .

echo "== Create customer =="
CUSTOMER=$(curl -sS -X POST "$BASE_URL/api/merchant/customers" \
  -H 'Content-Type: application/json' \
  -d '{
    "tenant_id":"tenant_demo",
    "name":"Test Owner",
    "email":"owner@test.com",
    "phone":"+85291231234"
  }')
echo "$CUSTOMER" | jq .

CUSTOMER_ID=$(echo "$CUSTOMER" | jq -r '.id')

echo "== Create product =="
PRODUCT=$(curl -sS -X POST "$BASE_URL/api/merchant/products" \
  -H 'Content-Type: application/json' \
  -d '{
    "tenant_id":"tenant_demo",
    "name":"Digestive Treat",
    "sku":"PW-DIG-01",
    "stock":22,
    "price_hkd":188
  }')
echo "$PRODUCT" | jq .

echo "== Create order =="
ORDER=$(curl -sS -X POST "$BASE_URL/api/merchant/orders" \
  -H 'Content-Type: application/json' \
  -d "{
    \"tenant_id\":\"tenant_demo\",
    \"customer_id\":\"$CUSTOMER_ID\",
    \"status\":\"paid\",
    \"amount_hkd\":188
  }")
echo "$ORDER" | jq .

echo "== Create visit =="
VISIT=$(curl -sS -X POST "$BASE_URL/api/merchant/visits" \
  -H 'Content-Type: application/json' \
  -d '{
    "tenant_id":"tenant_demo",
    "pet_id":"pet_123",
    "doctor_id":"doc_1",
    "appointment_id":"manual_link"
  }')
echo "$VISIT" | jq .

VISIT_ID=$(echo "$VISIT" | jq -r '.id')

echo "== Create prescription =="
curl -sS -X POST "$BASE_URL/api/merchant/prescriptions" \
  -H 'Content-Type: application/json' \
  -d "{\"visit_id\":\"$VISIT_ID\",\"medicine_name\":\"Antibiotic\",\"dosage\":\"5mg\",\"frequency\":\"BID\",\"duration_days\":5}" | jq .

echo "== Create follow-up =="
curl -sS -X POST "$BASE_URL/api/merchant/followups" \
  -H 'Content-Type: application/json' \
  -d "{\"visit_id\":\"$VISIT_ID\",\"due_at\":\"2026-03-03T09:00:00Z\",\"channel\":\"phone\"}" | jq .

echo "== Shopify OAuth placeholder =="
curl -sS -X POST "$BASE_URL/api/merchant/shopify/oauth/start" \
  -H 'Content-Type: application/json' \
  -d '{"shop_domain":"demo-store.myshopify.com"}' | jq .

echo "== Shopify webhook placeholder =="
curl -sS -X POST "$BASE_URL/api/merchant/shopify/webhooks" \
  -H 'Content-Type: application/json' \
  -d '{"topic":"orders/create","payload":{"id":"1001"}}' | jq .

echo "== List endpoints =="
curl -sS "$BASE_URL/api/merchant/dashboard?mode=shop" | jq .
curl -sS "$BASE_URL/api/merchant/dashboard?mode=clinic" | jq .
curl -sS "$BASE_URL/api/merchant/orders" | jq .
curl -sS "$BASE_URL/api/merchant/products" | jq .
curl -sS "$BASE_URL/api/merchant/customers" | jq .
curl -sS "$BASE_URL/api/merchant/appointments" | jq .
curl -sS "$BASE_URL/api/merchant/visits" | jq .
curl -sS "$BASE_URL/api/merchant/prescriptions" | jq .
curl -sS "$BASE_URL/api/merchant/followups" | jq .

echo "Smoke tests done."
