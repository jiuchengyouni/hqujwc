package jwc

type GetSessionIBody struct {
	Reportlet string `json:"reportlet"`
	XH        int    `json:"xh"`
	BBWID     string `json:"BBWID"`
}
