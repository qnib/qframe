package qtypes


import (
	"github.com/docker/docker/api/types"
	dc "github.com/fsouza/go-dockerclient"
	"math"
)

// Inspired by https://github.com/elastic/beats/blob/master/metricbeat/module/docker/cpu/helper.go
type MemoryStats struct {
	Base
	Container   *types.Container
	Failcnt   	float64
	Limit     	float64
	MaxUsage  	float64
	TotalRss  	float64
	TotalRssP 	float64
	Usage     	float64
	UsageP 		float64
}

func NewMemoryStats(src Base, stats *dc.Stats) MemoryStats {
	return MemoryStats{
		Base:      src,
		Failcnt:   float64(stats.MemoryStats.Failcnt),
		Limit:     float64(stats.MemoryStats.Limit),
		MaxUsage:  float64(stats.MemoryStats.MaxUsage),
		TotalRss:  float64(stats.MemoryStats.Stats.TotalRss),
		TotalRssP: calcUsage(float64(stats.MemoryStats.Stats.TotalRss), float64(stats.MemoryStats.Limit)),
		Usage:     float64(stats.MemoryStats.Usage),
		UsageP:    calcUsage(float64(stats.MemoryStats.Usage), float64(stats.MemoryStats.Limit)),
	}
}

func (ms *MemoryStats) ToMetrics(src string) []Metric {
	dim := AssembleDefaultDimensions(ms.Container)
	return []Metric{
		ms.NewExtMetric(src, "memory.usage.percent", Gauge, ms.UsageP, dim, ms.Time, true),
		ms.NewExtMetric(src, "memory.total_rss.percent", Gauge, ms.TotalRssP, dim, ms.Time, true),
		ms.NewExtMetric(src, "memory.total_rss.bytes", Gauge, ms.TotalRss, dim, ms.Time, true),
		ms.NewExtMetric(src, "memory.usage.bytes", Gauge, ms.Usage, dim, ms.Time, true),
		ms.NewExtMetric(src, "memory.failcnt", Gauge, ms.Failcnt, dim, ms.Time, true),
		ms.NewExtMetric(src, "memory.limit.bytes", Gauge, ms.Limit, dim, ms.Time, true),
	}
}

func calcUsage(frac, all float64) float64 {
	v := float64(frac / all)
	if math.IsNaN(v) {
		v = 0.0
	}
	return v
}
