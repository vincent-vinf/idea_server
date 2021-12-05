package idea

import (
	"github.com/gin-gonic/gin"
	v1 "idea_server/api/v1"
)

type IdeaCommentRouter struct {

}

func (e IdeaCommentRouter) InitIdeaCommentRouter(Router *gin.RouterGroup)  {
	customerRouter := Router.Group("idea")
	var ideaCommentApi = v1.ApiGroupApp.IdeaApiGroup.IdeaCommentApi
	{
		customerRouter.POST("create_comment", ideaCommentApi.CreateComment)
		customerRouter.POST("delete_comment", ideaCommentApi.DeleteComment)
		customerRouter.POST("get_comment_list", ideaCommentApi.GetCommentList)
	}
}