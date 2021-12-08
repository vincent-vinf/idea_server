package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/model/common/request"
	"idea_server/model/common/response"
	"idea_server/utils/constant"
)

type UserApi struct {
}

func (u *UserApi) GetMyInfo(c *gin.Context) {
	userInfo, _ := c.Get(constant.IdentityKey)
	if userInfo == nil {
		global.IDEA_LOG.Error("获取个人信息失败")
		response.Fail(c)
		return
	}
	response.OkWithDetailed(userInfo, "获取个人信息成功", c)
}


func (u *UserApi) GetUserInfo(c *gin.Context) {
	var ids request.IdsReq
	_ = c.ShouldBindJSON(&ids)
	if info, err := userService.GetUserInfo(ids.Ids); err != nil {
		global.IDEA_LOG.Error("获取用户信息失败", zap.Error(err))
		response.FailWithMessage("获取用户信息失败", c)
	} else  {
		response.OkWithDetailed(info, "获取用户信息成功", c)
	}
}
