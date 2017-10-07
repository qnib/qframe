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

func NewDockerEvent(base qtypes_messages.Base, eInfo types.Info, event events.Message) DockerEvent {
	return DockerEvent{
		Base: base,
		Message: fmt.Sprintf("%s.%s", event.Type, event.Action),
		Event: event,
		Engine: eInfo,
	}
}

func (de *DockerEvent) EventToJSON() (res map[string]interface{}) {
	res = de.Base.ToJSON()
	res["message"] = de.Message
	res["event"] = de.Event
	return
}


func (de *DockerEvent) EngineFlatJSON() (res map[string]interface{}) {
	res = de.Base.ToFlatJSON()
	res["engine_id"] = de.Engine.ID
	res["engine_name"] = de.Engine.Name
	eLab, err := qtypes_helper.PrefixFlatLabels(de.Engine.Labels, res, "engine_label")
	if err == nil {
		res = eLab
	}
	res["engine_arch"] = de.Engine.Architecture
	res["engine_kernel"] = de.Engine.KernelVersion
	res["engine_os"] = de.Engine.OperatingSystem
	return res
}

func (de *DockerEvent) SwarmNodeToFlatJSON() (res map[string]interface{}) {
	res = de.Base.ToFlatJSON()
	res["engine_id"] = de.Engine.ID
	res["swarm_cluster_id"] = de.Engine.Swarm.Cluster.ID
	res["swarm_node_address"] = de.Engine.Swarm.NodeAddr
	res["swarm_node_id"] = de.Engine.Swarm.NodeID
	res["swarm_node_state"] = de.Engine.Swarm.LocalNodeState
	res["swarm_node_error"] = de.Engine.Swarm.Error
	return
}

func (de *DockerEvent) EventToFlatJSON() (res map[string]interface{}) {
	res = de.Base.ToFlatJSON()
	res["engine_id"] = de.Engine.ID
	res["container_id"] = de.Event.Actor.ID
	res["msg_message"] = de.Message
	res["event_type"] = de.Event.Type
	res["event_action"] = de.Event.Action
	res["event_scope"] = de.Event.Scope
	actAtr, err := qtypes_helper.PrefixFlatKV(de.Event.Actor.Attributes, res, "event_actor_attr")
	if err == nil {
		res = actAtr
	}
	return
}
