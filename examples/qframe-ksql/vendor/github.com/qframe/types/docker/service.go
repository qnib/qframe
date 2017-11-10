package qtypes_docker

import (
	"github.com/docker/docker/api/types/swarm"
	"github.com/qframe/types/helper"
	"fmt"
)

type ServiceEvent struct {
	DockerEvent
	Service 	swarm.Service
}

func NewServiceEvent(de DockerEvent, srv swarm.Service) ServiceEvent {
	return ServiceEvent{
		DockerEvent: de,
		Service: srv,
	}
}

func (se *ServiceEvent) ServiceToJSON() (map[string]interface{}) {
	res := se.Base.ToJSON()
	res["message"] = se.Message
	res["service"] = se.Service
	return res
}

func (se *ServiceEvent) ServiceToFlatJSON() (map[string]interface{}) {
	res := se.Base.ToFlatJSON()
	res["swarm_cluster_id"] = se.Engine.Swarm.Cluster.ID
	res["service_id"] = se.Service.ID
	res["service_name"] = se.Service.Spec.Name
	res["service_version"] = fmt.Sprintf("%d", se.Service.Version.Index)
	sLab, err := qtypes_helper.PrefixFlatKV(se.Service.Spec.Labels, res, "service_label")
	if err == nil {
		res = sLab
	}
	return res
}

func (se *ServiceEvent) ServiceEventToFlatJSON() (map[string]interface{}) {
	res := se.Base.ToFlatJSON()
	res["service_id"] = se.Service.ID
	res["event_action"] = se.Event.Action
	actAtr, err := qtypes_helper.PrefixFlatKV(se.Event.Actor.Attributes, res, "event_actor_attr")
	if err == nil {
		res = actAtr
	}
	return res
}