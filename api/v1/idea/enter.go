package idea

import "idea_server/service"

type ApiGroup struct {
	IdeaApi
}

var ideaService = service.ServiceGroupApp.IdeaServiceGroup.IdeaService
