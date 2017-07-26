package qframe_collector_internal

import (
	"time"
	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"

)

const (
	version = "0.1.0"
	pluginTyp = "collector"
	pluginPkg = "internal"
)

type Plugin struct {
	qtypes.Plugin
}

func New(qChan qtypes.QChan, cfg config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
	}
	return p, err
}


func (p *Plugin) Run() {
	tickSec := p.CfgIntOr("ticker-sec", 1)
	p.Log("info", "Start internal collector v" + version)
	ticker := time.NewTicker(time.Duration(tickSec)*time.Second).C
	ims := qtypes.NewIntMemoryStats(p.Name)
	for {
		select {
		case <-ticker:
			ims.SnapShot()
			for _, m := range ims.ToMetrics(p.Name) {
				go p.QChan.Data.Send(m)
			}
		}
	}

}
