# PetWell 商家端实施 Checklist

> 规则：完成一项就勾选一项（`[x]`）。
> 说明：本清单用于快速确认“当前版本到底做到了什么”。

- [x] 线稿完成（`02_商家端线框稿.html`）
- [x] 高保真前端原型（Shop / Clinic 双视图 Dashboard）
- [x] 多页面功能实现（Orders/Appointments、Products/Visits、Customers/Prescriptions、Settings）
- [x] 医院端后端骨架（预约、就诊、处方、回访 API）
- [x] 商品端后端骨架（订单、商品、客户 API）
- [x] Shopify 对接后端占位（OAuth / Webhook / 同步任务）
- [x] 数据库初始化脚本（SQLite 核心表）
- [x] SQLite 持久化接入（默认 store）
- [x] 登录鉴权接入数据库（邮箱/手机号+密码）
- [x] Session 机制接入（`X-Session-ID`）
- [x] 租户隔离落地（按 `tenant_id` 强制过滤）
- [x] 处方/回访越权校验（visit 所属租户）
- [x] 前端中英切换（含表格列名）
- [x] 登录后 Guide 引导（右下角 Next）
- [x] 阶段性测试用例补充到可执行格式（API + UI）
- [x] API 冒烟脚本升级（自动登录 + session 请求）
- [x] 联调说明文档（与 PetWell App 对接字段与时序）
