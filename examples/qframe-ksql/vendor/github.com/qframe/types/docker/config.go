package qtypes_docker

import (
	"github.com/docker/docker/api/types/swarm"
	"github.com/qframe/types/helper"

)

type ConfigEvent struct {
	DockerEvent
	Config 	swarm.Config
}

func NewConfigEvent(de DockerEvent, config swarm.Config) ConfigEvent {
	return ConfigEvent{
		DockerEvent: de,
		Config: config,
	}
}

func (ce *ConfigEvent) DetailsToFlatJSON() (map[string]interface{}) {
	res := ce.Base.ToFlatJSON()
	res["config_id"] = ce.Config.ID
	res["config_name"] = ce.Config.Spec.Name
	res["config_annotations_name"] = ce.Config.Spec.Annotations.Name
	sLab, err := qtypes_helper.PrefixFlatKV(ce.Config.Spec.Labels, res, "config_label")
	if err == nil {
		res = sLab
	}
	return res
}

func (ce *ConfigEvent) EventToFlatJSON() (map[string]interface{}) {
	res := ce.Base.ToFlatJSON()
	res["config_id"] = ce.Config.ID
	res["event_action"] = ce.Event.Action
	res["event_from"] = ce.Event.From
	actAtr, err := qtypes_helper.PrefixFlatKV(ce.Event.Actor.Attributes, res, "event_actor_attr")
	if err == nil {
		res = actAtr
	}
	return res
}