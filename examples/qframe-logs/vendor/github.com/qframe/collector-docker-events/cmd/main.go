package main


import (
	"os"
	"sync"
	"log"
	"github.com/zpatrick/go-config"

	"github.com/qframe/types/qchannel"
	"github.com/qframe/handler-log"
	"github.com/qframe/collector-docker-events"
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
	}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(myCfg)})
	qChan := qtypes_qchannel.NewQChan()
	qChan.Broadcast()
	// Start log handler
	phl, err := qhandler_log.New(qChan, cfg, "log")
	check_err(phl.Name, err)
	go phl.Run()
	// start docker-events
	pe, err := qcollector_docker_events.New(qChan, cfg, "docker-events")
	check_err(pe.Name, err)
	go pe.Run()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}