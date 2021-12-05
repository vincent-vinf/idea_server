package response

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Num      int       `json:"num"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}
