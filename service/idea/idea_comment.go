package idea

import (
	"errors"
	"gorm.io/gorm"
	"idea_server/global"
	"idea_server/model/idea"
	ideaRes "idea_server/model/idea/response"
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
	// 事务
	err := global.IDEA_DB.Transaction(func(tx *gorm.DB) error {
		var comment idea.IdeaComment
		if err := tx.Find(&comment, id).Error ; err != nil {
			return err
		}

		if comment.CommentId == 0 { // 删除直接回复评论的子评论
			if err := tx.Where("comment_id = ?", comment.ID).Delete(&idea.IdeaComment{}).Error ; err != nil {
				return err
			}
		}

		if err := tx.Delete(&idea.IdeaComment{}, id).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
	return err
}

func (e *IdeaCommentService) GetComment(ideaId uint) (err error, comments []ideaRes.IdeaCommentResponse) {
	var c []idea.IdeaComment // 回复想法
	var c2 []idea.IdeaComment // 回复评论

	err =global.IDEA_DB.Where("idea_id = ? AND comment_id = ?", ideaId, 0).Find(&c).Error
	if err != nil {
		return err, make([]ideaRes.IdeaCommentResponse, 0, 1)
	}

	comments = make([]ideaRes.IdeaCommentResponse, len(c), cap(c))

	for index, _ := range c {
		comments[index].IdeaComment = c[index]
		err = global.IDEA_DB.Where("idea_id = ? AND comment_id = ?", ideaId, c[index].ID).Find(&c2).Error
		comments[index].Replys = c2
		// TODO 容错性
		if err != nil {
			return err, make([]ideaRes.IdeaCommentResponse, 0, 1)
		}
	}

	return nil, comments
}
