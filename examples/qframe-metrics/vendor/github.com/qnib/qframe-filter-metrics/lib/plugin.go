package qframe_filter_metrics

import (
	"fmt"
	"time"
	"strings"
	"strconv"
	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"

)

const (
	version   = "0.2.2"
	pluginTyp = "filter"
	pluginPkg = "metric"
)

var (
	containerStates = map[string]float64{
		"create": 1.0,
		"start": 2.0,
		"healthy": 3.0,
		"unhealthy": 4.0,
		"kill": -1.0,
		"die": -2.0,
		"stop": -3.0,
		"destroy": -4.0,
	}
)

type Plugin struct {
	qtypes.Plugin
}

func New(qChan qtypes.QChan, cfg *config.Config, name string) (p Plugin, err error) {
	p = Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
	}
	return
}

// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start plugin v%s", p.Version))
	dc := p.QChan.Data.Join()
	for {
		select {
		case val := <-dc.Read:
			switch val.(type) {
			case qtypes.Message:
				msg := val.(qtypes.Message)
				if p.StopProcessingMessage(msg, false) {
					continue
				}
				name, nok := msg.KV["name"]
				tval, tok := msg.KV["time"]
				value, vok := msg.KV["value"]
				if nok && tok && vok {
					mval, _ := strconv.ParseFloat(value, 64)
					tint, _ := strconv.Atoi(tval)
					dims := qtypes.AssembleJSONDefaultDimensions(&msg.Container)
					dims["source"] = msg.GetLastSource()
					met := qtypes.NewExt(p.Name, name, qtypes.Gauge, mval, dims, time.Unix(int64(tint), 0), true)
					tags, tagok := msg.KV["tags"]
					if tagok {
						for _, item := range strings.Split(tags, ",") {
							dim := strings.Split(item, "=")
							if len(dim) == 2 {
								met.Dimensions[dim[0]] = dim[1]
							}
						}
					}
					p.Log("trace", "send metric")
					p.QChan.Data.Send(met)
				}
			case qtypes.ContainerEvent:
				ce := val.(qtypes.ContainerEvent)
				if p.StopProcessingCntEvent(ce, false) {
					continue
				}
				switch ce.Event.Type {
				case "container":
					p.handleContainerEvent(ce)
				}
			}
		}
	}
}

func (p *Plugin) handleContainerEvent(ce qtypes.ContainerEvent) {
	if strings.HasPrefix(ce.Event.Action, "exec_") {
		p.MsgCount["execEvent"]++
		return
	}
	action := ce.Event.Action
	if strings.HasPrefix(ce.Event.Action, "health_status") {
		action = strings.Split(ce.Event.Action, ":")[1]
	}
	action = strings.Trim(action, " ")
	dims := qtypes.AssembleJSONDefaultDimensions(&ce.Container)
	dims["state-type"] = "container"
	p.Log("info", fmt.Sprintf("Action '%s' by %v", action, dims))
	if mval, ok := containerStates[action]; ok {
		met := qtypes.NewExt(p.Name, "state", qtypes.Gauge, mval, dims, ce.Time, false)
		p.QChan.Data.Send(met)
	} else {
		p.Log("warn", fmt.Sprintf("Could not fetch '%s' from containerState", action))
		met := qtypes.NewExt(p.Name, "state", qtypes.Gauge, 0.0, dims, ce.Time, false)
		p.QChan.Data.Send(met)
	}
}