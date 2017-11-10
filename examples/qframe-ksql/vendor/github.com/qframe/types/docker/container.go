package qtypes_docker

import (
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/qframe/types/helper"
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

// ContainerToJSON create a nested JSON object.
func (ce *ContainerEvent) ContainerToJSON() (map[string]interface{}) {
	res := ce.Base.ToJSON()
	res["msg_message"] = ce.Message
	res["container"] = ce.Container
	return res
}

// ContainerToFlatJSON create a key/val JSON map, which can be consumed by KSQL.
func (ce *ContainerEvent) ContainerToFlatJSON() (res map[string]interface{}) {
	res = ce.Base.ToFlatJSON()
	res["msg_message"] = ce.Message
	res["engine_id"] = ce.Engine.ID
	res["container_id"] = ce.Container.ID
	res["container_image"] = ce.Container.Image
	res["container_name"] = ce.GetContainerName()
	res["container_created"] = ce.Container.Created
	res["container_args"] = strings.Join(ce.Container.Args, " ")
	res["container_cfg_image"] = ce.Container.Config.Image
	rLab, err := qtypes_helper.PrefixFlatKV(ce.Container.Config.Labels, res, "container_label")
	if err == nil {
		res = rLab
	}
	return res
}
