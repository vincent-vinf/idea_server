package user

import (
	"fmt"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/model/common/request"
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

func (e *UserService) GetUserWeight(id uint) (weight uint, err error) {
	err = global.IDEA_DB.Model(&user.User{}).Where("id = ?", id).Select("weight").Find(&weight).Error
	return
}

func (e *UserService) Notice(pageInfo request.PageInfo, userId uint) (err error, list []userRes.NoticeResponse, num int) {
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)

	var results []userRes.NoticeField
	err = global.IDEA_DB.Raw("SELECT *\nFROM (\n\tSELECT id, 1 type, created_at FROM `like` WHERE user_id != ? AND idea_id in (SELECT id FROM ideas WHERE user_id = ? AND deleted_at IS NULL)\n\tUNION\n\tSELECT id, 2 type, created_at FROM `follow` WHERE followed_id = ?\n\tUNION\n\tSELECT id, 3 type, created_at FROM `comment` WHERE user_id != ? AND to_id != ? AND idea_id in (SELECT id FROM ideas WHERE user_id = ? AND deleted_at IS NULL)\n\tUNION\n\tSELECT id, 4 type, created_at FROM `comment` WHERE to_id = ?\n) AS a\nORDER BY created_at DESC\nLIMIT ?\nOFFSET ?", userId, userId, userId, userId, userId, userId, userId, limit, offset).Scan(&results).Error
	if err != nil {
		return err, make([]userRes.NoticeResponse, 0, 1), 0
	}
	list = make([]userRes.NoticeResponse, 0, len(results))
	for _, v := range results {
		switch v.Type {
		case 1: // 点赞你发的帖，取消点赞即 deleted_at 不为 NULL
			var data idea.IdeaLike
			err = global.IDEA_DB.Raw("SELECT * FROM `like` WHERE id = ?", v.ID).Scan(&data).Error
			if err == nil {
				list = append(list, userRes.NoticeResponse{
					NoticeField: v,
					Data:        data,
				})
			}
			break
		case 2: // 用户关注，与点赞同理
			var data user.UserFollow
			err = global.IDEA_DB.Raw("SELECT * FROM `follow` WHERE id = ?", v.ID).Scan(&data).Error
			if err == nil {
				list = append(list, userRes.NoticeResponse{
					NoticeField: v,
					Data:        data,
				})
			}
			break
		case 3, 4: // 3 回复你发的帖子的评论（作者回复自己的帖） 4 回复你的评论
			var data idea.IdeaComment
			err = global.IDEA_DB.Raw("SELECT * FROM `comment` WHERE id = ?", v.ID).Scan(&data).Error
			if err == nil {
				list = append(list, userRes.NoticeResponse{
					NoticeField: v,
					Data:        data,
				})
			}
			break
		}

	}
	return err, list, len(list)
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
			continue
		}
		pOstnum = float64(len(ideaIds))
		//fmt.Println("ideaIds", ideaIds)
		var pSum, rSum int64
		for _, v2 := range ideaIds {
			// 用户收到的点赞数
			if err := global.IDEA_DB.Model(&idea.IdeaLike{}).Where("idea_id = ?", v2).Count(&cnt).Error; err != nil {
				global.IDEA_LOG.Error("更新用户权值定时任务——获取 id "+strconv.Itoa(int(v))+"用户收到的点赞数失败！", zap.Error(err))
				continue
			}
			pSum += cnt
			//fmt.Println("cnt p", cnt)

			// 用户收到的评论数
			if err := global.IDEA_DB.Model(&idea.IdeaComment{}).Where("idea_id = ?", v2).Count(&cnt).Error; err != nil {
				global.IDEA_LOG.Error("更新用户权值定时任务——获取 id "+strconv.Itoa(int(v))+"用户收到的评论数失败！", zap.Error(err))
				continue
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
			continue
		}
		pSend = float64(cnt)
		//fmt.Println("pSend", pSend)

		// 用户给予的评论数
		if err := global.IDEA_DB.Model(&idea.IdeaComment{}).Where("user_id = ?", v).Count(&cnt).Error; err != nil {
			global.IDEA_LOG.Error("更新用户权值定时任务——获取 id "+strconv.Itoa(int(v))+"用户给予的评论数失败！", zap.Error(err))
			continue
		}
		rSend = float64(cnt)
		//fmt.Println("rSend", rSend)

		// 用户被关注的人数
		if err := global.IDEA_DB.Model(&user.UserFollow{}).Where("followed_id = ?", v).Count(&cnt).Error; err != nil {
			global.IDEA_LOG.Error("更新用户权值定时任务——获取 id "+strconv.Itoa(int(v))+"用户被关注的人数失败！", zap.Error(err))
			continue
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
