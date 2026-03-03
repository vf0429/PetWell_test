# Task Plan: PetWell商家端优化长计划（MVP -> 可上线）

## Goal
在现有可运行MVP基础上，把PetWell商家端（Shop + Clinic）从“联调演示版”升级到“可灰度上线版”，覆盖安全、集成、质量、可观测与发布流程。

## Current Phase
Phase 4（鉴权与权限体系强化）

## Phases

### Phase 1: 现状盘点与目标对齐
- [x] 盘点当前功能边界（前端、后端、文档、测试）
- [x] 明确“已落地 vs 占位”模块
- [x] 给出当前阶段判断（MVP/Alpha）
- [x] 输出优化优先级（P0/P1/P2）
- **Status:** complete

### Phase 2: 制定优化路线图（本次输出重点）
- [x] 形成8-10阶段长计划
- [x] 每阶段定义：目标、任务、交付物、验收标准
- [x] 补充依赖、风险、时间节奏
- [x] 与你确认“今天先做哪一块”
- **Status:** complete

### Phase 3: 稳定性与工程基线（P0）
- [x] 后端配置分层（dev/staging/prod）与环境变量治理
- [x] 统一错误结构（error code + message + request id）
- [x] 增加请求级日志上下文（request_id, tenant_id, user_id）
- [x] 增加健康检查细化（db connectivity, migration version）
- [x] 构建API契约文档（OpenAPI初版）
- **Status:** complete

### Phase 4: 鉴权与权限体系强化（P0）
- [ ] 真实Google OAuth替换Demo登录
- [ ] 增加`GET /me`、`POST /logout`、session撤销与过期清理
- [ ] session安全策略（TTL、并发登录策略、可选设备管理）
- [ ] RBAC初版（admin/doctor/staff）接口权限矩阵
- [ ] 敏感操作审计日志（登录、权限拒绝、关键状态变更）
- **Status:** in_progress

### Phase 5: 业务数据模型与API完善（P0/P1）
- [ ] 补齐PATCH类接口（订单/预约/就诊/回访状态流转）
- [ ] 引入输入校验层（必填、格式、枚举、边界值）
- [ ] 增加分页/筛选/排序（列表接口）
- [ ] 引入幂等键（关键写操作）
- [ ] 数据一致性优化（事务边界、外键策略、约束）
- **Status:** pending

### Phase 6: Shopify真实集成（P0）
- [ ] OAuth安装流程（connect-url/callback/token持久化）
- [ ] Webhook验签（HMAC）+ 去重（event_id）
- [ ] 事件处理器（orders/fulfillments/inventory）
- [ ] 失败重试与死信队列（最小可用版本）
- [ ] 对账任务（按日增量核对）
- **Status:** pending

### Phase 7: App同步与异步任务系统（P1）
- [ ] 设计`event_logs`与`app_sync_queue`数据结构
- [ ] 实现事件出站（订单/就诊摘要/处方摘要）
- [ ] 消费者重试策略（30s/2min/10min）
- [ ] 幂等消费与重复事件防抖
- [ ] 同步可视化状态页（成功率/失败详情）
- **Status:** pending

### Phase 8: 前端体验与运营效率优化（P1）
- [ ] 登录后态增强（用户信息、租户信息、session到期提示）
- [ ] 表单可用性优化（字段校验、错误提示、加载态）
- [ ] 列表体验优化（分页、筛选、空态、骨架屏）
- [ ] 多语言覆盖补全（遗漏字段、统一文案key）
- [ ] Settings页扩展为“集成配置中心”
- **Status:** pending

### Phase 9: 测试自动化与CI质量门禁（P0）
- [ ] 增加后端单测（store/server核心路径）
- [ ] 增加租户隔离与越权回归测试
- [ ] 将冒烟脚本接入CI（GitHub Actions）
- [ ] 增加契约测试（OpenAPI + schema校验）
- [ ] 制定发布前测试清单（冒烟/回归/性能）
- **Status:** pending

### Phase 10: 上线准备与灰度发布（P0）
- [ ] 生产配置与密钥管理（Secret Manager）
- [ ] 监控告警（成功率、错误率、延迟、队列堆积）
- [ ] 回滚预案（DB备份、版本回退、流量开关）
- [ ] 灰度策略（租户白名单/分批放量）
- [ ] 发布复盘与下一轮迭代计划
- **Status:** pending

## Key Questions
1. 近2周内最优先目标是“可联调稳定”还是“可上线安全”？
2. Shopify与Google OAuth哪个必须先落地（业务驱动顺序）？
3. 当前是否允许引入消息队列（如Redis/NATS）还是先SQLite任务表过渡？
4. 前端是否保持原生JS，还是允许迁移到框架（Vue/React）？
5. 计划周期是1周冲刺还是2周冲刺？

## Decisions Made
| Decision | Rationale |
|----------|-----------|
| 先做P0基础能力再扩功能 | 当前主要短板在安全、集成与质量门禁，不在页面数量 |
| 保留SQLite作为短期默认存储 | 现有代码已稳定运行，适合快速迭代与验证 |
| 先补自动化测试再做大规模重构 | 没有测试护栏会导致功能迭代风险高 |
| Shopify采用“先只读同步，后写入回推” | 降低首期集成风险，便于快速上线 |
| 优先实现统一错误响应与request_id日志 | 这是最低成本、最高收益的工程稳定性改进 |
| OpenAPI先覆盖“已落地接口”再扩展规划接口 | 保证文档与现网行为一致，减少联调歧义 |

## Errors Encountered
| Error | Attempt | Resolution |
|-------|---------|------------|
| 写入`progress.md`时shell heredoc解析失败 | 1 | 改为修复式补丁编辑（apply_patch）并避免反引号被shell执行 |
| 本地启动`go run .`监听8090时被沙箱拒绝 | 1 | 使用`GOCACHE=/tmp`完成`go test`编译级验证，运行级验证待非沙箱环境执行 |
| 处理OpenAPI `$ref`批量替换时Python一次性脚本语法错误 | 1 | 改为使用`perl`批量替换并重新做YAML语义校验 |

## Notes
- 下一步建议：先从Phase 3 + Phase 4中挑1个“今天可完成”的子任务开始动手。
