package user

// User 结构体
type User struct {
	// TODO recover
	//global.GVA_MODEL
	ID       uint   `gorm:"primarykey"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Passwd   string `json:"passwd"`
}

// TableName User 表名
func (u *User) TableName() string {
	return "users"
}
