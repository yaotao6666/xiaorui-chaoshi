package jsbank

// Config 江苏银行支付配置
type Config struct {
	Enabled        bool
	AppID          string
	MchID          string
	MchName        string
	MasterAccount  string
	PartnerID      string
	DeviceNo       string
	PostURL        string
	NotifyURL      string
	PFXPath        string
	PFXPassword    string
	PublicCertPath string
}

// MiniProgramPayRequest 小程序支付请求
type MiniProgramPayRequest struct {
	OrderNo string
	Amount  float64
	OpenID  string
	AppID   string
}

// MiniProgramPayParams 小程序支付参数
type MiniProgramPayParams struct {
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

// MiniProgramPayResponse 小程序支付结果
type MiniProgramPayResponse struct {
	RespCode string
	RespMsg  string
	Raw      map[string]string
	PayParams MiniProgramPayParams
}

// RefundRequest 退款请求
type RefundRequest struct {
	OrderNo  string
	RefundNo string
	Amount   float64
}

// RefundResponse 退款结果
type RefundResponse struct {
	RespCode    string
	RespMsg     string
	OrderStatus string
	Raw         map[string]string
}

// QRCodePayRequest 扫码支付请求
type QRCodePayRequest struct {
	OrderNo string
	Amount  float64
	MchIP   string
}

// QRCodePayResponse 扫码支付结果
type QRCodePayResponse struct {
	RespCode string
	RespMsg  string
	Raw      map[string]string
	PayURL   string
}
