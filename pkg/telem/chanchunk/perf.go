package chanchunk

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	log "github.com/sirupsen/logrus"
	"runtime"
	"sync"
	"time"
)

type Perf struct {
	mu     sync.Mutex
	Data   map[string]int64
	Client api.WriteAPIBlocking
	NodeID int
}

func (p *Perf) IncrementSample(sample string, value int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Data[sample] += value
}

func (p *Perf) WriteSamples() {
	t := time.NewTicker(1 * time.Second)
	for range t.C {
		p.runtime()
		err := p.Client.WritePoint(context.Background(), p.buildMeasurement())
		if err != nil {
			log.Warn(err)
		}
		p.Data = map[string]int64{}
	}
}

func (p *Perf) runtime() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	p.Data["mem_usage"] = int64(bToMb(m.Alloc))
	p.Data["num_go_routine"] = int64(runtime.NumGoroutine())
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func (p *Perf) buildMeasurement() *write.Point {
	point := influxdb2.NewPointWithMeasurement(fmt.Sprintf("node_%v_stats", p.NodeID))
	for k, v := range p.Data {
		point.AddField(k, v)
	}
	point.SetTime(time.Now())
	return point
}
