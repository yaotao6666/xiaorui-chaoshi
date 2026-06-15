package jsbank

import (
	"crypto/rsa"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"chaoshi_api/internal/config"
)

type Client struct {
	config     Config
	privateKey *rsa.PrivateKey
}

func NewClient() (*Client, error) {
	if config.Config == nil {
		return nil, fmt.Errorf("应用配置未初始化")
	}

	clientConfig := Config{
		Enabled:        config.Config.JsBankPay.Enabled,
		AppID:          strings.TrimSpace(config.Config.JsBankPay.AppID),
		MchID:          strings.TrimSpace(config.Config.JsBankPay.MchID),
		MchName:        strings.TrimSpace(config.Config.JsBankPay.MchName),
		MasterAccount:  strings.TrimSpace(config.Config.JsBankPay.MasterAccount),
		PartnerID:      strings.TrimSpace(config.Config.JsBankPay.PartnerID),
		DeviceNo:       strings.TrimSpace(config.Config.JsBankPay.DeviceNo),
		PostURL:        strings.TrimSpace(config.Config.JsBankPay.PostURL),
		NotifyURL:      strings.TrimSpace(config.Config.JsBankPay.NotifyURL),
		PFXPath:        strings.TrimSpace(config.Config.JsBankPay.PFXPath),
		PFXPassword:    config.Config.JsBankPay.PFXPassword,
		PublicCertPath: strings.TrimSpace(config.Config.JsBankPay.PublicCertPath),
	}

	if !clientConfig.Enabled {
		return nil, fmt.Errorf("江苏银行支付未启用")
	}

	if clientConfig.AppID == "" {
		clientConfig.AppID = strings.TrimSpace(config.Config.Wechat.AppID)
	}

	requiredFields := map[string]string{
		"JSBANK_PAY_APP_ID":          clientConfig.AppID,
		"JSBANK_PAY_PARTNER_ID":      clientConfig.PartnerID,
		"JSBANK_PAY_DEVICE_NO":       clientConfig.DeviceNo,
		"JSBANK_PAY_POST_URL":        clientConfig.PostURL,
		"JSBANK_PAY_NOTIFY_URL":      clientConfig.NotifyURL,
		"JSBANK_PAY_PFX_PATH":        clientConfig.PFXPath,
		"JSBANK_PAY_PFX_PASSWORD":    clientConfig.PFXPassword,
		"JSBANK_PAY_PUBLIC_CERT_PATH": clientConfig.PublicCertPath,
	}
	for key, value := range requiredFields {
		if strings.TrimSpace(value) == "" {
			return nil, fmt.Errorf("%s 未配置", key)
		}
	}

	privateKey, err := loadPrivateKeyFromPFX(clientConfig.PFXPath, clientConfig.PFXPassword)
	if err != nil {
		return nil, err
	}

	return &Client{
		config:     clientConfig,
		privateKey: privateKey,
	}, nil
}

func (c *Client) CreateMiniProgramPay(req MiniProgramPayRequest) (*MiniProgramPayResponse, error) {
	appID := strings.TrimSpace(req.AppID)
	if appID == "" {
		appID = c.config.AppID
	}
	if appID == "" {
		return nil, fmt.Errorf("支付 AppID 未配置")
	}
	if strings.TrimSpace(req.OpenID) == "" {
		return nil, fmt.Errorf("用户 openid 不能为空")
	}

	payload := c.buildCommonPayload()
	payload["service"] = "paymentWXPay"
	payload["totalFee"] = formatAmount(req.Amount)
	payload["tradeType"] = "JSAPI"
	payload["extfld2"] = strings.TrimSpace(req.OrderNo)
	payload["backUrl"] = c.config.NotifyURL
	payload["mchIp"] = "0.0.0.0"
	payload["extfld1"] = appID
	payload["openId"] = strings.TrimSpace(req.OpenID)

	rawResult, err := c.doRequest(payload)
	if err != nil {
		return nil, err
	}

	payParams := MiniProgramPayParams{
		TimeStamp: strings.TrimSpace(rawResult["timeStamp"]),
		NonceStr:  strings.TrimSpace(rawResult["nonceStr"]),
		Package:   normalizeMiniProgramPackage(rawResult),
		SignType:  getOrDefault(strings.TrimSpace(rawResult["signType"]), "RSA"),
		PaySign:   strings.TrimSpace(rawResult["paySign"]),
	}
	if payParams.TimeStamp == "" || payParams.NonceStr == "" || payParams.Package == "" || payParams.PaySign == "" {
		return nil, fmt.Errorf("江苏银行返回的小程序支付参数不完整")
	}

	return &MiniProgramPayResponse{
		RespCode: strings.TrimSpace(rawResult["respCode"]),
		RespMsg:  strings.TrimSpace(rawResult["respMsg"]),
		Raw:      rawResult,
		PayParams: payParams,
	}, nil
}

