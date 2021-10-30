package service

import "idea_server/service/user"

type ServiceGroup struct {
	UserServiceGroup user.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
