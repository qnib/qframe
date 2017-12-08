package qcollector_docker_stats

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"github.com/zpatrick/go-config"
	"github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"github.com/qframe/types/docker-events"
	"github.com/qframe/types/messages"
	"github.com/qframe/types/health"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/types/plugin"
	"github.com/qframe/types/metrics"
)

const (
	version = "0.2.1"
	pluginTyp = "collector"
	pluginPkg = "docker-stats"
	dockerAPI = "v1.29"

)

var (
	ctx = context.Background()
)

// struct to keep info and channels to goroutine
// -> get heartbeats so that we know it's still alive
// -> allow for gracefully shutdown of the supervisor
type ContainerSupervisor struct {
	CntID 	string 			 // ContainerID
	CntName string			 // sanatized name of container
	Container docker.Container
	Com 	chan interface{} // Channel to communicate with goroutine
	cli 	*docker.Client
	qChan 	qtypes_qchannel.QChan
}

func SplitLabels(labels []string) map[string]string {
	res := map[string]string{}
	for _, label := range labels {
		tupel := strings.Split(label, "=")
		res[tupel[0]] = tupel[1]
	}
	return res
}

func (cs ContainerSupervisor) Run() {
	log.Printf("[II] Start listener for: '%s' [%s]", cs.CntName, cs.CntID)
	// TODO: That is not realy straight forward...
	filter := map[string][]string{
		"id": []string{cs.CntID},
	}
	df := docker.ListContainersOptions{
		Filters: filter,
	}
	info, _ := cs.cli.Info()
	engineLabels := SplitLabels(info.Labels)
	cnt, _ := cs.cli.ListContainers(df)
	if len(cnt) != 1 {
		log.Printf("[EE] Could not found excatly one container with id '%s'", cs.CntID)
		return
	}
	errChannel := make(chan error, 1)
	statsChannel := make(chan *docker.Stats)

	opts := docker.StatsOptions{
		ID:     cs.CntID,
		Stats:  statsChannel,
		Stream: true,
	}

	go func() {
		errChannel <- cs.cli.Stats(opts)
	}()

	for {
		select {
		case msg := <-cs.Com:
			switch msg {
			case "died":
				log.Printf("[DD] Container [%s]->'%s' died -> BYE!", cs.CntID, cs.CntName)
				return
			default:
				log.Printf("[DD] Container [%s]->'%s' got message from cs.Com: %v\n", cs.CntID, cs.CntName, msg)
			}
		case stats, ok := <-statsChannel:
			if !ok {
				err := errors.New(fmt.Sprintf("Bad response getting stats for container: %s", cs.CntID))
				log.Println(err.Error())
				return
			}
			qs := qtypes_metrics.NewContainerStats("docker-stats", stats, cnt[0])
			for k, v := range engineLabels {
				qs.Container.Labels[k] = v
			}
			cs.qChan.SendData(qs)
		}
	}
}

type Plugin struct {
	*qtypes_plugin.Plugin
	cli *docker.Client
	mobyClient *client.Client
	sMap map[string]ContainerSupervisor
}

func New(qChan qtypes_qchannel.QChan, cfg *config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes_plugin.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
		sMap: map[string]ContainerSupervisor{},
	}
	return p, err
}

