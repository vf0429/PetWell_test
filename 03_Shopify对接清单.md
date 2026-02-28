# Shopify 对接清单（PetWell 商家端）

## 1) 推荐接口
- 主接口：Shopify Admin GraphQL API（订单/商品/库存/客户）
- 可选：Storefront API（若 App 端需直接商品浏览）
- 事件：Webhook（订单、发货、库存变更）

## 2) 你需要准备的内容
1. Shopify Partner 账号
2. 创建 PetWell 私有/公有 App（建议走 OAuth）
3. 获取 API Key / API Secret
4. 配置 App Redirect URL（PetWell 后端）
5. 配置 Webhook 回调 URL（PetWell 后端）

## 3) 首期建议权限（scopes）
- read_products, write_products（如需从商家端改商品）
- read_inventory, write_inventory（库存同步）
- read_orders, write_orders（订单状态）
- read_customers（用户数据最小读取）
- read_fulfillments, write_fulfillments（发货）

> 注意：最终 scopes 以你实际功能范围收敛，遵循最小权限。

## 4) 必接 Webhooks（首期）
- orders/create
- orders/paid
- orders/cancelled
- fulfillments/create
- fulfillments/update
- inventory_levels/update

## 5) PetWell 后端需要做什么
1. OAuth 授权流程
   - 发起安装 -> 回调换 token -> 保存店铺 token
2. Webhook 验签
   - 校验 HMAC，防止伪造请求
3. 数据标准化
   - Shopify 字段转为 PetWell 统一订单/商品模型
4. 异步任务
   - 大批量同步（初次导入）走队列
5. 失败重试
   - Webhook 处理失败自动重试并告警

## 6) 建议的对接策略
- Phase A：只读 + 状态同步（低风险）
- Phase B：支持商家端改库存/发货（写入能力）
- Phase C：推荐引擎联动（按宠物健康档案做个性化商品推荐）

## 7) 关键接口示例（语义）
- GET /merchant/shopify/connect-url
- GET /merchant/shopify/callback
- POST /merchant/webhooks/shopify/orders
- POST /merchant/webhooks/shopify/fulfillments
- POST /merchant/webhooks/shopify/inventory

