package user

import (
	"idea_server/global"
	"idea_server/model/user"
	userRes "idea_server/model/user/response"
)

var userFollowService = new(UserFollowService)

type UserService struct {
}

func (e *UserService) GetUserInfo(ids []int, userId uint) (infos []userRes.UserInfoResponse, err error) {
	var users []user.User
	err = global.IDEA_DB.Omit("email", "passwd", "weight").Find(&users, ids).Error
	for i, _ := range users {
		infos = append(infos, userRes.UserInfoResponse{
			User:     users[i],
			IsFollow: userFollowService.IsFollow(users[i].ID, userId),
		})
	}
	return
}
