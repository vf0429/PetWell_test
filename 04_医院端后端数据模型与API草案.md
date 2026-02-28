# 医院商家端：后端数据模型与 API 草案（V0.1）

## 1) 数据库表（首期）

## 多租户核心
- tenants(id, name, type[shop/clinic], status)
- staff_users(id, tenant_id, role, name, email, phone, password_hash)

## 宠物与主人关联
- pets(id, owner_user_id, name, species, breed, sex, dob, weight)
- pet_profiles(id, pet_id, allergies, chronic_conditions, vaccination_status)

## 预约就诊
- appointments(id, tenant_id, pet_id, doctor_id, scheduled_at, status, chief_complaint, source)
- visits(id, tenant_id, appointment_id, pet_id, doctor_id, checkin_at, checkout_at, status)
- diagnoses(id, visit_id, code, name, notes)
- treatments(id, visit_id, item_type, item_name, qty, unit_price)
- prescriptions(id, visit_id, medicine_name, dosage, frequency, duration_days, notes)
- visit_files(id, visit_id, file_type, file_url, uploaded_by)
- followups(id, visit_id, due_at, channel, status, result_notes)

## 审计与同步
- event_logs(id, tenant_id, entity_type, entity_id, action, actor_id, payload, created_at)
- app_sync_queue(id, event_type, payload, status, retry_count)

## 2) 状态流转
- appointment.status: pending -> confirmed -> checked_in -> completed / cancelled
- visit.status: in_progress -> diagnosed -> treated -> closed
- followup.status: pending -> done / skipped

## 3) API 草案（REST）

### 认证与组织
- POST /merchant/auth/login
- GET /merchant/me
- GET /merchant/tenants/:id

### 预约
- GET /merchant/appointments
- POST /merchant/appointments
- PATCH /merchant/appointments/:id/status

### 就诊
- GET /merchant/visits
- POST /merchant/visits
- GET /merchant/visits/:id
- PATCH /merchant/visits/:id

### 诊断/处方/文件
- POST /merchant/visits/:id/diagnoses
- POST /merchant/visits/:id/prescriptions
- POST /merchant/visits/:id/files

### 回访
- GET /merchant/followups
- POST /merchant/followups
- PATCH /merchant/followups/:id/status

### App 同步
- POST /internal/sync/pet-profile
- POST /internal/sync/visit-summary

## 4) 数据安全与合规
- 医疗文件需加访问控制（签名 URL）
- 关键表需审计日志
- 不同租户数据严格隔离（tenant_id 强约束）

## 5) 首期可交付最小集（MVP）
1. 预约管理
2. 就诊记录创建与更新
3. 处方记录
4. 病历文件上传
5. 同步摘要到 PetWell App

