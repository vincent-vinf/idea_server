package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/model/common/response"
	userReq "idea_server/model/user/request"
	"idea_server/utils"
)

type UserFollowApi struct {
}

func (u *UserFollowApi) CreatFollow(c *gin.Context) {
	var info userReq.GetByFollowedId
	_ = c.ShouldBindJSON(&info)
	if err := userFollowService.CreateFollow(info.Uint(), utils.GetJwtId(c)); err != nil {
		global.IDEA_LOG.Error("关注用户失败", zap.Error(err))
		response.FailWithMessage("关注用户失败", c)
	} else {
		response.OkWithMessage("关注用户成功", c)
	}
}

func (u *UserFollowApi) DeleteFollow(c *gin.Context) {
	var info userReq.GetByFollowedId
	_ = c.ShouldBindJSON(&info)
	if err := userFollowService.DeleteFollow(info.Uint(), utils.GetJwtId(c)); err != nil {
		global.IDEA_LOG.Error("取消用户关注失败", zap.Error(err))
		response.FailWithMessage("取消用户关注失败", c)
	} else {
		response.OkWithMessage("取消用户关注成功", c)
	}
}
