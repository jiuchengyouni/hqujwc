package types

type WxAccessRequest struct {
	Signature string `json:"signature"`
	Timestamp string `json:"timestamp"`
	Nonce     string `json:"nonce"`
	Echoster  string `json:"echoster"`
}
