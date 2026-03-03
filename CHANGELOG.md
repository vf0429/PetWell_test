# Changelog

All notable changes to this project are documented here.

## 2026-03-02

### Added
- Added OpenAPI initial contract at `backend/openapi.yaml`.
- Added structured API error response shape with `error.code`, `error.message`, and `request_id`.
- Added `X-Request-ID` request/response header support for traceability.
- Added environment profile config loader (`dev/staging/prod`) in `backend/internal/config.go`.
- Added config unit tests in `backend/internal/config_test.go`.

### Changed
- Upgraded `GET /health` to include store type, DB reachability, and schema version metadata.
- Upgraded request logging to include `request_id`, `tenant_id`, `user_id`, `status`, and `duration_ms`.
- Updated CORS headers to allow/expose `X-Request-ID`.

## 2026-03-01

### Added
- Added database-backed authentication and session model (`staff_users`, `sessions`).
- Added strict tenant isolation for business APIs using `X-Session-ID`.
- Added multilingual (EN/ZH) UI toggle and localized table headers.
- Added login UX updates: password fields for email/phone and demo account hint.
- Added final test spec document: `09_最终测试稿_登录与租户隔离.md`.
- Added documentation index: `README.md`.

### Changed
- Switched backend default store to SQLite persistence.
- Updated API smoke script to include login and authenticated requests.
- Updated CORS allow headers to include `X-Session-ID`.
- Updated project docs (`00`, `03`, `05`, `07`, `08`) to reflect current architecture and testing flow.
- Refactored docs `01` and `04` into table-style quick-reference format.

### Notes
- Google login is currently demo-mode (not full OAuth flow).
- Shopify integration remains placeholder-level (OAuth/Webhook endpoints are stubs).
