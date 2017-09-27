package main

import (
	"log"
	"os"
	"sync"
	"github.com/zpatrick/go-config"
	"github.com/codegangsta/cli"
	
	"github.com/qframe/handler-elasticsearch"
	"github.com/qframe/cache-inventory"
	"github.com/qframe/handler-influxdb"
	"github.com/qframe/filter-grok"
	"github.com/qframe/collector-docker-events"
	"github.com/qframe/collector-internal"
	"github.com/qframe/collector-docker-logs"
	"github.com/qframe/types/qchannel"
)

const (
	dockerHost = "unix:///var/run/docker.sock"
	dockerAPI = "v1.29"
)


func check_err(pname string, err error) {
	if err != nil {
		log.Printf("[EE] Failed to create %s plugin: %s", pname, err.Error())
		os.Exit(1)
	}
}

func Run(ctx *cli.Context) {
	// Create conf
	log.Printf("[II] Start Version: %s", ctx.App.Version)
	cfg := config.NewConfig([]config.Provider{})
	if _, err := os.Stat(ctx.String("config")); err == nil {
		log.Printf("[II] Use config file: %s", ctx.String("config"))
		cfg.Providers = append(cfg.Providers, config.NewYAMLFile(ctx.String("config")))
	} else {
		log.Printf("[II] No config file found")
	}
	cfg.Providers = append(cfg.Providers, config.NewCLI(ctx, false))
	qChan := qtypes_qchannel.NewQChan()
	qChan.Broadcast()
	// Start InfluxDB handler
	phi, err := qhandler_influxdb.New(qChan, *cfg, "influxdb")
	check_err(phi.Name, err)
	go phi.Run()
	// Start Elasticsearch handler
	phe, err := qhandler_elasticsearch.New(qChan, *cfg, "es_logstash")
	check_err(phe.Name, err)
	go phe.Run()
	// Inventory
	pfi, err := qcache_inventory.New(qChan, *cfg, "inventory")
	check_err(pfi.Name, err)
	go pfi.Run()
	// App Log filter
	pfg, err := qfilter_grok.New(qChan, *cfg, "app-log")
	check_err(pfg.Name, err)
	go pfg.Run()
	// Elasticsearch Log filter
	pfgEs, err := qfilter_grok.New(qChan, *cfg, "es-log")
	check_err(pfgEs.Name, err)
	go pfgEs.Run()
	// start docker-events
	pe, err := qcollector_docker_events.New(qChan, *cfg, "docker-events")
	check_err(pe.Name, err)
	go pe.Run()
	// Docker Logs collector
	pcdl, err := qcollector_docker_logs.New(qChan, *cfg, "docker-log")
	check_err(pcdl.Name, err)
	go pcdl.Run()
	// Internal metrics collector
	pci, err := qcollector_internal.New(qChan, *cfg, "internal")
	check_err(pci.Name, err)
	go pci.Run()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func main() {
	app := cli.NewApp()
	app.Name = "ETC event collector based on qframe, inspired by qcollect,logstash and fullerite"
	app.Usage = "qframe-events [options]"
	app.Version = "0.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "qframe.yml",
			Usage: "Config file, will overwrite flag default if present.",
		},
	}
	app.Action = Run
	app.Run(os.Args)
}
