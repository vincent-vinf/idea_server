package v1

import "idea_server/api/v1/user"

type ApiGroup struct {
	UserApiGroup user.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
