package qtypes

import (
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types"

	"strings"
)
type ContainerEvent struct {
	Base
	Message   	string
	Container 	types.ContainerJSON
	Event 		events.Message
	EngineInfo  types.Info
}

func NewContainerEvent(base Base, cnt types.ContainerJSON, event events.Message, info types.Info) ContainerEvent {
	return ContainerEvent{
		Base: base,
		Container: cnt,
		Event: event,
		EngineInfo: info,
	}
}


func (ce *ContainerEvent) GetContainerName() string {
	if ce.Container.Name != "" {
		return strings.Trim(ce.Container.Name, "/")
	} else {
		return "<none>"
	}
}