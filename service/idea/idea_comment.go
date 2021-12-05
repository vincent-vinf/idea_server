package idea

import (
	"errors"
	"gorm.io/gorm"
	"idea_server/global"
	"idea_server/model/common/request"
	"idea_server/model/idea"
)

type IdeaCommentService struct {
}

func (e *IdeaCommentService) CreateComment(info *idea.IdeaComment) error {
	db := global.IDEA_DB

	// 类似这种判断，示例
	if errors.Is(db.Where("id = ?", info.IdeaId).First(&idea.Idea{}).Error, gorm.ErrRecordNotFound) { // 判断想法是否存在
		return errors.New("想法不存在")
	}

	err := db.Create(info).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *IdeaCommentService) DeleteComment(id uint) error {
	//err := global.IDEA_DB.Where("idea_id = ? AND user_id = ? AND to_id = ?",  info.IdeaId, info.UserId, info.ToId).Delete(&idea.IdeaComment{}).Error
	err := global.IDEA_DB.Delete(&idea.IdeaComment{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *IdeaCommentService) GetCommentList(ideaCommentInfo idea.IdeaComment, pageInfo request.PageInfo, order string, desc bool) (err error, list interface{}, total int64, num int) {
	ideaComments := make([]idea.IdeaComment, 0, 1)
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	db := global.IDEA_DB.Model(&idea.IdeaComment{}).Where("idea_id = ?", ideaCommentInfo.IdeaId)

	if err = db.Count(&total).Error; err != nil {
		return err, ideaComments, total, 0
	}
	err = db.Limit(limit).Offset(offset).Order("created_at").Find(&ideaComments).Error
	return err, ideaComments, total, len(ideaComments)
}
