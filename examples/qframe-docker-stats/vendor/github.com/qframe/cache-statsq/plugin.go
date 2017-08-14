package qcache_statsq

import (
	"fmt"
	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
	"github.com/qnib/statsq/lib"
)

const (
	version   = "0.1.3"
	pluginTyp = qtypes.CACHE
	pluginPkg = "statsq"
)

type Plugin struct {
	qtypes.Plugin
	Statsq statsq.StatsQ
}


func New(qChan qtypes.QChan, cfg *config.Config, name string) (Plugin, error) {
	p := qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version)
	sdName := fmt.Sprintf("%s.%s", pluginTyp, name)
	sd := statsq.NewNamedStatsQ(sdName, cfg, p.QChan)
	return Plugin{Plugin: p,Statsq: sd}, nil
}

// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start plugin v%s", p.Version))
	dc := p.QChan.Data.Join()
	go p.Statsq.LoopChannel()
	for {
		select {
		case val := <-dc.Read:
			switch val.(type) {
			case qtypes.Message:
				msg := val.(qtypes.Message)
				if p.StopProcessingMessage(msg, false) {
					continue
				}
				p.Log("debug", fmt.Sprintf("Received Message: %s %v", msg.Message))
				p.Statsq.ParseLine(msg.Message)
			case *qtypes.StatsdPacket:
				sd := val.(*qtypes.StatsdPacket)
				p.Log("trace", fmt.Sprintf("Received StatsdPacket: %s %v", sd.Bucket, sd.ValFlt))
				p.Statsq.HandlerStatsdPacket(sd)
			// TODO: Add qframe/type messages
			// TODO: Add qframe/type StatsdPacket
			}
		}
	}
}
