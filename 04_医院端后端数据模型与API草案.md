# 医院商家端：后端数据模型与 API 草案（V0.2）

> 文档用途：给后端、前端、测试统一医院端数据与接口认知。  
> 当前实现说明：`appointments / visits / prescriptions / followups` 已落地；其余为规划。

## 1) 数据模型总览

| 模块 | 表/实体 | 用途 | 当前状态 |
|---|---|---|---|
| 多租户 | `tenants` | 租户基础信息 | 规划中 |
| 账号权限 | `staff_users` | 医院账号与角色 | 已落地（简化版） |
| 会话 | `sessions` | 登录会话与鉴权 | 已落地 |
| 宠物资料 | `pets`, `pet_profiles` | 宠物档案与健康信息 | 规划中 |
| 预约 | `appointments` | 预约与排班 | 已落地 |
| 就诊 | `visits` | 到诊、诊疗主流程 | 已落地 |
| 诊断 | `diagnoses` | 病种与诊断记录 | 规划中 |
| 处置 | `treatments` | 项目与收费项 | 规划中 |
| 处方 | `prescriptions` | 用药记录 | 已落地 |
| 回访 | `followups` | 复诊/回访任务 | 已落地 |
| 文件 | `visit_files` | 报告、影像、附件 | 规划中 |
| 审计 | `event_logs` | 操作日志与追溯 | 规划中 |
| 同步 | `app_sync_queue` | 对 App 异步同步队列 | 规划中 |

## 2) 状态流转建议

| 对象 | 状态流 |
|---|---|
| appointment.status | `pending -> confirmed -> checked_in -> completed / cancelled` |
| visit.status | `in_progress -> diagnosed -> treated -> closed` |
| followup.status | `pending -> done / skipped` |

## 3) API 草案（目标态）

### 认证与组织
- `POST /merchant/auth/login`
- `GET /merchant/me`
- `GET /merchant/tenants/:id`

### 预约
- `GET /merchant/appointments`
- `POST /merchant/appointments`
- `PATCH /merchant/appointments/:id/status`

### 就诊
- `GET /merchant/visits`
- `POST /merchant/visits`
- `GET /merchant/visits/:id`
- `PATCH /merchant/visits/:id`

### 诊断/处方/文件
- `POST /merchant/visits/:id/diagnoses`
- `POST /merchant/visits/:id/prescriptions`
- `POST /merchant/visits/:id/files`

### 回访
- `GET /merchant/followups`
- `POST /merchant/followups`
- `PATCH /merchant/followups/:id/status`

### App 同步
- `POST /internal/sync/pet-profile`
- `POST /internal/sync/visit-summary`

## 4) 当前已落地 API（以代码为准）

| 接口 | 状态 |
|---|---|
| `POST /api/merchant/auth/login` | 已落地（数据库鉴权） |
| `GET/POST /api/merchant/appointments` | 已落地 |
| `GET/POST /api/merchant/visits` | 已落地 |
| `GET/POST /api/merchant/prescriptions` | 已落地 |
| `GET/POST /api/merchant/followups` | 已落地 |

## 5) 安全与隔离（当前）
- 业务接口要求 `X-Session-ID`
- 服务端按 session 解析 `tenant_id`
- 租户数据强制隔离（读写过滤）
- 处方/回访会校验 `visit` 所属租户，防越权写入

## 6) 医疗数据合规建议（上线前）
- 文件访问改签名 URL（短时有效）
- 审计日志覆盖关键医疗动作
- 角色权限细化（管理员/医生/前台）
- 关键字段加密与脱敏策略

## 7) MVP（当前 + 下一步）
1. 预约管理（已）
2. 就诊记录创建与更新（已）
3. 处方记录（已）
4. 回访任务（已）
5. 文件上传（待）
6. App 摘要同步（待）
