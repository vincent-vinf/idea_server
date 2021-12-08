package idea

import (
	"fmt"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/model/common/request"
	"idea_server/model/idea"
	ideaRes "idea_server/model/idea/response"
	"math"
	"regexp"
	"time"
)

var ideaCommentService = new(IdeaCommentService)
var ideaLikeService = new(IdeaLikeService)

type mdRegexp struct {
	expr string
	repl string
}

var mdRegexps = []mdRegexp{
	// 全局匹配内粗体
	{
		expr: "(\\*\\*|__)(.*?)(\\*\\*|__)",
		repl: "${2}",
	},
	// 全局匹配图片
	{
		expr: "\\!\\[[\\s\\S]*?\\]\\([\\s\\S]*?\\)",
		repl: "",
	},
	// 全局匹配连接
	{
		expr: "\\[([\\s\\S]*?)\\]\\([\\s\\S]*?\\)",
		repl: "${1}",
	},
	// 全局匹配内 html 标签
	{
		expr: "<\\/?.+?\\/?>",
		repl: "",
	},
	// 全局匹配内联代码块
	{
		expr: "(\\*)(.*?)(\\*)",
		repl: "",
	},
	// 全局匹配内联代码块
	{
		expr: "`{1,2}[^`](.*?)`{1,2}",
		repl: "",
	},
	// 全局匹配代码块
	{
		expr: "```([\\s\\S]*?)```[\\s]*",
		repl: "",
	},
	// 全局匹配删除线
	{
		expr: "\\~\\~(.*?)\\~\\~",
		repl: "${1}",
	},
	// 全局匹配无序列表
	{
		expr: "[\\s]*[-\\*\\+]+(.*)",
		repl: "${1}",
	},
	// 全局匹配有序列表
	{
		expr: "[\\s]*[0-9]+\\.(.*)",
		repl: "${1}",
	},
	// 全局匹配标题
	{
		expr: "(#+)(.*)",
		repl: "${2}",
	},
	// 全局匹配摘要
	{
		expr: "(>+)(.*)",
		repl: "${2}",
	},
	// 全局匹配两次 \n
	{
		expr: "\\n\\n",
		repl: "\\n",
	},
	// 全局匹配换行
	//{
	//	expr: "\\r\\n",
	//	repl: "",
	//},
	//全局匹配换行
	{
		expr: "\\\\n",
		repl: "\n",
	},
	// 全局匹配空字符
	//{
	//	expr: "\\s",
	//	repl: "",
	//},
}

type IdeaService struct {
}

func (e *IdeaService) SimpleContent(content string) string {
	//part := []rune(content)[:50]
	for _, value := range mdRegexps {
		//fmt.Println("value", value)
		if r, err := regexp.Compile(value.expr); err != nil {
			fmt.Println("正则表示式编译错误", err)
			return ""
		} else {
			content = r.ReplaceAllString(content, value.repl)
		}
	}
	return content
}

