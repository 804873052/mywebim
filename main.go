package main

import (
	_ "MyWebIM/routers"
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

const (
	APP_VER = "1.0"
)

func main() {
	beego.Info(beego.BConfig.AppName, APP_VER)

	/*注册模板函数 国际化*/
	beego.AddFuncMap("i18n", i18n.Tr)

	// 开启热升级
	beego.BConfig.Listen.Graceful = true

	//测试 161 服务器

	beego.Run()
}
