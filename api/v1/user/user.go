package user

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/model/common/response"
	"idea_server/utils/constant"
	"strconv"
)

type UserApi struct {
}

func (u *UserApi) GetMyInfo(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	id, err := strconv.Atoi(claims[constant.IdentityKey].(string))
	if err != nil {
		global.IDEA_LOG.Error("获取个人信息失败", zap.Error(err))
		response.Fail(c)
		return
	}
	userInfo := userService.GetMyInfo(id)
	response.OkWithData(userInfo, c)
}
