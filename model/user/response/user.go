package response

import "idea_server/model/user"

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type UserInfoResponse struct {
	user.User
	IsFollow bool `json:"isFollow"`
}

type NoticeField struct {
	ID   uint `json:"id"`
	Type uint `json:"type"`
}

type NoticeResponse struct {
	NoticeField
	Data interface{} `json:"data"`
}
