package request

type Login struct {
	Email    string `json:"email" form:"email"`
	Passwd   string `json:"passwd" form:"passwd"`
	Code     string `json:"code" form:"code"`
}

type Register struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Passwd   string `json:"passwd" form:"passwd"`
	Code     string `json:"code" form:"code"`
}
