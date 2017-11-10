package main


import (
	"log"
	"os"
	"sync"
	"github.com/zpatrick/go-config"

	"github.com/qframe/types/qchannel"
	"github.com/qframe/collector-docker-events"
	"github.com/qframe/collector-docker-logs"
)

const (
	dockerHost = "unix:///var/run/docker.sock"
	dockerAPI  = "v1.29"
)

func check_err(pname string, err error) {
	if err != nil {
		log.Printf("[EE] Failed to create %s plugin: %s", pname, err.Error())
		os.Exit(1)
	}
}

func main() {
	// Create conf
	myCfg := map[string]string{
		"log.level": "debug",
		"cache.inventory.inputs": "docker-events",
	}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(myCfg)})
	qChan := qtypes_qchannel.NewQChan()
	qChan.Broadcast()
	// start docker-events
	pe, err := qcollector_docker_events.New(qChan, cfg, "docker-events")
	check_err(pe.Name, err)
	go pe.Run()
	// start docker-logs
	pl, err := qcollector_docker_logs.New(qChan, cfg, "docker-logs")
	check_err(pl.Name, err)
	go pl.Run()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}