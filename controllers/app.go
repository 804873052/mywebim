package controllers

import (
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"strings"
)

var langTypes []string //支持的语言

/*国际化*/
func init() {
	//初始化语言列表
	langTypes = strings.Split(beego.AppConfig.String("lang_types"), "|")

	//根据语言类型加载语言环境文件
	for _, lang := range langTypes {
		beego.Trace("Loading language: " + lang)
		if err := i18n.SetMessage(lang, "conf/"+"locale_"+lang+".ini"); err != nil {
			beego.Error("Fail to set message file: ", err)
			return
		}
	}
}

/*
	baseController	所有的controller都继承它 可以实现一些公共部分的功能
*/
type baseController struct {
	beego.Controller //具有接口存根实现的嵌入结构
	i18n.Locale      //用于处理数据和呈现模板时使用i18n
}

/*
	用于语言选项的 检测 和 设置
*/
func (this *baseController) Prepare() {
	this.Lang = "" //from i18n.Locale

	// 1.从 “Accept-Language” 获取语言信息
	al := this.Ctx.Request.Header.Get("Accept-Language")
	if len(al) > 4 {
		al = al[:5] //只要前面 5 个字母
		if i18n.IsExist(al) {
			this.Lang = al
		}
	}

	//2.默认语言是 英语
	if len(this.Lang) == 0 {
		this.Lang = "en-US"
	}

	//设置 模板级别 的语言选项
	this.Data["Lang"] = this.Lang

}

/* 处理欢迎页面，允许用户选择技术和用户名 */
type AppController struct {
	baseController // 嵌入以使用 baseController 中实现的方法
}

/*首页*/
func (this *AppController) Get() {
	this.TplName = "welcome.html"
}

/*进入 聊天室
	uname： 聊天名
	tech：	进入的方式 Long polling，WebSocket
重定向到，不同方式的 api
*/
func (this *AppController) Join() {
	// Get form value
	uname := this.GetString("uname")
	tech := this.GetString("tech")

	//Check value
	if len(uname) == 0 { //没设置 uname
		this.Redirect("/", 302) //重定向
		return
	}

	/*判断选择的聊天类型 并重定向到相应的 接口上去*/
	switch tech {
	case "longpolling": //长连接
		this.Redirect("/lp?uname="+uname, 302)
	case "websocket": //websocket
		this.Redirect("/ws?uname="+uname, 302)
	default:
		this.Redirect("/", 302)
	}

	//Usually put return after redirect
	return
}
