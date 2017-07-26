package qtypes

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/docker/docker/api/types"
)

type ContainerStats struct {
	Base
	Stats *docker.Stats
	Container docker.APIContainers
}

func NewContainerStats(src string, stats *docker.Stats, cnt docker.APIContainers) ContainerStats{
	return ContainerStats{
		Base: NewBase(src),
		Stats: stats,
		Container: cnt,
	}
}

func (cs *ContainerStats) GetContainer() *types.Container {
	return &types.Container{
		ID: cs.Container.ID,
		Names: cs.Container.Names,
		Command: cs.Container.Command,
		Created: cs.Container.Created,
		Image: cs.Container.Image,
		Labels: cs.Container.Labels,
	}
}

// Flat out copied from https://github.com/elastic/beats
func (cs *ContainerStats) GetCpuStats() CPUStats {
	cnt := cs.GetContainer()
	// TODO: Use NewCPUStats?
	return CPUStats{
		Base: cs.Base,
		Container:                   cnt,
		PerCpuUsage:                 perCpuUsage(cs.Stats),
		TotalUsage:                  totalUsage(cs.Stats),
		UsageInKernelmode:           cs.Stats.CPUStats.CPUUsage.UsageInKernelmode,
		UsageInKernelmodePercentage: usageInKernelmode(cs.Stats),
		UsageInUsermode:             cs.Stats.CPUStats.CPUUsage.UsageInUsermode,
		UsageInUsermodePercentage:   usageInUsermode(cs.Stats),
		SystemUsage:                 cs.Stats.CPUStats.SystemCPUUsage,
		SystemUsagePercentage:       systemUsage(cs.Stats),
	}
}

func (cs *ContainerStats) GetMemStats() MemoryStats {
	cnt := cs.GetContainer()
	// TODO: Use NewMemoryStats?
	return MemoryStats{
		Base: cs.Base,
		Container: cnt,
		Failcnt:   float64(cs.Stats.MemoryStats.Failcnt),
		Limit:     float64(cs.Stats.MemoryStats.Limit),
		MaxUsage:  float64(cs.Stats.MemoryStats.MaxUsage),
		TotalRss:  float64(cs.Stats.MemoryStats.Stats.TotalRss),
		TotalRssP: float64(cs.Stats.MemoryStats.Stats.TotalRss) / float64(cs.Stats.MemoryStats.Limit),
		Usage:     float64(cs.Stats.MemoryStats.Usage),
		UsageP:    float64(cs.Stats.MemoryStats.Usage) / float64(cs.Stats.MemoryStats.Limit),
	}
}

func (cs *ContainerStats) GetNetStats() NetStats {
	cnt := cs.GetContainer()
	// TODO: Use NewNetStats?
	return NetStats{
		Base: 		cs.Base,
		Container: 	cnt,
		RxBytes:    float64(cs.Stats.Network.RxBytes),
		RxDropped:  float64(cs.Stats.Network.RxDropped),
		RxErrors:   float64(cs.Stats.Network.RxErrors),
		RxPackets:  float64(cs.Stats.Network.RxPackets),
		TxBytes:    float64(cs.Stats.Network.TxBytes),
		TxDropped:  float64(cs.Stats.Network.TxDropped),
		TxErrors:   float64(cs.Stats.Network.TxErrors),
		TxPackets:  float64(cs.Stats.Network.TxPackets),
	}
}

func (cs *ContainerStats) GetNetPerIfaceStats(iface string) NetStats {
	cnt := cs.GetContainer()
	// TODO: Use NewNetStats?
	return NetStats{
		Base: 		cs.Base,
		Container: 	cnt,
		NameInterface: iface,
		RxBytes:    float64(cs.Stats.Networks[iface].RxBytes),
		RxDropped:  float64(cs.Stats.Networks[iface].RxDropped),
		RxErrors:   float64(cs.Stats.Networks[iface].RxErrors),
		RxPackets:  float64(cs.Stats.Networks[iface].RxPackets),
		TxBytes:    float64(cs.Stats.Networks[iface].TxBytes),
		TxDropped:  float64(cs.Stats.Networks[iface].TxDropped),
		TxErrors:   float64(cs.Stats.Networks[iface].TxErrors),
		TxPackets:  float64(cs.Stats.Networks[iface].TxPackets),
	}
}
