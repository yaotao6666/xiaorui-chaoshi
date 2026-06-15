# 商家助手 API 产品需求文档 (PRD)

## 0. 文档导航与同步规则

### 0.1 必读入口

`PRD.md` 是本项目的 **PRD 总入口、索引页和同步约束页**。后续无论是人工查阅还是大模型读取，**都必须先读本文件，再按下列顺序联读拆分文档**：

1. `PRD.md`
2. `docs/prd/PRD-功能说明.md`
3. `docs/prd/PRD-接口文档.md`
4. `docs/prd/PRD-测试与附录.md`
5. 按需继续读取：
   - `docs/prd/完整业务链路目录.md`
   - `docs/prd/测试报告模板.md`

### 0.2 各文档职责

| 文件                       | 职责                               |
| ------------------------ | -------------------------------- |
| `PRD.md`                 | 总索引、阅读顺序、同步规则、关键摘要、历史完整正文        |
| `docs/prd/PRD-功能说明.md`   | 产品角色、页面结构、业务功能、关键链路说明            |
| `docs/prd/PRD-接口文档.md`   | 接口分组、返回口径、空值约定 |
| `docs/prd/PRD-测试与附录.md`  | 测试账号、验收口径、部署/环境要点、文档维护规范         |
| `docs/prd/商家入驻协议草案目录.md` | 历史入驻协议目录，仅作废弃背景参考                |
| `docs/prd/完整业务链路目录.md`   | 按角色梳理完整业务流程与验证路径                 |
| `docs/prd/测试报告模板.md`     | 完整链路测试报告模板                       |

### 0.3 更新规则

1. 修改业务功能时：
   - 必须同步更新 `PRD.md` 中的目录摘要；
   - 必须同步更新 `docs/prd/PRD-功能说明.md` 对应章节。
2. 修改接口、字段、返回结构、空值口径时：
   - 必须同步更新 `PRD.md` 中的摘要或索引说明；
   - 必须同步更新 `docs/prd/PRD-接口文档.md` 对应章节。
3. 修改测试流程、测试账号、验收结论、部署要求时：
   - 必须同步更新 `docs/prd/PRD-测试与附录.md`；
   - 如形成一次完整回归，还应同步更新对应测试报告。
4. 新增协议、链路说明、模板类文档时：
   - 必须在 `PRD.md` 的本章节补充索引；
   - 必须标注用途与适用范围。
5. 更新任意 PRD 文档时：
   - 只记录功能、页面、接口、字段、流程、状态“是什么”；
   - 不展开实现原因、决策过程和冗长解释；
   - `PRD.md` 与 `docs/prd/*` 均遵循同一写法。

### 0.4 同步约束

- 后续读取 `PRD.md` 时，默认视为“需要继续联读上述拆分文档”。
- 若 `PRD.md` 与拆分文档存在冲突，以 **最近一次同步修改后的同组文档** 为准，并在下一次更新时完成收敛，避免长期分叉。
- 禁止只更新代码不更新 PRD；也禁止只更新拆分文档不回写 `PRD.md` 索引。

### 0.5 本轮已同步摘要

- 超市独立版改造：
  - 项目形态已从“服务商多商家”调整为“单超市主体 + 多门店”。
  - `api-chaoshi` 已切换本地文件存储，移除七牛依赖与上传凭证接口。
  - 总部后台接口前缀已切换为 `/api/v1/admin/*`，`web-admin` 构建基路径已切换为 `/admin/`。
  - 用户端主体已转为 H5，微信小程序保留用户登录壳与 `web-view` 容器；商家端仅保留账号密码登录。
  - 微信支付、在线退款、分账、服务商主体小程序配置已从当前版本主链路移除，等待后续新支付方案接入。
  - 后端代码模型与数据库启动补列逻辑已同步移除 `sub_mch_id`、`payment_config_status`、`profit_sharing_*` 等旧支付/分账字段。
- 商家端：
  - 商家端 H5 当前不再建立 WebSocket 握手连接，登录失效与主动退出仅清理本地登录态。
  - 订单详情与订单列表支持直接使用订单 `verify_code` 核销，工作台保留扫码/输入核销码的快速核销。
  - 已支付、已完成订单支持退款，退款金额默认按订单实付金额处理。
  - 商家工作台快捷功能补齐图标与辅助文案，今日概览中的待处理订单改为待核销订单，商品数量改为已上架商品总数。
  - 商家订单列表新增快捷范围与自定义日期筛选，支持按开始日期和结束日期过滤。
  - 商家分析页库存预警改为看板展示，包含汇总、风险分级、空态和商品管理入口。
  - 商家设置已移除微信快捷登录、微信绑定与分账历史入口。
- 服务商端：
  - 商家由服务商直接创建，不再存在小程序入驻、审核、进件流程。
  - 总部后台当前仅维护门店资料、管理员账号、图片资源、订单与经营分析，不再维护支付配置与分账参数。
  - 默认最小初始化账号已调整为 `admin / tm666666`。
  - 数据分析收口为商家访问率、下单率、下单金额、下单均价、日周月年订单量与排行榜维度切换。
  - 商家详情支持维护 `logo` 与 `cover_image`。
- 文档与测试：
  - 本轮接口、功能、测试结论已同步到拆分文档。
  - 完整回归结果单独沉淀在 `docs/prd/完整链路测试报告-20260513.md`。
  - `api-chaoshi/migrations/` 当前仅保留 `20240101000000_full_init.sql` 作为单文件初始化基线。
  - `api-chaoshi/scripts/regression-test.sh` 当前仅使用该单文件初始化脚本。
  - 新增 `xcx/` 小程序验证壳目录，用于验证“微信小程序无感登录 -> 带 token 进入现有 H5”的交付链路。
  - C 端微信登录接口返回结构统一为 `code/message/data.token/data.app_id/data.user`，前端缓存字段为 `user_token/userInfo/openid/user_login_app_id`，下单接口使用 `user_token` 鉴权。
  - C 端订单退款状态口径统一为：`status=5` 退款中、`status=6` 已退款。
  - C 端访问/行为埋点接口使用 `openid` 作为用户唯一标识（`/api/v1/store/:merchant_id/visit`、`/api/v1/store/:merchant_id/event`）。
  - C 端订单详情已独立为 `pages/store/order-detail`，支付完成、取消支付与“去购物”均保留当前 `merchant_id`，避免返回错误商家店铺。

### 0.6 历史归档说明

- 本文件中凡是出现 `sub_mch_id`、`payment_config_status`、`profit_sharing_*`、微信支付回调、服务商分账等描述，除非明确标注“当前生效”，否则均视为历史方案归档。
- 当前生效口径以以下文档为准：
  - `docs/prd/PRD-功能说明.md`
  - `docs/prd/PRD-接口文档.md`
  - `docs/prd/PRD-测试与附录.md`
- 本次已将最易混淆的旧服务商支付章节在根 PRD 中改为“历史归档 + 当前说明”，避免继续被误读为实施基线。

## 1. 项目概述

### 1.1 项目背景

开发一套面向单超市主体、多门店经营场景的管理系统，当前以总部后台、商家端 H5 和用户端 H5/小程序壳协同运行为核心形态。

### 1.2 小程序产品定位

**产品名称**：商家助手

**产品定位**：专为中小商家打造的统一私域经营平台

**核心原则**：小程序围绕**商家运营管理**为中心，默认服务于商家日常经营

**小程序整体架构**：

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            商家助手                                      │
│                        (统一小程序 + 多入口模式)                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                        │
                    ┌───────────────────┼───────────────────┐
                    │                   │                   │
                    ▼                   ▼                   ▼
          ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
          │   商家管理端     │ │   用户端         │ │  服务商PC后台    │
          │  (默认入口)      │ │  (扫码入口)      │ │   (Web入口)      │
          │                 │ │                 │ │                 │
          │ • 商家登录       │ │ • 进入商家店铺   │ │ • 后台登录       │
          │ • 店铺信息管理   │ │ • 商品浏览下单   │ │ • 商家统计       │
          │ • 商品管理       │ │ • 创建订单       │ │ • 订单统计       │
          │ • 订单管理       │ │ • 我的订单       │ │ • 数据分析       │
          │ • 数据分析       │ │                 │ │ • 商家管理       │
          │ • 商家二维码     │ │                 │ │                 │
          └─────────────────┘ └─────────────────┘ └─────────────────┘
```

**小程序设计风格**
为一个专为中小商家设计的私域经营平台设计图标，产品名称为“商家助手”。

设计原则：

- 核心：围绕“商家运营管理”这一中心
- 理念：统一、聚合、专业、可靠

视觉元素建议：

- 主图形：抽象的“店铺”+“数据”+“管家”元素融合
- 可选元素：盾牌（安全）、网格（统一入口）、指南针（私域向导）
- 风格：简洁扁平、现代商务、辨识度高

配色方案：

- 主色：深蓝 + 金色/橙色点缀（深蓝代表专业可靠，金色/橙色代表商家的财富增长）
- 辅助色：白色用于图形/文字

输出：6个不同方向的SVG设计变体，包含方形图标和圆形图标两种形式。

### 1.3 商家管理端（默认入口）

**说明**：商家管理端是小程序的**默认入口**，商家通过账号密码登录后进入，专注于日常店铺经营管理。

**核心功能**：

| 模块   | 功能描述                     |
| ---- | ------------------------ |
| 店铺管理 | 商家信息编辑、营业状态设置、店铺二维码生成    |
| 商品管理 | 商品上架/下架、库存管理、价格调整        |
| 订单管理 | 订单查看、订单处理、退款管理           |
| 数据分析 | 今日订单、交易额、客单价、趋势分析        |
| 设置与员工 | 配送设置、密码修改与员工管理 |

**导航结构（商家端）**：

- 底部 Tab：工作台、订单、数据分析、设置。
- tabBar 仅在 Tab 页面展示；登录等非 Tab 页面不展示。
- 工作台快捷入口当前已精简，不再重复提供“订单管理”“店铺设置”入口。
- 工作台在“今日概览”上方展示系统公告滚动栏，支持点击查看详情与点击关闭，关闭状态本地记忆。

**商家二维码**：

- 每个商家拥有独立的店铺二维码
- 二维码用于C端用户扫码访问商家店铺
- 二维码携带商家ID参数，实现店铺隔离

### 1.4 用户端（扫码入口）

**说明**：用户通过扫描商家二维码直接进入指定商家的店铺页面，进行商品浏览、下单购买等操作。

**访问路径**：

```
商家二维码 → 扫码 → 直接进入商家店铺首页 → 操作指引弹窗 → 商品浏览 → 确认订单 → 创建订单
```

**核心功能**：

| 模块   | 功能描述                    |
| ---- | ----------------------- |
| 店铺首页 | 展示商家信息、店铺氛围、商品分类、操作指引弹窗 |
| 商品详情 | 商品图片、价格、规格选择            |
| 购物车  | 商品数量修改、结算               |
| 下单确认 | 确认订单、创建订单（当前不启用在线支付） |
| 我的订单 | 订单列表、订单详情、退款申请（支持按商家筛选） |

**用户订单记录**：

- 用户登录后可查看历史订单记录
- 订单记录按商家维度展示
- 支持订单状态筛选（待支付/待发货/已完成/退款中）
- 支持按商家筛选，仅显示当前商家的订单

**店铺首页优化功能**：

- **操作指引弹窗**：用户首次进入店铺首页时，自动显示操作指引弹窗，包含四步购物流程说明；关闭后不再自动弹出
- **我的订单入口**：在店铺首页商品区上方提供「我的订单」快捷入口，方便用户快速查看当前商家的订单记录
- **底部购物车结算栏**：店铺首页底部固定展示购物车汇总与去购物车/去结算动作，统一与商品详情页底部的视觉风格
- **加购弹窗**：在店铺首页点击商品「+」时弹出规格与数量选择，确认后加入购物车并累计金额
- **自动授权登录**：用户进入店铺首页时自动完成微信授权登录，获取openid作为唯一身份标识，无需手动登录
  - 说明：openid 同时用于访问/行为埋点归因；登录需做并发保护，避免进店时重复触发多次登录请求

**用户授权登录流程**：

1. 用户进入店铺首页
2. 前端自动调用 `uni.login()` 获取微信授权码
3. 调用后端微信登录接口，传入授权码
4. 后端根据授权码获取openid
5. 查询或创建用户记录（openid作为唯一标识）
6. 返回 `token + user(openid)`，前端写入 `user_token` 与 `openid` 并用于后续下单与订单查询
7. 返回JWT token和用户信息
8. 前端保存token，后续请求携带token
9. 用户二次访问时，直接通过openid识别用户，无需重新授权

**openid管理**：

- openid是用户在微信生态下的唯一标识
- 不同小程序的openid不同，但同一小程序的openid不变
- 用户二次访问时，直接通过openid识别用户，无需重新授权

**隐私保护**：

- 仅获取用户openid，不获取手机号等敏感信息
- 用户昵称默认显示为"微信用户"
- 后续可扩展用户主动完善个人信息

### 1.5 服务商管理端（Web/PC 后台）

**说明**：总部后台通过独立的 Web/PC 入口登录，统一管理门店、订单、公告与经营分析能力，与商家端和 C 端隔离。

**当前模式**：

- 系统采用单总部主体模式，不存在平台多租户角色。
- 后台能力已迁移到 `web-admin/`，不再提供服务商小程序入口。
- 后台当前不维护旧支付配置、分账历史或子商户参数。

**访问控制**：

- 后台通过 Web 登录页访问
- 需要后台账户登录验证
- 登录成功后进入总部后台

**核心功能**：

| 模块    | 功能描述               |
| ----- | ------------------ |
| 数据看板  | 门店总数、今日订单、交易额、客单价  |
| 门店管理  | 创建门店、门店信息查看、资料维护 |
| 商家统计  | 商家行业分布、状态分布、趋势分析   |
| 订单统计  | 订单量统计、订单趋势、退款统计    |
| 金额统计  | 交易额统计、TOP门店排行、行业分析 |
| 系统公告接口 | 小程序页面已下线，后端接口保留供后续 PC 端复用 |
| 商家二维码 | 查看/生成商家小程序码        |

**系统公告功能**：

- 服务商公告管理页面已从小程序下线，不再作为当前小程序端交付范围
- 后端保留系统公告相关接口，供后续 PC 端公告管理后台复用
- 商家端仍可继续读取并展示已发布公告

**数据权限**：

- 服务商可查看所有商家的汇总数据
- 支持按商家筛选查看单个商家详情
- 数据实时更新，支持多维度分析

### 1.6 图片存储方案

当前版本统一使用**服务端接收 + 本地文件存储**方案：

```
┌─────────────────────────────────────────────────────────────┐
│                      图片上传流程                            │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
  ┌───────────┐         ┌───────────┐         ┌───────────┐
  │ 商家Logo   │         │ 商品图片   │         │ 营业执照   │
  │ 门头照片   │         │ 商品详情   │         │ 身份证照片 │
  └───────────┘         └───────────┘         └───────────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │  本地文件目录     │
                    │  (按 scope 分类) │
                    └─────────────────┘
