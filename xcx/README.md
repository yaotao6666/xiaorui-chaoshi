# xcx 小程序验证壳

该目录用于验证以下交付链路：

- 小程序端先调用 `Taro.login`
- 将 `code` 提交到 `POST /api/v1/auth/user/wechat-login`
- 成功后拿到 `token/openid/app_id`
- 再通过 `web-view` 打开现有 H5，并把上述参数拼入 H5 地址

## 当前页面

- `pages/index/index`：启动页，负责无感登录并进入 H5
- `pages/settings/index`：配置页，维护 API 基址、H5 入口、门店 ID 与 H5 预览调试参数
- `pages/debug/index`：调试页，查看本地登录态与最终拼好的 H5 地址
- `pages/webview/index`：H5 容器页

## 默认配置

- API 基址：`http://127.0.0.1:8081`
- H5 入口：`http://localhost:3000/#/pages/store/home`
- 默认门店：`merchant_id=1`

## 测试说明

- 微信真机环境会实际调用 `Taro.login -> /api/v1/auth/user/wechat-login`
- H5 预览环境无法拿到真实微信 `code`，可在配置页填写 `debugToken/debugOpenid` 验证地址拼装与 `web-view` 跳转
- 微信小程序 `web-view` 无法直接访问 `localhost`，如需真机验证嵌入，请替换为可公网访问且已配置业务域名的 HTTPS H5 地址
