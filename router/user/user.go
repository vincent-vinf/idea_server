package user

import (
	"github.com/gin-gonic/gin"
	v1 "idea_server/api/v1"
)

type UserFollowRouter struct {

}

func (e UserFollowRouter) InitUserFollowRouter(Router *gin.RouterGroup)  {
	customerRouter := Router.Group("user")
	var userFollowApi = v1.ApiGroupApp.UserApiGroup.UserFollowApi
	{
		customerRouter.POST("create_follow", userFollowApi.CreatFollow)
		customerRouter.POST("delete_follow", userFollowApi.DeleteFollow)
	}

}