func (e *IdeaService) CreateIdea(userId uint, content string) (bool, error) {
	life := getLife(0, 0, 0)
	simple := e.SimpleContent(content)
	//fmt.Println("simple", simple)
	idea := idea.Idea{
		UserId:  userId,
		Simple:  simple,
		Content: content,
		Life:    life,
		Level:   1,
	}

	result := global.IDEA_DB.Create(&idea)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (e *IdeaService) GetIdeaInfo(info *request.GetById, userId uint) (interface{}, error) {
	err, comments := ideaCommentService.GetComment(info.Uint())
	if err != nil {
		return nil, err
	}

	//i := idea.Idea{}
	//result := global.IDEA_DB.First(&i, info.Uint())
	//if result.Error != nil {
	//	return nil, result.Error
	//}
	return &ideaRes.IdeaInfoResponse{
		//Idea: i,
		Comments: comments,
		IsLike:   ideaLikeService.IsLike(userId, info.Uint()),
	}, nil
}

func (e *IdeaService) GetIdeaList(ideaInfo idea.Idea, pageInfo request.PageInfo, order string, desc bool, userId uint) (err error, list interface{}, total int64, num int) {
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	//db := global.IDEA_DB.Model(&idea.Idea{}).Omit("content")
	db := global.IDEA_DB.Debug().Model(&idea.Idea{})
	var ideas []idea.Idea
	ideaListResponses := make([]ideaRes.IdeaListResponse, 0, pageInfo.PageSize)

	// 添加一些条件
	//if ideaInfo.Content != "" {
	//	db = db.Where("content LIKE ?", "%"+ideaInfo.Content+"%")
	//}

	err = db.Where("level > 0").Count(&total).Error

	if err != nil {
		return err, ideaListResponses, total, len(ideaListResponses)
	} else {
		db = db.Limit(limit).Offset(offset)
		if order != "" {
			var OrderStr string
			// 设置有效排序key 防止sql注入
			// 感谢 Tom4t0 提交漏洞信息
			orderMap := make(map[string]bool, 2)
			orderMap["life"] = true
			orderMap["updated_at"] = true
			if orderMap[order] {
				if desc {
					OrderStr = order + " desc"
				} else {
					OrderStr = order
				}
			}

			db = db.Order(OrderStr)
			err = db.Error
		}
		if err == nil {
			err = db.Order("created_at desc").Find(&ideas).Error
		}
	}
	for _, v := range ideas {
		r := []rune(v.Simple)
		if len(r) > 40 {
			v.Simple = string(r[:40])
		} else {
			v.Simple = string(r)
		}
		v.Life = math.Trunc(v.Life*1e4+0.5) * 1e-4
		response := ideaRes.IdeaListResponse{
			Idea:      v,
			IsLike:    ideaLikeService.IsLike(userId, v.ID),
			LikeCount: ideaLikeService.GetLikeCount(v.ID),
		}
		ideaListResponses = append(ideaListResponses, response)
	}
	return err, ideaListResponses, total, len(ideaListResponses)
}

func (e *IdeaService) GetSimilarIdeasByText(text string) (similarIdeas []ideaRes.SimilarIdea, err error) {
	similarIdeas = make([]ideaRes.SimilarIdea, 0, 5)
	for i := 0; i < 5; i++ {
		var idea idea.Idea
		err = global.IDEA_DB.Find(&idea, i+1).Error
		if err != nil {
			return make([]ideaRes.SimilarIdea, 0, 1), err
		}
		similarIdeas = append(similarIdeas, ideaRes.SimilarIdea{
			Idea:       idea,
			Similarity: float64(i) * 0.1,
		})
	}
	return
}

func getLife(p, r, t float64) float64 {
	g := 1.194
	score := (p + 1.5*r + 20) / (math.Pow(t+2, g))
	return score
}

func (e *IdeaService) LifeCronFunc() {
	var ideas []idea.Idea
	if err := global.IDEA_DB.Where("level > 0").Find(&ideas).Error; err != nil {
		global.IDEA_LOG.Error("更新生命值定时任务——获取想法列表失败！", zap.Error(err))
	}
	now := time.Now()
	for _, v := range ideas {
		// 距离发帖的时间
		t := now.Sub(v.CreatedAt).Minutes() / 60
		// p 点赞数，r 评论数
		var p, r int64
		if err := global.IDEA_DB.Model(&idea.IdeaLike{}).Where("idea_id = ?", v.ID).Count(&p).Error; err != nil {
			global.IDEA_LOG.Error("更新生命值定时任务——统计想法点赞数失败！", zap.Error(err))
		}
		if err := global.IDEA_DB.Model(&idea.IdeaComment{}).Where("idea_id = ?", v.ID).Count(&r).Error; err != nil {
			global.IDEA_LOG.Error("更新生命值定时任务——统计想法评论数失败！", zap.Error(err))
		}
		//fmt.Println("t", t)
		//fmt.Println("p", p)
		//fmt.Println("r", r)
		//fmt.Println("life: " + strconv.FormatFloat(getLife(float64(p), float64(r), t), 'f', 2, 64))
		life := getLife(float64(p), float64(r), t)
		level := v.Level
		// TODO level 规则
		if life > 11 {
			level = 2
		} else if life < 0.5 && t > 1*24*60 {
			level = 0
		}
		if err := global.IDEA_DB.Model(&v).Update("life", life).Update("level", level).Error; err != nil {
			global.IDEA_LOG.Error("更新生命值定时任务——更新生命值失败！", zap.Error(err))
		}
	}
	global.IDEA_LOG.Info("更新生命值定时任务——成功！")
	fmt.Println("更新生命值定时任务——成功！")
}
