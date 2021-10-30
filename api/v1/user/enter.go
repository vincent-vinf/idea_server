package user

import "idea_server/service"

type ApiGroup struct {
	UserApi
}

var userBaseService = service.ServiceGroupApp.UserServiceGroup.UserService
