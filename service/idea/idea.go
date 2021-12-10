package idea

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"idea_server/global"
	"idea_server/model/common/request"
	"idea_server/model/idea"
	ideaRes "idea_server/model/idea/response"
	"idea_server/service/user"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

var ideaCommentService = new(IdeaCommentService)
var ideaLikeService = new(IdeaLikeService)
var userService = new(user.UserService)

type RepRegexp struct {
	expr string
	repl string
}

var mdRegexps = []RepRegexp{
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

var escapeRegexps = []RepRegexp{
	// \，注意要 4 个，一定要在 \n 之前
	{
		expr: "\\\\",
		repl: "",
	},
	// 空格
	{
		expr: "\\s",
		repl: "",
	},
	// 回车
	{
		expr: "\\n",
		repl: "",
	},
	{
		expr: "\\r\\n",
		repl: "",
	},
	// 双引号
	{
		expr: "\"",
		repl: "\\\"",
	},
}

type IdeaService struct {
}

// 解决转义
func resolveEscape(text string) (string, error) {
	for _, value := range escapeRegexps {
		//fmt.Println("value", value)
		if r, err := regexp.Compile(value.expr); err != nil {
			fmt.Println("正则表示式编译错误", err)
			return "", errors.New("正则表达式编译错误")
		} else {
			text = r.ReplaceAllString(text, value.repl)
		}
	}
	return text, nil
}

func (e *IdeaService) GetClassification(text string) (typeId uint, err error) {
	text, err = resolveEscape(text)
	if err != nil {
		return 0, errors.New("获取分类失败：" + err.Error())
	}
	jsonData := "{\"text\": [\"" + text + "\"]}"
	//fmt.Println("jsonData", jsonData)
	res, _ := http.Post("http://hlj.vinf.top/model_classfication", "application/json", bytes.NewBuffer([]byte(jsonData)))
	defer res.Body.Close()
	if res.StatusCode == 200 {
		data, _ := ioutil.ReadAll(res.Body)
		var m []uint
		_ = json.Unmarshal(data, &m)
		//fmt.Println("m", m)
		return m[0], nil
	}
	return 0, errors.New("http 请求失败")
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

func (e *IdeaService) AuditContent(text string) (err error) {
	text, err = resolveEscape(text)
	if err != nil {
		return errors.New("内容审核失败：" + err.Error())
	}
	jsonData := "{\"text\": [\"" + text + "\"]}"
	res, err := http.Post("http://127.0.0.1:9998/audit_content", "application/json", bytes.NewBuffer([]byte(jsonData)))
	if res.StatusCode != 200 {
		// 服务出错
		return err
	}
	data, _ := ioutil.ReadAll(res.Body)
	if string(data) == "合规" {
		return nil
	}
	return errors.New(string(data))
}

func noticeAddIndexes(text string) {
	text, _ = resolveEscape(text)
	jsonData := "{\"text\": \"" + text + "\"}"
	http.Post("http://hlj.vinf.top/model_addsentence", "application/json", bytes.NewBuffer([]byte(jsonData)))
}

func noticeDelIndexes(text string) {
	text, _ = resolveEscape(text)
	jsonData := "{\"text\": \"" + text + "\"}"
	http.Post("http://hlj.vinf.top/model_delsentence", "application/json", bytes.NewBuffer([]byte(jsonData)))
}

func (e *IdeaService) CreateIdea(userId uint, content string) (bool, error) {
	if content == "" {
		return false, errors.New("内容不能为空")
	}
	life := getLife(0, 0, 0, 1)
	simple := e.SimpleContent(content)
	//fmt.Println("simple", simple)

	// 低俗辱骂过滤
	err := e.AuditContent(simple)
	if err != nil {
		return false, err
	}

	// 通知添加索引
	go noticeAddIndexes(simple)

	typeId, _ := e.GetClassification(simple) // 失败则默认 0，无类型
	i := idea.Idea{
		UserId:  userId,
		Simple:  simple,
		Content: content,
		Life:    life,
		Level:   1,
		TypeId:  typeId, // TODO 修改成固定值
	}

	result := global.IDEA_DB.Create(&i)
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
	//db := global.IDEA_DB.Debug().Model(&idea.Idea{}) // debug
	db := global.IDEA_DB.Model(&idea.Idea{})
	var ideas []idea.Idea
	ideaListResponses := make([]ideaRes.IdeaListResponse, 0, pageInfo.PageSize)

	// 添加一些条件
	//if ideaInfo.Content != "" {
	//	db = db.Where("content LIKE ?", "%"+ideaInfo.Content+"%")
	//}
	if ideaInfo.UserId != 0 {
		db = db.Where("user_id = ?", ideaInfo.UserId)
	}

	err = db.Where("level > 0 AND content != \"\"").Count(&total).Error

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
		if len(r) > 60 {
			v.Simple = string(r[:60])
		} else {
			v.Simple = string(r)
		}
		v.Life = math.Trunc(v.Life*1e4+0.5) * 1e-4
		response := ideaRes.IdeaListResponse{
			Idea:      v,
			IsLike:    ideaLikeService.IsLike(userId, v.ID),
			LikeCount: ideaLikeService.GetLikeCount(v.ID),
			TypeName:  e.GetIdeaTypeName(v.TypeId),
		}
		ideaListResponses = append(ideaListResponses, response)
	}
	return err, ideaListResponses, total, len(ideaListResponses)
}

func (e *IdeaService) GetSimilarIdeasByText(text string) (similarIdeas []ideaRes.SimilarIdea, err error) {
	jsonData := "{\"text\": \"" + text + "\"}"
	res, _ := http.Post("http://hlj.vinf.top/model_findsimilar", "application/json", bytes.NewBuffer([]byte(jsonData)))
	defer res.Body.Close()
	if res.StatusCode == 200 {
		data, _ := ioutil.ReadAll(res.Body)
		var result []ideaRes.SimilarModelResponse
		_ = json.Unmarshal(data, &result)
		fmt.Println("result", result)
		if len(result) > 3 {
			result = result[:3]
		}
		for _, v := range result {
			var i idea.Idea
			// 防止相似服务那未成功删除
			if !errors.Is(global.IDEA_DB.Where("level > 0").First(&i, v.IdeaId+1).Error, gorm.ErrRecordNotFound) {
				r := []rune(i.Simple)
				if len(r) > 60 {
					i.Simple = string(r[:60])
				} else {
					i.Simple = string(r)
				}
				similarIdeas = append(similarIdeas, ideaRes.SimilarIdea{
					Idea:       i,
					Similarity: v.Similarity,
					TypeName:   e.GetIdeaTypeName(i.TypeId),
				})
			}
		}
		return
	}
	similarIdeas = make([]ideaRes.SimilarIdea, 0, 3)
	for j := 0; j < 3; j++ {
		var i idea.Idea
		err = global.IDEA_DB.Find(&i, j+1).Error
		if err != nil {
			return make([]ideaRes.SimilarIdea, 0, 1), err
		}
		r := []rune(i.Simple)
		if len(r) > 40 {
			i.Simple = string(r[:60])
		} else {
			i.Simple = string(r)
		}
		similarIdeas = append(similarIdeas, ideaRes.SimilarIdea{
			Idea:       i,
			Similarity: float64(j) * 0.1,
		})
	}
	return
}

func (e *IdeaService) GetIdeaTypeName(typeId uint) string {
	if typeId == 0 {
		return ""
	}
	var t idea.IdeaType
	global.IDEA_DB.Find(&t, typeId)
	return t.Name
}

func (e *IdeaService) DeleteIdea(id uint) (err error) {
	var simple string
	if errors.Is(global.IDEA_DB.Debug().Model(&idea.Idea{}).Select("simple").Where("id = ?", id).First(&simple).Error, gorm.ErrRecordNotFound) {
		return gorm.ErrRecordNotFound
	}
	if simple != "" {
		err = global.IDEA_DB.Delete(&idea.Idea{}, id).Error
		if err != nil {
			return err
		}

		// 确保数据库删除成功后再做删除
		go noticeDelIndexes(simple)
	}
	return
}

func getLife(p, r, t, w float64) float64 {
	g := 1.194
	score := (p + 1.5*r + 20) / (math.Pow(t+2, g/w))
	return score
}

type LifeCronField struct {
	ID        uint
	UserId    uint
	Level     uint
	CreatedAt time.Time
}

type CountGroupByUser struct {
	UserId uint
	Cnt    uint
}

func LifeCronFunc() {
	//fmt.Println(getLife(0, 0, 0, 1))
	//fmt.Println(getLife(10, 5, 1, 8))
	//fmt.Println(getLife(0, 0, 20, 1))
	//fmt.Println(getLife(0, 0, 40, 1))

	var ids []LifeCronField
	if err := global.IDEA_DB.Model(&idea.Idea{}).Where("level > 0").Find(&ids).Error; err != nil {
		global.IDEA_LOG.Error("更新生命值定时任务——获取想法 id 列表失败！", zap.Error(err))
		return
	}
	now := time.Now()
	//fmt.Println("ids", ids)
	for _, v := range ids {
		// 点赞计算
		var likeIds []uint
		if err := global.IDEA_DB.Model(&idea.IdeaLike{}).Select("user_id").Where("idea_id = ?", v.ID).Find(&likeIds).Error; err != nil {
			global.IDEA_LOG.Error("更新生命值定时任务——获取 id "+strconv.Itoa(int(v.ID))+" 想法点赞列表失败！", zap.Error(err))
			continue
		}
		var pSum uint
		for _, v2 := range likeIds {
			weight, err := userService.GetUserWeight(v2)
			if err != nil {
				global.IDEA_LOG.Error("更新生命值定时任务——获取 id "+strconv.Itoa(int(v.ID))+" 想法点赞用户权值失败！", zap.Error(err))
				continue
			}
			pSum += weight
		}
		p := float64(pSum)

		// 评论计算 (Group)
		var comments []CountGroupByUser
		if err := global.IDEA_DB.Model(&idea.IdeaComment{}).Select("user_id, COUNT(*) as cnt").Where("idea_id = ?", v.ID).Group("user_id").Find(&comments).Error; err != nil {
			global.IDEA_LOG.Error("更新生命值定时任务——获取 id "+strconv.Itoa(int(v.ID))+" 想法评论列表失败！", zap.Error(err))
			continue
		}
		//fmt.Println("comments", comments)
		var rSum uint
		for _, v2 := range comments {
			weight, err := userService.GetUserWeight(v2.UserId)
			if err != nil {
				global.IDEA_LOG.Error("更新生命值定时任务——获取 id "+strconv.Itoa(int(v.ID))+" 想法评论用户权值失败！", zap.Error(err))
				continue
			}
			rSum += v2.Cnt * weight
		}
		r := float64(rSum)

		// 距离发帖的时间
		t := now.Sub(v.CreatedAt).Minutes() / 60
		w, err := userService.GetUserWeight(v.UserId)
		if err != nil {
			global.IDEA_LOG.Error("更新生命值定时任务——获取 id "+strconv.Itoa(int(v.ID))+" 想法发布者权值失败！", zap.Error(err))
			continue
		}
		//fmt.Println("p", p)
		//fmt.Println("r", r)
		//fmt.Println("t", t)
		//fmt.Println("w", w)
		life := getLife(p, r, t, float64(w))
		level := v.Level
		// TODO level 规则
		if life > getLife(15, 5, 1, float64(w)) {
			level = 2
		} else if life < getLife(0, 0, 20, float64(w)) && t < 1*24 {
			level = 0
		}
		//fmt.Println("life", life)
		//fmt.Println("level", level)
		if err := global.IDEA_DB.Model(&idea.Idea{}).Where("id = ?", v.ID).Update("life", life).Update("level", level).Error; err != nil {
			global.IDEA_LOG.Error("更新生命值定时任务——更新生命值失败！", zap.Error(err))
		}
	}
	global.IDEA_LOG.Info("更新生命值定时任务——成功！")
	fmt.Println("更新生命值定时任务——成功！")
}

func (e *IdeaService) Convert() {
	var count int64
	global.IDEA_DB.Model(&idea.Idea{}).Count(&count)
	for j := 1015; j <= int(count); j++ {
		// 一定要每次初始化一个新对象
		var i idea.Idea
		err := global.IDEA_DB.First(&i, j).Error
		simple := e.SimpleContent(i.Content)
		typeId, err := e.GetClassification(simple)
		if err != nil {
			fmt.Println("错误转换 id", j)
			os.Exit(1)
		}
		err = global.IDEA_DB.Model(&idea.Idea{}).Where("id = ?", j).Updates(map[string]interface{}{"simple": simple, "life": 0, "level": 1, "type_id": typeId}).Error
		if err != nil {
			fmt.Println("错误转换 id", j)
			os.Exit(1)
		}
		//fmt.Println("成功转换 id", j)
	}
	fmt.Println("转换成功")
}
