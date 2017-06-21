package main

import (
	"log"
	"os"
	"sync"
	"github.com/zpatrick/go-config"
	"github.com/codegangsta/cli"

	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-collector-docker-events/lib"
	"github.com/qnib/qframe-handler-influxdb/lib"
	"github.com/qnib/qframe-collector-internal/lib"
	"github.com/qnib/qframe-filter-inventory/lib"
	"github.com/qnib/qframe-filter-grok/lib"
	"github.com/qnib/qframe-filter-metrics/lib"
	"github.com/qnib/qframe-collector-docker-stats/lib"
	"github.com/qnib/qframe-filter-docker-stats/lib"
	"github.com/qnib/qframe-filter-statsq/lib"
	"github.com/qnib/qframe-collector-tcp/lib"
	"github.com/qnib/qframe-handler-elasticsearch/lib"
	"github.com/qnib/qframe-collector-docker-log/lib"
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
	qChan := qtypes.NewQChan()
	qChan.Broadcast()
	//////// Handlers
	// Start InfluxDB
	phi, err := qframe_handler_influxdb.New(qChan, cfg, "influxdb")
	check_err(phi.Name, err)
	go phi.Run()
	// Start Elasticsearch handler to push to logstash
	phe, err := qframe_handler_elasticsearch.New(qChan, cfg, "es_logstash")
	check_err(phe.Name, err)
	go phe.Run()
	// Start Elasticsearch handler to push to events
	phee, err := qframe_handler_elasticsearch.New(qChan, cfg, "es_events")
	check_err(phee.Name, err)
	go phee.Run()
	//////// Filters
	// GROK
	pfm, err := qframe_filter_grok.New(qChan, cfg, "opentsdb")
	check_err(pfm.Name, err)
	go pfm.Run()
	// StatsD
	pfs, err := qframe_filter_statsq.New(qChan, cfg, "statsq")
	check_err(pfs.Name, err)
	//go pfs.Run()
	// Container Stats
	pfcs, err := qframe_filter_docker_stats.New(qChan, cfg, "container-stats")
	check_err(pfcs.Name, err)
	go pfcs.Run()
	// Inventory
	pfi, err := qframe_filter_inventory.New(qChan, cfg, "inventory")
	check_err(pfi.Name, err)
	go pfi.Run()
	// Metrics
	pfmet, err := qframe_filter_metrics.New(qChan, cfg, "metrics")
	check_err(pfmet.Name, err)
	go pfmet.Run()
	// Event filter
	pfg, err := qframe_filter_grok.New(qChan, cfg, "app-event")
	check_err(pfg.Name, err)
	go pfg.Run()
	// App Log filter
	pfal, err := qframe_filter_grok.New(qChan, cfg, "app-log")
	check_err(pfal.Name, err)
	go pfal.Run()
	// Elasticsearch Log filter
	pfgEs, err := qframe_filter_grok.New(qChan, cfg, "es-log")
	check_err(pfgEs.Name, err)
	go pfgEs.Run()
	//////// Collectors
	// Internal metrics
	pci, err := qframe_collector_internal.New(qChan, cfg, "internal")
	check_err(pci.Name, err)
	go pci.Run()
	// start docker-events
	pe, err := qframe_collector_docker_events.New(qChan, cfg, "docker-events")
	check_err(pe.Name, err)
	go pe.Run()
	// start docker-stats
	pds, err := qframe_collector_docker_stats.New(qChan, cfg, "docker-stats")
	check_err(pds.Name, err)
	go pds.Run()
	// start docker-logs
	pdl, err := qframe_collector_docker_log.New(qChan, cfg, "docker-log")
	check_err(pdl.Name, err)
	go pdl.Run()
	// TCP
	pct, err := qframe_collector_tcp.New(qChan, cfg, "tcp")
	check_err(pct.Name, err)
	go pct.Run()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func main() {
	app := cli.NewApp()
	app.Name = "ETC collector based on qframe, inspired by qcollect,logstash and fullerite"
	app.Usage = "qframe [options]"
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

