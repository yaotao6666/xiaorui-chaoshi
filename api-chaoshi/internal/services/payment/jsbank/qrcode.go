package jsbank

import (
	"fmt"
	"strings"
)

func (c *Client) CreateQRCodePay(req QRCodePayRequest) (*QRCodePayResponse, error) {
	payload := c.buildCommonPayload()
	payload["service"] = "dPay"
	payload["outTradeNo"] = strings.TrimSpace(req.OrderNo)
	payload["totalFee"] = formatAmount(req.Amount)
	payload["qrValidTime"] = "15"
	payload["backUrl"] = c.config.NotifyURL
	payload["mchIp"] = strings.TrimSpace(req.MchIP)
	payload["deviceNo"] = c.config.DeviceNo

	rawResult, err := c.doRequest(payload)
	if err != nil {
		return nil, err
	}

	payURL := strings.TrimSpace(rawResult["qrCode"])
	if payURL == "" {
		return nil, fmt.Errorf("江苏银行返回的二维码地址为空")
	}

	return &QRCodePayResponse{
		RespCode: strings.TrimSpace(rawResult["respCode"]),
		RespMsg:  strings.TrimSpace(rawResult["respMsg"]),
		Raw:      rawResult,
		PayURL:   payURL,
	}, nil
}
