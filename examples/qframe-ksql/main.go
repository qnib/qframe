package main

import (
	"log"
	"os"
	"sync"
	"github.com/zpatrick/go-config"
	"github.com/codegangsta/cli"

	"github.com/qframe/collector-docker-events"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/handler-kafka"
)

const (
	dockerHost = "unix:///var/run/docker.sock"
	dockerAPI = "v1.31"
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
	// Start kafka handler
	phi, err := qhandler_kafka.New(qChan, cfg, "kafka")
	check_err(phi.Name, err)
	go phi.Run()
	// start docker-events
	pe, err := qcollector_docker_events.New(qChan, cfg, "docker-events")
	check_err(pe.Name, err)
	go pe.Run()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func main() {
	app := cli.NewApp()
	app.Name = "ETC event collector based on qframe, inspired by qcollect,logstash and fullerite"
	app.Usage = "qframe-ksql [options]"
	app.Version = "0.1.2"
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
