package qtypes


import (
	"github.com/docker/docker/api/types"
)


// Inspired by https://github.com/elastic/beats/blob/master/metricbeat/module/docker/net/helper.go
type NetStats struct {
	Base
	Container     *types.Container
	NameInterface string
	RxBytes       float64
	RxDropped     float64
	RxErrors      float64
	RxPackets     float64
	TxBytes       float64
	TxDropped     float64
	TxErrors      float64
	TxPackets     float64
}

func AggregateNetStats(iface string, s1, s2 NetStats) NetStats {
	// Adds the NetStats counters and names the interface 'iface'
	return NetStats{
		Base: s1.Base,
		Container: s1.Container,
		NameInterface: iface,
		RxBytes:    s1.RxBytes + s2.RxBytes,
		RxDropped:  s1.RxDropped + s2.RxDropped,
		RxErrors:   s1.RxErrors + s2.RxErrors,
		RxPackets:  s1.RxPackets + s2.RxPackets,
		TxBytes:    s1.TxBytes + s2.TxBytes,
		TxDropped:  s1.TxDropped + s2.TxDropped,
		TxErrors:   s1.TxErrors + s2.TxErrors,
		TxPackets:  s1.TxPackets + s2.TxPackets,
	}
}

func NewNetStats(base Base, cnt *types.Container) NetStats {
	return NetStats{
		Base: base,
		Container: cnt,
		NameInterface: "none",
		RxBytes:    0.0,
		RxDropped:  0.0,
		RxErrors:   0.0,
		RxPackets:  0.0,
		TxBytes:    0.0,
		TxDropped:  0.0,
		TxErrors:   0.0,
		TxPackets:  0.0,
	}
}


func (ns *NetStats) ToMetrics(src string) []Metric {
	dim := AssembleDefaultDimensions(ns.Container)
	iface := "global"
	if ns.NameInterface != "" {
		iface = ns.NameInterface
		dim["network_iface"] = iface
	}
	return []Metric{
		ns.NewExtMetric(src, "network."+iface+".rx.bytes", Gauge, ns.RxBytes, dim, ns.Time, true),
		ns.NewExtMetric(src, "network."+iface+".rx.dropped", Gauge, ns.RxDropped, dim, ns.Time, true),
		ns.NewExtMetric(src, "network."+iface+".rx.errors", Gauge, ns.RxErrors, dim, ns.Time, true),
		ns.NewExtMetric(src, "network."+iface+".rx.packets", Gauge, ns.RxPackets, dim, ns.Time, true),
		ns.NewExtMetric(src, "network."+iface+".tx.bytes", Gauge, ns.TxBytes, dim, ns.Time, true),
		ns.NewExtMetric(src, "network."+iface+".tx.dropped", Gauge, ns.TxDropped, dim, ns.Time, true),
		ns.NewExtMetric(src, "network."+iface+".tx.errors", Gauge, ns.TxErrors, dim, ns.Time, true),
		ns.NewExtMetric(src, "network."+iface+".tx.packets", Gauge, ns.TxPackets, dim, ns.Time, true),
	}
}


