package main

import (
	"github.com/codegangsta/cli"
	"github.com/zpatrick/go-config"
	"log"
	"os"
	"sync"

	"github.com/qnib/qframe-collector-docker-events/lib"
	"github.com/qnib/qframe-collector-internal/lib"
	"github.com/qnib/qframe-collector-tcp/lib"
	"github.com/qnib/qframe-filter-grok/lib"
	"github.com/qnib/qframe-filter-inventory/lib"
	"github.com/qnib/qframe-handler-elasticsearch/lib"
	"github.com/qnib/qframe-handler-influxdb/lib"
	"github.com/qnib/qframe-types"
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
	qChan := qtypes.NewQChan()
	qChan.Broadcast()
	// Start InfluxDB handler
	phi, err := qframe_handler_influxdb.New(qChan, *cfg, "influxdb")
	check_err(phi.Name, err)
	go phi.Run()
	// Start Elasticsearch handler
	phe := qframe_handler_elasticsearch.NewElasticsearch(qChan, *cfg, "es_logstash")
	check_err(phi.Name, err)
	go phe.Run()
	// Inventory
	pfi := qframe_filter_inventory.New(qChan, *cfg, "inventory")
	go pfi.Run()
	// Event filter
	pfg, err := qframe_filter_grok.New(qChan, *cfg, "app-event")
	check_err(pfg.Name, err)
	go pfg.Run()
	// start docker-events
	pe, err := qframe_collector_docker_events.New(qChan, *cfg, "docker-events")
	check_err(pe.Name, err)
	go pe.Run()
	// TCP collector
	pct, err := qframe_collector_tcp.New(qChan, *cfg, "tcp")
	check_err(pct.Name, err)
	go pct.Run()
	// Internal metrics collector
	pci, err := qframe_collector_internal.New(qChan, *cfg, "internal")
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
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "gcollect.yml",
			Usage: "Config file, will overwrite flag default if present.",
		},
	}
	app.Action = Run
	app.Run(os.Args)
}
