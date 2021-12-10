package user

import (
	"fmt"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/model/idea"
	"idea_server/model/user"
	userRes "idea_server/model/user/response"
	"math"
	"strconv"
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

func getWeight(pGet, pSend, rGet, rSent, pOstnum, fOcusnum float64) uint {
	p := 2*math.Log10(pGet+1) + math.Log10(pSend+1)
	r := 2*math.Log10(rGet+1) + math.Log10(rSent+1)
	num := math.Log10(pOstnum+1) + math.Log10(fOcusnum+1)
	x := -(1 / (math.Log10(0.5*p+0.75*r+math.Log10(num+1)+1) + 0.001))
	w := 8 * math.Exp(x) / 3.81 * 8
	if int(w) == 0 {
		return 1
	} else if int(w) >= 8 {
		return 8
	}
	return uint(math.Ceil(w))
}

func UserWeightCronFunc() {
	//fmt.Println(getWeight(300000, 600000, 100000, 300000, 100000, 100000))
	//fmt.Println(getWeight(3000, 6000, 1000, 3000, 1000, 1000))
	//fmt.Println(getWeight(300, 600, 100, 300, 100, 100))
	//fmt.Println(getWeight(30, 60, 10, 30, 10, 10))
	//fmt.Println(getWeight(0, 0, 0, 0, 0, 0))

	var ids []uint
	if err := global.IDEA_DB.Model(&user.User{}).Select("id").Find(&ids).Error; err != nil {
		global.IDEA_LOG.Error("更新用户权值定时任务——获取用户 id 列表失败！", zap.Error(err))
		return
	}
	//fmt.Println("ids", ids)

	// 默认第一个用户不更新，thanks
	for _, v := range ids[1:] {
		var pGet, pSend, rGet, rSend, pOstnum, fOcusnum float64
		var cnt int64
		// 用户发布的想法数
		var ideaIds []uint
		if err := global.IDEA_DB.Model(&idea.Idea{}).Select("id").Where("user_id = ?", v).Find(&ideaIds).Error; err != nil {
			global.IDEA_LOG.Error("更新用户权值定时任务——获取 id "+strconv.Itoa(int(v))+" 用户发布想法失败！", zap.Error(err))
		}
		pOstnum = float64(len(ideaIds))
		//fmt.Println("ideaIds", ideaIds)
		var pSum, rSum int64
		for _, v2 := range ideaIds {
			// 用户收到的点赞数
			if err := global.IDEA_DB.Model(&idea.IdeaLike{}).Where("idea_id = ?", v2).Count(&cnt).Error; err != nil {
				global.IDEA_LOG.Error("更新用户权值定时任务——获取 id "+strconv.Itoa(int(v))+"用户收到的点赞数失败！", zap.Error(err))
			}
			pSum += cnt
			//fmt.Println("cnt p", cnt)

			// 用户收到的评论数
			if err := global.IDEA_DB.Model(&idea.IdeaComment{}).Where("idea_id = ?", v2).Count(&cnt).Error; err != nil {
				global.IDEA_LOG.Error("更新用户权值定时任务——获取 id "+strconv.Itoa(int(v))+"用户收到的评论数失败！", zap.Error(err))
			}
			rSum += cnt
			//fmt.Println("cnt r", cnt)
		}
		pGet = float64(pSum)
		rGet = float64(rSum)
		//fmt.Println("pGet", pGet)
		//fmt.Println("rGet", rGet)

		// 用户给予的点赞数
		if err := global.IDEA_DB.Model(&idea.IdeaLike{}).Where("user_id = ?", v).Count(&cnt).Error; err != nil {
			global.IDEA_LOG.Error("更新用户权值定时任务——获取 id "+strconv.Itoa(int(v))+"用户给予的点赞数失败！", zap.Error(err))
		}
		pSend = float64(cnt)
		//fmt.Println("pSend", pSend)

		// 用户给予的评论数
		if err := global.IDEA_DB.Model(&idea.IdeaComment{}).Where("user_id = ?", v).Count(&cnt).Error; err != nil {
			global.IDEA_LOG.Error("更新用户权值定时任务——获取 id "+strconv.Itoa(int(v))+"用户给予的评论数失败！", zap.Error(err))
		}
		rSend = float64(cnt)
		//fmt.Println("rSend", rSend)

		// 用户被关注的人数
		if err := global.IDEA_DB.Model(&user.UserFollow{}).Where("followed_id = ?", v).Count(&cnt).Error; err != nil {
			global.IDEA_LOG.Error("更新用户权值定时任务——获取 id "+strconv.Itoa(int(v))+"用户被关注的人数失败！", zap.Error(err))
		}
		fOcusnum = float64(cnt)
		//fmt.Println("fOcusnum", fOcusnum)

		// 更新权值
		if err := global.IDEA_DB.Model(&user.User{}).Where("id = ?", v).Update("weight", getWeight(pGet, pSend, rGet, rSend, pOstnum, fOcusnum)).Error; err != nil {
			global.IDEA_LOG.Error("更新用户权值定时任务——获取 id "+strconv.Itoa(int(v))+"用户更新权值失败！", zap.Error(err))
		}
	}
	global.IDEA_LOG.Info("更新用户权值定时任务——成功！")
	fmt.Println("更新用户权值定时任务——成功！")
}
