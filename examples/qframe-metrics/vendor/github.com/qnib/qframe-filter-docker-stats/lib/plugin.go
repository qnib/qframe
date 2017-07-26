package qframe_filter_docker_stats

import (
	"C"
	"fmt"
	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
)

const (
	version = "0.1.3"
	pluginTyp = "filter"
	pluginPkg = "docker-stats"
)

type Plugin struct {
	qtypes.Plugin
}

func New(qChan qtypes.QChan, cfg *config.Config, name string) (p Plugin, err error) {
	p = Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
	}
	return p, err
}

// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start docker-stats filter v%s", p.Version))
	dc := p.QChan.Data.Join()
	go p.DispatchMsgCount()
	for {
		select {
		case val := <- dc.Read:
			switch val.(type) {
			case qtypes.ContainerStats:
				qcs := val.(qtypes.ContainerStats)
				if p.StopProcessingCntStats(qcs, false) {
					continue
				}
				// Process ContainerStats and create send multiple qtypes.Metrics
				go p.GetCpuMetrics(qcs)
				go p.GetMemoryMetrics(qcs)
				go p.GetNetworkMetrics(qcs)
			}
		}
	}
}

func (p *Plugin) GetCpuMetrics(qcs qtypes.ContainerStats) {
	stat := qcs.GetCpuStats()
	for _, m := range stat.ToMetrics(p.Name) {
		p.QChan.Data.Send(m)
	}
}

func (p *Plugin) GetMemoryMetrics(qcs qtypes.ContainerStats) {
	stat := qcs.GetMemStats()
	for _, m := range stat.ToMetrics(p.Name) {
		p.QChan.Data.Send(m)
	}
}

func (p *Plugin) GetNetworkMetrics(qcs qtypes.ContainerStats) {
	stat := qcs.GetNetStats()
	for _, m := range stat.ToMetrics(p.Name) {
		p.QChan.Data.Send(m)
	}
	aggStats := qtypes.NewNetStats(qcs.Base, qcs.GetContainer())
	for iface, _ := range qcs.Stats.Networks {
		stats := qcs.GetNetPerIfaceStats(iface)
		for _, m := range stats.ToMetrics(p.Name) {
			p.QChan.Data.Send(m)
		}
		aggStats = qtypes.AggregateNetStats("total", aggStats, stats)
	}
	for _, m := range aggStats.ToMetrics(p.Name) {
		p.QChan.Data.Send(m)
	}
}
