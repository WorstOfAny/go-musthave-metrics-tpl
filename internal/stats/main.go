package stats

import(
	"runtime"
	"math/rand"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
	"iter"
)

type Stats struct {
	memstats *runtime.MemStats
	metrics [29]models.Metrics

	Alloc *float64
	BuckHashSys *float64
	Frees *float64
	GCCPUFraction *float64
	GCSys *float64
	HeapAlloc *float64
	HeapIdle *float64
	HeapInuse *float64
	HeapObjects *float64
	HeapReleased *float64
	HeapSys *float64
	LastGC *float64
	Lookups *float64
	MCacheInuse *float64
	MCacheSys *float64
	MSpanInuse *float64
	MSpanSys *float64
	Mallocs *float64
	NextGC *float64
	NumForcedGC *float64
	NumGC *float64
	OtherSys *float64
	PauseTotalNs *float64
	StackInuse *float64
	StackSys *float64
	Sys *float64
	TotalAlloc *float64
	RandomValue *float64
	PollCount *int64
}

func NewStats() (*Stats) {
	stats := &Stats{
		memstats: &runtime.MemStats{},
		Alloc: new(float64),
		BuckHashSys: new(float64),
		Frees: new(float64),
		GCCPUFraction: new(float64),
		GCSys: new(float64),
		HeapAlloc: new(float64),
		HeapIdle: new(float64),
		HeapInuse: new(float64),
		HeapObjects: new(float64),
		HeapReleased: new(float64),
		HeapSys: new(float64),
		LastGC: new(float64),
		Lookups: new(float64),
		MCacheInuse: new(float64),
		MCacheSys: new(float64),
		MSpanInuse: new(float64),
		MSpanSys: new(float64),
		Mallocs: new(float64),
		NextGC: new(float64),
		NumForcedGC: new(float64),
		NumGC: new(float64),
		OtherSys: new(float64),
		PauseTotalNs: new(float64),
		StackInuse: new(float64),
		StackSys: new(float64),
		Sys: new(float64),
		TotalAlloc: new(float64),
		RandomValue: new(float64),
		PollCount: new(int64),
	}

	stats.metrics = [29]models.Metrics{
		models.Metrics{ID: "Alloc", MType: models.Gauge, Value: stats.Alloc},
		models.Metrics{ID: "BuckHashSys", MType: models.Gauge, Value: stats.BuckHashSys},
		models.Metrics{ID: "Frees", MType: models.Gauge, Value: stats.Frees},
		models.Metrics{ID: "GCCPUFraction", MType: models.Gauge, Value: stats.GCCPUFraction},
		models.Metrics{ID: "GCSys", MType: models.Gauge, Value: stats.GCSys},
		models.Metrics{ID: "HeapAlloc", MType: models.Gauge, Value: stats.HeapAlloc},
		models.Metrics{ID: "HeapIdle", MType: models.Gauge, Value: stats.HeapIdle},
		models.Metrics{ID: "HeapInuse", MType: models.Gauge, Value: stats.HeapInuse},
		models.Metrics{ID: "HeapObjects", MType: models.Gauge, Value: stats.HeapObjects},
		models.Metrics{ID: "HeapReleased", MType: models.Gauge, Value: stats.HeapReleased},
		models.Metrics{ID: "HeapSys", MType: models.Gauge, Value: stats.HeapSys},
		models.Metrics{ID: "LastGC", MType: models.Gauge, Value: stats.LastGC},
		models.Metrics{ID: "Lookups", MType: models.Gauge, Value: stats.Lookups},
		models.Metrics{ID: "MCacheInuse", MType: models.Gauge, Value: stats.MCacheInuse},
		models.Metrics{ID: "MCacheSys", MType: models.Gauge, Value: stats.MCacheSys},
		models.Metrics{ID: "MSpanInuse", MType: models.Gauge, Value: stats.MSpanInuse},
		models.Metrics{ID: "MSpanSys", MType: models.Gauge, Value: stats.MSpanSys},
		models.Metrics{ID: "Mallocs", MType: models.Gauge, Value: stats.Mallocs},
		models.Metrics{ID: "NextGC", MType: models.Gauge, Value: stats.NextGC},
		models.Metrics{ID: "NumForcedGC", MType: models.Gauge, Value: stats.NumForcedGC},
		models.Metrics{ID: "NumGC", MType: models.Gauge, Value: stats.NumGC},
		models.Metrics{ID: "OtherSys", MType: models.Gauge, Value: stats.OtherSys},
		models.Metrics{ID: "PauseTotalNs", MType: models.Gauge, Value: stats.PauseTotalNs},
		models.Metrics{ID: "StackInuse", MType: models.Gauge, Value: stats.StackInuse},
		models.Metrics{ID: "StackSys", MType: models.Gauge, Value: stats.StackSys},
		models.Metrics{ID: "Sys", MType: models.Gauge, Value: stats.Sys},
		models.Metrics{ID: "TotalAlloc", MType: models.Gauge, Value: stats.TotalAlloc},
		models.Metrics{ID: "RandomValue", MType: models.Gauge, Value: stats.RandomValue},
		models.Metrics{ID: "PollCount", MType: models.Counter, Delta: stats.PollCount},
	}

	return stats
}

func(s *Stats) AllMetrics() iter.Seq[*models.Metrics] {
	return func(yield func(*models.Metrics) bool) {
		for _, v := range s.metrics {
			if !yield(&v) { return }
		}
	}
}

func(s *Stats) Update() {
	runtime.ReadMemStats(s.memstats)

	*s.Alloc = float64(s.memstats.Alloc)
	*s.BuckHashSys = float64(s.memstats.BuckHashSys)
	*s.Frees = float64(s.memstats.Frees)
	*s.GCCPUFraction = float64(s.memstats.GCCPUFraction)
	*s.GCSys = float64(s.memstats.GCSys)
	*s.HeapAlloc = float64(s.memstats.HeapAlloc)
	*s.HeapIdle = float64(s.memstats.HeapIdle)
	*s.HeapInuse = float64(s.memstats.HeapInuse)
	*s.HeapObjects = float64(s.memstats.HeapObjects)
	*s.HeapReleased = float64(s.memstats.HeapReleased)
	*s.HeapSys = float64(s.memstats.HeapSys)
	*s.LastGC = float64(s.memstats.LastGC)
	*s.Lookups = float64(s.memstats.Lookups)
	*s.MCacheSys = float64(s.memstats.MCacheInuse)
	*s.MCacheSys = float64(s.memstats.MCacheSys)
	*s.MSpanInuse = float64(s.memstats.MSpanInuse)
	*s.MSpanSys = float64(s.memstats.MSpanSys)
	*s.Mallocs = float64(s.memstats.Mallocs)
	*s.NextGC = float64(s.memstats.NextGC)
	*s.NumForcedGC = float64(s.memstats.NumForcedGC)
	*s.NumGC = float64(s.memstats.NumGC)
	*s.OtherSys = float64(s.memstats.OtherSys)
	*s.PauseTotalNs = float64(s.memstats.PauseTotalNs)
	*s.StackInuse = float64(s.memstats.StackInuse)
	*s.StackSys = float64(s.memstats.StackSys)
	*s.Sys = float64(s.memstats.Sys)
	*s.TotalAlloc = float64(s.memstats.TotalAlloc)

	*s.PollCount += 1
	*s.RandomValue = rand.Float64()
}
