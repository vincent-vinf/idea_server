package idea

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/model/common/request"
	"idea_server/model/common/response"
	"idea_server/model/idea"
	"idea_server/utils"
)

type IdeaCommentApi struct {
}

func (e *IdeaCommentApi) CreateComment(c *gin.Context) {
	var info idea.IdeaComment
	_ = c.ShouldBindJSON(&info)
	info.UserId = utils.GetJwtId(c)
	if err := utils.Verify(info, utils.CreateCommentVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := ideaCommentService.CreateComment(&info); err != nil {
		global.IDEA_LOG.Error("创建评论失败", zap.Error(err))
		response.FailWithDetailed(err.Error(),"创建评论失败", c)
	} else {
		response.OkWithMessage("创建评论成功", c)
	}

}

func (e *IdeaCommentApi) DeleteComment(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := ideaCommentService.DeleteComment(idInfo.Uint()); err != nil {
		global.IDEA_LOG.Error("删除评论失败", zap.Error(err))
		response.FailWithMessage("删除评论失败", c)
	} else {
		response.OkWithMessage("删除评论成功", c)
	}
}
