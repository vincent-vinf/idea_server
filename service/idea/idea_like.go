package idea

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"idea_server/global"
	"idea_server/model/idea"
)

type IdeaLikeService struct {
}

func (e *IdeaLikeService) CreateLike(userId, ideaId uint) (err error) {
	if !errors.Is(global.IDEA_DB.Where("user_id = ? AND idea_id = ?", userId, ideaId).First(&idea.IdeaLike{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("重复点赞")
	}

	err = global.IDEA_DB.Create(&idea.IdeaLike{IdeaId: ideaId, UserId: userId}).Error
	//TODO nil
	fmt.Println("err", err)
	return
}

func (e *IdeaLikeService) DeleteLike(userId, ideaId uint) (err error) {
	var likeId uint
	err = global.IDEA_DB.Model(&idea.IdeaLike{}).Select("id").Where("user_id = ? AND idea_id = ?", userId, ideaId).Row().Scan(&likeId)
	if err != nil {
		return
	}
	err = global.IDEA_DB.Delete(&idea.IdeaLike{}, likeId).Error
	return
}

func (e *IdeaLikeService) IsLike(userId, ideaId uint) bool {
	if errors.Is(global.IDEA_DB.Where("user_id = ? AND idea_id = ?", userId, ideaId).First(&idea.IdeaLike{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
