package main

import (
	"log"
	"github.com/zpatrick/go-config"
	"github.com/qframe/collector-docker-events"
	"github.com/qframe/types/plugin"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/handler-elasticsearch"
)

func main() {
	qChan := qtypes_qchannel.NewQChan()
	qChan.Broadcast()
	cfgMap := map[string]string{
		"log.level": "info",
		"handler.es_logstash.inputs": "events",
	}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(cfgMap)})
	// Create Health Cache\
	b := qtypes_plugin.NewBase(qChan, cfg)
	p, err := qhandler_elasticsearch.New(b, "es_logstash")
	if err != nil {
		log.Fatalf("[EE] Failed to create logstash: %v", err)
	}
	go p.Run()
	// Create docker events collector
	pde, err := qcollector_docker_events.New(qChan, cfg, "events")
	if err != nil {
		log.Fatalf("[EE] Failed to create docker-events: %v", err)
	}
	go pde.Run()
	bg := p.QChan.Data.Join()
	for {
		val := <- bg.Read
		switch val.(type) {
		default:
			continue
		}
	}
}

