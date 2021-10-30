package user

import (
	"github.com/gin-gonic/gin"
	"idea_server/model/common/response"
	"idea_server/model/user"
)

type UserApi struct {
}

func (u *UserApi) IsExistEmail(c *gin.Context) {
	var msg user.User
	_ = c.ShouldBindJSON(&msg)
	if isExist, err := userBaseService.IsExistEmail(msg.Email); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		if isExist {
			response.FailWithMessage("email 已使用", c)
		} else {
			response.OkWithMessage("email 未使用", c)
		}
	}
}
