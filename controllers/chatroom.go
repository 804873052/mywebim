package controllers

import (
	"MyWebIM/models"
	"container/list"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"time"
)

type Subscription struct {
	Archive []models.Event      //档案中的所有事件
	New     <-chan models.Event //新事件
}

type Subscriber struct {
	Name string
	Conn *websocket.Conn //只适用于WebSocket用户;否则nil。
}

var (
	//新加入用户的 通道 入口
	subscribe = make(chan Subscriber, 10)
	//出口用户的通道
	unsubscribe = make(chan string, 10)
	//发送事件到这里发布它们。
	publish = make(chan models.Event, 10)
	//长轮候名单。
	waitingList = list.New()
	subscribers = list.New()
)

/*新建一个 事件*/
func newEvent(ep models.EventType, user, msg string) models.Event {
	return models.Event{ep, user, int(time.Now().Unix()), msg}
}

/*加入 用户 新通道 */
func Join(user string, ws *websocket.Conn) {
	subscribe <- Subscriber{Name: user, Conn: ws}
}

/* 用户 出口 */
func Leave(user string) {
	unsubscribe <- user
}

/*此函数处理所有传入的chan消息*/
func chatroom() {
	for {
		/*
			select 是 Go 中的一个控制结构，类似于用于通信的 switch 语句。每个 case 必须是一个通信操作，要么是发送要么是接收。
			select 随机执行一个可运行的 case。如果没有 case 可运行，它将阻塞，直到有 case 可运行。一个默认的子句应该总是可运行的。
		*/
		select {
		case sub := <-subscribe:
			if !isUserExist(subscribers, sub.Name) { //还未加入聊天室
				subscribers.PushBack(sub) //添加到list 末尾
				//发布一个连接事件。
				publish <- newEvent(models.EVENT_JOIN, sub.Name, "")
				beego.Info("New user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			} else {
				beego.Info("Old user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			}
		case event := <-publish:
			//通知名单。
			for ch := waitingList.Back(); ch != nil; ch = ch.Prev() {
				ch.Value.(chan bool) <- true
				waitingList.Remove(ch)
			}

			//广播到 websocket
			broadcastWebSocket(event)
			models.NewArchive(event)

			if event.Type == models.EVENT_MESSAGE {
				beego.Info("Message from", event.User, ";Content:", event.Content)
			}
		case unsbu := <-unsubscribe:
			for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Subscriber).Name == unsbu {
					subscribers.Remove(sub)
					//关闭 连接
					ws := sub.Value.(Subscriber).Conn
					if ws != nil {
						ws.Close()
						beego.Error("WebSocket closed:", unsbu)
					}
					publish <- newEvent(models.EVENT_LEAVE, unsbu, "")
					break
				}
			}

		}
	}
}

func init() {
	go chatroom()
}

/*判断用户是否在 用户列表中
通过用户 name 确定
*/
func isUserExist(subscribers *list.List, user string) bool {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).Name == user {
			return true
		}
	}
	return false
}
