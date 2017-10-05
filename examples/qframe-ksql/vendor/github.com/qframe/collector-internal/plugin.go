package qcollector_internal

import (
	"time"
	"github.com/zpatrick/go-config"
	"net/http"
	"net/http/pprof"

	"github.com/qnib/qframe-types"

	"fmt"
)

const (
	version = "0.1.3"
	pluginTyp = "collector"
	pluginPkg = "internal"
)

type Plugin struct {
	qtypes.Plugin
}

func New(qChan qtypes.QChan, cfg *config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
	}
	return p, err
}

func (p *Plugin) StartHttp(port string) {
	r := http.NewServeMux()
	// Register pprof handlers
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	p.Log("info", fmt.Sprintf("Start pprof http server on :%s", port))
	http.ListenAndServe(fmt.Sprintf(":%s",port), r)
}

func (p *Plugin) Run() {
	tickSec := p.CfgIntOr("ticker-sec", 1)
	p.Log("notice", "Start internal collector v" + version)
	ticker := time.NewTicker(time.Duration(tickSec)*time.Second).C
	pprofEnabled := p.CfgBoolOr("pperf.enabled", false)
	if pprofEnabled {
		pport := p.CfgStringOr("pperf.port", "8080")
		go p.StartHttp(pport)
	}
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
