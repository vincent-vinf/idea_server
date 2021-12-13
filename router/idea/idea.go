package idea

import (
	"github.com/gin-gonic/gin"
	v1 "idea_server/api/v1"
)

type IdeaRouter struct {
}

func (e IdeaRouter) InitIdeaRouter(Router *gin.RouterGroup) {
	customerRouter := Router.Group("idea")
	var ideaApi = v1.ApiGroupApp.IdeaApiGroup.IdeaApi
	{
		customerRouter.POST("create_idea", ideaApi.CreateIdea)
		customerRouter.POST("delete_idea", ideaApi.DeleteIdea)
		customerRouter.POST("get_idea_info", ideaApi.GetIdeaInfo)
		customerRouter.POST("get_idea_list", ideaApi.GetIdeaList)
		customerRouter.POST("get_follow_idea_list", ideaApi.GetFollowIdeaList)
		customerRouter.POST("get_my_idea_list", ideaApi.GetMyIdeaList)
		customerRouter.POST("get_similar_ideas", ideaApi.GetSimilarIdeas)
	}

}
