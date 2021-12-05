package idea

import "idea_server/service"

type ApiGroup struct {
	IdeaApi
	IdeaCommentApi
}

var ideaService = service.ServiceGroupApp.IdeaServiceGroup.IdeaService
var ideaCommentService = service.ServiceGroupApp.IdeaServiceGroup.IdeaCommentService
