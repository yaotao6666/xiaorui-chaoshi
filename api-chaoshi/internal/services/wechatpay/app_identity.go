package wechatpay

import (
	"fmt"
	"strings"

	"chaoshi_api/internal/config"
)

type AppIdentity struct {
	AppID     string
	AppSecret string
}

func GetActiveAppIdentity() (*AppIdentity, error) {
	if config.Config == nil {
		return nil, fmt.Errorf("应用配置未初始化")
	}

	appID := strings.TrimSpace(config.Config.Wechat.AppID)
	appSecret := strings.TrimSpace(config.Config.Wechat.AppSecret)
	if appID == "" || appSecret == "" {
		return nil, fmt.Errorf("微信小程序配置不完整")
	}

	return &AppIdentity{
		AppID:     appID,
		AppSecret: appSecret,
	}, nil
}
