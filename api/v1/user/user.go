package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/model/common/response"
	"idea_server/model/user"
)

type UserApi struct {
}

func (u *UserApi) IsExistEmail(c *gin.Context) {
	var msg user.User
	_ = c.ShouldBindJSON(&msg)
	if isExist, err := userBaseService.IsExistEmail(msg.Email); err != nil {
		global.IDEA_LOG.Error("判断是否存在邮箱失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		if isExist {
			response.FailWithMessage("email 已使用", c)
		} else {
			response.OkWithMessage("email 未使用", c)
		}
	}
}
