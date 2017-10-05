package qtypes_docker_events

import (

	"github.com/qframe/types/messages"
	"github.com/docker/docker/api/types/events"
	"fmt"
	"github.com/qframe/types/helper"
	"github.com/docker/docker/api/types"
)

type DockerEvent struct {
	qtypes_messages.Base
	Message   	string
	Event 		events.Message
	Engine 		types.Info
}

func NewDockerEvent(base qtypes_messages.Base, event events.Message) DockerEvent {
	return DockerEvent{
		Base: base,
		Message: fmt.Sprintf("%s.%s", event.Type, event.Action),
		Event: event,
	}
}

func (de *DockerEvent) SetEngineInfo(e types.Info) {
	de.Engine = e
}

func (de *DockerEvent) EventToJSON() (res map[string]interface{}) {
	res = de.Base.ToJSON()
	res["message"] = de.Message
	res["event"] = de.Event
	return
}


func (de *DockerEvent) AddEngineFlatJSON(res map[string]interface{}) map[string]interface{} {
	res["engine_id"] = de.Engine.ID
	res["engine_name"] = de.Engine.Name
	eLab, err := qtypes_helper.PrefixFlatLabels(de.Engine.Labels, res, "engine_label")
	if err == nil {
		res = eLab
	}
	res["engine_arch"] = de.Engine.Architecture
	res["engine_kernel"] = de.Engine.KernelVersion
	res["engine_os"] = de.Engine.OperatingSystem
	res["swarm_node_address"] = de.Engine.Swarm.NodeAddr
	res["swarm_node_id"] = de.Engine.Swarm.NodeID
	res["swarm_cluster_id"] = de.Engine.Swarm.Cluster.ID
	return res
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
