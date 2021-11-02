package user

import (
	"idea_server/global"
	"idea_server/model/user"
	"idea_server/model/user/response"
)

type UserService struct {

}

func (u *UserService) GetMyInfo(id int) (userInfo response.UserResponse) {
	global.IDEA_DB.Model(&user.User{}).Where("id = ?", id).Find(&userInfo)

	return userInfo
}
