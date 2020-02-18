package events

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/Strum355/log"
	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
)

var (
	scripts = make([]*interp.Interpreter, 0)
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

func LoadScripts(custom map[string]map[string]reflect.Value) {
	files, err := ioutil.ReadDir(".cache/Cards/scripts")
	if err != nil {
		log.WithError(err).Error("Could not load scripts")
		return
	}

	for _, file := range files {
		interpreter := interp.New(interp.Options{})
		interpreter.Use(stdlib.Symbols)

		interpreter.Use(custom)
		files, err := ioutil.ReadDir(".cache/Cards/scripts/" + file.Name())
		read, err := ioutil.ReadFile(".cache/Cards/scripts/" + file.Name() + "/" + files[0].Name())
		if err != nil {
			log.WithError(err).Error("Could not read script")
			return
		}
		_, err = interpreter.Eval(string(read))
		if err != nil {
			log.WithError(err).Error("Could not run script")
			return
		}

		scripts = append(scripts, interpreter)
	}
}

func PushEvent(event Event) Event {
	for _, script := range scripts {
		v, err := script.Eval("Handle" + string(event.EventID))
		if err != nil {
			continue
		}
		function := v.Interface().(func(Event) (Event, bool))
		eve, changed := function(event)
		fmt.Println("ran function")
		if changed {
			return eve
		}
	}
	return event
}
