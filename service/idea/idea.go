package idea

import (
	"fmt"
	"idea_server/global"
	"idea_server/model/common/request"
	"idea_server/model/idea"
	ideaRes "idea_server/model/idea/response"
	"regexp"
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
	// TODO life

	simple := e.SimpleContent(content)
	//fmt.Println("simple", simple)
	idea := idea.Idea{
		UserId:  userId,
		Simple:  simple,
		Content: content,
		Life:    0,
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
		IsLike: ideaLikeService.IsLike(userId, info.Uint()),
	}, nil
}

func (e *IdeaService) GetIdeaList(ideaInfo idea.Idea, pageInfo request.PageInfo, order string, desc bool, userId uint) (err error, list interface{}, total int64, num int) {
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	//db := global.IDEA_DB.Model(&idea.Idea{}).Omit("content")
	db := global.IDEA_DB.Model(&idea.Idea{})
	ideaList := make([]ideaRes.IdeaListResponse, 0, 1)

	// 添加一些条件
	//if ideaInfo.Content != "" {
	//	db = db.Where("content LIKE ?", "%"+ideaInfo.Content+"%")
	//}

	err = db.Count(&total).Error

	if err != nil {
		return err, ideaList, total, len(ideaList)
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
			err = db.Order("created_at desc").Find(&ideaList).Error
		}
	}
	for index := range ideaList {
		r := []rune(ideaList[index].Simple)
		if len(r) > 40 {
			ideaList[index].Simple = string(r[:40])
		} else {
			ideaList[index].Simple = string(r)
		}
		// 减少传输字节
		//ideaList[index].Content = ""
		ideaList[index].IsLike = ideaLikeService.IsLike(userId, ideaList[index].ID)
		ideaList[index].LikeCount = ideaLikeService.GetLikeCount(ideaList[index].ID)
	}
	return err, ideaList, total, len(ideaList)
}

func (e *IdeaService) GetSimilarIdeasByText(text string) (similarIdeas []ideaRes.SimilarIdea, err error) {
	similarIdeas = make([]ideaRes.SimilarIdea, 0, 5)
	for i := 0; i < 5; i ++ {
		var idea idea.Idea
		err = global.IDEA_DB.Find(&idea, i + 1).Error
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