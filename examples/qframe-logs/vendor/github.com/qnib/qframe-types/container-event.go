package qtypes

import (
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types"

)
type ContainerEvent struct {
	Base
	Message   string
	Container types.ContainerJSON
	Event events.Message
}

func NewContainerEvent(base Base, cnt types.ContainerJSON, event events.Message) ContainerEvent {
	return ContainerEvent{
		Base: base,
		Container: cnt,
		Event: event,
	}
}
