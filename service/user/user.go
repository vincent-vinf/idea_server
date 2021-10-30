package user

import (
	"errors"
	"gorm.io/gorm"
	"idea_server/global"
	"idea_server/model/user"
)

type UserService struct {

}

func (u *UserService) IsExistEmail(email string) (bool, error) {
	if err := global.IDEA_DB.Where("email = ?", email).First(&user.User{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
