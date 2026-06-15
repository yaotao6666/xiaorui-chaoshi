# PRD-接口文档

## 0. 2026-06 当前接口摘要

- 后台接口前缀：`/api/v1/admin/*`
- 后台页面文案统一使用“门店”，后端字段仍使用 `merchant_id` 承载门店实体
- 文件上传接口为服务端直传本地存储：
  - `POST /api/v1/upload/file`
  - 请求方式：`multipart/form-data`
  - 返回字段：`path`、`url`、`filename`、`scope`
- 商家侧认证方式为账号密码登录；员工消息推送标识使用 `push_openid` 管理
- 订单创建成功响应不直接返回 `pay_params`，而是返回支付下一步动作
- 当前版本已接入江苏银行微信小程序支付主链路
- 当前版本提供：
  - `POST /api/v1/user/orders/:order_id/pay/prepare`
  - `POST /api/v1/payments/jsbank/notify`
- 当前版本已支持“商家端发起退款（江苏银行）”，暂不提供分账历史相关接口
- 当前数据库与代码模型不包含 `sub_mch_id`、`payment_config_status`、`profit_sharing_*` 等旧支付字段

## 1. 文档定位

本文件用于承接 `PRD.md` 中与“接口、返回结构、字段口径、空值约定”相关的内容。

读取顺序：

1. 先读根目录 `PRD.md`
2. 再读本文件
3. 如需功能背景，联读 `docs/prd/PRD-功能说明.md`

## 2. 接口分组

### 2.1 认证相关接口

- 商家账号密码登录：`POST /api/v1/auth/merchant/login`
- 总部后台登录：`POST /api/v1/admin/auth/login`
- 总部后台退出：`POST /api/v1/admin/auth/logout`
- C 端微信登录：`POST /api/v1/auth/user/wechat-login`

关键约束：

- 登录态字段、缓存字段、角色标识必须在前后端统一
- 商家端 H5 当前不建立 WebSocket 长连接

#### C 端微信登录

- **接口**：`POST /api/v1/auth/user/wechat-login`
- **请求**：`{ code: string }`
- **成功响应**：
  - `code`: `0`
  - `message`: `success`
  - `data.token`: C 端用户 token
  - `data.app_id`: 当前生效的小程序 `appid`
  - `data.user`：
    - `id`
    - `openid`
    - `nickname`
- **前端缓存字段**：
  - `user_token`：对应 `data.token`
  - `userInfo`：对应 `data.user`
  - `openid`：对应 `data.user.openid`
  - `user_login_app_id`：对应 `data.app_id`
- **鉴权使用**：
  - `Authorization: Bearer {user_token}` 用于 `GET/POST /api/v1/user/*`
  - `Authorization: Bearer {user_token}` 用于 `POST /api/v1/store/:merchant_id/orders`
  - `/api/v1/store/:merchant_id/home`、`/products`、`/products/:product_id`、`/delivery-rules` 为公开接口，不要求登录
  - `pages/store/home`、`pages/store/product`、`pages/store/confirm` 统一支持从 `merchant_id` 或 `scene` 解析商家入口参数
  - `GET /api/v1/store/:merchant_id/home` 与 `GET /api/v1/store/:merchant_id/delivery-rules` 进入页面后可直接发起，不等待登录完成

#### C 端访问与埋点

- **访问埋点**：`POST /api/v1/store/:merchant_id/visit`
  - 请求：`{ openid: string, source?: string }`
- **行为埋点**：`POST /api/v1/store/:merchant_id/event`
  - 请求：
    - `openid: string`
    - `event_type: page_view | product_view | submit_order | pay_success`
    - `page?: string`
    - `product_id?: number`
    - `order_id?: number`
    - `source?: string`
    - `payload?: object`

#### 后台公告接口说明

- 公告管理由 `web-admin/` 通过 `/api/v1/admin/announcements*` 调用
- 商家端首页公告展示继续读取已发布公告

### 2.2 总部后台接口

