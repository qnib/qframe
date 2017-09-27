package qhandler_influxdb

import (
	"fmt"
	"time"
	"reflect"
	"sync"
	"github.com/zpatrick/go-config"
	"github.com/influxdata/influxdb/client/v2"
	"strings"
	"github.com/qframe/types/plugin"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/types/metrics"
	"github.com/qframe/types/ticker"
)

const (
	version = "0.2.1"
	pluginTyp = "handler"
	pluginPkg = "influxdb"
)

type Plugin struct {
    *qtypes_plugin.Plugin
	cli client.Client
	metricCount int
	mutex sync.Mutex
	SanitizeLabels bool

}

func New(qChan qtypes_qchannel.QChan, cfg *config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes_plugin.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
		metricCount: 0,
	}
	p.SanitizeLabels = p.CfgBoolOr("sanitize-labels", false)
	if p.SanitizeLabels {
		p.Log("debug", "Replace '.' in container labels with '_' to play nicer with grafana")
	}
	return p, err
}

// Connect creates a connection to InfluxDB
func (p *Plugin) Connect() {
	host := p.CfgStringOr("host", "localhost")
	port := p.CfgStringOr("port", "8086")
	username := p.CfgStringOr("username", "root")
	password := p.CfgStringOr("password", "root")
	addr := fmt.Sprintf("http://%s:%s", host, port)
	cli := client.HTTPConfig{
		Addr:     addr,
		Username: username,
		Password: password,
	}
	var err error
	p.cli, err = client.NewHTTPClient(cli)
	if err != nil {
		p.Log("error", fmt.Sprintf("Error during connection to InfluxDB '%s': %v", addr, err))
	} else {
		p.Log("info", fmt.Sprintf("Established connection to '%s", addr))
	}
}

func (p *Plugin) NewBatchPoints() client.BatchPoints {
	dbName := p.CfgStringOr("database", "qframe")
	dbPrec := p.CfgStringOr("precision", "s")
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  dbName,
		Precision: dbPrec,
	})
	if err != nil {
		p.Log("error", fmt.Sprintf("Not able to create BatchPoints: %v", err))
	}
	return bp

}

func (p *Plugin) WriteBatch(points client.BatchPoints) client.BatchPoints {
	err := p.cli.Write(points)
	if err != nil {
		p.Log("error", fmt.Sprintf("Not able to write BatchPoints: %v", err))
	}
	return p.NewBatchPoints()
}

func (p *Plugin) MetricsToBatchPoint(m qtypes_metrics.Metric) (pt *client.Point, err error) {
	fields := map[string]interface{}{
		"value": m.Value,
	}
	dims := map[string]string{}
	if p.SanitizeLabels {
		for k,v := range m.Dimensions {
			dims[strings.Replace(k, ".", "_", -1)] = v
		}
	} else {
		dims = m.Dimensions
	}
	pt, err = client.NewPoint(m.Name, dims, fields, m.Time)
	return
}

// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start handler %sv%s", p.Name, version))
	batchSize := p.CfgIntOr("batch-size", 100)
	tick := p.CfgIntOr("ticker-msec", 1000)
	p.Connect()
	dc, _, tc := p.JoinChannels()
	bp := p.NewBatchPoints()
	p.StartTicker("influxdb", tick)
	dims := map[string]string{
		"version": version,
		"plugin": p.Name,
	}
	// Initialise lastTick with time of a year ago
	lastTick := time.Now().AddDate(0,0,-1)
	for {
		select {
		case val := <- dc.Read:
			switch val.(type) {
			case qtypes_metrics.Metric:
				m := val.(qtypes_metrics.Metric)
				pt, err := p.MetricsToBatchPoint(m)
				if err != nil {
					p.Log("error", fmt.Sprintf("%v", err))
					continue
				}
				bp.AddPoint(pt)
				if len(bp.Points()) >= batchSize {
					now := time.Now()
					bLen := len(bp.Points())
					p.Log("debug", fmt.Sprintf("%d >= %d: Write batch",bLen, batchSize))
					p.metricCount += bLen+1
					pt, _ = p.MetricsToBatchPoint(qtypes_metrics.NewExt(p.Name, "influxdb.batch.size", qtypes_metrics.Gauge, float64(bLen+1), dims, time.Now(), false))
					bp.AddPoint(pt)
					bp = p.WriteBatch(bp)
					//took := time.Now().Sub(now)
					//p.QChan.Data.Send(qtypes_metrics.NewStatsdPacket("influxdb.batch.write.ns",  fmt.Sprintf("%d", took.Nanoseconds()), "ms"))
					lastTick = now
				}
			}
		case val := <-tc.Read:
			switch val.(type) {
			case qtypes_ticker.Ticker:
				tick := val.(qtypes_ticker.Ticker)
				tickDiff, skipTick := tick.SkipTick(lastTick)
				if skipTick {
					msg := fmt.Sprintf("tick '%s' | Last tick %s ago (< %s)", tick.Name, tickDiff.String(), tick.Duration.String())
					p.Log("trace", msg)
					continue
				}
				now := time.Now()
				lastTick = now
				// Might take some time
				bLen := len(bp.Points())
				p.Log("trace", fmt.Sprintf("tick '%s' | Last tick %s ago ([some wiggel room] >= %s) - Write batch of %d", tick.Name, tickDiff.String(), tick.Duration.String(), bLen))
				pt, _ := p.MetricsToBatchPoint(qtypes_metrics.NewExt(p.Name, "influxdb.batch.size", qtypes_metrics.Gauge, float64(bLen+1), dims, time.Now(), false))
				bp.AddPoint(pt)
				p.metricCount += bLen+1
				bp = p.WriteBatch(bp)
				//took := time.Now().Sub(now)
				//p.QChan.Data.Send(qtypes_metrics.NewStatsdPacket("influxdb.batch.write.ns",  fmt.Sprintf("%d", took.Nanoseconds()), "ms"))
			default:
				p.Log("warn", fmt.Sprintf("Received Tick of type %s", reflect.TypeOf(val)))
			}
		}
	}
}
