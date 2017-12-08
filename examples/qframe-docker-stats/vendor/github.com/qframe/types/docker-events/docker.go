package qtypes_docker_events

import (

	"github.com/qframe/types/messages"
	"github.com/docker/docker/api/types/events"
	"fmt"
)

type DockerEvent struct {
	qtypes_messages.Base
	Message   	string
	Event 		events.Message
}

func NewDockerEvent(base qtypes_messages.Base, event events.Message) DockerEvent {
	return DockerEvent{
		Base: base,
		Message: fmt.Sprintf("%s.%s", event.Type, event.Action),
		Event: event,
	}
}

func (de *DockerEvent) EventToJSON() (res map[string]interface{}) {
	res = de.Base.ToJSON()
	res["message"] = de.Message
	res["event"] = de.Event
	return
}