- 当前调用端：`web-admin/` 总部后台
- 仪表盘：`GET /api/v1/admin/dashboard`
- 创建门店：`POST /api/v1/admin/stores`
- 门店列表：`GET /api/v1/admin/stores/list`
- 门店详情：`GET /api/v1/admin/stores/:merchant_id`
- 更新门店资料：`PUT /api/v1/admin/stores/:merchant_id`
- 重置门店管理员密码：`POST /api/v1/admin/stores/:merchant_id/admin/reset-password`
- 门店图片资产更新：`PUT /api/v1/admin/stores/:merchant_id/assets`
- 门店分类、商品、规格、自提点管理：统一通过 `/api/v1/admin/stores/:merchant_id/*` 调用
- 订单列表：`GET /api/v1/admin/orders`
- 订单详情：`GET /api/v1/admin/orders/:order_id`
- 订单分析：`GET /api/v1/admin/orders/analytics`
- 金额分析：`GET /api/v1/admin/amount/analytics`
- 门店排行：`GET /api/v1/admin/amount/top-stores`
- 门店分布分析：`GET /api/v1/admin/stores/analytics/distribution`
- 门店二维码：`GET /api/v1/admin/stores/:merchant_id/qrcode`
- 系统公告：`/api/v1/admin/announcements*`
- 平台活动：`/api/v1/admin/activities*`
- 修改后台密码：`POST /api/v1/admin/account/change-password`

当前重点约束：

- 总部后台数据分析接口需按以下维度输出：
  - 门店访问率
  - 门店下单率
  - 门店下单金额
  - 门店下单均价
  - 日 / 周 / 月 / 年订单量
  - 门店排行维度切换
- `GET /api/v1/admin/stores/analytics/distribution` 返回 `merchants + totals` 结构
- `GET /api/v1/admin/orders` 需支持 `merchant_id`、`status`、`start_date`、`end_date`、`keyword`、`page`、`page_size` 筛选，并返回符合条件的门店订单
- `GET /api/v1/admin/orders/:order_id` 返回结构与商家侧订单详情保持一致，包含商品明细、用户信息、配送信息、支付单号、核销码、核销时间与核销人
- `GET /api/v1/admin/orders/analytics` 返回 `day/week/month/year` 四组订单量桶
- `GET /api/v1/admin/amount/top-stores` 支持 `metric` 参数切换排行维度

### 2.3 商家管理接口

- 商家资料查询 / 更新：`GET/PUT /api/v1/merchant/profile`
- 商家设置：`GET/PUT /api/v1/merchant/settings`
- 修改密码：`POST /api/v1/merchant/account/change-password`
- 商家营业状态：`POST /api/v1/merchant/status`
- 商家二维码：`GET /api/v1/merchant/qrcode`
- 配送设置：`GET/PUT /api/v1/merchant/delivery-settings`
- 自提点管理：`GET/POST/PUT/DELETE /api/v1/merchant/pickup-points*`
- 商家满减规则：`GET/PUT /api/v1/merchant/full-reduction-rules`
- 员工消息推送订阅：`GET/PUT /api/v1/merchant/subscriptions`
- 员工管理：`GET/POST/PUT/DELETE /api/v1/merchant/staff*`
- 员工重置密码：`POST /api/v1/merchant/staff/:id/reset-password`
- 系统公告查看：`GET /api/v1/merchant/announcements`、`GET /api/v1/merchant/announcements/:id`
- 商家打印机管理：`GET/POST/PUT/DELETE /api/v1/merchant/printers*`
- 打印测试：`POST /api/v1/merchant/printers/:printer_id/test`
- 打印日志：`GET /api/v1/merchant/print-logs`

当前重点约束：

- 商家资料更新需支持 `logo` 与背景图字段
- `GET /api/v1/merchant/settings` 需返回 `takeout_enabled`、`dine_in_enabled`、`pickup_enabled`
- `PUT /api/v1/merchant/settings` 需支持更新上述三个开关
- `GET /api/v1/store/:merchant_id/delivery-rules` 除配送费结构外，还需返回 `takeout_enabled`、`dine_in_enabled`、`pickup_enabled`
- `GET /api/v1/merchant/full-reduction-rules` 返回 `rules` 与 `active_rules`，单档规则至少包含 `threshold_amount`、`discount_amount`、`status`、`sort`
- `PUT /api/v1/merchant/full-reduction-rules` 最多支持 5 档规则，`discount_amount` 必须小于 `threshold_amount`
- `GET /api/v1/store/:merchant_id/full-reduction-rules` 为公开接口，仅返回当前商家启用中的满减规则
- `delivery_settings.enabled` 只表示配送费规则是否生效，确认页是否展示“配送”以后端返回的 `takeout_enabled` 为准
- `GET /api/v1/merchant/printers` 返回打印机列表时，需返回 `type`、`status`、`auto_print`、`is_default`、`print_count`、`last_print_at`、`has_api_key`、`has_feie_ukey`
- 飞鹅打印机请求字段至少包含 `feie_user`、`feie_ukey`、`feie_sn`
- `GET /api/v1/merchant/qrcode` 必须固定生成指向 `pages/store/home` 的小程序码，且 `scene` 需使用 `merchant_id={当前商家ID}`
- `GET /api/v1/admin/stores/:merchant_id/qrcode` 作为后台代查看门店二维码接口，返回内容需与商家侧二维码一致
- 门店二维码只在接口请求时动态生成，不写入 `merchants.qrcode_url` 之类的持久化字段

