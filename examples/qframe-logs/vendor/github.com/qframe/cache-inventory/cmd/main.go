package main


import (
	"fmt"
	"log"
	"os"
	"time"
	"github.com/zpatrick/go-config"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/cache-inventory"
	"github.com/qframe/collector-docker-events"

	"github.com/qframe/types/plugin"
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
		"cache.inventory.ticker-ms": "500",
	}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(myCfg)})
	qChan := qtypes_qchannel.NewQChan()
	qChan.Broadcast()
	// Start inventory cache
	phi, err := qcache_inventory.New(qChan, cfg, "inventory")
	check_err(phi.Name, err)
	go phi.Run()
	time.Sleep(time.Second)
	// start docker-events
	pe, err := qcollector_docker_events.New(qChan, cfg, "docker-events")
	check_err(pe.Name, err)
	go pe.Run()
	timeout := time.NewTimer(time.Duration(2000)*time.Millisecond).C
	if len(os.Args) != 2 {
		phi.Log("error", "Please provide name of container to search for `go run main.go some_name`")
		os.Exit(1)
	}
	cntName := os.Args[1]
	phi.Log("info", fmt.Sprintf("Create query for container by Name '%s' : qcache_inventory.NewNameContainerRequest('q1', '%s')", cntName, cntName))
	req := qcache_inventory.NewNameContainerRequest("q1", cntName)
	qChan.Data.Send(req)
	select {
	case resp := <- req.Back:
		if resp.Error != nil {
			phi.Log("error", resp.Error.Error())
		} else {
			phi.Log("info", fmt.Sprintf("Got InventoryResponse: Container '%s' has ID-digest:%s", cntName, resp.Container.ID[:13]))
			secondQuery(phi.Plugin, resp)
		}
	case <- timeout:
		phi.Log("error", fmt.Sprintf("Experience timeout for Name %s.", cntName))
		os.Exit(1)
	}
}

func secondQuery(phi *qtypes_plugin.Plugin, resp qcache_inventory.Response) {
	// Second query
	ipAddr := ""
	for k,v := range resp.Container.NetworkSettings.Networks {
		ipAddr = v.IPAddress
		phi.Log("info", fmt.Sprintf("Use IP from first query response container (network:%s) to generate another query: qcache_inventory.NewIPContainerRequest('q2', '%s')", k, ipAddr))
		break
	}
	if ipAddr == "" {
		phi.Log("error", "Seems like the container does not have networks attached to it (maybe only bridge?)")
		os.Exit(1)
	}
	timeout := time.NewTimer(time.Duration(2000)*time.Millisecond).C
	req := qcache_inventory.NewIPContainerRequest("q2", ipAddr)
	phi.QChan.Data.Send(req)
	select {
	case resp := <- req.Back:
		if resp.Error != nil {
			phi.Log("error", resp.Error.Error())
		} else {
			phi.Log("info", fmt.Sprintf("Got InventoryResponse: Container w/ IP %s has Digest:%s", ipAddr, resp.Container.ID[:13]))
		}
	case <- timeout:
		phi.Log("error", fmt.Sprintf("Experience timeout for IP %s.", ipAddr))
	}
}