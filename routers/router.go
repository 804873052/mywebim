package routers

import (
	"MyWebIM/controllers"
	"github.com/astaxie/beego"
)

func init() {

	// 首页
	beego.Router("/", &controllers.AppController{})

	//有一个 post 连接请求 确定需要跳转到那个 路由
	beego.Router("/join", &controllers.AppController{}, "post:Join")

	// long polling
	beego.Router("/lp", &controllers.LongPollingController{}, "get:Join")
	beego.Router("/lp/post", &controllers.LongPollingController{})               // 请求路径为：/lp/post 请求方式为：post
	beego.Router("/lp/fetch", &controllers.LongPollingController{}, "get:Fetch") // 获取 聊天信息。每3秒访问一次，看看有没信息

	//WebSocket
	beego.Router("/ws", &controllers.WebSocketController{})
	beego.Router("/ws/join", &controllers.WebSocketController{}, "get:Join")

	/*什么鬼o*/
	//	??
}
