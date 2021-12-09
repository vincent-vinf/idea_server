package idea

import "idea_server/global"

// IdeaType 结构体
type IdeaType struct {
	global.IDEA_MODEL
	Name string `json:"name"`
}

// TableName User 表名
func (u *IdeaType) TableName() string {
	return "type"
}
