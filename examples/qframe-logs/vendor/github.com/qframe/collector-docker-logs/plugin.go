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
	"reflect"
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
	defer delete(p.sMap, ce.Actor.ID)
	s.Run()
}

func (p *Plugin) StartSupervisorCE(ce qtypes_docker_events.ContainerEvent) {
	go p.StartSupervisor(ce.Event, ce.Container, p.info)
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
		logEnv := p.CfgStringOr("enable-log-env", "LOG_CAPTURE_ENABLED")
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
			// Skip those with the Env:
			logCnt, err := SkipContainer(&cjson, logEnv)
			if err != nil {
				p.Log("debug", err.Error())
			}
			if ! logCnt {
				p.Log("info", fmt.Sprintf("Skip subscribing to logs of '%s' as environment variable '%s' was not found", cnt.Names, logEnv))
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
	logEnv := p.CfgStringOr("enable-log-env", "LOG_CAPTURE_ENABLED")
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
	dc := p.QChan.Data.Join()
	for {
		select {
		case msg := <-dc.Read:
			switch msg.(type) {
			case qtypes_docker_events.ContainerEvent:
				ce := msg.(qtypes_docker_events.ContainerEvent)
				if ce.Event.Type == "container" && strings.HasPrefix(ce.Event.Action, "exec_") {
					continue
				}
				p.Log("trace", fmt.Sprintf("Receied ContainerEvent: %s.%s", ce.Event.Type, ce.Event.Action))
				skipCnt, err := SkipContainer(&ce.Container, logEnv)
				if err != nil {
					p.Log("debug", err.Error())
				}
				switch ce.Event.Type {
				case "container":
					switch ce.Event.Action {
					case "start":
						if ! skipCnt {
							p.sendHealthhbeat("routine.logSkip", ce, "start")
							continue
						}
						p.sendHealthhbeat("routine.log", ce, "start")
						p.StartSupervisorCE(ce)
					case "die":
						if ! skipCnt {
							p.sendHealthhbeat("routine.logSkip", ce, "stop")
							continue
						}
						p.sendHealthhbeat("routine.log", ce, "stop")
						p.sMap[ce.Event.Actor.ID].Com <- ce.Event.Action
					}
				}
			default:
				p.Log("trace", fmt.Sprintf("Dunno what to do with type: %s", reflect.TypeOf(msg)))
			}
		}
	}
}


func (p *Plugin) sendHealthBeats(hbs []qtypes_health.HealthBeat) (err error) {
	for _, h := range hbs {
		p.QChan.SendData(h)
	}
	return
}
func (p *Plugin) sendHealthhbeat(rName string, ce qtypes_docker_events.ContainerEvent, action string) {
	if ce.Container.HostConfig.LogConfig.Type != "json-file" {
		b := qtypes_messages.NewTimedBase(p.Name, ce.Time)
		h := qtypes_health.NewHealthBeat(b, "routine.logWrongType", ce.Container.ID[:12], action)
		p.QChan.SendData(h)
		return
	}
	hbs := createHealthhbeats(p.Name, rName, action, ce)
	p.sendHealthBeats(hbs)
}