```

**上传方式：**

1. 前端以 `multipart/form-data` 方式调用 `POST /api/v1/upload/file`。
2. 服务端按当前登录身份解析 `scope`，写入配置文件指定的本地目录。
3. 接口直接返回 `path`、`url`、`filename`、`scope`，前端使用返回结果回显并保存。

**存储空间规划：**

| 目录范围   | 用途                    | 说明                   |
| ------ | --------------------- | -------------------- |
| admin  | 后台上传的图片与资料文件         | 对应总部后台登录态上传         |
| store/{id} | 门店 Logo、封面、商品图片 | 对应商家登录态上传，按门店隔离    |
| user/{id} | 用户侧补充资料文件        | 仅在有实际上传场景时使用        |
| common | 通用文件                  | 未识别登录身份时的兜底目录       |

**相关接口：**

```
POST /api/v1/upload/file      # 上传文件到本地存储
```

**返回字段：**

- `path`：服务端保存的相对路径
- `url`：当前文件的可访问地址
- `filename`：原始文件名
- `scope`：当前上传归属目录

**商品图片持久化与展示规则：**

- 数据库存储当前接口返回的稳定路径或 URL，不引入上传凭证、回调或动态签名地址。
- 商品列表、商品详情、C 端店铺商品接口统一返回可直接展示的图片地址。
- 更新商品图片时，前端直接提交上一次保存成功的图片路径，不需要二次拼接。

**图片上传完成后的处理说明：**

- 前端上传成功后，需等待业务保存接口调用成功，再将图片视为已保存状态。
- 若仅上传文件但未完成业务保存，图片不应直接出现在商品或门店资料中。
- 文件物理存储目录、访问前缀和大小限制以 `api-chaoshi` 配置文件为准。

### 1.7 门店接入与总部后台

当前系统由总部后台直接创建门店及管理员账号，不再维护旧服务商支付配置。

#### 1.7.1 总部后台说明

当前提供独立的 `web-admin/` Web/PC 后台，通过浏览器访问，与商家端和用户端分离。

**访问路径**：

```text
浏览器打开总部后台 → 进入登录页 → 使用 admin 账号密码登录 → 进入后台管理端
```

**核心功能**：

| 功能模块   | 描述                           |
| ------ | ---------------------------- |
| 后台认证   | 管理员账号密码登录，独立权限体系             |
| 首页数据看板 | 门店总数、订单量、交易额、趋势与概览           |
| 门店管理   | 创建门店、维护资料、维护管理员账号与图片资源       |
| 经营分析   | 门店分布、经营统计、排行榜与订单分析           |
| 订单管理   | 查看门店订单、详情、筛选与汇总              |
| 系统公告接口 | 小程序页面已下线，后台继续承接公告管理能力        |
| 后台设置   | 修改密码、查看后台主体信息                |

#### 1.7.2 门店接入流程

```text
后台登录 → 创建门店账号与基础资料 → 商家登录 → 维护商品/分类/配送设置 → 用户扫码或进入小程序壳 → 无感进入 H5 店铺
```

#### 1.7.3 当前支付口径

- 当前版本已接入“江苏银行微信小程序支付 + xcx 壳承接原生支付”主链路。
- C 端创建订单后不直接返回 `pay_params`，而是返回 `payment.next_action` 与 `payment.prepare_url`。
- 真正的小程序支付参数由 `POST /api/v1/user/orders/:order_id/pay/prepare` 返回，并只供 `xcx` 原生页调用。
- 支付回调通过 `POST /api/v1/payments/jsbank/notify` 更新订单状态。

#### 1.7.4 C端用户下单流程

**访问路径**：

```
用户扫描商家二维码或进入小程序壳 → 完成登录 → 跳转 H5 店铺首页 → 浏览商品 → 选择商品加入购物车 → 确认订单信息 → 选择已开启的下单方式（配送 / 堂食 / 自提）→ 创建订单 → H5 跳回 xcx 支付页 → 小程序完成支付 → 返回 H5 订单详情
```

**订单状态流转**：

| 值 | 状态        | 说明  |
| - | --------- | --- |
| 1 | pending   | 待支付 |
| 2 | paid      | 已支付 |
| 3 | completed | 已完成 |
| 4 | cancelled | 已取消 |
| 5 | refunding | 退款中 |
| 6 | refunded  | 已退款 |

### 1.8 历史支付方案归档

- 原“微信支付服务商模式”、子商户结算、支付回调、自动分账、服务商抽佣等说明均属于历史方案。
- 当前仓库已不再以该方案作为实现基线，相关代码、字段和入口已在本轮改造中逐步清理。
- 若后续接入新的支付方案，应新建独立章节或独立文档，不再沿用本节内容。
- 当前订单、退款与资料文档只保留现行状态定义，不再展开旧支付签名、回调验签、`session_key`、`refresh_token` 等历史细节。

## 2. 功能模块

### 2.1 总部后台模块（Web/PC）

| 功能     | 描述                       | 优先级 |
| ------ | ------------------------ | --- |
| 后台登录   | 管理员账号密码登录               | P0  |
| 首页数据看板 | 门店总数、今日订单、交易额、门店分布与趋势   | P0  |
| 门店管理   | 直接创建门店、维护资料和管理员账号       | P0  |
| 门店分析   | 门店分布、行业统计、排行榜与经营概览      | P0  |
| 订单分析   | 多时间维度统计、订单趋势图           | P0  |
| 门店列表   | 查看所有门店基本信息及运营状况         | P1  |
| 门店详情   | 查看单个门店详细数据              | P1  |
| 系统公告接口 | 小程序页面已下线，后端接口保留供后台复用    | P1  |
| 后台设置   | 修改密码、查看后台主体信息           | P0  |

**系统公告功能说明**：

- 服务商公告管理页面已从小程序下线，不再提供发布、编辑、删除页面
- 后端继续保留 `/api/v1/admin/announcements*` 接口，供后续 PC 端复用
- 商家端可在首页查看最新公告

### 2.2 商家管理模块

| 功能     | 描述                   | 优先级 |
| ------ | -------------------- | --- |
| 商家资料   | 查看和维护商家基础资料          | P0  |
| 商家信息   | 商家基本信息管理             | P0  |
| 商家设置   | 设置首页、快捷入口与账号安全管理     | P0  |
| 满减营销   | 配置多档满多少减多少规则         | P0  |
| 商家二维码  | 生成商家专属小程序码           | P0  |
| 商家状态   | 开启/关闭店铺              | P0  |
| 系统公告查看 | 查看服务商发布的平台公告         | P0  |
| 商家员工   | 员工账号管理               | P2  |

### 2.3 商品分类模块

| 功能   | 描述            | 优先级 |
| ---- | ------------- | --- |
| 分类创建 | 商家创建自己的商品分类   | P0  |
| 分类编辑 | 修改分类名称、排序     | P0  |
| 分类删除 | 删除分类（需检查关联商品） | P0  |
| 分类排序 | 调整分类显示顺序      | P1  |

### 2.4 商品管理模块

| 功能    | 描述        | 优先级 |
| ----- | --------- | --- |
| 商品创建  | 创建商品信息    | P0  |
| 商品编辑  | 修改商品信息    | P0  |
| 商品上下架 | 控制商品是否可售  | P0  |
| 商品删除  | 删除商品（软删除） | P0  |
| 库存管理  | 商品库存数量管理  | P1  |
| 规格管理  | 商品多规格支持   | P2  |

### 2.5 订单管理模块

| 功能   | 描述                  | 优先级 |
| ---- | ------------------- | --- |
| 创建订单 | 用户下单                | P0  |
| 创建后状态 | 创建订单后按 `payment.next_action` 继续支付；在 `xcx` 壳内跳原生支付页，普通 H5 仅提示去小程序支付 | P0  |
| 订单列表 | 商家查看订单列表            | P0  |
| 订单详情 | 查看订单详细信息            | P0  |
| 订单状态 | 待支付/已支付/已完成/已取消/已退款 | P0  |
| 订单核销 | 线下核销订单              | P1  |
| 订单退款 | 处理退款申请              | P1  |

### 2.6 数据分析模块

| 功能   | 描述          | 优先级 |
| ---- | ----------- | --- |
| 销售统计 | 销售额、订单量统计   | P0  |
| 商品分析 | 商品销量排行、库存预警 | P1  |
| 时段分析 | 分时段销售趋势     | P1  |
| 用户分析 | 用户消费行为分析    | P2  |

### 2.7 C端用户模块（扫码进入商家店铺）

| 功能   | 描述             | 优先级 |
| ---- | -------------- | --- |
| 微信登录 | 用户授权登录         | P0  |
| 扫码进店 | 扫描商家二维码进入店铺首页  | P0  |
| 店铺首页 | 展示商家信息、公告、商品分类 | P0  |
| 商品浏览 | 浏览当前商家的商品列表    | P0  |
| 商品详情 | 查看商品详细信息、规格选择  | P0  |
| 购物车  | 加入购物车、修改数量     | P1  |
| 下单确认 | 确认订单、满减优惠、创建订单（当前不启用在线支付） | P0  |
| 我的订单 | 查看个人订单（按商家分组）  | P0  |
| 订单详情 | 查看订单详情、申请退款    | P0  |

### 2.8 历史 WebSocket 能力归档

- 商家端 H5 当前不再建立 WebSocket 长连接，也不展示连接状态与联调入口。
- `/api/v1/ws/merchant`、`/api/v1/dev/order-notify`、`/api/v1/dev/store-visit-notify` 仅保留为历史能力归档，不作为当前实现基线。
- 商家端页面功能以常规 HTTP 接口为准，不依赖实时推送完成主链路。

### 2.9 云打印模块

| 功能     | 描述           | 优先级 |
| ------ | ------------ | --- |
| 打印机配置  | 添加/编辑/删除云打印机 | P1  |
| 打印机管理  | 查看打印机列表、状态、默认设备 | P1  |
| 自动打印开关 | 商家开启/关闭自动打印  | P1  |
| 飞鹅参数配置 | 维护飞鹅账号、UKey、终端号 | P1  |
| 打印模板设置 | 设置小票打印格式     | P2  |
| 打印记录   | 查看历史打印记录     | P2  |
| 打印测试   | 测试打印机连接      | P1  |

**云打印流程**：

```
┌────────┐    ┌────────────┐    ┌────────────┐    ┌────────────┐
│ 订单创建 │ ──▶│  触发打印  │ ──▶│ 云打印API  │ ──▶│  打印机   │
│        │    │  条件判断   │    │  (易联云等) │    │  打印小票  │
└────────┘    └────────────┘    └────────────┘    └────────────┘
```

## 3. API 接口设计

### 3.1 认证相关接口

#### 3.1.1 商家管理员登录

```
POST /api/v1/auth/merchant/login
```

**请求参数：**

```json
{
  "username": "商家账号",
  "password": "密码"
}
```

**开发环境测试账号：**

- `username`: `merchant`
- `password`: `merchant123`

#### 3.1.2 上传文件

```
POST /api/v1/upload/file
Authorization: Bearer {token}
```

**请求方式：**

- `multipart/form-data`
- 文件字段：`file`

**响应：**

```json
{
  "code": 0,
  "data": {
    "path": "store/1/product/20260614-demo.jpg",
    "url": "/uploads/store/1/product/20260614-demo.jpg",
    "filename": "demo.jpg",
    "scope": "store/1"
  }
}
```

**说明：**

- 当前上传由服务端直接写入本地存储，不提供上传凭证与回调接口。
- 前端直接使用返回的 `path` 或 `url` 参与后续业务保存。

#### 3.1.3 C端用户微信登录

```
POST /api/v1/user/auth/wechat-login
```

**请求参数：**

```json
{
  "code": "微信登录凭证",
  "nickname": "用户昵称",
  "avatar": "头像URL"
}
```

### 3.2 总部后台接口（前缀 /api/v1/admin）

> **说明**：本节保留后台管理接口总览，当前生效前缀为 `/api/v1/admin/*`；若个别历史示例仍使用旧命名，以拆分文档中的最新接口说明为准。

#### 3.2.1 后台登录

```
POST /api/v1/admin/auth/login
```

**请求参数：**

```json
{
  "username": "后台账号",
  "password": "密码"
}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "token": "JWT Token",
    "admin": {
      "id": 1,
      "name": "超市总部",
      "display_name": "总部管理员"
    }
  }
}
```

#### 3.2.2 服务商首页数据看板

```
GET /api/v1/admin/dashboard
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "summary": {
      "total_merchants": 128,
      "today_orders": 1256,
      "today_amount": 58960.00,
      "avg_order_amount": 46.95
    },
    "merchant_distribution": {
      "categories": [
        {"name": "餐饮", "count": 58, "percentage": 45.3},
        {"name": "零售", "count": 38, "percentage": 29.7},
        {"name": "服务", "count": 19, "percentage": 14.8},
        {"name": "其他", "count": 13, "percentage": 10.2}
      ]
    },
    "order_trend": {
      "dates": ["2024-01-01", "2024-01-02", "..."],
      "orders": [120, 145, "..."]
    },
    "pending_tasks": {
      "merchant_approvals": 3,
      "order_issues": 12
    }
  }
}
```

#### 3.2.3 创建商家

```
POST /api/v1/admin/stores
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "name": "美味餐厅",
  "contact_name": "张三",
  "contact_phone": "13800138000",
  "contact_email": "merchant@example.com",
  "address": "XX市XX区XX路XX号",
  "business_category": "餐饮",
  "business_hours": "09:00-22:00",
  "announcement": "欢迎光临",
  "username": "merchant001",
  "password": "merchant123"
}
```

**说明：**

- 门店由总部后台直接创建，不存在商家自助入驻申请。
- 当前创建接口仅维护基础资料与管理员账号，不包含支付配置字段。

#### 3.2.4 商家详情查看

```
GET /api/v1/admin/stores/{merchant_id}
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "美味餐厅",
    "contact_name": "张三",
    "contact_phone": "13800138000",
    "business_category": "餐饮",
    "address": "店铺地址",
    "status": 1,
    "qrcode_url": "",
    "created_at": "2024-01-01T10:00:00Z",
    "settings": {
      "store_images": []
    },
    "total_orders": 0,
    "total_amount": 0,
    "total_users": 0
  }
}
```

#### 3.2.5 更新商家资料

```
PUT /api/v1/admin/stores/{merchant_id}
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "name": "美味餐厅旗舰店",
  "contact_name": "张三",
  "contact_phone": "13800138000",
  "contact_email": "merchant@example.com",
  "address": "XX市XX区XX路XX号",
  "business_category": "餐饮",
  "business_hours": "10:00-22:00",
  "announcement": "欢迎光临"
}
```

**响应：**

```json
{
  "code": 0,
  "message": "更新成功"
}
```

#### 3.2.6 历史支付配置接口归档

- `PUT /api/v1/admin/stores/{merchant_id}/payment-config` 对应的旧支付配置能力已停用。
- 当前总部后台不再维护 `sub_mch_id`、支付配置状态、分账开关或抽佣比例。
- 后续如接入新支付方案，须新增独立接口而不是恢复本节旧定义。

#### 3.2.7 商家数据分析

```
GET /api/v1/admin/stores/analytics/distribution
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "by_category": [
      {"name": "餐饮", "count": 58, "percentage": 45.3},
      {"name": "零售", "count": 38, "percentage": 29.7},
      {"name": "服务", "count": 19, "percentage": 14.8},
      {"name": "其他", "count": 13, "percentage": 10.2}
    ],
    "by_status": [
      {"name": "营业中", "count": 100, "percentage": 78.1},
      {"name": "休息中", "count": 15, "percentage": 11.7},
      {"name": "已关闭", "count": 13, "percentage": 10.2}
    ],
    "by_month": [
      {"month": "2024-01", "count": 15},
      {"month": "2024-02", "count": 22}
    ]
  }
}
```

#### 3.2.8 商家列表

```
GET /api/v1/admin/stores/list
Authorization: Bearer {token}
```

**请求参数：**

| 参数         | 类型     | 必填 | 描述    |
| ---------- | ------ | -- | ----- |
| page       | int    | 否  | 页码    |
| page\_size | int    | 否  | 每页数量  |
| keyword    | string | 否  | 搜索关键词 |
| category   | string | 否  | 行业分类  |

**响应：**

```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "name": "美味餐厅",
        "business_category": "餐饮",
        "status": "active",
        "total_orders": 1256,
        "total_amount": 58960.00,
        "created_at": "2024-01-01"
      }
    ],
    "total": 128,
    "page": 1,
    "page_size": 10
  }
}
```

#### 3.2.9 订单数据分析

```
GET /api/v1/admin/orders/analytics
Authorization: Bearer {token}
```

**请求参数：**

| 参数          | 类型     | 必填 | 描述                              |
| ----------- | ------ | -- | ------------------------------- |
| period      | string | 否  | 周期：today/week/month/year/custom |
| start\_date | string | 否  | 开始日期                            |
| end\_date   | string | 否  | 结束日期                            |

**响应：**

```json
{
  "code": 0,
  "data": {
    "summary": {
      "total_orders": 125600,
      "total_amount": 5896000.00,
      "today_orders": 1256,
      "today_amount": 58960.00,
      "avg_order_amount": 46.95
    },
    "trend": [
      {"date": "2024-01-01", "orders": 1200, "amount": 56400.00}
    ],
    "by_delivery_type": [
      {"type": "配送", "count": 800, "percentage": 63.7},
      {"type": "堂食", "count": 300, "percentage": 23.9},
      {"type": "自提", "count": 156, "percentage": 12.4}
    ]
  }
}
```

#### 3.2.10 金额数据分析

```
GET /api/v1/admin/amount/analytics
Authorization: Bearer {token}
```

**请求参数：**

| 参数          | 类型     | 必填 | 描述                       |
| ----------- | ------ | -- | ------------------------ |
| period      | string | 否  | 周期：today/week/month/year |
| start\_date | string | 否  | 开始日期                     |
| end\_date   | string | 否  | 结束日期                     |

**响应：**

```json
{
  "code": 0,
  "data": {
    "summary": {
      "total_amount": 5896000.00,
      "today_amount": 58960.00,
      "week_amount": 412720.00,
      "month_amount": 1766400.00,
      "growth_rate": 15.5
    },
    "trend": [
      {"date": "2024-01-01", "amount": 56400.00}
    ],
    "by_merchant": [
      {"merchant_id": 1, "merchant_name": "美味餐厅", "amount": 589600.00, "percentage": 10.0}
    ],
    "by_category": [
      {"category": "餐饮", "amount": 2948000.00, "percentage": 50.0}
    ]
  }
}
```

#### 3.2.11 TOP商家排行

```
GET /api/v1/admin/amount/top-merchants
Authorization: Bearer {token}
```

**请求参数：**

| 参数     | 类型     | 必填 | 描述        |
| ------ | ------ | -- | --------- |
| period | string | 否  | 周期        |
| limit  | int    | 否  | 返回数量，默认10 |

**响应：**

```json
{
  "code": 0,
  "data": [
    {"rank": 1, "merchant_id": 1, "merchant_name": "美味餐厅", "amount": 589600.00},
    {"rank": 2, "merchant_id": 2, "merchant_name": "隔壁小馆", "amount": 412720.00}
  ]
}
```

#### 3.2.12 历史分账记录接口归档

- `GET /api/v1/admin/profit-sharing-records` 所对应的旧分账记录能力已停用。
- 当前版本总部后台不再展示或查询分账历史。
- 根 PRD 不再维护该接口的请求/响应样例，避免与现有实现冲突。

#### 3.2.13 商家年费管理

```
GET /api/v1/admin/stores/{id}/fee
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "merchant_id": 123,
    "merchant_name": "美味餐厅",
    "fees": [
      {
        "year": 2024,
        "amount": 365.00,
        "status": "paid",
        "pay_time": "2024-01-01 00:00:00",
        "free_reason": ""
      },
      {
        "year": 2025,
        "amount": 0.00,
        "status": "free",
        "pay_time": null,
        "free_reason": "限时优惠免年费"
      }
    ]
  }
}
```

#### 3.2.14 服务商退出登录

```
POST /api/v1/admin/auth/logout
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "message": "退出成功"
}
```

#### 3.2.15 获取商家二维码

```
GET /api/v1/admin/stores/{id}/qrcode
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "merchant_id": 123,
    "merchant_name": "美味餐厅",
    "qrcode_url": "data:image/png;base64,...",
    "page_path": "pages/store/home",
    "scene": "merchant_id=123"
  }
}
```

#### 3.2.16 系统公告列表

说明：服务商公告管理页面已从当前小程序下线，以下接口保留供后续 PC 端公告管理后台复用；商家端公告展示链路仍可继续读取已发布公告。

```
GET /api/v1/admin/announcements
Authorization: Bearer {token}
```

**请求参数：**

| 参数         | 类型  | 必填 | 描述   |
| ---------- | --- | -- | ---- |
| page       | int | 否  | 页码   |
| page\_size | int | 否  | 每页数量 |

**响应：**

```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "title": "系统升级通知",
        "content": "平台将于本周六进行系统升级",
        "status": 1,
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "total": 10,
      "page": 1,
      "page_size": 10
    }
  }
}
```

#### 3.2.17 系统公告详情

```
GET /api/v1/admin/announcements/{id}
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "title": "系统升级通知",
    "content": "平台将于本周六进行系统升级",
    "status": 1,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### 3.2.18 创建系统公告

```
POST /api/v1/admin/announcements
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "title": "系统公告标题",
  "content": "公告内容",
  "status": 1
}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "message": "创建成功"
  }
}
```

#### 3.2.19 更新系统公告

```
PUT /api/v1/admin/announcements/{id}
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "title": "更新后的标题",
  "content": "更新后的内容",
  "status": 1
}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "title": "更新后的标题",
    "content": "更新后的内容",
    "status": 1,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### 3.2.20 删除系统公告

```
DELETE /api/v1/admin/announcements/{id}
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "message": "删除成功"
}
```

#### 3.2.21 退款订单查询(服务商)

```
GET /api/v1/admin/orders/refunds
Authorization: Bearer {token}
```

**请求参数：**

| 参数           | 类型     | 必填 | 描述   |
| ------------ | ------ | -- | ---- |
| page         | int    | 否  | 页码   |
| page\_size   | int    | 否  | 每页数量 |
| merchant\_id | int    | 否  | 商家ID |
| status       | string | 否  | 退款状态 |
| start\_time  | string | 否  | 开始时间 |
| end\_time    | string | 否  | 结束时间 |

**响应：**

```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "order_id": "ORD202401150001",
        "merchant_name": "美味餐厅",
        "amount": 100.00,
        "refund_amount": 100.00,
        "status": "success",
        "create_time": "2024-01-15 10:30:00",
        "complete_time": "2024-01-15 10:35:00"
      }
    ],
    "total": 20
  }
}
```

#### 3.2.22 服务商设置

```
GET /api/v1/admin/settings
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "name": "服务商名称",
    "sp_name": "服务商姓名",
    "contact_phone": "13800138000",
    "contact_email": "",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 3.2.23 更新服务商设置

```
PUT /api/v1/admin/settings
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "contact_name": "新联系人",
  "contact_phone": "13900139000"
}
```

**响应：**

```json
{
  "code": 0,
  "message": "设置成功"
}
```

#### 3.2.24 服务商修改密码

```
POST /api/v1/admin/account/change-password
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "old_password": "旧密码",
  "new_password": "新密码"
}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "message": "修改成功"
  }
}
```

#### 3.2.25 商家图片资产更新

```
PUT /api/v1/admin/stores/{merchant_id}/assets
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "logo": "https://cdn.example.com/merchant/logo.png",
  "cover_image": "https://cdn.example.com/merchant/cover.png"
}
```

**说明：**

- 仅允许服务商更新其名下商家的图片资源。
- 图片更新后需同步影响商家端与 C 端展示。

#### 3.2.26 历史支付配置说明归档

- 本节对应的旧服务商支付配置口径已废弃。
- 当前总部后台不再维护支付配置状态、子商户号或分账比例。

### 3.3 C端用户接口（前缀 /api/v1/user）

#### 3.3.1 用户收货地址列表

```
GET /api/v1/user/addresses
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "name": "张三",
      "phone": "13800138000",
      "province": "广东省",
      "city": "深圳市",
      "district": "南山区",
      "detail": "科技园路88号",
      "is_default": 1
    }
  ]
}
```

#### 3.3.2 添加收货地址

```
POST /api/v1/user/addresses
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "name": "李四",
  "phone": "13900139000",
  "province": "广东省",
  "city": "广州市",
  "district": "天河区",
  "detail": "珠江新城XX号",
  "is_default": 1
}
```

**响应：**

```json
{
  "code": 0,
  "data": {"id": 2}
}
```

#### 3.3.3 更新收货地址

```
PUT /api/v1/user/addresses/{id}
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "name": "李四",
  "phone": "13900139000",
  "province": "广东省",
  "city": "广州市",
  "district": "天河区",
  "detail": "珠江新城XX号",
  "is_default": 1
}
```

**响应：**

```json
{
  "code": 0,
  "message": "更新成功"
}
```

#### 3.3.4 删除收货地址

```
DELETE /api/v1/user/addresses/{id}
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "message": "删除成功"
}
```

### 3.4 商家管理接口（前缀 /api/v1/merchant）

#### 3.4.1 获取商家信息

```
GET /api/v1/merchant/profile
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "美味餐厅",
    "logo": "店铺Logo",
    "contact_name": "张三",
    "contact_phone": "13800138000",
    "address": "XX市XX区XX路XX号",
    "business_category": "餐饮",
    "status": "active",
    "settings": {
      "announcement": "今日特惠：全场8折",
      "business_hours": "09:00-22:00",
      "min_order_amount": 20.00,
      "delivery_fee": 5.00
    },
    "qrcode_url": "",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 3.4.2 更新商家基本信息

```
PUT /api/v1/merchant/profile
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "name": "商家名称",
  "logo": "店铺Logo URL",
  "contact_name": "联系人",
  "contact_phone": "联系电话",
  "address": "店铺地址"
}
```

#### 3.4.3 商家自定义设置（获取/更新）

```
GET /api/v1/merchant/settings
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "announcement": "门店公告内容",
    "business_hours": "09:00-22:00",
    "min_order_amount": 20.00,
    "takeout_enabled": true,
    "dine_in_enabled": true,
    "pickup_enabled": true,
    "notify_enabled": true,
    "browse_notify_enabled": true,
    "push_openid": "push_openid_xxx",
    "delivery_settings": {
      "enabled": true,
      "base_fee": 5.00,
      "free_delivery_amount": 50.00,
      "distance_rules": [
        {"min_distance": 0, "max_distance": 2, "fee": 0},
        {"min_distance": 2, "max_distance": 5, "fee": 3.00},
        {"min_distance": 5, "max_distance": 10, "fee": 6.00}
      ],
      "max_distance": 10
    }
  }
}
```

```
PUT /api/v1/merchant/settings
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "takeout_enabled": true,
  "dine_in_enabled": true,
  "pickup_enabled": true,
  "notify_enabled": true,
  "browse_notify_enabled": true
}
```

**说明：**

- 设置页包含"满减营销"、"打印机管理"、"修改密码"和"员工管理"入口（仅店主可见）。
- `takeout_enabled`、`dine_in_enabled`、`pickup_enabled` 均为商家主表维度开关，分别控制配送、堂食、自提是否可下单。
- `delivery_settings.enabled` 只表示配送费规则是否生效，确认页是否展示“配送”仍以后端返回的 `takeout_enabled` 为准。
- 本轮不包含优惠券、发券和复杂营销活动，但新增商家满减营销配置。
- 满减规则通过 `GET/PUT /api/v1/merchant/full-reduction-rules` 维护，用户确认订单页通过 `GET /api/v1/store/{merchant_id}/full-reduction-rules` 获取当前启用规则。

**配送距离规则说明：**

| 字段                     | 类型      | 描述         |
| ---------------------- | ------- | ---------- |
| enabled                | boolean | 是否开启配送     |
| base\_fee              | decimal | 基础配送费      |
| free\_delivery\_amount | decimal | 满额免配送费金额   |
| distance\_rules        | array   | 按距离收费规则    |
| max\_distance          | int     | 最大配送距离（公里） |

**distance\_rules 结构：**

| 字段            | 类型      | 描述           |
| ------------- | ------- | ------------ |
| min\_distance | int     | 最小距离（公里），包含  |
| max\_distance | int     | 最大距离（公里），不包含 |
| fee           | decimal | 该距离范围内的配送费   |

- 配送距离是"商家提供给用户选择的服务范围档位"，不是地图定位距离。
- C 端用户下单时手动选择距离档位；超出商家支持范围时仅提示，且不接入第三方地图。

**响应：**

```json
{
  "code": 0,
  "message": "设置更新成功"
}
```

#### 3.4.4 配送设置（获取/更新）

```
GET /api/v1/merchant/delivery-settings
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "enabled": true,
    "base_fee": 5.00,
    "free_delivery_amount": 50.00,
    "distance_rules": [
      {"min_distance": 0, "max_distance": 2, "fee": 0},
      {"min_distance": 2, "max_distance": 5, "fee": 3.00},
      {"min_distance": 5, "max_distance": 10, "fee": 6.00}
    ],
    "max_distance": 10
  }
}
```

```
PUT /api/v1/merchant/delivery-settings
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "enabled": true,
  "base_fee": 5.00,
  "free_delivery_amount": 50.00,
  "distance_rules": [
    {"min_distance": 0, "max_distance": 2, "fee": 0},
    {"min_distance": 2, "max_distance": 5, "fee": 3.00},
    {"min_distance": 5, "max_distance": 10, "fee": 6.00}
  ],
  "max_distance": 10
}
```

**说明：**

- distance\_rules 允许传空数组，用于清空"按距离收费"的所有规则。
- 所有规则必须满足 `max_distance > min_distance`、`fee >= 0`，且规则之间不能重叠。
- 每条规则的 `max_distance` 不能超过 `max_distance` 字段本身。
- 商家端保存时需要前后端双重校验；C 端仅消费商家保存后的档位结果。

**响应：**

```json
{
  "code": 0,
  "message": "更新成功"
}
```

#### 3.4.5 开启/关闭店铺

```
POST /api/v1/merchant/status
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "status": 1
}
```

**说明：**

- `status = 1` 表示营业中，`status = 0` 表示休息中。
- 商家切换为休息后，用户无法再创建新订单。
- 已进入支付流程的订单不因商家切换休息而支付失败，支付成功后仍提示下单成功。

#### 3.4.6 历史分账历史接口归档

- 商家端分账历史页面与 `GET /api/v1/merchant/profit-sharing-records` 已从当前版本移除。
- 当前商家设置不再承接分账历史入口，避免与后续新支付方案的账务设计冲突。

#### 3.4.7 获取商家小程序码

```
GET /api/v1/merchant/qrcode
Authorization: Bearer {token}
```

**说明：**

- 该接口固定为商家生成 C 端店铺首页二维码，页面路径为 `pages/store/home`。
- 二维码 `scene` 使用 `merchant_id={当前商家ID}`，以便 `parseStoreEntryOptions()` 按现有规则解析商家入口。
- 返回体中的 `scene` 与 `page` 为调试字段，需与实际二维码生成配置保持一致。
- 二维码为接口动态生成结果，不再写回 `merchants` 表持久化保存。

**响应：**

```json
{
  "code": 0,
  "data": {
    "qrcode_url": "小程序码图片URL",
    "scene": "merchant_id=123",
    "page": "pages/store/home",
    "placeholder": false,
    "message": "微信小程序码生成成功",
    "expire_time": null
  }
}
```

#### 3.4.8 获取系统公告列表

```
GET /api/v1/merchant/announcements
Authorization: Bearer {token}
```

**请求参数：**

| 参数         | 类型  | 必填 | 描述   |
| ---------- | --- | -- | ---- |
| page       | int | 否  | 页码   |
| page\_size | int | 否  | 每页数量 |

**响应：**

```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "title": "系统升级通知",
        "content": "平台将于本周六进行系统升级",
        "create_time": "2024-01-15 10:30:00"
      }
    ],
    "total": 10
  }
}
```

#### 3.4.9 获取系统公告详情

```
GET /api/v1/merchant/announcements/{id}
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "title": "系统升级通知",
    "content": "平台将于本周六进行系统升级，请各位商家提前做好相关准备。",
    "create_time": "2024-01-15 10:30:00"
  }
}
```

#### 3.4.10 历史 WebSocket 能力归档

- `/api/v1/ws/merchant`、`/api/v1/dev/order-notify`、`/api/v1/dev/store-visit-notify` 已从当前实现中移除。
- 当前商家端 H5 不建立 WebSocket 长连接，也不提供联调测试页。
- 若后续恢复实时提醒能力，应以新的接口协议重新定义。

```json
{
  "type": "store_visit_notify",
  "payload": {
    "merchant_id": 1,
    "visitor_openid": "wx_xxx",
    "source": "scan"
  }
}
```

#### 3.4.11 历史联调接口归档

- `POST /api/v1/dev/order-notify` 与 `POST /api/v1/dev/store-visit-notify` 已从当前实现中移除。
- 当前版本不再提供实时提醒联调入口。

#### 3.4.12 修改当前员工密码

```
POST /api/v1/merchant/account/change-password
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "old_password": "merchant123",
  "new_password": "merchant456"
}
```

**说明：**

- 修改密码成功后，前端需强制当前员工退出并重新登录。
- 新旧密码不能相同，新密码最少 6 位。

#### 3.4.13 员工推送 OpenID 说明

- 当前版本不再提供商家员工微信快捷登录、绑定、解绑接口。
- `merchant_staffs.push_openid` 仅作为订阅消息或提醒推送标识使用，不承担登录鉴权语义。
- 员工登录方式统一为账号密码登录。

#### 3.4.14 获取员工列表

```
GET /api/v1/merchant/staff
Authorization: Bearer {token}
```

**请求参数：**

| 参数         | 类型  | 必填 | 描述   |
| ---------- | --- | -- | ---- |
| page       | int | 否  | 页码   |
| page\_size | int | 否  | 每页数量 |

**响应：**

```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "username": "zhangsan",
        "name": "张三",
        "phone": "13800138000",
        "role": "staff",
        "push_openid": "push_openid_xxx",
        "notify_enabled": true,
        "browse_notify_enabled": true,
        "last_login_at": "2024-01-02T10:00:00Z",
        "status": 1,
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 5
  }
}
```

#### 3.4.15 添加员工

```
POST /api/v1/merchant/staff
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "name": "员工姓名",
  "phone": "手机号",
  "username": "登录账号",
  "password": "初始密码",
  "role": "staff"
}
```

**响应：**

```json
{
  "code": 0,
  "message": "添加成功",
  "data": {
    "id": 1
  }
}
```

#### 3.4.18 更新员工信息

```
PUT /api/v1/merchant/staff/{id}
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "name": "员工姓名",
  "phone": "手机号",
  "role": "staff",
  "notify_enabled": true,
  "browse_notify_enabled": true,
  "status": 1
}
```

**响应：**

```json
{
  "code": 0,
  "message": "更新成功"
}
```

#### 3.4.19 删除员工

```
DELETE /api/v1/merchant/staff/{id}
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "message": "删除成功"
}
```

**说明：**

- 员工管理首版仅 `owner` 可见和可操作。
- 不允许删除当前登录账号，也不允许删除最后一个 `owner`。

#### 3.4.20 重置员工密码

```
POST /api/v1/merchant/staff/{id}/reset-password
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "new_password": "merchant456"
}
```

#### 3.4.21 获取云打印机列表

```
GET /api/v1/merchant/printers
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "name": "前台打印机",
      "type": "yilianyun",
      "device_no": "设备编号",
      "api_url": "https://api.example.com",
      "feie_user": "",
      "feie_sn": "",
      "status": 1,
      "auto_print": true,
      "is_default": true,
      "print_count": 156,
      "has_api_key": true,
      "has_feie_ukey": false
    }
  ]
}
```

#### 3.4.22 添加云打印机

```
POST /api/v1/merchant/printers
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "name": "前台打印机",
  "type": "yilianyun",
  "device_no": "设备编号",
  "api_key": "API密钥",
  "api_url": "https://api.example.com",
  "print_types": ["order"],
  "auto_print": true,
  "is_default": false
}
```

**飞鹅打印机附加字段：**

```json
{
  "type": "feie",
  "device_no": "打印机设备编号",
  "feie_user": "飞鹅账号",
  "feie_ukey": "飞鹅UKey",
  "feie_sn": "飞鹅终端号"
}
```

#### 3.4.23 更新云打印机

```
PUT /api/v1/merchant/printers/{printer_id}
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "name": "打印机名称",
  "auto_print": true,
  "is_default": true,
  "status": 1
}
```

#### 3.4.24 删除云打印机

```
DELETE /api/v1/merchant/printers/{printer_id}
Authorization: Bearer {token}
```

#### 3.4.25 测试云打印机

```
POST /api/v1/merchant/printers/{printer_id}/test
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "success": true,
    "message": "打印测试成功"
  }
}
```

#### 3.4.26 获取打印记录

```
GET /api/v1/merchant/print-logs
Authorization: Bearer {token}
```

**请求参数：**

| 参数          | 类型     | 必填 | 描述   |
| ----------- | ------ | -- | ---- |
| page        | int    | 否  | 页码   |
| page\_size  | int    | 否  | 每页数量 |
| start\_date | string | 否  | 开始日期 |
| end\_date   | string | 否  | 结束日期 |

### 3.5 商品分类接口（前缀 /api/v1/merchant）

#### 3.5.1 获取分类列表

```
GET /api/v1/merchant/categories
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "name": "热销推荐",
      "sort": 1,
      "product_count": 10,
      "status": "active",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### 3.5.2 创建分类

```
POST /api/v1/merchant/categories
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "name": "分类名称",
  "sort": 1
}
```

#### 3.5.3 更新分类

```
PUT /api/v1/merchant/categories/{category_id}
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "name": "新分类名称",
  "sort": 2
}
```

#### 3.5.4 删除分类

```
DELETE /api/v1/merchant/categories/{category_id}
Authorization: Bearer {token}
```

#### 3.5.5 批量排序分类

```
POST /api/v1/merchant/categories/sort
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "orders": [
    {"id": 1, "sort": 1},
    {"id": 2, "sort": 2}
  ]
}
```

### 3.6 商品管理接口（前缀 /api/v1/merchant）

#### 3.6.1 商品列表

```
GET /api/v1/merchant/products
Authorization: Bearer {token}
```

**请求参数：**

| 参数           | 类型     | 必填 | 描述                    |
| ------------ | ------ | -- | --------------------- |
| page         | int    | 否  | 页码                    |
| page\_size   | int    | 否  | 每页数量                  |
| category\_id | int    | 否  | 分类ID                  |
| status       | string | 否  | 状态：on\_sale/off\_sale |
| keyword      | string | 否  | 搜索关键词                 |

**响应：**

```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "name": "招牌红烧肉",
        "images": ["图片URL"],
        "price": 58.00,
        "original_price": 68.00,
        "stock": 100,
        "sales": 256,
        "category_id": 1,
        "category_name": "热销推荐",
        "status": "on_sale",
        "sort": 1,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10
  }
}
```

#### 3.6.2 商品详情

```
GET /api/v1/merchant/products/{product_id}
Authorization: Bearer {token}
```

**响应特点：**

- 商品详情返回结构与创建商品、更新商品成功后的回读结构保持一致。
- 图片字段统一为 `images` 数组，返回值为可直接显示的图片地址。
- 规格字段统一为结构化 `specs` 数组，便于商家端详情页、编辑页直接复用。

**错误语义：**

- 当商品不存在、不属于当前商家或已删除时，接口返回 `404`
- 当分类、规格等关系预加载失败，或服务端在详情查询/回读阶段发生内部异常时，接口返回 `500`
- 创建商品、更新商品在写入成功后的详情回读阶段，如果回读失败，同样返回 `500`

#### 3.6.3 创建商品

```
POST /api/v1/merchant/products
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "name": "商品名称",
  "description": "商品描述",
  "images": ["图片URL1", "图片URL2"],
  "category_id": 1,
  "price": 58.00,
  "original_price": 68.00,
  "stock": 100,
  "unit": "份",
  "sort": 1,
  "specs": [
    {
      "name": "规格",
      "options": [
        {"name": "小份", "price": 48.00, "stock": 50},
        {"name": "大份", "price": 68.00, "stock": 50}
      ]
    }
  ]
}
```

#### 3.6.4 更新商品

```
PUT /api/v1/merchant/products/{product_id}
Authorization: Bearer {token}
```

#### 3.6.5 商品上架

```
POST /api/v1/merchant/products/{product_id}/on-sale
Authorization: Bearer {token}
```

#### 3.6.6 商品下架

```
POST /api/v1/merchant/products/{product_id}/off-sale
Authorization: Bearer {token}
```

#### 3.6.7 批量上下架

```
POST /api/v1/merchant/products/batch-status
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "product_ids": [1, 2, 3],
  "status": "on_sale"
}
```

#### 3.6.8 删除商品

```
DELETE /api/v1/merchant/products/{product_id}
Authorization: Bearer {token}
```

**说明：**

- 删除商品采用软删除语义，已删除商品不会继续出现在商家商品列表与 C 端可售商品列表中。

#### 3.6.9 更新库存

```
PUT /api/v1/merchant/products/{product_id}/stock
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "stock": 100,
  "action": "set"
}
```

| action   | 描述     |
| -------- | ------ |
| set      | 设置为指定值 |
| add      | 增加指定数量 |
| subtract | 减少指定数量 |

#### 3.6.10 获取商品规格

```
GET /api/v1/merchant/products/{product_id}/specs
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "specs": [
      {
        "id": 1,
        "name": "颜色",
        "values": ["红色", "蓝色", "绿色"]
      },
      {
        "id": 2,
        "name": "尺码",
        "values": ["S", "M", "L", "XL"]
      }
    ],
    "skus": [
      {
        "id": 1,
        "specs": [{"spec_id": 1, "value": "红色"}, {"spec_id": 2, "value": "S"}],
        "price": 29.90,
        "stock": 100,
        "sku_code": "RED-S-001",
        "status": "active"
      }
    ]
  }
}
```

实现说明：

- 当前版本暂未实现 `skus` 维度的持久化与校验，接口会返回 `skus: []`，并在保存时忽略请求中的 `skus` 字段。

#### 3.6.11 创建/更新商品规格

```
PUT /api/v1/merchant/products/{product_id}/specs
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "specs": [
    {"id": 1, "name": "颜色", "values": ["红色", "蓝色"]},
    {"id": 2, "name": "尺码", "values": ["S", "M", "L"]}
  ],
  "skus": [
    {
      "specs": [{"spec_id": 1, "value": "红色"}, {"spec_id": 2, "value": "S"}],
      "price": 29.90,
      "stock": 100,
      "sku_code": "RED-S-001"
    },
    {
      "specs": [{"spec_id": 1, "value": "红色"}, {"spec_id": 2, "value": "M"}],
      "price": 29.90,
      "stock": 80,
      "sku_code": "RED-M-001"
    }
  ]
}
```

**响应：**

```json
{
  "code": 0,
  "message": "保存成功"
}
```

#### 3.6.12 删除商品规格

```
DELETE /api/v1/merchant/products/{product_id}/specs
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "spec_ids": [1, 2]
}
```

**响应：**

```json
{
  "code": 0,
  "message": "删除成功"
}
```

### 3.7 订单管理接口（前缀 /api/v1/merchant）

#### 3.7.1 订单列表

```
GET /api/v1/merchant/orders
Authorization: Bearer {token}
```

**请求参数：**

| 参数          | 类型     | 必填 | 描述    |
| ----------- | ------ | -- | ----- |
| page        | int    | 否  | 页码    |
| page\_size  | int    | 否  | 每页数量  |
| status      | string | 否  | 订单状态  |
| start\_date | string | 否  | 开始日期  |
| end\_date   | string | 否  | 结束日期  |
| order\_no   | string | 否  | 订单号搜索 |

**响应：**

```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "order_no": "202401010001",
        "user": {
          "id": 1,
          "nickname": "用户昵称",
          "avatar": "头像URL"
        },
        "total_amount": 116.00,
        "pay_amount": 106.00,
        "discount_amount": 10.00,
        "status": "paid",
        "items": [
          {
            "product_id": 1,
            "product_name": "招牌红烧肉",
            "image": "图片URL",
            "price": 58.00,
            "quantity": 2,
            "specs": "大份"
          }
        ],
        "remark": "少放辣",
        "created_at": "2024-01-01T12:00:00Z",
        "paid_at": "2024-01-01T12:01:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

#### 3.7.2 订单详情

```
GET /api/v1/merchant/orders/{order_id}
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "order_no": "202401010001",
    "user": {
      "id": 1,
      "nickname": "用户昵称",
      "avatar": "头像URL",
      "phone": "138****8000"
    },
    "merchant": {
      "id": 1,
      "name": "美味餐厅",
      "address": "店铺地址",
      "phone": "店铺电话"
    },
    "items": [
      {
        "product_id": 1,
        "product_name": "招牌红烧肉",
        "image": "图片URL",
        "price": 58.00,
        "quantity": 2,
        "specs": "大份",
        "subtotal": 116.00
      }
    ],
    "total_amount": 116.00,
    "pay_amount": 106.00,
    "discount_amount": 10.00,
    "delivery_fee": 0,
    "status": "paid",
    "remark": "少放辣",
    "transaction_id": "支付交易号",
    "created_at": "2024-01-01T12:00:00Z",
    "paid_at": "2024-01-01T12:01:00Z",
    "completed_at": null,
    "completed_by_name": null,
    "refunded_at": null
  }
}
```

#### 3.7.3 订单核销

```
POST /api/v1/merchant/orders/{order_id}/complete
Authorization: Bearer {token}
```

> **重要说明**：商家只能核销**属于自己店铺**的订单，不能跨商家核销。系统会验证订单的 merchant\_id 是否与当前登录商家一致。

**请求参数：**

```json
{
  "verify_code": "核销码"
}
```

**核销码规则：**

- verify\_code 必须为 6 位数字字符串（示例：`123456`）。

**响应：**

```json
{
  "code": 0,
  "message": "核销成功",
  "data": {
    "id": 1,
    "order_no": "202401010001",
    "verify_code": "123456",
    "status": 3,
    "completed_at": "2024-01-01T14:00:00Z",
    "completed_by_name": "收银员A"
  }
}
```

#### 3.7.4 按核销码快速核销

```
POST /api/v1/merchant/orders/quick-complete
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "verify_code": "123456"
}
```

**说明：**

- 用于商家工作台的"扫一扫 / 输入核销码快速核销"入口。
- 系统会在当前登录商家下查找匹配核销码且状态为"已支付"的订单并完成核销。

**错误情况：**

| 错误码  | 描述             |
| ---- | -------------- |
| 5001 | 订单不存在          |
| 5002 | 订单状态错误（非已支付状态） |
| 5005 | 无权操作此订单（跨商家核销） |
| 5006 | 核销码错误          |

#### 3.7.5 订单退款

```
POST /api/v1/merchant/orders/{order_id}/refund
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "refund_amount": 106.00,
  "refund_reason": "退款原因"
}
```

#### 3.7.6 订单统计

```
GET /api/v1/merchant/orders/statistics
Authorization: Bearer {token}
```

**请求参数：**

| 参数          | 类型     | 必填 | 描述   |
| ----------- | ------ | -- | ---- |
| start\_date | string | 否  | 开始日期 |
| end\_date   | string | 否  | 结束日期 |

**响应：**

```json
{
  "code": 0,
  "data": {
    "total_orders": 1000,
    "total_sales": 50000.00,
    "pending_payment": 10,
    "pending_complete": 50,
    "completed": 900,
    "refunded": 40,
    "cancelled": 0
  }
}
```

### 3.8 数据分析接口（前缀 /api/v1/merchant）

> **说明**：当前商家数据看板以订单数据和 `user_behavior_events` 最小行为事件表为基础，重点覆盖访客、浏览、下单、支付成功等核心经营指标，不扩展复杂 BI 分析能力。

#### 3.8.1 销售概览

```
GET /api/v1/merchant/analytics/overview
Authorization: Bearer {token}
```

**请求参数：**

| 参数     | 类型     | 必填 | 描述                       |
| ------ | ------ | -- | ------------------------ |
| period | string | 否  | 周期：today/week/month/year |

**响应：**

```json
{
  "code": 0,
  "data": {
    "total_sales": 50000.00,
    "total_orders": 1000,
    "total_customers": 500,
    "avg_order_amount": 50.00,
    "sales_growth": 15.5,
    "orders_growth": 10.2,
    "customers_growth": 8.3,
    "visit_count": 1800,
    "visit_users": 900,
    "pay_success_users": 320
  }
}
```

#### 3.8.2 销售趋势

```
GET /api/v1/merchant/analytics/sales-trend
Authorization: Bearer {token}
```

**请求参数：**

| 参数          | 类型     | 必填 | 描述                |
| ----------- | ------ | -- | ----------------- |
| start\_date | string | 是  | 开始日期              |
| end\_date   | string | 是  | 结束日期              |
| granularity | string | 否  | 粒度：day/week/month |

**响应：**

```json
{
  "code": 0,
  "data": [
    {
      "date": "2024-01-01",
      "sales": 5000.00,
      "orders": 100,
      "customers": 80,
      "visit_users": 80,
      "submit_order_users": 45
    }
  ]
}
```

#### 3.8.3 商品销量排行

```
GET /api/v1/merchant/analytics/product-ranking
Authorization: Bearer {token}
```

**请求参数：**

| 参数          | 类型     | 必填 | 描述                |
| ----------- | ------ | -- | ----------------- |
| start\_date | string | 否  | 开始日期              |
| end\_date   | string | 否  | 结束日期              |
| limit       | int    | 否  | 返回数量，默认10         |
| sort\_by    | string | 否  | 排序字段：sales/amount |

**响应：**

```json
{
  "code": 0,
  "data": [
    {
      "rank": 1,
      "product_id": 1,
      "product_name": "招牌红烧肉",
      "image": "图片URL",
      "sales_count": 256,
      "sales_amount": 14848.00
    }
  ]
}
```

#### 3.8.4 时段分析

```
GET /api/v1/merchant/analytics/hourly
Authorization: Bearer {token}
```

**请求参数：**

| 参数   | 类型     | 必填 | 描述      |
| ---- | ------ | -- | ------- |
| date | string | 否  | 日期，默认今天 |

**响应：**

```json
{
  "code": 0,
  "data": [
    {
      "hour": 9,
      "orders": 15,
      "sales": 750.00
    },
    {
      "hour": 10,
      "orders": 25,
      "sales": 1250.00
    }
  ]
}
```

#### 3.8.5 库存预警

```
GET /api/v1/merchant/analytics/stock-alert
Authorization: Bearer {token}
```

**请求参数：**

| 参数        | 类型  | 必填 | 描述        |
| --------- | --- | -- | --------- |
| threshold | int | 否  | 库存阈值，默认10 |

**响应：**

```json
{
  "code": 0,
  "data": [
    {
      "product_id": 1,
      "product_name": "招牌红烧肉",
      "image": "可直接显示的图片地址",
      "images": ["可直接显示的图片地址1", "可直接显示的图片地址2"],
      "stock": 5,
      "status": "low_stock"
    }
  ]
}
```

- 库存预警接口返回的 `image` 与 `images` 均为可直接展示的图片地址。

#### 3.8.6 商家用户分析（商家维度）

```
GET /api/v1/merchant/analytics/customers
Authorization: Bearer {token}
```

> **说明**：此接口从商家维度统计顾客与访客行为，当前版本以最小闭环为主，返回可直接用于商家数据看板页面的核心指标。

**请求参数：**

| 参数          | 类型     | 必填 | 描述   |
| ----------- | ------ | -- | ---- |
| start\_date | string | 否  | 开始日期 |
| end\_date   | string | 否  | 结束日期 |

**响应：**

```json
{
  "code": 0,
  "data": {
    "total_customers": 500,
    "new_customers": 50,
    "repeat_rate": 38.0,
    "visit_users": 900,
    "visit_count": 1800,
    "submit_order_users": 420,
    "pay_success_users": 320
  }
}
```

#### 3.8.7 用户消费趋势

```
GET /api/v1/merchant/analytics/customer-trend
Authorization: Bearer {token}
```

**请求参数：**

| 参数          | 类型     | 必填 | 描述                |
| ----------- | ------ | -- | ----------------- |
| start\_date | string | 是  | 开始日期              |
| end\_date   | string | 是  | 结束日期              |
| granularity | string | 否  | 粒度：day/week/month |

**响应：**

```json
{
  "code": 0,
  "data": [
    {
      "date": "2024-01-01",
      "new_customers": 10,
      "returning_customers": 40,
      "total_customers": 50
    }
  ]
}
```

### 3.9 C端店铺接口（扫码进入商家店铺，前缀 /api/v1/store）

> **说明**：用户通过扫描商家二维码进入统一小程序，小程序携带商家ID参数，所有操作都在当前商家店铺上下文中进行。

#### 3.9.1 微信登录

```
POST /api/v1/user/auth/wechat-login
```

**请求参数：**

```json
{
  "code": "微信登录凭证",
  "nickname": "用户昵称",
  "avatar": "头像URL"
}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "token": "JWT Token",
    "user": {
      "id": 1,
      "nickname": "用户昵称",
      "avatar": "头像URL"
    }
  }
}
```

#### 3.9.2 店铺首页（扫码进入）

```
GET /api/v1/store/{merchant_id}/home
```

> 用户扫描商家二维码后，小程序调用此接口获取店铺首页数据

- `pages/store/home`、`pages/store/product`、`pages/store/confirm` 统一支持从 `merchant_id` 或 `scene` 解析商家入口参数。
- 店铺首页公开数据请求进入页面后立即发起，不依赖先完成登录。

**响应：**

```json
{
  "code": 0,
  "data": {
    "merchant": {
      "id": 1,
      "name": "美味餐厅",
      "logo": "店铺Logo",
      "images": ["店铺图片"],
      "address": "店铺地址",
      "phone": "店铺电话",
      "business_hours": "09:00-22:00",
      "announcement": "今日特惠：全场8折",
      "status": "open",
      "rating": 4.8,
      "sales_count": 1000
    },
    "categories": [
      {
        "id": 1,
        "name": "热销推荐",
        "sort": 1,
        "product_count": 10
      }
    ],
    "hot_products": [
      {
        "id": 1,
        "name": "招牌红烧肉",
        "images": ["可直接显示的图片地址"],
        "price": 58.00,
        "original_price": 68.00,
        "sales": 256,
        "stock": 100
      }
    ],
    "delivery_settings": {
      "enabled": true,
      "base_fee": 5.00,
      "free_delivery_amount": 50.00,
      "max_distance": 10
    }
  }
}
```

**说明：**

- `hot_products` 当前返回商品结构化对象，图片字段为 `images` 数组。
- 商品图片地址为当前可直接展示的访问地址。

#### 3.9.3 商家商品列表

```
GET /api/v1/store/{merchant_id}/products
```

**请求参数：**

| 参数           | 类型     | 必填 | 描述               |
| ------------ | ------ | -- | ---------------- |
| category\_id | int    | 否  | 分类ID，不传则返回全部分类商品 |
| keyword      | string | 否  | 商品名称关键词          |
| page         | int    | 否  | 页码，默认 1          |
| page\_size   | int    | 否  | 每页数量，默认 20，最大 50 |

**响应：**

```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "name": "招牌红烧肉",
        "images": ["可直接显示的图片地址"],
        "price": 58.00,
        "original_price": 68.00,
        "sales": 256,
        "stock": 100,
        "specs": [
          {
            "id": 1,
            "name": "规格",
            "options": [
              {"name": "小份", "price": 48.00, "stock": 50},
              {"name": "大份", "price": 68.00, "stock": 50}
            ]
          }
        ]
      }
    ],
    "merchant": {
      "min_order_amount": 20.00,
      "takeout_enabled": true,
      "dine_in_enabled": true,
      "pickup_enabled": true
    },
    "pagination": {
      "total": 1,
      "page": 1,
      "page_size": 20
    }
  }
}
```

**说明：**

- 当前 C 端商品列表接口返回扁平商品列表与分页信息，不再按分类分组返回。
- 图片字段统一为 `images` 数组，前端直接取第一张图作为缩略图即可。
- 商品图片返回值已可直接用于页面展示。

#### 3.9.4 商品详情

```
GET /api/v1/store/{merchant_id}/products/{product_id}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "招牌红烧肉",
    "images": ["可直接显示的图片地址1", "可直接显示的图片地址2"],
    "description": "商品描述",
    "price": 58.00,
    "original_price": 68.00,
    "sales": 256,
    "stock": 100,
    "unit": "份",
    "specs": [
      {
        "id": 1,
        "name": "规格",
        "options": [
          {"id": 1, "name": "小份", "price": 48.00, "stock": 50},
          {"id": 2, "name": "大份", "price": 68.00, "stock": 50}
        ]
      }
    ]
  }
}
```

**说明：**

- 商品详情图片为可直接展示的访问地址。

#### 3.9.5 获取配送费规则

```
GET /api/v1/store/{merchant_id}/delivery-rules
```

> **说明**：获取商家的配送费规则，前端基于商家返回的配送档位供用户手动选择，实际配送费仍以后端创建订单时的校验和计算结果为准。

- `pages/store/confirm` 进入页面后先加载配送规则，再执行登录与埋点流程。

**响应：**

```json
{
  "code": 0,
  "data": {
    "enabled": true,
    "base_fee": 5.00,
    "free_delivery_amount": 50.00,
    "max_distance": 10,
    "takeout_enabled": true,
    "dine_in_enabled": true,
    "pickup_enabled": true,
    "distance_rules": [
      {"min_distance": 0, "max_distance": 2, "fee": 0},
      {"min_distance": 2, "max_distance": 5, "fee": 3.00},
      {"min_distance": 5, "max_distance": 10, "fee": 6.00}
    ]
  }
}
```

**说明**：

- 用户选择的是商家配置的配送距离档位，不是实时定位距离。
- 前端可根据 `distance_rules` 展示预计配送费区间
- 如果订单金额 >= `free_delivery_amount`，预计配送费可展示为 0
- `takeout_enabled`、`dine_in_enabled`、`pickup_enabled` 用于确认页动态展示配送、堂食、自提三种下单方式，未开启的方式不显示。
- `enabled` 仍表示配送费规则是否开启；即使存在配送规则，确认页是否展示“配送”也以后端返回的 `takeout_enabled` 为准。
- 超出商家支持范围时，前端仅提示不可提交，不接入第三方地图。
- 创建订单时，服务端会根据 `delivery_distance`、`distance_rules` 和满免门槛重新计算并校验最终配送费

#### 3.9.6 创建订单

```
POST /api/v1/store/{merchant_id}/orders
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "items": [
    {
      "product_id": 1,
      "spec_info": "大份",
      "quantity": 2
    }
  ],
  "delivery_type": 1,
  "delivery_distance": 3,
  "delivery_address": "XX市XX区XX路XX号",
  "contact_name": "张三",
  "contact_phone": "13800138000",
  "remark": "少放辣"
}
```

**delivery\_type 说明：**

| 值 | 描述                         |
| - | -------------------------- |
| 1 | 配送（需要填写配送距离、配送地址、联系人、联系电话） |
| 2 | 堂食（不需要配送费）                 |
| 3 | 自提（不需要配送费）                 |

**响应：**

```json
{
  "code": 0,
  "data": {
    "order": {
      "id": 1,
      "order_no": "202401010001",
      "total_amount": 116.00,
      "delivery_fee": 3.00,
      "pay_amount": 119.00,
      "delivery_type": 1,
      "delivery_distance": 3,
      "delivery_address": "收货地址",
      "status": 1,
      "remark": "少放辣"
    }
  }
}
```

> **说明**：当前下单接口不直接返回 `pay_params`，而是返回 `payment.next_action` 与 `payment.prepare_url`，再由 `xcx` 原生支付页单独请求支付参数。

#### 3.9.7 我的订单列表

```
GET /api/v1/user/orders
Authorization: Bearer {token}
```

**请求参数：**

| 参数           | 类型  | 必填 | 描述                                     |
| ------------ | --- | -- | -------------------------------------- |
| page         | int | 否  | 页码，默认1                                 |
| page\_size   | int | 否  | 每页数量，默认10                              |
| merchant\_id | int | 否  | 商家ID，用于筛选特定商家的订单（配合「我的订单」按钮使用）         |
| status       | int | 否  | 订单状态：0全部 1待支付 2已支付 3已完成 4已取消 5退款中 6已退款 |

**使用场景**：

- **全部订单**：不传 merchant\_id 参数，获取用户所有订单
- **商家订单**：传入 merchant\_id 参数，获取用户在特定商家的订单（用于店铺首页「我的订单」按钮跳转）

**响应：**

```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "order_no": "202401010001",
        "merchant": {
          "id": 1,
          "name": "美味餐厅",
          "logo": "店铺Logo"
        },
        "items": [
          {
            "product_name": "招牌红烧肉",
            "image": "图片URL",
            "quantity": 2
          }
        ],
        "total_amount": 116.00,
        "pay_amount": 120.00,
        "delivery_fee": 4.00,
        "status": 2,
        "status_text": "已支付",
        "verify_code": "123456",
        "created_at": "2024-01-01T12:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10
  }
}
```

#### 3.9.8 订单详情

```
GET /api/v1/user/orders/{order_id}
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "order_no": "202401010001",
    "merchant": {
      "id": 1,
      "name": "美味餐厅",
      "logo": "店铺Logo",
      "address": "店铺地址",
      "phone": "店铺电话"
    },
    "items": [
      {
        "product_id": 1,
        "product_name": "招牌红烧肉",
        "image": "图片URL",
        "price": 58.00,
        "quantity": 2,
        "spec_info": "大份",
        "subtotal": 116.00
      }
    ],
    "total_amount": 116.00,
    "delivery_fee": 3.00,
    "pay_amount": 119.00,
    "discount_amount": 0,
    "delivery_info": {
      "type": "delivery",
      "address": "收货地址",
      "contact_name": "联系人",
      "contact_phone": "联系电话",
      "distance": 3.5
    },
    "status": "paid",
    "status_text": "已支付",
    "remark": "少放辣",
    "verify_code": "123456",
    "created_at": "2024-01-01T12:00:00Z",
    "paid_at": "2024-01-01T12:01:00Z"
  }
}
```

#### 3.9.9 取消订单

```
POST /api/v1/user/orders/{order_id}/cancel
Authorization: Bearer {token}
```

**响应：**

```json
{
  "code": 0,
  "message": "订单已取消"
}
```

#### 3.9.10 申请退款

```
POST /api/v1/user/orders/{order_id}/refund
Authorization: Bearer {token}
```

**请求参数：**

```json
{
  "refund_reason": "退款原因"
}
```

**响应：**

```json
{
  "code": 0,
  "data": {
    "refund_id": 1,
    "refund_no": "RF202401010001",
    "status": "processing"
  }
}
```

### 3.10 历史支付方案归档

- 旧版“微信支付服务商模式”接口说明已归档，不再作为当前实现基线。
- 当前版本中：
  - `POST /api/v1/store/{merchant_id}/orders` 仅负责创建订单；
  - 不再返回 `pay_params`；
  - 不再定义微信支付回调、子商户回填、自动分账等接口为当前有效协议。
- 当前有效接口口径请以 `docs/prd/PRD-接口文档.md` 为准。

## 4. 数据模型设计

### 4.1 服务商表 (service\_providers)

| 字段               | 类型           | 描述            |
| ---------------- | ------------ | ------------- |
| id               | BIGINT       | 主键            |
| name             | VARCHAR(128) | 服务商名称         |
| contact\_name    | VARCHAR(64)  | 联系人           |
| contact\_phone   | VARCHAR(20)  | 联系电话          |
| mch\_id          | VARCHAR(32)  | 服务商商户号        |
| api\_key         | VARCHAR(128) | API密钥（加密存储）   |
| api\_v3\_key     | VARCHAR(128) | APIv3密钥（加密存储） |
| cert\_serial\_no | VARCHAR(64)  | 证书序列号         |
| private\_key     | TEXT         | 商户私钥（加密存储）    |
| public\_key      | TEXT         | 平台公钥          |
| callback\_url    | VARCHAR(256) | 回调地址          |
| status           | TINYINT      | 状态：0禁用 1正常    |
| created\_at      | DATETIME     | 创建时间          |
| updated\_at      | DATETIME     | 更新时间          |

### 4.2 服务商表 (service\_provider\_sps)

| 字段                    | 类型           | 描述                |
| --------------------- | ------------ | ----------------- |
| id                    | BIGINT       | 主键                |
| service\_provider\_id | BIGINT       | 服务商ID             |
| username              | VARCHAR(64)  | 用户名               |
| password              | VARCHAR(128) | 密码（加密存储）          |
| name                  | VARCHAR(64)  | 姓名                |
| phone                 | VARCHAR(20)  | 手机号               |
| role                  | VARCHAR(32)  | 角色：sp/operator |
| status                | TINYINT      | 状态                |
| last\_login\_at       | DATETIME     | 最后登录时间            |
| created\_at           | DATETIME     | 创建时间              |
| updated\_at           | DATETIME     | 更新时间              |

### 4.3 商家支付配置字段补充说明

本节原用于描述服务商支付模式下的门店支付配置字段。

当前口径：

- `sub_mch_id`
- `payment_config_status`
- `profit_sharing_*`
- `merchant_profit_sharing_records`

以上字段与表均已从当前实现基线中移除，仅保留历史归档意义。

### 4.4 商家表 (merchants)

| 字段                       | 类型            | 描述               |
| ------------------------ | ------------- | ---------------- |
| id                       | BIGINT        | 主键               |
| service\_provider\_id    | BIGINT        | 服务商ID            |
| name                     | VARCHAR(128)  | 商家名称             |
| logo                     | VARCHAR(512)  | 店铺Logo           |
| contact\_name            | VARCHAR(64)   | 联系人              |
| contact\_phone           | VARCHAR(20)   | 联系电话             |
| contact\_email           | VARCHAR(128)  | 联系邮箱             |
| address                  | VARCHAR(256)  | 店铺地址             |
| lat                      | DECIMAL(10,6) | 纬度               |
| lng                      | DECIMAL(10,6) | 经度               |
| business\_category       | VARCHAR(64)   | 经营类目             |
| business\_hours          | VARCHAR(64)   | 营业时间             |
| announcement             | TEXT          | 门店公告             |
| min\_order\_amount       | DECIMAL(10,2) | 最低起送金额           |
| takeout\_enabled         | BOOLEAN       | 是否支持外卖           |
| dine\_in\_enabled        | BOOLEAN       | 是否支持堂食           |
| pickup\_enabled          | BOOLEAN       | 是否支持自提           |
| status                   | TINYINT       | 状态：0关闭 1营业中      |
| rating                   | DECIMAL(2,1)  | 评分               |
| sales\_count             | INT           | 销量               |
| created\_at              | DATETIME      | 创建时间             |
| updated\_at              | DATETIME      | 更新时间             |

### 4.5 商家配送设置表 (merchant\_delivery\_settings)

| 字段                     | 类型            | 描述                    |
| ---------------------- | ------------- | --------------------- |
| id                     | BIGINT        | 主键                    |
| merchant\_id           | BIGINT        | 商家ID                  |
| enabled                | BOOLEAN       | 是否开启配送                |
| base\_fee              | DECIMAL(10,2) | 基础配送费                 |
| free\_delivery\_amount | DECIMAL(10,2) | 满额免配送费金额              |
| max\_distance          | INT           | 最大配送距离（公里，作为用户可选档位上限） |
| distance\_rules        | JSON          | 按距离收费规则               |
| created\_at            | DATETIME      | 创建时间                  |
| updated\_at            | DATETIME      | 更新时间                  |

**distance\_rules JSON 结构：**

```json
[
  {"min_distance": 0, "max_distance": 2, "fee": 0},
  {"min_distance": 2, "max_distance": 5, "fee": 3.00},
  {"min_distance": 5, "max_distance": 10, "fee": 6.00}
]
```

- 该表只描述商家允许用户手动选择的配送档位，不参与真实地图定位。
- `distance_rules` 需满足区间不重叠、结束距离大于起始距离、且不超过 `max_distance`。

### 4.6 历史支付配置状态归档

- `payment_config_status` 与分账开关说明已失效，不再作为当前数据库与接口设计的一部分。
- 若旧库仍存在历史字段，以清理迁移脚本和当前代码模型为准。

### 4.7 商家员工表 (merchant\_staffs)

| 字段                      | 类型           | 描述                     |
| ----------------------- | ------------ | ---------------------- |
| id                      | BIGINT       | 主键                     |
| merchant\_id            | BIGINT       | 商家ID                   |
| username                | VARCHAR(64)  | 用户名                    |
| password                | VARCHAR(128) | 密码                     |
| name                    | VARCHAR(64)  | 姓名                     |
| phone                   | VARCHAR(20)  | 手机号                    |
| push\_openid            | VARCHAR(64)  | 员工消息推送 OpenID           |
| role                    | VARCHAR(32)  | 角色：owner/manager/staff |
| notify\_enabled         | TINYINT(1)   | 订单提示音开关，1开启 0关闭        |
| browse\_notify\_enabled | TINYINT(1)   | 浏览提示音开关，1开启 0关闭        |
| status                  | TINYINT      | 状态                     |
| last\_login\_at         | DATETIME     | 最后登录时间                 |
| created\_at             | DATETIME     | 创建时间                   |
| updated\_at             | DATETIME     | 更新时间                   |

- 员工管理首版由 `owner` 角色负责，不引入细粒度权限模型。

### 4.8 商品分类表 (categories)

| 字段           | 类型          | 描述         |
| ------------ | ----------- | ---------- |
| id           | BIGINT      | 主键         |
| merchant\_id | BIGINT      | 商家ID       |
| name         | VARCHAR(64) | 分类名称       |
| sort         | INT         | 排序权重       |
| status       | TINYINT     | 状态：0禁用 1启用 |
| created\_at  | DATETIME    | 创建时间       |
| updated\_at  | DATETIME    | 更新时间       |

### 4.9 商品表 (products)

| 字段              | 类型            | 描述         |
| --------------- | ------------- | ---------- |
| id              | BIGINT        | 主键         |
| merchant\_id    | BIGINT        | 商家ID       |
| category\_id    | BIGINT        | 分类ID       |
| name            | VARCHAR(128)  | 商品名称       |
| description     | TEXT          | 商品描述       |
| images          | JSON          | 图片数组       |
| price           | DECIMAL(10,2) | 售价         |
| original\_price | DECIMAL(10,2) | 原价         |
| stock           | INT           | 库存         |
| unit            | VARCHAR(16)   | 单位         |
| sales           | INT           | 销量         |
| sort            | INT           | 排序权重       |
| status          | TINYINT       | 状态：0下架 1上架 |
| deleted\_at     | DATETIME      | 删除时间       |
| created\_at     | DATETIME      | 创建时间       |
| updated\_at     | DATETIME      | 更新时间       |

### 4.10 商品规格表 (product\_specs)

| 字段          | 类型          | 描述                                                                                           |
| ----------- | ----------- | -------------------------------------------------------------------------------------------- |
| id          | BIGINT      | 主键                                                                                           |
| product\_id | BIGINT      | 商品ID                                                                                         |
| name        | VARCHAR(64) | 规格名称（如：份量）                                                                                   |
| options     | JSON        | 规格选项（存储规格和对应加价，最终售价=商品基础价+所选规格加价；如：\[{"name":"小份","price":0.00},{"name":"大份","price":3.00}]） |
| created\_at | DATETIME    | 创建时间                                                                                         |
| updated\_at | DATETIME    | 更新时间                                                                                         |

### 4.11 C端用户表 (users)

| 字段               | 类型            | 描述        |
| ---------------- | ------------- | --------- |
| id               | BIGINT        | 主键        |
| openid           | VARCHAR(64)   | 微信OpenID  |
| union\_id        | VARCHAR(64)   | 微信UnionID |
| nickname         | VARCHAR(64)   | 昵称        |
| avatar           | VARCHAR(512)  | 头像        |
| phone            | VARCHAR(20)   | 手机号       |
| status           | TINYINT       | 状态        |
| first\_visit\_at | DATETIME      | 首次访问时间    |
| last\_visit\_at  | DATETIME      | 最近访问时间    |
| visit\_count     | INT           | 访问次数      |
| has\_ordered     | TINYINT(1)    | 是否下过单     |
| total\_orders    | INT           | 累计订单数     |
| total\_spent     | DECIMAL(10,2) | 累计消费金额    |
| has\_paid        | TINYINT(1)    | 是否发生过支付   |
| first\_paid\_at  | DATETIME      | 首次支付时间    |
| created\_at      | DATETIME      | 创建时间      |
| updated\_at      | DATETIME      | 更新时间      |

### 4.12 订单表 (orders)

| 字段                         | 类型            | 描述                    |
| -------------------------- | ------------- | --------------------- |
| id                         | BIGINT        | 主键                    |
| order\_no                  | VARCHAR(32)   | 订单号                   |
| user\_id                   | BIGINT        | 用户ID                  |
| merchant\_id               | BIGINT        | 商家ID                  |
| total\_amount              | DECIMAL(10,2) | 商品总金额                 |
| delivery\_fee              | DECIMAL(10,2) | 配送费                   |
| discount\_amount           | DECIMAL(10,2) | 优惠金额                  |
| pay\_amount                | DECIMAL(10,2) | 实付金额                  |
| delivery\_type             | TINYINT       | 配送类型：1配送 2堂食 3自提      |
| delivery\_distance         | DECIMAL(5,2)  | 配送距离（公里）              |
| delivery\_address          | VARCHAR(256)  | 配送地址                  |
| contact\_name              | VARCHAR(64)   | 联系人姓名                 |
| contact\_phone             | VARCHAR(20)   | 联系人电话                 |
| status                     | TINYINT       | 状态                    |
| remark                     | VARCHAR(256)  | 备注                    |
| verify\_code               | VARCHAR(16)   | 核销码                   |
| transaction\_id            | VARCHAR(64)   | 支付交易号（当前未启用在线支付时为空） |
| paid\_at                   | DATETIME      | 支付时间                  |
| completed\_at              | DATETIME      | 完成时间                  |
| completed\_by\_name        | VARCHAR(64)   | 核销人                   |
| cancelled\_at              | DATETIME      | 取消时间                  |
| refunded\_at               | DATETIME      | 退款时间                  |
| created\_at                | DATETIME      | 创建时间                  |
| updated\_at                | DATETIME      | 更新时间                  |

**订单状态流转：**

```
pending_payment(待支付) → paid(已支付) → completed(已完成)
        ↓                    ↓
   cancelled(已取消)    refunding(退款中) → refunded(已退款)
```

### 4.13 订单商品表 (order\_items)

| 字段            | 类型            | 描述            |
| ------------- | ------------- | ------------- |
| id            | BIGINT        | 主键            |
| order\_id     | BIGINT        | 订单ID          |
| merchant\_id  | BIGINT        | 商家ID          |
| product\_id   | BIGINT        | 商品ID          |
| product\_name | VARCHAR(128)  | 商品名称          |
| image         | VARCHAR(512)  | 商品图片          |
| price         | DECIMAL(10,2) | 单价            |
| quantity      | INT           | 数量            |
| spec\_info    | JSON          | 规格信息（如：大份/小份） |
| subtotal      | DECIMAL(10,2) | 小计            |
| created\_at   | DATETIME      | 创建时间          |

### 4.14 退款记录表 (refunds)

| 字段             | 类型            | 描述              |
| -------------- | ------------- | --------------- |
| id             | BIGINT        | 主键              |
| order\_id      | BIGINT        | 订单ID            |
| refund\_no     | VARCHAR(32)   | 退款单号            |
| refund\_amount | DECIMAL(10,2) | 退款金额            |
| refund\_reason | VARCHAR(256)  | 退款原因            |
| status         | TINYINT       | 状态：0处理中 1成功 2失败 |
| refund\_id     | VARCHAR(64)   | 微信退款单号          |
| refunded\_at   | DATETIME      | 退款完成时间          |
| created\_at    | DATETIME      | 创建时间            |
| updated\_at    | DATETIME      | 更新时间            |

### 4.15 系统公告表 (announcements)

| 字段                    | 类型           | 描述         |
| --------------------- | ------------ | ---------- |
| id                    | BIGINT       | 主键         |
| service\_provider\_id | BIGINT       | 服务商ID      |
| title                 | VARCHAR(128) | 公告标题       |
| content               | TEXT         | 公告内容       |
| status                | TINYINT      | 状态：0禁用 1启用 |
| created\_at           | DATETIME     | 创建时间       |
| updated\_at           | DATETIME     | 更新时间       |

### 4.16 云打印机表 (cloud\_printers)

| 字段           | 类型           | 描述                      |
| ------------ | ------------ | ----------------------- |
| id           | BIGINT       | 主键                      |
| merchant\_id | BIGINT       | 商家ID                    |
| name         | VARCHAR(64)  | 打印机名称                   |
| type         | VARCHAR(32)  | 打印机类型：yilianyun/feie/xp |
| device\_no   | VARCHAR(64)  | 设备编号                    |
| api\_key     | VARCHAR(128) | API密钥                   |
| api\_url     | VARCHAR(256) | API接口地址                 |
| status       | TINYINT      | 状态：0离线 1在线              |
| auto\_print  | BOOLEAN      | 自动打印开关                  |
| print\_count | INT          | 累计打印次数                  |
| is\_default  | BOOLEAN      | 是否默认打印机                 |
| created\_at  | DATETIME     | 创建时间                    |
| updated\_at  | DATETIME     | 更新时间                    |

### 4.17 打印记录表 (print\_logs)

| 字段             | 类型           | 描述         |
| -------------- | ------------ | ---------- |
| id             | BIGINT       | 主键         |
| merchant\_id   | BIGINT       | 商家ID       |
| printer\_id    | BIGINT       | 打印机ID      |
| order\_id      | BIGINT       | 订单ID       |
| status         | TINYINT      | 状态：0失败 1成功 |
| error\_message | VARCHAR(256) | 错误信息       |
| print\_time    | DATETIME     | 打印时间       |
| created\_at    | DATETIME     | 创建时间       |

### 4.18 历史分账记录表归档

- `merchant_profit_sharing_records` 为旧服务商支付模式下的历史表。
- 当前版本已从主链路、代码模型和数据库基线中移除该表，不再维护字段明细。

### 4.19 商家年费表 (merchant\_fees)

| 字段           | 类型            | 描述                  |
| ------------ | ------------- | ------------------- |
| id           | BIGINT        | 主键                  |
| merchant\_id | BIGINT        | 商家ID                |
| year         | YEAR          | 年费年度                |
| amount       | DECIMAL(10,2) | 年费金额                |
| status       | VARCHAR(32)   | 状态：unpaid/paid/free |
| pay\_time    | DATETIME      | 支付时间                |
| free\_reason | VARCHAR(256)  | 免年费原因（如限时优惠）        |
| created\_at  | DATETIME      | 创建时间                |
| updated\_at  | DATETIME      | 更新时间                |

### 4.20 用户收货地址表 (user\_addresses)

| 字段          | 类型           | 描述         |
| ----------- | ------------ | ---------- |
| id          | BIGINT       | 主键         |
| user\_id    | BIGINT       | 用户ID       |
| name        | VARCHAR(64)  | 收货人姓名      |
| phone       | VARCHAR(20)  | 联系电话       |
| province    | VARCHAR(64)  | 省份         |
| city        | VARCHAR(64)  | 城市         |
| district    | VARCHAR(64)  | 区县         |
| detail      | VARCHAR(512) | 详细地址       |
| is\_default | TINYINT      | 是否默认：0否 1是 |
| status      | TINYINT      | 状态：0禁用 1正常 |
| created\_at | DATETIME     | 创建时间       |
| updated\_at | DATETIME     | 更新时间       |

### 4.21 用户行为事件表 (user\_behavior\_events)

| 字段           | 类型          | 描述                  |
| ------------ | ----------- | ------------------- |
| id           | BIGINT      | 主键                  |
| merchant\_id | BIGINT      | 商家ID                |
| user\_id     | BIGINT      | 用户ID                |
| openid       | VARCHAR(64) | 微信OpenID            |
| event\_type  | VARCHAR(32) | 事件类型                |
| page         | VARCHAR(64) | 页面标识                |
| product\_id  | BIGINT      | 商品ID，可空             |
| order\_id    | BIGINT      | 订单ID，可空             |
| source       | VARCHAR(32) | 来源，如 scan/store/dev |
| payload      | JSON        | 扩展数据                |
| created\_at  | DATETIME    | 创建时间                |

**当前事件类型：**

- `store_visit`
- `page_view`
- `product_view`
- `submit_order`
- `pay_success`

### 4.22 优惠券/活动表 (coupons)

| 字段               | 类型            | 描述                      |
| ---------------- | ------------- | ----------------------- |
| id               | BIGINT        | 主键                      |
| merchant\_id     | BIGINT        | 商家ID（0表示平台活动）           |
| name             | VARCHAR(128)  | 活动名称                    |
| type             | VARCHAR(32)   | 类型：discount/coupon/gift |
| discount         | DECIMAL(5,2)  | 折扣比例（如0.9表示9折）          |
| min\_amount      | DECIMAL(10,2) | 最低消费金额                  |
| discount\_amount | DECIMAL(10,2) | 优惠金额                    |
| total\_count     | INT           | 发放总量                    |
| used\_count      | INT           | 已使用数量                   |
| start\_time      | DATETIME      | 开始时间                    |
| end\_time        | DATETIME      | 结束时间                    |
| status           | TINYINT       | 状态：0禁用 1启用              |
| created\_at      | DATETIME      | 创建时间                    |
| updated\_at      | DATETIME      | 更新时间                    |

> **说明**：优惠券/活动能力不是本轮实施范围，当前仅保留数据结构预留说明。

### 4.23 优惠券领取记录表 (coupon\_records)

| 字段            | 类型          | 描述                     |
| ------------- | ----------- | ---------------------- |
| id            | BIGINT      | 主键                     |
| coupon\_id    | BIGINT      | 优惠券ID                  |
| user\_id      | BIGINT      | 用户ID                   |
| order\_id     | BIGINT      | 使用订单ID（NULL表示未使用）      |
| status        | VARCHAR(32) | 状态：unused/used/expired |
| receive\_time | DATETIME    | 领取时间                   |
| used\_time    | DATETIME    | 使用时间                   |
| created\_at   | DATETIME    | 创建时间                   |

## 5. 技术架构

### 5.1 技术栈

| 层级    | 技术选型         | 说明              |
| ----- | ------------ | --------------- |
| Web框架 | Gin          | Go语言高性能Web框架    |
| ORM   | GORM         | Go语言ORM库        |
| 数据库   | MySQL 8.0    | 主数据库            |
| 缓存    | Redis        | 缓存、Session、分布式锁 |
| 认证    | JWT          | 用户认证            |
| 文档    | Swagger      | API文档自动生成       |
| 文件存储 | Local Storage | 本地文件存储         |

### 5.2 项目结构

#### 后端项目 (api-chaoshi/)

> **说明**：当前后端顶层目录为 `api-chaoshi/`；下方 `cmd/server/` 仅表示启动入口子目录，不代表旧的顶层 `server/` 项目路径。

```
api-chaoshi/                     # Go后端服务
├── cmd/
│   └── server/
│       └── main.go            # API 启动入口
├── internal/
│   ├── config/                # 配置管理
│   ├── handlers/               # 处理器层
│   │   ├── sp/                # 服务商管理
│   │   ├── merchant/           # 商家管理
│   │   └── user/              # 用户功能
│   ├── middleware/            # 中间件
│   ├── models/                # 数据模型
│   └── utils/                 # 工具函数
├── pkg/
│   ├── database/              # 数据库连接
│   └── response/              # 统一响应
├── migrations/                # 单文件初始化基线目录，仅保留 20240101000000_full_init.sql
├── .env                       # 环境配置
├── .env.example
├── docker-compose.yml
└── go.mod
```

#### 小程序项目 (miniprogram/)

```
miniprogram/                    # 微信小程序
├── src/
│   ├── App.vue
│   ├── main.ts
│   ├── pages.json
│   ├── manifest.json
│   └── ...
├── package.json
└── ...
```

#### 小程序验证壳 (xcx/)

```text
xcx/                            # 微信小程序验证壳（Taro）
├── src/
│   ├── pages/
│   │   ├── index/              # 启动页：无感登录后进入 H5
│   │   ├── settings/           # 配置页：维护 API/H5/merchant_id
│   │   ├── debug/              # 调试页：查看 token/openid/最终 H5 地址
│   │   └── webview/            # H5 容器页
│   ├── services/               # 登录、配置、请求与本地存储封装
│   ├── types/                  # 壳工程类型定义
│   └── utils/                  # H5 地址拼装等工具
└── README.md                   # 验证链路说明
```

### 5.3 分层架构

```
┌─────────────────────────────────────────────┐
│              Handlers (Controller)           │  ← 处理HTTP请求/响应
├─────────────────────────────────────────────┤
│              Services (Business)             │  ← 业务逻辑处理
├─────────────────────────────────────────────┤
│           Repositories (Data Access)         │  ← 数据访问层
├─────────────────────────────────────────────┤
│              Models (Entity)                 │  ← 数据模型定义
└─────────────────────────────────────────────┘
```

### 5.4 历史支付架构归档

- 原“微信支付服务商模式架构”已归档。
- 当前技术架构以本地文件上传、总部后台、商家端 H5、用户端 H5/小程序壳协同为主。

### 5.5 数据库索引设计

**核心表索引规划：**

| 表名                                 | 索引字段                                | 索引类型   | 用途            |
| ---------------------------------- | ----------------------------------- | ------ | ------------- |
| merchants                          | status, created\_at                 | BTREE  | 商家状态筛选、创建时间排序 |
| merchants                          | business\_category                  | BTREE  | 按行业分类统计       |
| orders                             | merchant\_id, created\_at           | BTREE  | 商家订单查询、时间范围统计 |
| orders                             | order\_no                           | UNIQUE | 订单号快速查询       |
| orders                             | status, created\_at                 | BTREE  | 订单状态筛选        |
| users                              | openid                              | UNIQUE | 微信用户快速查询      |
| products                           | merchant\_id, status                | BTREE  | 商家商品查询        |

### 5.6 Redis缓存策略

**缓存类型：**

| 缓存键                    | 数据类型   | 过期时间 | 用途          |
| ---------------------- | ------ | ---- | ----------- |
| token:{token}          | STRING | 2小时  | 用户登录Token验证 |
| refresh\_token:{token} | STRING | 7天   | Token刷新     |
| merchant:{id}          | HASH   | 5分钟  | 商家基本信息缓存    |
| merchant:{id}:products | LIST   | 5分钟  | 商家商品列表缓存    |
| dashboard:{sp\_id}     | STRING | 1分钟  | 服务商数据看板缓存   |
| rate\_limit:{ip}       | STRING | 1分钟  | 接口限流计数      |

**缓存更新策略：**

- **写操作后更新缓存**：商家信息、商品信息更新时同步更新缓存
- **定时刷新**：数据看板类数据设置较短过期时间，定期刷新
- **缓存击穿处理**：使用互斥锁避免缓存击穿

### 5.7 接口限流策略

**限流规则：**

| 接口类别   | 限制频率       | 说明      |
| ------ | ---------- | ------- |
| 登录接口   | 5次/分钟/IP   | 防止暴力破解  |
| 普通业务接口 | 100次/分钟/IP | 常规限流    |
| 数据导出接口 | 1次/分钟/用户   | 资源密集型接口 |
| 历史通知接口 | 1000次/分钟   | 旧支付通知接口（当前未启用） |

**限流实现：**

1. 使用Redis实现分布式限流
2. 采用令牌桶算法或滑动窗口算法
3. 返回429状态码表示限流触发
4. 记录限流日志便于监控

**限流响应：**

```json
{
  "code": 429,
  "message": "请求过于频繁，请稍后重试",
  "retry_after": 60
}
```

## 6. 开发计划

### 6.1 第一阶段：基础架构 (P0)

- [ ] 项目初始化与依赖配置
- [ ] 数据库连接与迁移脚本
- [ ] 基础错误处理与响应格式
- [ ] JWT认证中间件
- [ ] 日志系统

### 6.2 第二阶段：服务商与商家管理 (P0)

- [ ] 服务商登录（服务端文件配置服务商信息）
- [ ] 服务商直接创建商家
- [ ] 商家信息管理与门店基础配置
- [ ] 商家信息管理
- [ ] 商家自定义设置
- [ ] 系统公告管理（服务商发布，商家查看）

### 6.3 第三阶段：商品管理 (P0)

- [ ] 商品分类CRUD
- [ ] 商品CRUD
- [ ] 商品上下架
- [ ] 商品规格管理
- [ ] 库存管理

### 6.4 第四阶段：订单与支付预留 (P0)

- [ ] 订单创建
- [ ] 订单状态管理
- [ ] 订单核销
- [ ] 退款处理
- [ ] 新支付方案预留接入设计

### 6.5 第五阶段：数据分析 (P1)

- [ ] 销售统计
- [ ] 商品销量排行
- [ ] 时段分析
- [ ] 库存预警

### 6.6 第六阶段：优化迭代 (P2)

- [ ] 员工管理
- [ ] 缓存优化
- [ ] 性能调优
- [ ] 监控告警

### 6.7 小程序开发计划

#### 技术选型

- **开发框架**：uni-app（Vue3 + TypeScript）
- **UI组件**：uView / Vant Weapp
- **状态管理**：Pinia
- **样式规范**：简洁实用型，支持商家自定义模板
- **发布平台**：微信小程序（可扩展至H5/APP）

#### 开发原则

1. **简洁实用** - 第一版本专注于功能打通，不追求复杂UI
2. **商家模板预留** - 组件设计预留样式扩展接口，便于后续商家店铺模板
3. **接口优先** - 先确保API对接，再优化界面

#### 页面优先级（按批次）

**第一批次：认证与设置**

| 页面                           | 功能描述       | 对接API                                       |
| ---------------------------- | ---------- | ------------------------------------------- |
| auth/login                   | 商家登录（账号密码） | POST /api/v1/auth/merchant/login            |
| merchant/settings            | 商家设置与员工能力  | GET/PUT /api/v1/merchant/settings           |

**第二批次：商户管理中心首页**

| 页面                          | 功能描述              | 对接API                                   |
| --------------------------- | ----------------- | --------------------------------------- |
| merchant/home               | 商户后台首页、快捷入口、待处理事项 | GET /api/v1/merchant/profile            |
| merchant/orders/statistics  | 订单统计卡片            | GET /api/v1/merchant/orders/statistics  |
| merchant/analytics/overview | 数据概览              | GET /api/v1/merchant/analytics/overview |

**第三批次：商品管理**

| 页面                     | 功能描述        | 对接API                                |
| ---------------------- | ----------- | ------------------------------------ |
| merchant/products/list | 商品列表、搜索、上下架 | GET /api/v1/merchant/products        |
| merchant/products/edit | 商品编辑、规格设置   | POST/PUT /api/v1/merchant/products   |
| merchant/categories    | 分类管理、排序     | GET/POST /api/v1/merchant/categories |

**第四批次：订单管理**

| 页面                     | 功能描述      | 对接API                            |
| ---------------------- | --------- | -------------------------------- |
| merchant/orders/list   | 订单列表、状态筛选 | GET /api/v1/merchant/orders      |
| merchant/orders/detail | 订单详情、核销操作 | GET /api/v1/merchant/orders/{id} |

**第五批次：商家设置与系统公告**

| 页面                         | 功能描述      | 对接API                                      |
| -------------------------- | --------- | ------------------------------------------ |
| merchant/settings          | 商家设置首页、快捷入口 | GET /api/v1/merchant/settings               |
| merchant/marketing         | 满减营销配置    | GET/PUT /api/v1/merchant/full-reduction-rules |
| merchant/printers          | 打印机管理     | GET/POST/PUT/DELETE /api/v1/merchant/printers* |
| merchant/delivery-settings | 配送设置      | GET/PUT /api/v1/merchant/delivery-settings |
| merchant/announcements     | 系统公告列表    | GET /api/v1/merchant/announcements         |

**第六批次：数据分析**

| 页面                           | 功能描述              | 对接API                                          |
| ---------------------------- | ----------------- | ---------------------------------------------- |
| merchant/analytics/sales     | 经营趋势图表（浏览人数/下单人数） | GET /api/v1/merchant/analytics/sales-trend     |
| merchant/analytics/products  | 商品销量排行            | GET /api/v1/merchant/analytics/product-ranking |
| merchant/analytics/customers | 客户分析              | GET /api/v1/merchant/analytics/customers       |

**第七批次：C端店铺购物（扫码进入）**

| 页面             | 功能描述      | 对接API                                    |
| -------------- | --------- | ---------------------------------------- |
| store/home     | 店铺首页、分类商品 | GET /api/v1/store/{id}/home              |
| store/products | 商品列表      | GET /api/v1/store/{id}/products          |
| store/product  | 商品详情      | GET /api/v1/store/{id}/products/{id}     |
| store/confirm  | 确认订单、创建订单 | POST /api/v1/store/{merchant\_id}/orders |

**第八批次：订单创建与退款联调**

| 功能   | 描述                 | 对接API                                    |
| ---- | ------------------ | ---------------------------------------- |
| 创建订单 | 提交订单，返回支付下一步动作 | POST /api/v1/store/{merchant\_id}/orders |
| 江苏银行支付 | `xcx` 支付页拉起江苏银行微信小程序支付 | POST /api/v1/user/orders/{id}/pay/prepare |
| 支付回调 | 江苏银行回调更新订单状态   | POST /api/v1/payments/jsbank/notify |
| 退款处理 | 申请退款               | POST /api/v1/merchant/orders/{id}/refund |

## 7. 非功能性需求

### 7.1 性能要求

- API 响应时间 < 200ms (P95)
- 支持并发 1000 QPS
- 数据库查询优化，避免 N+1 问题
- 合理使用缓存，热点数据缓存命中率 > 90%

### 7.2 安全要求

- 所有接口使用 HTTPS
- 敏感信息加密存储（密码、支付密钥等）
- SQL 注入防护
- XSS 防护
- 接口限流保护
- 支付签名验证

### 7.3 可维护性

- 完善的日志记录
- API 文档自动生成
- 单元测试覆盖率 > 70%
- 错误码规范化
- 数据库迁移版本管理

## 8. 附录

### 8.1 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 8.2 错误码定义

| 错误码  | 描述             |
| ---- | -------------- |
| 0    | 成功             |
| 1001 | 参数错误           |
| 1002 | 未授权            |
| 1003 | 禁止访问           |
| 1004 | 资源不存在          |
| 2001 | 用户不存在          |
| 2002 | 用户已存在          |
| 2003 | 密码错误           |
| 3001 | 商家不存在          |
| 3002 | 商家基础配置未完成      |
| 3003 | 商家已禁用          |
| 3004 | 商家未开通配送        |
| 3005 | 超出配送范围         |
| 4001 | 商品不存在          |
| 4002 | 商品已下架          |
| 4003 | 库存不足           |
| 5001 | 订单不存在          |
| 5002 | 订单状态错误         |
| 5003 | 订单已支付          |
| 5004 | 订单已取消          |
| 5005 | 无权操作此订单（跨商家操作） |
| 5006 | 核销码错误          |
| 6001 | 创建订单失败         |
| 6002 | 退款失败           |
| 6003 | 订单状态处理失败       |
| 6004 | 历史分账能力未启用      |
| 7001 | 分类不存在          |
| 7002 | 分类下存在商品        |
| 8001 | 上传配置错误         |
| 8002 | 文件上传失败         |
| 9001 | 服务器内部错误        |

### 8.3 数据库枚举与状态定义

#### 8.3.1 订单状态枚举 (orders.status)

| 值 | 状态标识       | 描述  | 说明  |
| - | ---------- | --- | --- |
| 1 | pending\_payment | 待支付 | 订单已创建，未支付 |
| 2 | paid | 已支付 | 支付成功 |
| 3 | completed | 已完成 | 已核销/完成 |
| 4 | cancelled | 已取消 | 订单已取消 |
| 5 | refunding | 退款中 | 退款处理中 |
| 6 | refunded | 已退款 | 退款完成 |

#### 8.3.2 历史分账状态枚举归档

- `orders.profit_sharing_status` 与 `merchant_profit_sharing_records.status` 已移出当前实现基线。
- 本节枚举不再作为当前数据库设计依据。

#### 8.3.3 商家状态枚举 (merchants.status)

| 值 | 状态标识       | 描述  |
| - | ---------- | --- |
| 0 | disabled | 已禁用 |
| 1 | active | 营业中/正常 |

#### 8.3.4 历史支付配置状态枚举归档

- `merchants.payment_config_status` 已从当前实现基线移除。
- 本节仅保留为历史概念说明，不再参与现有接口与数据库校验。

#### 8.3.5 商品状态枚举 (products.status)

| 值 | 状态标识       | 描述  |
| - | ---------- | --- |
| 0 | offline | 已下架 |
| 1 | online | 已上架 |

#### 8.3.6 商品分类状态枚举 (categories.status)

| 值 | 状态标识       | 描述  |
| - | ---------- | --- |
| 0 | inactive | 停用 |
| 1 | active | 正常 |

#### 8.3.7 服务商状态枚举 (service\_providers.status)

| 值 | 状态标识       | 描述  |
| - | ---------- | --- |
| 0 | disabled | 已禁用 |
| 1 | active | 正常 |

#### 8.3.8 服务商状态枚举 (service\_provider\_sps.status)

| 值 | 状态标识       | 描述  |
| - | ---------- | --- |
| 0 | disabled | 已禁用 |
| 1 | active | 正常 |

#### 8.3.9 商家员工状态枚举 (merchant\_staffs.status)

| 值 | 状态标识       | 描述  |
| - | ---------- | --- |
| 0 | disabled | 已禁用 |
| 1 | active | 正常 |

#### 8.3.10 用户状态枚举 (users.status)

| 值 | 状态标识       | 描述  |
| - | ---------- | --- |
| 0 | disabled | 已禁用 |
| 1 | active | 正常 |

#### 8.3.11 系统公告状态枚举 (announcements.status)

| 值 | 状态标识       | 描述  |
| - | ---------- | --- |
| 0 | draft | 草稿 |
| 1 | published | 已发布 |

#### 8.3.12 云打印机状态枚举 (cloud\_printers.status)

| 值 | 状态标识       | 描述  |
| - | ---------- | --- |
| 0 | inactive | 停用 |
| 1 | active | 正常 |

#### 8.3.13 打印记录状态枚举 (print\_logs.status)

| 值 | 状态标识       | 描述  |
| - | ---------- | --- |
| 0 | pending | 待打印 |
| 1 | success | 打印成功 |
| 2 | failed | 打印失败 |

#### 8.3.14 退款记录状态枚举 (refunds.status)

| 值 | 状态标识       | 描述  |
| - | ---------- | --- |
| 0 | pending | 待处理 |
| 1 | processing | 处理中 |
| 2 | success | 退款成功 |
| 3 | failed | 退款失败 |

#### 8.3.15 配送类型枚举 (orders.delivery\_type)

| 值 | 类型标识       | 描述  |
| - | ---------- | --- |
| 1 | dine\_in | 堂食 |
| 2 | takeaway | 外卖 |
| 3 | self\_pickup | 自提 |

### 8.4 分页参数

```json
{
  "page": 1,
  "page_size": 10,
  "total": 100,
  "total_pages": 10
}
```

### 8.5 历史微信支付回调归档

- 微信支付回调报文不再属于当前版本有效协议。
- 若后续接入新的支付方案，应在新的接口文档中重新定义通知结构。

## 9. 小程序测试指南

### 9.1 测试环境说明

当前开发环境包含三类测试入口，分别对应不同角色与用途：

1. **总部后台**：通过 `web-admin/` 登录后进入管理端，用于创建门店、维护资料、查看统计与维护后台设置。
2. **商家端**：商家员工登录后进入经营后台，用于商品、订单、设置与员工能力联调。
3. **C 端店铺端**：用户浏览店铺并下单的链路测试，由于开发阶段无法总是通过扫码直达指定商家店铺，因此提供测试入口页与编译模式两种方式。

开发测试时，需明确区分以下入口：

- `web-admin/`：总部后台 Web/PC 项目入口
- `pages/auth/login`：商家登录入口
- `pages/auth/agreement`：商家服务协议入口页，整合原隐私政策内容
- `pages/store/test-entry`：C 端店铺测试入口，不是服务商入口
- `pages/store/home?merchant_id=...`：C 端店铺首页直达方式，不是服务商入口

### 9.1.1 双小程序发布约定

- 小程序前端支持通过一套源码切换发布两个微信小程序。
- 品牌配置文件：
  - `miniprogram/config/brands/xunmeng.json`
  - `miniprogram/config/brands/caixu.json`
- 品牌环境文件：
  - `miniprogram/.env.production.xunmeng`
  - `miniprogram/.env.production.caixu`
- 预构建脚本：`miniprogram/scripts/apply-brand-config.mjs`
- 发布命令：
  - `npm run build:mp-weixin:xunmeng`
  - `npm run build:mp-weixin:caixu`
- 执行品牌构建脚本时，会自动写入 `src/manifest.json` 的应用名称与 `appid`，并同步对应的 `.env.production`。
- 执行品牌构建脚本时，也会同步写入 `src/pages.json` 的全局导航标题。
- 页面运行时统一通过 `src/config/env.ts` 中的 `APP_NAME` 读取品牌名称，页面文案不再手工硬编码品牌名。
- `versionName` 与 `versionCode` 继续统一维护在 `miniprogram/src/manifest.json`。

### 9.2 测试账号总表（开发环境）

| 角色     | 账号         | 密码            | 说明                 |
| ------ | ---------- | ------------- | ------------------ |
| 后台管理员 | `admin` | `tm666666`    | 用于总部后台登录与功能验证 |
| 商家员工   | `merchant` | `merchant123` | 用于商家工作台、商品、订单与设置联调 |

### 9.3 总部后台测试

#### 入口与登录

- 项目目录：`web-admin/`
- 登录页：`/login`
- 登录后首页：`/dashboard`
- 建议通过浏览器访问本地开发地址进行测试

#### 使用步骤

1. 进入 `web-admin/` 目录并启动本地开发环境
2. 打开浏览器访问总部后台登录页
3. 使用测试账号 `admin / tm666666` 登录
4. 登录成功后进入后台工作台 `/dashboard`
5. 继续验证后台相关页面与接口联调

#### 数据初始化口径

- 后端初始化目录：`api-chaoshi/migrations/`
- 当前仅保留 `20240101000000_full_init.sql` 一个有效初始化文件
- 回归脚本：`api-chaoshi/scripts/regression-test.sh init`
- 默认最小初始化账号：
  - 后台管理员：`admin / tm666666`
  - 商家员工：`merchant / merchant123`

#### 可测试功能

| 模块   | 功能点            | 页面/入口                                                 |
| ---- | -------------- | ----------------------------------------------------- |
| 工作台  | 数据看板、快捷入口      | `/dashboard`                                          |
| 门店管理 | 门店列表、门店详情、门店编辑 | `/merchants` / `/merchants/:id` / `/merchants/:id/edit` |
| 订单管理 | 按商家筛选查看订单与详情   | `/orders` / `/orders/:id`                             |
| 数据分析 | 商家统计           | `/analytics`                                          |
| 设置   | 修改密码、查看后台资料   | `/settings`                                           |

### 9.4 商家端测试

#### 入口与登录

- 入口页面：`pages/auth/login`
- 登录后首页：`pages/merchant/home`
- 使用测试账号：`merchant / merchant123`

#### 使用步骤

1. 在微信开发者工具中打开小程序项目
2. 选择"商家登录"编译模式，或直接进入 `pages/auth/login`
3. 使用测试账号 `merchant / merchant123` 登录
4. 登录成功后进入商家工作台
5. 按需继续验证商品、订单、设置与员工能力

#### 商家端自测清单

| 模块   | 功能点            | 说明                                                                                                                           |
| ---- | -------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| 工作台  | 系统公告滚动栏        | 可查看公告详情、可关闭（本地记忆）                                                                                                            |
| 工作台  | 快速核销入口         | 支持扫一扫和输入 6 位核销码快速核销，核销成功后可直接查看订单详情                                                                                           |
| 商品管理 | 商品详情/编辑/删除     | 列表进入详情页正常，编辑保存成功，删除后列表隐藏                                                                                                     |
| 商品管理 | 商品图片上传         | 调用 `/api/v1/upload/file` 直传到后端本地存储，上传成功后页面回显                                                                                |
| 商品管理 | 私有图片显示         | 商品列表、商品详情、编辑页回显使用接口返回的可访问图片地址                                                                                                |
| 订单   | 核销码校验          | 输入 6 位数字核销码，核销成功/失败提示正确，订单详情可查看核销码、核销时间与核销人                                                                                  |
| 设置   | 配送设置-按距离收费     | 支持新增/删除规则、区间合法性校验，保存后重新进入仍生效                                                                                                 |
| 设置   | 账号与员工能力        | 支持修改密码、仅店主可见的员工管理入口                                                                                                          |

### 9.5 C 端店铺测试

#### 9.5.1 测试方法一：测试入口页面（推荐）

##### 功能说明

创建了专门的测试入口页面 `pages/store/test-entry.vue`，用于在小程序发布前测试用户下单流程。

##### 页面路径

- 路由：`/pages/store/test-entry`
- 导航标题：`商家店铺测试`
- 背景色：`#667eea`（紫色渐变）

##### 功能特性

1. **商家ID输入**：支持手动输入商家ID
2. **商家列表选择**：提供预设商家列表，点击即可快速选择
3. **直接跳转店铺**：输入商家ID后直接跳转至店铺首页

##### 使用步骤

1. 在微信开发者工具中打开小程序项目
2. 进入"商家店铺测试"页面
3. 在输入框中输入商家ID（默认：1）
4. 或直接点击商家列表中的商家
5. 点击"进入商家店铺"按钮
6. 进入店铺首页后可进行商品浏览、加入购物车、创建订单等操作

##### 测试商家列表

| 商家ID | 商家名称  | 状态  |
| ---- | ----- | --- |
| 1    | 美味餐厅  | 营业中 |
| 2    | 示例商家B | 营业中 |
| 3    | 示例商家C | 营业中 |

##### 注意事项

- 订单数据会真实写入数据库
- 此页面仅用于 C 端店铺测试，不用于服务商登录

#### 9.5.3 测试方法三：编译模式直接进入

##### 功能说明

利用微信开发者工具的编译模式功能，直接编译到指定页面并携带参数。

##### 配置步骤

1. 打开微信开发者工具
2. 点击顶部"编译模式"下拉框
3. 选择"添加编译模式"
4. 配置以下参数：

| 参数   | 值                  | 说明     |
| ---- | ------------------ | ------ |
| 编译模式 | `pages/store/home` | 店铺首页   |
| 启动参数 | `merchant_id=1`    | 商家ID参数 |

##### 启动参数说明

| 参数名          | 类型     | 必填 | 说明                 |
| ------------ | ------ | -- | ------------------ |
| merchant\_id | number | 是  | 商家ID，用于加载对应商家的店铺数据 |
| scene | string | 否 | 扫码场景值，支持解析 `merchant_id` 与业务扩展参数 |

##### 使用示例

```text
启动参数：merchant_id=1
```

访问商家ID为1的店铺首页

```text
启动参数：merchant_id=2
```

访问商家ID为2的店铺首页

#### 9.5.4 可测试功能模块

| 模块   | 功能           | 测试方法一 | 测试方法二 |
| ---- | ------------ | ----- | ----- |
| 店铺浏览 | 查看店铺信息、公告    | ✅     | ✅     |
| 商品展示 | 商品列表、分类筛选    | ✅     | ✅     |
| 商品详情 | 查看商品详情、规格选择  | ✅     | ✅     |
| 购物车  | 添加商品、修改数量、删除 | ✅     | ✅     |
| 下单流程 | 确认订单、选择配送方式  | ✅     | ✅     |
| 下单创建 | 创建订单并跳转订单列表  | ✅     | ✅     |
| 订单管理 | 查看订单列表、订单详情  | ✅     | ✅     |
| 退款申请 | 申请退款、查看退款状态  | ✅     | ✅     |

### 9.6 环境与部署说明

#### 9.6.1 开发环境配置

1. 打开微信开发者工具
2. 进入"详情" → "本地设置"
3. 勾选"不校验合法域名"（开发阶段）
4. 勾选"不校验HTTPS证书"
5. 小程序前端统一通过环境变量文件配置地址：
   - 开发环境使用 `miniprogram/.env.development`
   - 生产环境使用 `miniprogram/.env.production`
   - 示例见 `miniprogram/.env.example`
6. `VITE_API_BASE_URL` 管理 HTTP 接口基址
7. 业务代码禁止再直接写死 `localhost:8080`
8. H5 开发代理统一由 `miniprogram/vite.config.ts` 读取环境变量生成

#### 9.6.2 生产环境要求

- 需要配置已备案的域名
- 需要配置SSL证书（HTTPS）
- 需要在微信小程序后台添加服务器域名配置
- 需要配置上传目录与静态资源访问路径
- 如需未来接入新支付方案，单独追加部署说明

### 9.7 调试技巧

#### 查看接口请求

1. 打开微信开发者工具
2. 切换到"Network"面板
3. 筛选接口请求
4. 查看请求参数和响应数据

#### 查看控制台日志

1. 在关键代码位置添加 `console.log`
2. 在"Console"面板查看日志输出
3. 检查接口返回的数据结构

#### 数据检查

- 使用后端接口直接查看数据库状态
- 使用 Postman 或 curl 测试 API 接口
- 检查订单状态是否正确更新
