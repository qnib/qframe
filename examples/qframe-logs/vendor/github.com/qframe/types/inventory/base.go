package qtypes_inventory

import (
	"time"
	"encoding/json"
	"github.com/qframe/types/messages"
	"github.com/qframe/types/docker-events"
)

const nanoX = 1000000000

func SplitUnixNano(t int64) (sec, nano int64) {
	sec = t/nanoX
	nano = t - sec*nanoX
	return
}

type Base struct {
	qtypes_messages.Base
	Time 			time.Time
	TimeUnixNano	int64				`json:"time"`
	Subject			interface{} 		`json:"subject"`	// Subject of what is going on (e.g. container)
	Action			interface{}			`json:"action"`
	Object  		interface{}     	`json:"object"` 	// Passive object
	Tags 			map[string]string 	`json:"tags"` 		// Tags that should be applied to the action
}


func NewBaseFromJson(qb qtypes_messages.Base, str string) (b Base, err error) {
	b.Base = qb
	err = json.Unmarshal([]byte(str), &b)
	s,n := SplitUnixNano(b.TimeUnixNano)
	b.Time = time.Unix(s, n)
	return
}

func NewBase(b qtypes_messages.Base, subject,action,object interface{}, tags map[string]string) (Base, error) {
	invBase, err := NewEmptyBase(b)
	invBase.Subject = subject
	invBase.Action = action
	invBase.Object = object
	invBase.Tags = tags
	return invBase, err
}

func NewEmptyBase(b qtypes_messages.Base) (Base, error) {
	var err error
	invBase := Base{
		Base: b,
		Time: b.Time,
		TimeUnixNano: b.Time.UnixNano(),
	}
	return invBase, err
}

func NewBaseFromContainerEvent(ce qtypes_docker_events.ContainerEvent) (Base, error) {
	var err error
	b := Base{
		Base: ce.Base,
		Time: ce.Time,
		TimeUnixNano: ce.Time.UnixNano(),
	}
	switch ce.Event.Type {
	case "container":
		b.EnrichContainer(ce)
	}
	return b, err
}


func (b *Base) EnrichContainer(ce qtypes_docker_events.ContainerEvent) {
	// TODO: Has to change to ID
	//b.Subject = ce.EngineInfo.Name
	b.Action = ce.Event.Action
	b.Object = ce.Container.Name
}