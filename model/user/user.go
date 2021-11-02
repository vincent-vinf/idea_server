package user

import "idea_server/global"

// User 结构体
type User struct {
	// TODO recover
	global.IDEA_MODEL
	Username string `json:"username"`
	Email    string `json:"email"`
	Passwd   string `json:"passwd"`
}

// TableName User 表名
func (u *User) TableName() string {
	return "users"
}
