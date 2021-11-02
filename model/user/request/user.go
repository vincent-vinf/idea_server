package request

type Login struct {
	Email    string `json:"email"`
	Passwd   string `json:"passwd"`
	Code     string `json:"code"`
}

type Register struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Passwd   string `json:"passwd"`
	Code     string `json:"code"`
}
