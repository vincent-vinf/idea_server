package v1

import (
	"idea_server/api/v1/idea"
	"idea_server/api/v1/user"
)

type ApiGroup struct {
	UserApiGroup user.ApiGroup
	IdeaApiGroup idea.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
