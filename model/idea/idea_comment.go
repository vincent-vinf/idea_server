package idea

import "idea_server/global"

// IdeaComment 结构体
type IdeaComment struct {
	global.IDEA_MODEL
	IdeaId    uint   `json:"ideaId"`
	UserId    uint   `json:"userId"`
	ToId      uint   `json:"toId"`
	CommentId uint   `json:"commentId"`
	Content   string `json:"content"`
}

// TableName User 表名
func (u *IdeaComment) TableName() string {
	return "comment"
}
