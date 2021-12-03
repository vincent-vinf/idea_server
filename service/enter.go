package service

import (
	"idea_server/service/idea"
	"idea_server/service/user"
)

type ServiceGroup struct {
	UserServiceGroup user.ServiceGroup
	IdeaServiceGroup idea.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
