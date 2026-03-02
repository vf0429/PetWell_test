# Findings & Decisions

## Requirements
- 用户希望基于当前仓库现状，判断项目处于什么阶段。
- 用户希望继续优化项目，并先提供一份“较长、较详细”的执行计划。
- 计划需要可落地，后续可据此决定“今天先做什么”。

## Research Findings
- 文档与代码一致表明：项目已完成可运行MVP（登录、session、租户隔离、Shop/Clinic核心CRUD）。
- `README.md`与`CHANGELOG.md`显示当前时间点（2026-03）重点是完成SQLite持久化与鉴权升级。
- `backend/internal/server.go`里Shopify OAuth/Webhook仍是占位实现，未做真实OAuth与验签。
- `backend/internal/store.go`已实现会话查询、租户过滤、处方/回访越权校验。
- `web/app.js`为原生JS单页控制台，支持EN/ZH、登录、Guide与基础数据操作。
- 当前自动化测试主要是Shell冒烟脚本；Go单测、CI流程尚未看到。
- 后端统一错误结构已改为：`error.code + error.message + request_id`（用于联调与排障）。
- 日志中已增加`request_id / tenant_id / user_id / status / duration_ms`上下文，便于请求级追踪。
- 健康检查已升级：`/health`会返回`store.type`、`db.reachable`和`db.schema_version`，DB异常时返回`503`。
- OpenAPI初版已落地到`backend/openapi.yaml`，覆盖当前已实现接口与统一错误响应结构。
- 环境配置已从散落环境变量读取，升级为`LoadConfig`分层策略（`dev/staging/prod` + 显式覆盖）。

## Technical Decisions
| Decision | Rationale |
|----------|-----------|
| 将当前阶段定义为“可联调MVP/Alpha” | 功能链路可跑通，但集成、安全、测试体系未达上线标准 |
| 优化路线采用“P0先稳定再扩展” | 先补底座（鉴权/测试/监控）才能安全迭代业务功能 |
| 计划按10个Phase拆分 | 便于与你按“今天做什么”逐步执行 |
| 第一优先落地“错误标准化 + request_id” | 可直接提升问题定位效率，且改造风险低 |
| 配置分层先做“轻量Profile + 环境变量覆盖” | 不引入额外依赖，保持MVP演进速度 |

## Issues Encountered
| Issue | Resolution |
|-------|------------|
| 无阻塞问题 | 代码结构清晰，信息完整，可直接进入计划阶段 |
| 沙箱环境禁止监听本地端口8090 | 改用`go test`完成编译级验证，待真实环境补跑冒烟脚本 |

## Resources
- `README.md`
- `CHANGELOG.md`
- `08_商家端总览与后续拓展建议.md`
- `09_最终测试稿_登录与租户隔离.md`
- `backend/internal/server.go`
- `backend/internal/store.go`
- `backend/internal/config.go`
- `backend/internal/config_test.go`
- `backend/openapi.yaml`
- `web/app.js`
- `06_API冒烟测试脚本.sh`

## Visual/Browser Findings
- 本轮未使用浏览器/图片类工具。
