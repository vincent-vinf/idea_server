package idea

import "idea_server/global"

// IdeaLike 结构体
type IdeaLike struct {
	global.IDEA_MODEL
	IdeaId uint `json:"ideaId"`
	UserId uint `json:"userId"`
}

// TableName User 表名
func (u *IdeaLike) TableName() string {
	return "like"
}
