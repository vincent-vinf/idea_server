package user

import "idea_server/global"

// UserFollow 结构体
type UserFollow struct {
	global.IDEA_MODEL
	FollowedId uint `json:"followedId"`
	FollowId   uint `json:"followId"`
}

// TableName UserFollow 表名
func (u *UserFollow) TableName() string {
	return "follow"
}
