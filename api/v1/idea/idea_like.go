package idea

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/model/common/response"
	ideaReq "idea_server/model/idea/request"
	"idea_server/utils"
)

type IdeaLikeApi struct {
}

func (e *IdeaLikeApi) CreateLike(c *gin.Context) {
	var info ideaReq.GetByIdeaId
	_ = c.ShouldBindJSON(&info)
	if err := ideaLikeService.CreateLike(utils.GetJwtId(c), info.Uint()); err != nil {
		global.IDEA_LOG.Error("点赞失败", zap.Error(err))
		response.FailWithMessage("点赞失败", c)
	} else {
		response.OkWithMessage("点赞成功", c)
	}
}

func (e *IdeaLikeApi) DeleteLike(c *gin.Context) {
	var info ideaReq.GetByIdeaId
	_ = c.ShouldBindJSON(&info)
	if err := ideaLikeService.DeleteLike(utils.GetJwtId(c), info.Uint()); err != nil {
		global.IDEA_LOG.Error("取消点赞失败", zap.Error(err))
		response.FailWithMessage("取消点赞失败", c)
	} else {
		response.OkWithMessage("取消点赞成功", c)
	}
}