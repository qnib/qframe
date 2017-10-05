package qtypes_docker_events

import (
	"github.com/docker/docker/api/types/swarm"
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