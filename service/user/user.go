package user

import (
	"idea_server/global"
	"idea_server/model/user"
)

type UserService struct {

}

func (e *UserService) GetUserInfo(ids []int) (infos []user.User, err error) {
	err = global.IDEA_DB.Omit("email", "passwd", "weight").Find(&infos, ids).Error
	return
}