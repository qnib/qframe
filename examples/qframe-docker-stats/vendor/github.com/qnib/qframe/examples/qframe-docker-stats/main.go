package main

import (
	"github.com/codegangsta/cli"
	"github.com/zpatrick/go-config"
	"log"
	"os"
	"sync"


	"github.com/qframe/handler-influxdb"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/cache-inventory"
	"github.com/qframe/collector-docker-stats"
	"github.com/qframe/collector-internal"
	"github.com/qframe/collector-docker-events"
)

const (
	dockerHost = "unix:///var/run/docker.sock"
	dockerAPI  = "v1.30"
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
	//////// Handlers
	// Start InfluxDB
	phi, err := qhandler_influxdb.New(qChan, cfg, "influxdb")
	check_err(phi.Name, err)
	go phi.Run()
	/*
	//////// Filters
	// GROK
	pfm, err := qfilter_grok.New(qChan, cfg, "opentsdb")
	check_err(pfm.Name, err)
	go pfm.Run()
	// StatsQ
	pfs, err := qcache_statsq.New(qChan, cfg, "statsq")
	check_err(pfs.Name, err)
	go pfs.Run()
	*/
	// Container Stats
	pfcs, err := qcollector_docker_stats.New(qChan, cfg, "container-stats")
	check_err(pfcs.Name, err)
	go pfcs.Run()
	// Inventory
	pfi, err := qcache_inventory.New(qChan, cfg, "inventory")
	check_err(pfi.Name, err)
	go pfi.Run()
	/*
	// Metrics
	pfmet, err := qfilter_metrics.New(qChan, cfg, "metrics")
	check_err(pfmet.Name, err)
	go pfmet.Run()
	*/
	//////// Collectors
	// Internal metrics
	pci, err := qcollector_internal.New(qChan, cfg, "internal")
	check_err(pci.Name, err)
	go pci.Run()
	// start docker-events
	pe, err := qcollector_docker_events.New(qChan, cfg, "docker-events")
	check_err(pe.Name, err)
	go pe.Run()
	/*
	// TCP
	pct, err := qcollector_tcp.New(qChan, cfg, "tcp")
	check_err(pct.Name, err)
	go pct.Run()
	*/
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func main() {
	app := cli.NewApp()
	app.Name = "ETC event collector based on qframe, inspired by qcollect,logstash and fullerite"
	app.Usage = "qframe-metrics [options]"
	app.Version = "0.0.1"
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
