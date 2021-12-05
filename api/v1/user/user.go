package user

import (
	"github.com/gin-gonic/gin"
	"idea_server/global"
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
