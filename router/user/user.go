package user

import (
	"github.com/gin-gonic/gin"
	v1 "idea_server/api/v1"
)

type UserRouter struct {

}

func (e UserRouter) InitUserBaseRouter(Router *gin.RouterGroup)  {
	customerRouter := Router.Group("user")
	var userBaseApi = v1.ApiGroupApp.UserApiGroup.UserApi
	{
		customerRouter.POST("test", userBaseApi.IsExistEmail)   // 邮箱是否存在
	}
}
