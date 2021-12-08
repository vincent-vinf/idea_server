package idea

import (
	"github.com/gin-gonic/gin"
	v1 "idea_server/api/v1"
)

type IdeaLikeRouter struct {

}

func (e IdeaLikeRouter) InitIdeaLikeRouter(Router *gin.RouterGroup)  {
	customerRouter := Router.Group("idea")
	var ideaLikeApi = v1.ApiGroupApp.IdeaApiGroup.IdeaLikeApi
	{
		customerRouter.POST("create_like", ideaLikeApi.CreateLike)
		customerRouter.POST("delete_like", ideaLikeApi.DeleteLike)
	}
}