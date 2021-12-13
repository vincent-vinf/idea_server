package user

import (
	"errors"
	"gorm.io/gorm"
	"idea_server/global"
	"idea_server/model/user"
)

type UserFollowService struct {
}

func (e *UserFollowService) CreateFollow(followedId, followId uint) (err error) {
	if !errors.Is(global.IDEA_DB.Where("followed_id = ? AND follow_id = ?", followedId, followId).First(&user.UserFollow{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("重复关注")
	}

	err = global.IDEA_DB.Create(&user.UserFollow{
		FollowedId: followedId,
		FollowId:   followId,
	}).Error
	return
}

func (e *UserFollowService) DeleteFollow(followedId, followId uint) (err error) {
	var id uint
	err = global.IDEA_DB.Model(&user.UserFollow{}).Select("id").Where("followed_id = ? AND follow_id = ?", followedId, followId).Row().Scan(&id)
	if err != nil {
		return
	}
	err = global.IDEA_DB.Delete(&user.UserFollow{}, id).Error
	return
}

func (e *UserFollowService) IsFollow(followedId, followId uint) bool {
	if errors.Is(global.IDEA_DB.Where("followed_id = ? AND follow_id = ?", followedId, followId).First(&user.UserFollow{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (e *UserFollowService) GetFollowList(userId uint) (list []uint, err error) {
	err = global.IDEA_DB.Model(&user.UserFollow{}).Where("follow_id = ?", userId).Select("followed_id").Find(&list).Error
	if err != nil {
		return make([]uint, 0, 1), err
	}
	return list, err
}
