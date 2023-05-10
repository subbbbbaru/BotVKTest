package events

type EventType string

const (
	EventMessageNew = "message_new"
)

// type GroupEvent struct {
// 	Type    EventType       `json:"type"`
// 	Object  json.RawMessage `json:"object"`
// 	GroupID int             `json:"group_id"`
// 	EventID string          `json:"event_id"`
// 	V       string          `json:"v"`
// 	Secret  string          `json:"secret"`
// }

// // MessageNewObject struct.
// type MessageNewObject struct {
// 	Message    object.MessagesMessage `json:"message"`
// 	ClientInfo object.ClientInfo      `json:"client_info"`
// }

// type FuncList struct {
// 	messageNew []func(context.Context, MessageNewObject)
// 	special    map[EventType][]func(context.Context, GroupEvent)
// 	eventsList []EventType

// 	goroutine bool
// }

// func NewFuncList() *FuncList {
// 	return &FuncList{
// 		special: make(map[EventType][]func(context.Context, GroupEvent)),
// 	}
// }

// // MessageNew handler.
// func (fl *FuncList) MessageNew(f func(context.Context, MessageNewObject)) {
// 	fl.messageNew = append(fl.messageNew, f)
// 	fl.eventsList = append(fl.eventsList, EventMessageNew)
// }
