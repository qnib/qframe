package qtypes_container_stats

import (
	"strconv"
	"github.com/elastic/beats/libbeat/common"
	"github.com/docker/docker/api/types"
	dc "github.com/fsouza/go-dockerclient"
	"github.com/qframe/types/messages"
	//"github.com/qframe/types/metrics"
)

// Inspired by https://github.com/elastic/beats/blob/master/metricbeat/module/docker/cpu/helper.go
type CPUStats struct {
	qtypes_messages.Base
	Container                   *types.Container
	PerCpuUsage                 common.MapStr
	TotalUsage                  float64
	UsageInKernelmode           uint64
	UsageInKernelmodePercentage float64
	UsageInUsermode             uint64
	UsageInUsermodePercentage   float64
	SystemUsage                 uint64
	SystemUsagePercentage       float64
}

func NewCPUStats(src qtypes_messages.Base, stats *dc.Stats) CPUStats {
	return CPUStats{
		Base: src,
		PerCpuUsage: perCpuUsage(stats),
		TotalUsage: totalUsage(stats),
		UsageInKernelmode: stats.CPUStats.CPUUsage.UsageInKernelmode,
		UsageInKernelmodePercentage: usageInKernelmode(stats),
		UsageInUsermode: stats.CPUStats.CPUUsage.UsageInUsermode,
		UsageInUsermodePercentage: usageInUsermode(stats),
		SystemUsage: stats.CPUStats.SystemCPUUsage,
		SystemUsagePercentage: systemUsage(stats),
	}
}


/*func (cs *CPUStats) ToMetrics(src string) []qtypes_metrics.Metric {
	//dim := qtypes_helper.AssembleDefaultDimensions(cs.Container)
	return []qtypes_metrics.Metric{
		cs.NewExtMetric(src, "cpu.usage.kernel.percent", qtypes_metrics.Gauge, cs.UsageInKernelmodePercentage, dim, cs.Time, true),
		cs.NewExtMetric(src, "cpu.usage.user.percent", qtypes_metrics.Gauge, cs.UsageInUsermodePercentage, dim, cs.Time, true),
		cs.NewExtMetric(src, "cpu.system.usage.percent", qtypes_metrics.Gauge, cs.SystemUsagePercentage, dim, cs.Time, true),

	}
}*/


func perCpuUsage(stats *dc.Stats) common.MapStr {
	var output common.MapStr
	if len(stats.CPUStats.CPUUsage.PercpuUsage) == len(stats.PreCPUStats.CPUUsage.PercpuUsage) {
		output = common.MapStr{}
		for index := range stats.CPUStats.CPUUsage.PercpuUsage {
			cpu := common.MapStr{}
			cpu["pct"] = calculateLoad(stats.CPUStats.CPUUsage.PercpuUsage[index] - stats.PreCPUStats.CPUUsage.PercpuUsage[index])
			cpu["ticks"] = stats.CPUStats.CPUUsage.PercpuUsage[index]
			output[strconv.Itoa(index)] = cpu
		}
	}
	return output
}

func totalUsage(stats *dc.Stats) float64 {
	return calculateLoad(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
}

func usageInKernelmode(stats *dc.Stats) (val float64) {
	if stats.PreCPUStats.CPUUsage.UsageInKernelmode == 0 {
		return
	}
	return calculateLoad(stats.CPUStats.CPUUsage.UsageInKernelmode - stats.PreCPUStats.CPUUsage.UsageInKernelmode)
}

func usageInUsermode(stats *dc.Stats) (val float64) {
	if stats.PreCPUStats.CPUUsage.UsageInUsermode == 0 {
		return
	}
	return calculateLoad(stats.CPUStats.CPUUsage.UsageInUsermode - stats.PreCPUStats.CPUUsage.UsageInUsermode)
}

func systemUsage(stats *dc.Stats) (val float64) {
	if stats.PreCPUStats.SystemCPUUsage == 0 {
		return
	}
	return calculateLoad(stats.CPUStats.SystemCPUUsage - stats.PreCPUStats.SystemCPUUsage)
}

func calculateLoad(value uint64) float64 {
	return float64(value) / float64(1000000000)
}

// \beats

