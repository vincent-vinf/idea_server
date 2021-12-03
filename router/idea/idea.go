package idea

import (
	"github.com/gin-gonic/gin"
	v1 "idea_server/api/v1"
)

type IdeaRouter struct {

}

func (e IdeaRouter) InitIdeaRouter(Router *gin.RouterGroup)  {
	customerRouter := Router.Group("idea")
	var ideaApi = v1.ApiGroupApp.IdeaApiGroup.IdeaApi
	{
		customerRouter.POST("create_idea", ideaApi.CreateIdea)
		customerRouter.POST("get_idea_info", ideaApi.GetIdeaInfo)
	}

}