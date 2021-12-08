package user

import "idea_server/service"

type ApiGroup struct {
	UserApi
	UserBaseApi
	UserFollowApi
}

var userService = service.ServiceGroupApp.UserServiceGroup.UserService
var userBaseService = service.ServiceGroupApp.UserServiceGroup.UserBaseService
var userFollowService = service.ServiceGroupApp.UserServiceGroup.UserFollowService
