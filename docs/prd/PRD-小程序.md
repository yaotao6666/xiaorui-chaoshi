# PRD-小程序

## 1. 文档定位

本文件用于单独沉淀“小程序壳登录 + 跳转 H5”相关说明，作为后续以下工作的统一基线：

- 小程序登录接入说明
- 小程序 `web-view` 嵌入 H5 交付说明
- H5 接收登录态透传参数说明
- 后续支付跳转、支付完成回跳、支付结果对账说明

当前版本先覆盖：

- 小程序使用微信登录换取系统用户 token
- 登录成功后跳转现有 H5 页面
- 小程序与 H5 之间的参数透传约定

## 2. 当前目标

当前 `xcx/` 工程不是完整业务小程序，而是一个用于交付验证的“小程序壳”。

当前交付目标为：

1. 小程序启动后先调用微信登录
2. 使用微信 `code` 调用 Go 服务登录接口
3. Go 服务返回当前用户 `token/openid/app_id`
4. 小程序将这些信息拼到目标 H5 地址中
5. 小程序通过 `web-view` 打开 H5，后续业务继续在 H5 内完成

当前 Go 服务统一使用：

- `WECHAT_APP_ID`
- `WECHAT_APP_SECRET`


## 3. 当前页面与职责

当前 `xcx/` 目录包含以下页面：

- `pages/index/index`
  - 启动页
  - 负责调用登录接口并拼接 H5 地址
- `pages/webview/index`
  - H5 容器页
  - 负责承载目标 H5
- `pages/settings/index`
  - 配置页
  - 维护 API 基址、H5 入口、默认门店、调试 token、调试 openid
- `pages/debug/index`
  - 调试页
  - 用于查看当前缓存登录态与最终生成的 H5 地址

## 4. 交互链路

### 4.1 标准链路

```text
小程序启动
  -> 调用 wx.login / Taro.login
  -> 获取临时 code
  -> POST /api/v1/auth/user/wechat-login
  -> 返回 token/openid/app_id/user
  -> 小程序本地缓存登录态
  -> 根据 merchant_id + token + openid + source + app_id 生成 H5 地址
  -> 跳转 pages/webview/index
  -> web-view 打开 H5
```

### 4.2 本地调试链路

当不是微信小程序运行环境时，`xcx/` 当前支持调试模式：

1. 在配置页填写：
   - `debugToken`
   - `debugOpenid`
2. 本地不调用真实 `Taro.login`
3. 直接生成调试登录态
4. 继续验证 H5 地址拼装和 `web-view` 跳转

该模式仅用于联调壳验证，不等同于真实微信登录。

## 5. 小程序配置项

当前小程序壳依赖以下配置：

- `apiBaseUrl`
  - Go 服务 API 基址
  - 例如：`http://127.0.0.1:8081`
- `h5EntryUrl`
  - H5 入口地址
  - 例如：`http://localhost:3000/#/pages/store/home`
- `merchantId`
  - 默认进入的门店 ID
- `debugToken`
  - 非微信环境调试 token
- `debugOpenid`
  - 非微信环境调试 openid

真机环境补充要求：

- 小程序 `web-view` 访问的 H5 地址必须为已配置业务域名的 HTTPS 地址
- 小程序登录依赖真实微信 `code`，不能直接用浏览器环境模拟

## 6. H5 参数透传约定

当前小程序跳转 H5 时，会在 URL 中透传以下字段：

- `merchant_id`
  - 当前门店 ID
- `token`
  - C 端登录 token
- `openid`
  - 当前用户微信 openid
- `source`
  - 来源标识
  - 小程序壳当前默认值：`xcx_shell`（用于支付链路识别）
- `app_id`
  - 当前小程序 `appid`

### 6.1 示例

```text
https://h5.example.com/#/pages/store/home?merchant_id=2&token=xxx&openid=ooo&source=xcx_shell&app_id=wxea66472a9f8b25af
```

### 6.2 H5 侧要求

H5 接到上述参数后，应按以下方式处理：