func (c *Client) Refund(req RefundRequest) (*RefundResponse, error) {
	orderNo := strings.TrimSpace(req.OrderNo)
	refundNo := strings.TrimSpace(req.RefundNo)
	if orderNo == "" {
		return nil, fmt.Errorf("原订单号不能为空")
	}
	if refundNo == "" {
		return nil, fmt.Errorf("退款单号不能为空")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("退款金额必须大于0")
	}

	payload := c.buildCommonPayload()
	payload["service"] = "payRefund"
	payload["outTradeNo"] = orderNo
	payload["outRefundNo"] = refundNo
	payload["refundAmt"] = formatAmount(req.Amount)

	rawResult, err := c.doRequest(payload)
	if err != nil {
		return nil, err
	}

	return &RefundResponse{
		RespCode:    strings.TrimSpace(rawResult["respCode"]),
		RespMsg:     strings.TrimSpace(rawResult["respMsg"]),
		OrderStatus: strings.TrimSpace(rawResult["orderStatus"]),
		Raw:         rawResult,
	}, nil
}

func (c *Client) VerifyNotify(params map[string]string) error {
	publicKey, err := loadPublicKeyFromCert(c.config.PublicCertPath)
	if err != nil {
		return err
	}
	return verifyPayload(params, publicKey)
}

func (c *Client) doRequest(payload map[string]string) (map[string]string, error) {
	signature, err := signPayload(payload, c.privateKey)
	if err != nil {
		return nil, err
	}

	payload["sign"] = signature
	payload["signType"] = "RSA"

	form := url.Values{}
	for key, value := range payload {
		form.Set(key, value)
	}

	httpClient := &http.Client{Timeout: 60 * time.Second}
	response, err := httpClient.Post(
		c.config.PostURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("请求江苏银行失败: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取江苏银行响应失败: %w", err)
	}

	result := parseKVResponse(string(body))
	if len(result) == 0 {
		return nil, fmt.Errorf("江苏银行响应为空: %s", string(body))
	}
	if code := strings.TrimSpace(result["respCode"]); code != "" && code != "000000" {
		return nil, fmt.Errorf("江苏银行返回错误: %s", strings.TrimSpace(result["respMsg"]))
	}

	return result, nil
}

func (c *Client) buildCommonPayload() map[string]string {
	now := time.Now()
	return map[string]string{
		"createDate":    now.Format("20060102"),
		"createTime":    now.Format("150405"),
		"bizDate":       now.Format("20060102"),
		"msgID":         now.Format("20060102150405.000000000"),
		"svrCode":       "",
		"partnerId":     c.config.PartnerID,
		"channelNo":     "m",
		"publicKeyCode": "00",
		"version":       "v1.0.0",
		"charset":       "utf-8",
		"signType":      "RSA",
		"sign":          "",
		"deviceNo":      c.config.DeviceNo,
	}
}

func parseKVResponse(body string) map[string]string {
	result := make(map[string]string)
	for _, item := range strings.Split(strings.TrimSpace(body), "&") {
		if strings.TrimSpace(item) == "" {
			continue
		}
		key, value, found := strings.Cut(item, "=")
		if !found {
			continue
		}
		decodedValue, err := url.QueryUnescape(value)
		if err != nil {
			result[key] = value
			continue
		}
		result[key] = decodedValue
	}
	return result
}

func formatAmount(amount float64) string {
	normalized := math.Round(amount*100) / 100
	return fmt.Sprintf("%.2f", normalized)
}

func normalizeMiniProgramPackage(raw map[string]string) string {
	candidates := []string{
		strings.TrimSpace(raw["package"]),
		strings.TrimSpace(raw["packAge"]),
	}
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if strings.HasPrefix(candidate, "prepay_id=") {
			return candidate
		}
		return "prepay_id=" + candidate
	}
	return ""
}

func getOrDefault(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
