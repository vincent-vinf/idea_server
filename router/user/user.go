package user

import (
	"github.com/gin-gonic/gin"
	v1 "idea_server/api/v1"
)

type UserRouter struct {
}

func (e UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	customerRouter := Router.Group("user")
	var userApi = v1.ApiGroupApp.UserApiGroup.UserApi
	{
		customerRouter.POST("get_my_info", userApi.GetMyInfo)
		customerRouter.POST("get_user_info", userApi.GetUserInfo)
		customerRouter.POST("notice", userApi.Notice)
	}

}
