package qtypes_docker_events

import (
	"strings"

	"github.com/docker/docker/api/types"
)

type ContainerEvent struct {
	DockerEvent
	Container 	types.ContainerJSON
}

func NewContainerEvent(de DockerEvent, cnt types.ContainerJSON) ContainerEvent {
	return ContainerEvent{
		DockerEvent: de,
		Container: cnt,
	}
}

func (ce *ContainerEvent) GetContainerName() string {
	if ce.Container.Name != "" {
		return strings.Trim(ce.Container.Name, "/")
	} else {
		return "<none>"
	}
}

func (ce *ContainerEvent) ContainerToJSON() (map[string]interface{}) {
	res := ce.Base.ToJSON()
	res["message"] = ce.Message
	res["container"] = ce.Container
	return res
}