### 2.4 商品与分类接口

- 分类列表 / 新增 / 修改 / 删除
- 分类排序
- 商品列表 / 详情 / 新增 / 更新 / 删除
- 商品上架 / 下架
- 商品批量状态更新
- 商品库存更新
- 商品规格查询 / 更新 / 删除

当前重点约束：

- 分类排序接口需支持按数组顺序批量更新
- 商品状态更新后，店铺首页、商品列表、购物车联动口径必须一致
- 商品规格更新后，前端规格展示、购物车下单与订单明细口径必须一致

### 2.5 订单接口

- C 端创建订单：`POST /api/v1/store/:merchant_id/orders`
- C 端订单列表 / 详情 / 取消 / 申请退款：`/api/v1/user/orders*`
- 商家订单列表 / 详情：`GET /api/v1/merchant/orders*`
- 商家快捷核销：`POST /api/v1/merchant/orders/quick-complete`
- 商家订单核销：`POST /api/v1/merchant/orders/:order_id/complete`
- 商家订单退款：`POST /api/v1/merchant/orders/:order_id/refund`
- 商家订单统计：`GET /api/v1/merchant/orders/statistics`

当前重点约束：

- 商家订单详情应返回：
  - `verify_code`
  - `completed_at`
  - `completed_by_name`
- 商家订单列表接口支持：
  - `status`
  - `start_date`
  - `end_date`
- 商家对当前订单执行核销时，接口以订单上下文完成核销，不要求重复输入核销码
- 商家退款接口需与前端请求参数保持一致，至少明确：
  - `reason`
  - `refund_amount`
- 商家退款接口在未传 `refund_amount` 或传入 `<= 0` 时，默认按订单实付金额处理
- C 端创建订单接口需按 `delivery_type` 与三字段一一校验：
  - `delivery_type=1` 校验 `takeout_enabled`
  - `delivery_type=2` 校验 `dine_in_enabled`
  - `delivery_type=3` 校验 `pickup_enabled`
- 当对应方式未开启时，接口分别返回“商家暂未开启配送 / 堂食 / 自提”
- C 端确认订单页展示的满减优惠仅作提示，下单结果必须以后端返回的 `discount_amount`、`pay_amount` 为准
- 创建订单接口需按商家启用中的满减规则重算：
  - `discount_amount`
  - `pay_amount`
- **退款状态口径**：
  - 订单 `status=5`：退款中（已发起退款处理，等待最终结果）
  - 订单 `status=6`：已退款（退款完成后写入）
  - `refunded_at`：在订单进入 `status=6` 时写入
- 创建订单成功后返回：
  - `payment.enabled`
  - `payment.status`
  - `payment.message`
  - `payment.next_action`
  - `payment.prepare_url`
- 当 `pay_amount > 0` 且来源为 `xcx_shell` 时，前端需按 `payment.next_action = open_xcx_payment` 跳回小程序支付页
- 当 `pay_amount > 0` 且不是 `xcx_shell` 时，前端创建待支付订单后仅展示“请在小程序中完成支付”引导
- 当 `pay_amount <= 0` 时，订单会直接进入已支付状态，前端跳订单详情即可

### 2.6 支付接口

- 支付发起：`POST /api/v1/user/orders/:order_id/pay/prepare`
- 支付回调：`POST /api/v1/payments/jsbank/notify`
- 当前在线支付只支持“小程序壳内发起的江苏银行微信小程序支付”
- 当前在线退款支持“商家端发起退款（江苏银行）”
- 当前版本不提供分账接口

#### `POST /api/v1/merchant/orders/:order_id/refund`

- **用途**：商家对已支付/已完成订单发起退款（江苏银行）
- **鉴权**：商家登录态
- **请求**：

```json
{
  "refund_amount": 18.6,
  "reason": "用户申请退款"
}
```

说明：

- `refund_amount` 为空或 `<= 0` 时，默认按订单实付金额全额退款
- 仅允许对 `status=2(已支付)` 或 `status=3(已完成)` 的订单发起退款
- 发起后订单会进入 `status=5(退款中)`；退款完成后进入 `status=6(已退款)` 并写入 `refunded_at`

- **成功响应**：