1. 读取 `merchant_id`
   - 作为当前门店上下文
2. 读取 `token`
   - 写入 C 端登录态缓存
3. 读取 `openid`
   - 写入用户标识缓存
4. 读取 `app_id`
   - 写入当前登录来源的小程序标识
5. 读取 `source`
   - 用于来源识别、埋点和后续支付来源判断

说明：

- 当前 H5 与小程序联调的关键是“参数透传 + token 复用”
- 当前版本不要求小程序内重复承载完整商城页面

## 7. 接口文档

### 7.1 小程序微信登录

- **接口**：`POST /api/v1/auth/user/wechat-login`
- **用途**：用微信登录 `code` 换取系统用户登录态
- **调用端**：`xcx/` 小程序壳
- **鉴权**：无需登录 token
- **Content-Type**：`application/json`

#### 请求参数

```json
{
  "code": "wx-login-code"
}
```

字段说明：

- `code`
  - 必填
  - 来源于 `wx.login` / `Taro.login`

#### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "jwt-token",
    "app_id": "wxea66472a9f8b25af",
    "user": {
      "id": 1,
      "openid": "oXXXXXX",
      "nickname": "微信用户"
    }
  }
}
```

字段说明：

- `data.token`
  - C 端用户 token
  - 后续 H5 请求用户相关接口时使用
- `data.app_id`
  - 当前生效的小程序 `appid`
- `data.user.id`
  - 系统用户 ID
- `data.user.openid`
  - 当前用户微信 openid
- `data.user.nickname`
  - 当前用户昵称

#### 失败响应

常见失败返回：

```json
{
  "code": 400,
  "message": "参数错误"
}
```

```json
{
  "code": 500,
  "message": "微信登录失败"
}
```

#### 服务端处理说明

当前服务端处理逻辑如下：

1. 校验 `code` 是否为空
2. 使用 `WECHAT_APP_ID + WECHAT_APP_SECRET` 调用微信 `jscode2session`
3. 获取 `openid/unionid`
4. 按 `openid` 查询系统用户
5. 若用户不存在，则自动创建 C 端用户
6. 更新最近访问时间与访问次数
7. 生成系统登录 token 并返回

#### 当前使用范围

成功拿到 `token` 后，可用于：

- `GET /api/v1/user/*`
- `POST /api/v1/user/*`
- `POST /api/v1/store/:merchant_id/orders`

公开接口仍可直接访问：

- `GET /api/v1/store/:merchant_id/home`
- `GET /api/v1/store/:merchant_id/products`
- `GET /api/v1/store/:merchant_id/products/:product_id`
- `GET /api/v1/store/:merchant_id/delivery-rules`

## 8. 小程序侧登录态缓存约定

当前小程序本地缓存统一保存以下字段：

- `token`
- `openid`
- `appId`
- `nickname`
- `obtainedAt`

其中：

- `token/openid/appId` 会继续透传给 H5
- `obtainedAt` 用于判断本地会话获取时间

## 9. 与 H5 的职责边界

当前职责分工如下：

### 9.1 小程序负责

- 获取微信登录 `code`
- 调后端登录接口
- 保存最小登录态
- 拼接 H5 地址
- 使用 `web-view` 打开 H5

### 9.2 H5 负责

- 承接商城主业务流程
- 读取小程序透传参数并恢复登录态
- 门店首页、商品页、下单页、订单页等业务页面渲染
- 用户下单、后续支付、订单查询等业务执行

## 10. 联调注意事项

### 10.1 域名要求

- 微信小程序 `web-view` 不支持直接访问本地 `localhost`
- 真机联调时必须使用可公网访问的 HTTPS 地址
- 该域名需提前配置到微信小程序业务域名中

### 10.2 参数要求

- `merchant_id` 必须有效
- `token` 必须是当前后端签发的有效 C 端 token
- H5 需要兼容 hash 路由地址参数追加

### 10.3 登录失败处理

建议小程序壳统一处理以下失败提示：

- 未拿到微信 `code`
- 登录接口返回错误
- H5 地址为空或非法
- `web-view` 打开失败

## 11. 小程序支付跳转链路

当前版本已按“`H5` 下单、`xcx` 原生页发起支付、支付后回到 `H5` 订单详情”的方式接入江苏银行微信小程序支付。

### 11.1 支付前置约定

- `H5` 负责创建订单，不直接拉起微信支付。
- 当订单 `pay_amount > 0` 且当前来源是 `xcx_shell` 时，`H5` 通过 `wx.miniProgram.navigateTo` 跳回 `xcx/pages/payment/index`。
- `xcx` 支付页再调用 `POST /api/v1/user/orders/:order_id/pay/prepare` 获取江苏银行支付参数，并执行 `Taro.requestPayment`。

### 11.2 支付跳转链路

```text
xcx 启动页登录
  -> web-view 打开 H5
  -> H5 提交订单
  -> H5 根据 payment.next_action 跳回 xcx 原生支付页
  -> xcx 请求 /api/v1/user/orders/:order_id/pay/prepare
  -> xcx 执行 Taro.requestPayment
  -> 支付成功后重新打开 H5 订单详情页
```

### 11.3 前端对接要点（H5 侧）

H5 侧需要做到两件事：

1. 创建订单时把来源标识传给后端
2. 创建订单成功后，根据后端返回的 `payment.next_action` 做跳转/引导

#### 11.3.1 创建订单请求

- **接口**：`POST /api/v1/store/:merchant_id/orders`
- **调用端**：H5（在小程序壳 `web-view` 内运行）
- **鉴权**：需要登录 token
- **请求体新增字段**：`source`

请求示例（只展示与支付相关字段；其他下单字段保持原有实现）：

```json
{
  "source": "xcx_shell"
}
```

#### 11.3.2 创建订单响应（新增 payment 字段）

成功响应示例（只展示与支付相关字段）：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "order": {
      "id": 123,
      "merchant_id": 2,
      "pay_amount": 19.9
    },
    "payment": {
      "enabled": true,
      "status": "pending",
      "message": "订单已创建，请继续完成支付",
      "next_action": "open_xcx_payment",
      "prepare_url": "/api/v1/user/orders/123/pay/prepare"
    }
  }
}
```

字段说明：

- `payment.enabled`
  - 是否需要在线支付（`pay_amount > 0` 时为 true）
- `payment.next_action`
  - `view_order_detail`：无需支付（0 元），H5 直接跳订单详情
  - `open_xcx_payment`：需要支付且来源为 `xcx_shell`，H5 需要跳回 xcx 原生支付页
  - `show_xcx_pay_guide`：需要支付但不在小程序壳环境，H5 需要提示“请在小程序内支付”
- `payment.prepare_url`
  - 预下单接口路径（相对路径），主要给排查与调试使用；H5 通常不直接调用该接口

#### 11.3.3 H5 跳回 xcx 原生支付页

当 `payment.next_action === "open_xcx_payment"` 时：

- H5 使用 `wx.miniProgram.navigateTo` 跳转：
  - 目标页：`/pages/payment/index`
  - 携带参数：
    - `orderId`：订单 ID
    - `merchantId`：门店 ID
    - `returnTarget`：支付完成后回到 H5 的 hash 路径（例如 `/pages/store/order-detail`）

当前项目已封装工具函数（建议直接复用）：

- [miniProgramBridge.ts](file:///e:/yt/xiaorui-chaoshi/miniprogram/src/utils/miniProgramBridge.ts) 的 `openXcxPaymentPage({ orderId, merchantId, returnTarget })`

### 11.4 前端对接要点（xcx 侧）

xcx 原生页负责两步：

1. 调后端 `prepare` 接口获取 `pay_params`
2. 将 `pay_params` 透传给 `Taro.requestPayment` 拉起微信支付

#### 11.4.1 进入参数

页面：`xcx/pages/payment/index`

URL 参数：

- `orderId`：订单 ID（必填）
- `merchantId`：门店 ID（可选，主要用于回跳拼参）
- `returnTarget`：回跳 H5 的 hash 路径（例如 `/pages/store/order-detail`）

#### 11.4.2 发起支付（prepare 接口）

- **接口**：`POST /api/v1/user/orders/:order_id/pay/prepare`
- **调用端**：xcx 原生页
- **鉴权**：需要登录 token（Authorization）
- **Content-Type**：`application/json`

请求示例：

```json
{
  "source": "xcx_shell",
  "return_path": "/pages/store/order-detail"
}
```

说明：

- `source` 必须为 `xcx_shell`，否则后端会拒绝（用于防止非壳环境直接发起支付）
- `return_path` 是 H5 内的 hash 路径，后端会补齐 `order_id/merchant_id` 并返回 `return_target`

成功响应示例：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "order_id": 123,
    "order_no": "202606150001",
    "merchant_id": 2,
    "pay_amount": 19.9,
    "channel": "jsbank_wechat_mini",
    "pay_params": {
      "timeStamp": "1710000000",
      "nonceStr": "xxxx",
      "package": "prepay_id=xxxx",
      "signType": "RSA",
      "paySign": "xxxx"
    },
    "return_target": "/pages/store/order-detail?order_id=123&merchant_id=2"
  }
}
```

#### 11.4.3 拉起微信支付

xcx 将 `pay_params` 映射到 `Taro.requestPayment` 参数：

- `pay_params.timeStamp` -> `timeStamp`
- `pay_params.nonceStr` -> `nonceStr`
- `pay_params.package` -> `package`
- `pay_params.signType` -> `signType`
- `pay_params.paySign` -> `paySign`

项目内现成实现参考：

- [payment/index.tsx](file:///e:/yt/xiaorui-chaoshi/xcx/src/pages/payment/index.tsx)

### 11.5 关键字段

- `order_id`：订单 ID
- `merchant_id`：当前门店 ID
- `source=xcx_shell`：标识当前链路来自小程序壳
- `return_path`：支付完成后希望回到的 H5 hash 路径
- `return_target`：后端组装好的订单详情回跳路径
- `pay_params`：仅供 `xcx` 原生页传给 `Taro.requestPayment`

### 11.6 支付回跳说明

- 支付成功：`xcx` 自动重新打开 `H5` 订单详情页
- 支付取消：停留在 `xcx` 支付页，允许用户重试或返回 `H5`
- 支付失败：停留在 `xcx` 支付页，展示错误信息并允许重试
- `H5` 订单详情页会按订单接口刷新状态，展示支付时间和核销码

### 11.7 常见错误与排查（给前端）

- `POST /api/v1/user/orders/:id/pay/prepare` 返回 401
  - token 未带或已过期：确认 xcx 已走登录并缓存 token（或点“清缓存并重新登录”）
- 返回 400 且提示“当前仅支持小程序壳发起支付”
  - `source` 不是 `xcx_shell`：确认 H5 下单请求体带了 `source=xcx_shell`，且 xcx 调 prepare 时也带了同样的 source
- 返回 400 且提示“当前订单状态不可发起支付 / 当前订单无需支付”
  - 订单不是待支付（status!=1）或金额为 0：应回到 H5 订单详情刷新状态
- 本地能拉起支付，但订单不自动变为已支付
  - 银行异步回调走 `JSBANK_PAY_NOTIFY_URL`（通常是线上域名），本地环境不一定能收到回调；线上联调以回调日志和订单状态更新为准

## 12. 当前结论

当前小程序壳的正式基线是：

- 小程序只负责登录、参数透传、壳内打开 H5
- 小程序原生页还负责承接支付动作
- H5 继续承接实际业务页面
- Go 服务统一以 `WECHAT_APP_ID + WECHAT_APP_SECRET` 对接微信小程序登录
- Go 服务通过江苏银行接口为 `xcx` 支付页生成支付参数，并接收支付回调更新订单状态
- 本文件作为后续“小程序壳 + H5 + 支付跳转”统一对接文档持续维护
