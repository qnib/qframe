package qcollector_docker_logs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/zpatrick/go-config"

	"regexp"
	"github.com/qframe/types/health"
	"github.com/qframe/types/messages"
	"github.com/qframe/types/docker-events"
	"github.com/qframe/types/plugin"
	"github.com/qframe/types/qchannel"
)

const (
	version = "0.4.2"
	pluginTyp = "collector"
	pluginPkg = "docker-logs"
	dockerAPI = "v1.29"
)

var (
	ctx = context.Background()
)

type Plugin struct {
	*qtypes_plugin.Plugin
	cli *client.Client
	info types.Info
	sMap map[string]ContainerSupervisor
	TimeRegex   *regexp.Regexp
}

func (p *Plugin) StartSupervisor(ce events.Message, cnt types.ContainerJSON, info types.Info) {
	s := ContainerSupervisor{
		Plugin: *p,
		Action: ce.Action,
		CntID: ce.Actor.ID,
		CntName: ce.Actor.Attributes["name"],
		Info: &info,
		Container: &cnt,
		Com: make(chan interface{}),
		cli: p.cli,
		qChan: p.QChan,
	}
	s.TimeRegex = regexp.MustCompile(p.CfgStringOr("time-regex", `2\d{3}.*`))
	if p.CfgBoolOr("disable-reparse-logs", false) {
		s.TailRunning = p.CfgStringOr("tail-logs-since", "1m")
	}
	p.sMap[ce.Actor.ID] = s
	go s.Run()
}

func (p *Plugin) StartSupervisorCE(ce qtypes_docker_events.ContainerEvent) {
	p.StartSupervisor(ce.Event, ce.Container, p.info)
}


func New(qChan qtypes_qchannel.QChan, cfg *config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes_plugin.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg,  name, version),
		sMap: map[string]ContainerSupervisor{},
	}
	return p, err
}

func (p *Plugin) SubscribeRunning() {
	cnts, err := p.cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		p.Log("error", fmt.Sprintf("Failed to list containers: %s", err.Error()))
	} else {
		logLabel := p.CfgStringOr("enable-container-label", "org.qnib.qframe.enable-log")
		for _, cnt := range cnts {
			cjson, err := p.cli.ContainerInspect(ctx, cnt.ID)
			if err != nil {
				continue
			}
			event := events.Message{
				Type:   "container",
				Action: "running",
				Actor: events.Actor{
					ID: cnt.ID,
					Attributes: map[string]string{"name": strings.Trim(cnt.Names[0],"/")},
				},
			}
			// Skip those with the label:
			logCnt := false
			for label, _ := range cjson.Config.Labels {
				if label == logLabel {
					p.Log("info", fmt.Sprintf("Subscribing to logs of '%s' as label '%s' is set", cnt.Names, logLabel))
					logCnt = true
					break

				}
			}
			if ! logCnt {
				p.Log("info", fmt.Sprintf("Skip subscribing to logs of '%s' as label '%s' is set", cnt.Names, logLabel))
				b := qtypes_messages.NewTimedBase(p.Name, time.Unix(cnt.Created, 0))
				de := qtypes_docker_events.NewDockerEvent(b, event)
				ce := qtypes_docker_events.NewContainerEvent(de, cjson)
				h := qtypes_health.NewHealthBeat(b, "routine.logSkip", ce.Container.ID[:12], "start")
				p.Log("info", "Send logSkip-HealthBeat for "+h.Actor)
				p.QChan.SendData(h)
				continue
			}
			if cjson.HostConfig.LogConfig.Type != "json-file" {
				b := qtypes_messages.NewTimedBase(p.Name, time.Unix(cnt.Created, 0))
				h := qtypes_health.NewHealthBeat(b, "routine.logWrongType", cjson.ID[:12], "start")
				p.QChan.SendData(h)
				continue
			}

			b := qtypes_messages.NewTimedBase(p.Name, time.Unix(cnt.Created, 0))
			de := qtypes_docker_events.NewDockerEvent(b, event)
			ce := qtypes_docker_events.NewContainerEvent(de, cjson)
			h := qtypes_health.NewHealthBeat(b, "routine.log", ce.Container.ID[:12], "start")
			p.QChan.SendData(h)
			p.StartSupervisorCE(ce)
		}
	}
}

func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start v%s", p.Version))

	var err error
	dockerHost := p.CfgStringOr("docker-host", "unix:///var/run/docker.sock")
	p.cli, err = client.NewClient(dockerHost, dockerAPI, nil, nil)
	if err != nil {
		p.Log("error", fmt.Sprintf("Could not connect docker/docker/client to '%s': %v", dockerHost, err))
		return
	}
	p.info, err = p.cli.Info(ctx)
	if err != nil {
		p.Log("error", fmt.Sprintf("Error during Info(): %v >err> %s", p.info, err))
		return
	} else {
		p.Log("info", fmt.Sprintf("Connected to '%s' / v'%s' (SWARM: %s)", p.info.Name, p.info.ServerVersion, p.info.Swarm.LocalNodeState))
	}
	// need to start listener for all containers
	skipRunning := p.CfgBoolOr("skip-running", false)
	if ! skipRunning {
		p.Log("info", fmt.Sprintf("Start listeners for already running containers: %d", p.info.ContainersRunning))
		p.SubscribeRunning()
	}
	inputs := p.GetInputs()
	srcSuccess := p.CfgBoolOr("source-success", true)
	dc := p.QChan.Data.Join()
	for {
		select {
		case msg := <-dc.Read:
			switch msg.(type) {
			case qtypes_docker_events.ContainerEvent:
				ce := msg.(qtypes_docker_events.ContainerEvent)
				if len(inputs) != 0 && ! ce.InputsMatch(inputs) {
					continue
				}
				if ce.SourceSuccess != srcSuccess {
					continue
				}
				if ce.Event.Type == "container" && (strings.HasPrefix(ce.Event.Action, "exec_create") || strings.HasPrefix(ce.Event.Action, "exec_start")) {
					continue
				}
				p.Log("debug", fmt.Sprintf("Received: %s", ce.Message))
				switch ce.Event.Type {
				case "container":
					switch ce.Event.Action {
					case "start":
						p.sendHealthhbeat(ce, "start")
						p.StartSupervisorCE(ce)
					case "die":
						p.sendHealthhbeat(ce, "stop")
						p.sMap[ce.Event.Actor.ID].Com <- ce.Event.Action
					}
				}
			}
		}
	}
}

func (p *Plugin) sendHealthhbeat(ce qtypes_docker_events.ContainerEvent, action string) {
	skipLabel := p.CfgStringOr("skip-container-label", "org.qnib.qframe.skip-log")
	b := qtypes_messages.NewTimedBase(p.Name, ce.Time)
	// Skip those with the label:
	routineName := "routine.log"
	for label, val := range ce.Container.Config.Labels {
		if label == skipLabel && val == "true" {
			routineName = "routine.logSkip"
			break
		}
	}
	if ce.Container.HostConfig.LogConfig.Type != "json-file" {
		b := qtypes_messages.NewTimedBase(p.Name, ce.Time)
		h := qtypes_health.NewHealthBeat(b, "routine.logWrongType", ce.Container.ID[:12], "start")
		p.QChan.SendData(h)
		return
	}

	h := qtypes_health.NewHealthBeat(b, routineName, ce.Container.ID[:12], action)
	p.QChan.SendData(h)
	h = qtypes_health.NewHealthBeat(b, "vitals", p.Name, fmt.Sprintf("%s.%s", ce.Container.ID[:12], action))
	p.QChan.SendData(h)
}