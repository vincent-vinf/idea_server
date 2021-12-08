package idea

import "idea_server/global"

// IdeaLike 结构体
type Idea struct {
	global.IDEA_MODEL
	UserId  uint    `json:"userId"`
	Simple  string  `json:"simple"`
	Content string  `json:"content"`
	Life    float64 `json:"life"`
}

// TableName User 表名
func (u *Idea) TableName() string {
	return "ideas"
}
