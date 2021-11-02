package user

import (
	"github.com/gin-gonic/gin"
	v1 "idea_server/api/v1"
)

type UserBaseRouter struct {

}

func (e UserBaseRouter) InitUserBaseRouter(Router *gin.RouterGroup)  {
	customerRouter := Router.Group("")
	var userBaseApi = v1.ApiGroupApp.UserApiGroup.UserBaseApi
	// 开启一个新的作用域，目前只是代码视觉分割的作用
	{
		customerRouter.POST("register", userBaseApi.Register)           // 邮箱是否存在
		customerRouter.POST("get_email_code", userBaseApi.GetEmailCode) // 生成邮箱验证码
	}
}
