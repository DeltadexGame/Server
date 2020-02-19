package events

import (
	"github.com/Strum355/log"
)

var (
	Handlers map[EventID][]func(map[string]interface{}) (map[string]interface{}, bool) = make(map[EventID][]func(map[string]interface{}) (map[string]interface{}, bool))
)

type EventID string

const (
	MonsterDieEvent    EventID = "MonsterDieEvent"
	MonsterDamageEvent EventID = "MonsterDamageEvent"
)

type Event struct {
	EventID   EventID
	EventInfo map[string]interface{}
}

func RegisterHandler(id EventID, f func(map[string]interface{}) (map[string]interface{}, bool)) {
	_, ok := Handlers[id]
	if !ok {
		Handlers[id] = make([]func(map[string]interface{}) (map[string]interface{}, bool), 0)
	}

	Handlers[id] = append(Handlers[id], f)
}

func PushEvent(event Event) Event {
	log.WithFields(log.Fields{
		"event": event.EventID,
		"info":  event.EventInfo,
	}).Info("Event sent")
	for _, fun := range Handlers[event.EventID] {
		eve, changed := fun(event.EventInfo)
		if changed {
			event.EventInfo = eve
			return event
		}
	}
	return event
}
