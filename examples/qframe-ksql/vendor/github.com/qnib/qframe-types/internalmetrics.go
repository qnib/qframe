package qtypes

import (
	"runtime"
	"time"
)

type IntMemoryStats struct {
	Base
	Stats *runtime.MemStats
	NumGoroutine int
}

func NewIntMemoryStats(src string) IntMemoryStats {
	return IntMemoryStats{
		Base: NewBase(src),
		Stats: new(runtime.MemStats),
		NumGoroutine: runtime.NumGoroutine(),
	}
}

func (s *IntMemoryStats) SnapShot() {
	runtime.ReadMemStats(s.Stats)
	s.NumGoroutine = runtime.NumGoroutine()
	s.Time = time.Now()
}

func (s *IntMemoryStats) ToMetrics(src string) ([]Metric) {
	var dims map[string]string
	return []Metric{
		NewExt(src, "internal.goroutine.count", Counter, float64(s.NumGoroutine), dims, s.Time, false),
		NewExt(src, "internal.memory.alloc.total", Counter, float64(s.Stats.TotalAlloc), dims, s.Time, false),
		NewExt(src, "internal.memory.lookups", Counter, float64(s.Stats.Lookups), dims, s.Time, false),
		NewExt(src, "internal.memory.mallocs", Counter, float64(s.Stats.Mallocs), dims, s.Time, false),
		NewExt(src, "internal.memory.frees", Counter, float64(s.Stats.Frees), dims, s.Time, false),
		NewExt(src, "internal.memory.pause.total.ns", Counter, float64(s.Stats.PauseTotalNs), dims, s.Time, false),
		NewExt(src, "internal.memory.gc.count", Counter, float64(s.Stats.NumGC), dims, s.Time, false),
		NewExt(src, "internal.memory.alloc.bytes", Gauge, float64(s.Stats.Alloc), dims, s.Time, false),
		NewExt(src, "internal.memory.sys.bytes", Gauge, float64(s.Stats.Sys), dims, s.Time, false),
		NewExt(src, "internal.memory.heap.alloc.bytes", Gauge, float64(s.Stats.HeapAlloc), dims, s.Time, false),
		NewExt(src, "internal.memory.heap.sys.bytes", Gauge, float64(s.Stats.HeapSys), dims, s.Time, false),
		NewExt(src, "internal.memory.heap.idle.bytes", Gauge, float64(s.Stats.HeapIdle), dims, s.Time, false),
		NewExt(src, "internal.memory.heap.inuse.bytes", Gauge, float64(s.Stats.HeapInuse), dims, s.Time, false),
		NewExt(src, "internal.memory.heap.objects.count", Gauge, float64(s.Stats.HeapObjects), dims, s.Time, false),
	}
}