func (p *Plugin) Run() {
	var err error
	dockerHost := p.CfgStringOr("docker-host", "unix:///var/run/docker.sock")
	p.mobyClient, err = client.NewClient(dockerHost, dockerAPI, nil, nil)
	if err != nil {
		p.Log("error", fmt.Sprintf("Could not connect docker/docker/client to '%s': %v", dockerHost, err))
		return
	}
	// Filter start/stop event of a container
	p.cli, err = docker.NewClient(dockerHost)
	if err != nil {
		p.Log("error", fmt.Sprintf("Could not connect fsouza/go-dockerclient to '%s': %v", dockerHost, err))
		return
	}
	info, err := p.cli.Info()
	if err != nil {
		p.Log("error", fmt.Sprintf("Error during Info(): %v >err> %s", info, err))
		return
	} else {
		p.Log("info", fmt.Sprintf("Connected to '%s' / v'%s' (SWARM: %s)", info.Name, info.ServerVersion, info.Swarm.LocalNodeState))
	}
	// List of current containers
	p.Log("info", fmt.Sprintf("Currently running containers: %d", info.ContainersRunning))
	// Dispatch Msg Count
	//go p.DispatchMsgCount()
	// Start listener for each container
	cnts, _ := p.cli.ListContainers(docker.ListContainersOptions{})
	for _,cnt := range cnts {
		event := events.Message{
			Type:   "container",
			Action: "running",
			Actor: events.Actor{
				ID: cnt.ID,
				Attributes: map[string]string{"name": strings.Trim(cnt.Names[0],"/")},
			},
		}
		cjson, err := p.mobyClient.ContainerInspect(ctx, cnt.ID)
		if err != nil {
			continue
		}
		b := qtypes_messages.NewTimedBase(p.Name, time.Unix(cnt.Created, 0))
		de := qtypes_docker_events.NewDockerEvent(b, event)
		ce := qtypes_docker_events.NewContainerEvent(de, cjson)
		h := qtypes_health.NewHealthBeat(b, "routine.stats", ce.Container.ID[:13], "start")
		p.Log("info", "Send routine.stats HealthBeat for "+h.Actor)
		p.QChan.SendData(h)
		p.StartSupervisor(cnt.ID, strings.TrimPrefix(cnt.Names[0], "/"))
	}
	dc,_,_ := p.JoinChannels()
	p.MsgCount["execEvent"] = 0
	for {
		select {
		case msg := <-dc.Read:
			switch msg.(type) {
			case qtypes_docker_events.ContainerEvent:
				ce := msg.(qtypes_docker_events.ContainerEvent)
				if ce.StopProcessing(p.Plugin, false) {
					continue
				}
				if ce.Event.Type == "container" && strings.HasPrefix(ce.Event.Action, "exec_") {
					p.MsgCount["execEvent"]++
					continue
				}
				switch ce.Event.Type {
				case "container":
					switch ce.Event.Action {
					case "start":
						b := qtypes_messages.NewTimedBase(p.Name, ce.Time)
						h := qtypes_health.NewHealthBeat(b, "vitals", p.Name, fmt.Sprintf("%s.%s", ce.Container.ID[:13], ce.Event.Action))
						p.QChan.SendData(h)
						hb := qtypes_health.NewHealthBeat(b, "routine.stats", ce.Container.ID[:13], "start")
						p.QChan.SendData(hb)
						p.StartSupervisorCe(ce)
					case "die":
						b := qtypes_messages.NewTimedBase(p.Name, ce.Time)
						h := qtypes_health.NewHealthBeat(b, "vitals", p.Name, fmt.Sprintf("%s.%s", ce.Container.ID[:13], ce.Event.Action))
						p.QChan.SendData(h)
						hb := qtypes_health.NewHealthBeat(b, "routine.stats", ce.Container.ID[:13], "stop")
						p.QChan.SendData(hb)
						p.sMap[ce.Event.Actor.ID].Com <- ce.Event.Action
					}

				}
			}
		}
	}
}


func (p *Plugin) StartSupervisor(CntID, CntName string) {
	s := ContainerSupervisor{
		CntID: CntID,
		CntName: CntName,
		Com: make(chan interface{}),
		cli: p.cli,
		qChan: p.QChan,
	}
	p.sMap[CntID] = s
	go s.Run()
}

func (p *Plugin) StartSupervisorCe(ce qtypes_docker_events.ContainerEvent) {
	p.StartSupervisor(ce.Event.Actor.ID, ce.Event.Actor.Attributes["name"])
}


