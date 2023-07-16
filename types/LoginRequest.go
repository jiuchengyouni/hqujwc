package types

type LoginRequestBody struct {
	Stunum   string `json:"stunum"`
	Password string `json:"password"`
}

type LoginInfo struct {
	Stunum  string `json:"stunum"`
	Stupass string `json:"stupass"`
}
