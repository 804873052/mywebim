package controllers

import (
	"MyWebIM/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"net/http"
)

type WebSocketController struct {
	baseController
}

func (this *WebSocketController) Get() {
	uname := this.GetString("uname")
	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}

	this.TplName = "websocket.html"
	this.Data["IsWebSocket"] = true
	this.Data["UserName"] = uname
}

/*连接方法 处理WebSocket请求*/
func (this *WebSocketController) Join() {
	uname := this.GetString("uname")
	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}
	//从http请求升级到 WebSocket。
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Connot setup WebSocket connection:", err)
		return
	}
	//加入到 chat room中
	Join(uname, ws)
	defer Leave(uname)

	//消息接收循环
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		publish <- newEvent(models.EVENT_MESSAGE, uname, string(p))
	}
}

/*广播消息给WebSocket用户。*/
func broadcastWebSocket(event models.Event) {
	data, err := json.Marshal(event)
	if err != nil {
		beego.Error("Fail to marshal event: ", err)
		return
	}

	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		//立即发送事件到WebSocket用户。
		ws := sub.Value.(Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// 用户断开连接
				unsubscribe <- sub.Value.(Subscriber).Name
			}
		}
	}
}