```json
{
  "code": 0,
  "message": "退款已提交",
  "data": null
}
```

#### `POST /api/v1/user/orders/:order_id/pay/prepare`

- **用途**：为待支付订单生成江苏银行小程序支付参数
- **鉴权**：`Authorization: Bearer {user_token}`
- **请求**：

```json
{
  "return_path": "/pages/store/order-detail",
  "source": "xcx_shell"
}
```

- **成功响应**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "order_id": 123,
    "order_no": "202606140001",
    "merchant_id": 2,
    "pay_amount": 18.6,
    "channel": "jsbank_wechat_mini",
    "pay_params": {
      "timeStamp": "1710000000",
      "nonceStr": "random",
      "package": "prepay_id=xxx",
      "signType": "RSA",
      "paySign": "signature"
    },
    "return_target": "/pages/store/order-detail?order_id=123&merchant_id=2"
  }
}
```

当前重点约束：

- 仅允许当前登录用户为自己的 `status=1` 待支付订单发起支付
- `source` 当前固定为 `xcx_shell`
- 返回的 `pay_params` 仅供 `xcx` 原生页调用 `Taro.requestPayment`
- H5 页面不直接调用小程序支付

#### `POST /api/v1/payments/jsbank/notify`

- **用途**：接收江苏银行支付成功回调并更新订单状态
- **请求格式**：`application/x-www-form-urlencoded`
- **核心处理**：
  - 使用江苏银行公钥证书验签
  - 根据商户订单号映射 `orders.order_no`
  - 幂等更新订单 `status/transaction_id/paid_at/pay_notify_payload`

当前重点约束：

- 回调成功后，订单状态从 `1=待支付` 更新为 `2=已支付`
- 重复回调必须幂等
- 与支付相关的历史分账字段和历史表不属于当前数据库基线

### 2.7 数据分析接口

#### 商家端

- 今日概览
- 订单趋势
- 客户分析
- 商品排行
- 库存预警

#### 总部后台

- 门店转化分析
- 门店排行榜
- 周期订单分析

当前重点约束：

- 所有列表型分析接口必须保证空数组兜底
- 所有对象型统计字段在无数据时应返回 `0` 或明确的空对象结构
- `GET /api/v1/merchant/analytics/stock-alert` 默认使用阈值查询低库存商品，当前前端按 `threshold = 10` 使用
- 库存预警接口返回结果按库存升序处理，前端据此做风险分组和看板展示

### 2.8 上传接口

- 上传接口：`POST /api/v1/upload/file`
- 上传方式：`multipart/form-data`
- 存储方式：服务端直传本地存储
- 返回字段：`path`、`url`、`filename`、`scope`

当前重点约束：

- 前端按接口返回的 `url` 与 `path` 使用上传结果
- 图片上传成功后的持久化字段口径必须在前后端统一

### 2.9 历史 WebSocket 能力归档

- `/api/v1/ws/merchant` 与开发联调接口不再作为当前商家端 H5 的有效协议。
- 当前商家端 H5 不主动建立握手连接，也不依赖实时推送完成页面功能。
- 若后续恢复实时提醒能力，应以新的协议说明重新定义。

## 3. 返回与空值约定

### 3.1 通用返回

- 成功响应统一使用：
  - `code`
  - `message`
  - `data`
- 失败响应统一使用：
  - `code != 0`
  - `message`

### 3.2 统计返回

- 统计对象无数据时返回结构化对象，不直接返回 `null`
- 数值字段默认返回 `0`

### 3.3 图表返回

- 趋势、排行、分布类图表字段必须返回数组
- 维度切换接口无数据时返回空数组，由前端展示空态

## 4. 当前重点同步要求

1. 功能改动涉及接口时，必须同时更新：
   - 根 `PRD.md`
   - 本文件对应章节
2. 若字段口径调整，必须同步更新：
   - 请求参数
   - 响应示例
   - 空值约定
3. 若只是 UI 调整但不改协议，可只更新功能说明文档，不必改本文件

## 5. 与根 PRD 的对应关系

- 认证相关接口：对应 `PRD.md` `3.1`
- 总部后台接口：对应 `PRD.md` `3.2`
- 商家管理接口：对应 `PRD.md` `3.3`
- 商品分类 / 商品接口：对应 `PRD.md` `3.4`、`3.5`
- 订单接口：对应 `PRD.md` `3.6`
- 数据分析接口：对应 `PRD.md` `3.7`
- C 端接口：对应 `PRD.md` `3.8`
- 支付接口：对应 `PRD.md` `3.9`
