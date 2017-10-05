package qtypes_docker_events

import (

	"github.com/qframe/types/messages"
	"github.com/docker/docker/api/types/events"
	"fmt"
	"github.com/qframe/types/helper"
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

func (de *DockerEvent) EventToFlatJSON() (res map[string]interface{}) {
	res = de.Base.ToFlatJSON()
	res["msg_message"] = de.Message
	res["event_type"] = de.Event.Type
	res["event_action"] = de.Event.Action
	res["event_scope"] = de.Event.Scope
	res["event_container_id"] = de.Event.Actor.ID
	actAtr, err := qtypes_helper.PrefixFlatKV(de.Event.Actor.Attributes, res, "event_actor_attr")
	if err == nil {
		res = actAtr
	}
	return
}
