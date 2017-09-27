package qhandler_elasticsearch

import (
	"encoding/json"
	"log"
)

type Logstash struct {
	Settings interface{} `json:"settings"`
	Mappings interface{} `json:"mappings"`
}

func NewLogstash(shards, replicas int) Logstash {
	return Logstash{
		Settings: map[string]interface{}{
			"index.number_of_shards":   shards,
			"index.number_of_replicas": replicas,
		},
		Mappings: map[string]interface{}{
			"_default_": map[string]interface{}{
				"_source": map[string]interface{}{
					"enabled": true,
				},
				"_all": map[string]interface{}{
					"enabled": false,
				},
			},
		},
	}
}

func (l *Logstash) GetConfig() (interface{}, error) {
	marshalled, err := json.Marshal(l)
	if err != nil {
		log.Printf("[WW] Failed to marshall indexCfg: %v >> %v", l, err)
		return nil, err
	}
	return marshalled, nil
}
