package events

import (
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
	MonsterAttackEvent EventID = "MonsterAttackEvent"
)

type Event struct {
	EventID   EventID
	EventInfo map[string]interface{}
}

func LoadScripts(custom map[string]map[string]reflect.Value) {
	files, err := ioutil.ReadDir("scripts")
	if err != nil {
		log.WithError(err).Error("Could not load scripts")
		return
	}

	for _, file := range files {
		interpreter := interp.New(interp.Options{})
		interpreter.Use(stdlib.Symbols)

		interpreter.Use(custom)
		read, err := ioutil.ReadFile("scripts/" + file.Name())
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

func PushEvent(event Event) {
	for _, script := range scripts {
		v, err := script.Eval("Handle" + string(event.EventID))
		if err != nil {
			log.WithError(err).Error("Couldn't push event")
			continue
		}
		function := v.Interface().(func(Event))
		function(event)
	}
}
