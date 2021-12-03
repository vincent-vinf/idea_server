package idea

import (
	"fmt"
	"idea_server/global"
	"idea_server/model/common/request"
	"idea_server/model/idea"
	"regexp"
)

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
	// 全局匹配换行
	//{
	//	expr: "\\n",
	//	repl: "",
	//},
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
	r := []rune(content)
	if len(r) > 50 {
		return string(r[:50])
	}
	return string(r)
}

func (e *IdeaService) CreateIdea(userId uint, content string) (bool, error) {
	// TODO life

	simple := e.SimpleContent(content)
	//fmt.Println("simple", simple)
	idea := idea.Idea{
		UserId: userId,
		Simple:     simple,
		Content:    content,
		Life:       0,
	}

	result := global.IDEA_DB.Create(&idea)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (e *IdeaService) GetIdeaInfo(info *request.GetById) (interface{}, error) {
	i := idea.Idea{}
	result := global.IDEA_DB.First(&i, info.Uint())
	if result.Error != nil {
		return nil, result.Error
	}
	return i, nil
}
