# PetWell 商家端 <-> App 联调说明（V0.1）

## 1) 联调范围
- Shop：订单状态、发货状态、售后状态同步到 App
- Clinic：预约状态、就诊摘要、处方摘要、回访任务同步到 App

## 2) 字段映射（首期）

## Shop 订单
- merchant_order.id -> app_order.external_order_id
- merchant_order.status -> app_order.status
- merchant_order.fulfillment_status -> app_order.shipping_status
- merchant_order.updated_at -> app_order.synced_at

## Clinic 就诊摘要
- visit.id -> app_medical_record.visit_id
- visit.status -> app_medical_record.status
- diagnoses[].name -> app_medical_record.primary_diagnosis
- prescriptions[] -> app_medical_record.meds_summary
- followups[].due_at -> app_medical_record.next_followup_at

## 3) 时序建议
1. 商家端更新状态
2. 商家端后端写 event_logs
3. 推送 app_sync_queue
4. 同步服务消费队列，调用 App 内部同步接口
5. App 返回确认，更新队列状态

## 4) 幂等规则
- 以 (entity_type, entity_id, action, updated_at) 做去重键
- 同步接口支持幂等请求头：X-Idempotency-Key

## 5) 错误重试
- 首次失败：30s
- 第二次失败：2min
- 第三次失败：10min
- 超过阈值告警并进入人工处理队列

## 6) 联调检查项
- App 是否正确展示商家端最新状态
- 同一事件重复推送是否去重
- 跨租户数据是否隔离

