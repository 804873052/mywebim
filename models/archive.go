package models

import "container/list"

type EventType int

/*事件*/
const (
	EVENT_JOIN = iota
	EVENT_LEAVE
	EVENT_MESSAGE
)

type Event struct {
	Type      EventType //join leave message
	User      string
	Timestamp int // Unix timestamp (secs) 时间戳 秒
	Content   string
}

/*档案*/
const archivesSize = 20

// Event archives.
var archive = list.New() //list对象

// NewArchive saves new event to archive list
func NewArchive(event Event) {
	if archive.Len() >= archivesSize { // 大于等于 20
		archive.Remove(archive.Front()) //移除第一个
	}
	archive.PushBack(event) //插入
}

//
func GetEvents(lastReceived int) []Event {
	events := make([]Event, 0, archive.Len())
	for event := archive.Front(); event != nil; event = event.Next() {
		e := event.Value.(Event)
		if e.Timestamp > int(lastReceived) { //插入在时间戳比自己大的地方 这个时间戳后面的信息吗？？
			events = append(events, e)
		}
	}
	return events
}
