package qframe_handler_influxdb

import (
	"fmt"
	"time"
	"reflect"
	"sync"
	"github.com/zpatrick/go-config"
	"github.com/influxdata/influxdb/client/v2"

	"github.com/qnib/qframe-types"
)

const (
	version = "0.1.1"
	pluginTyp = "handler"
	pluginPkg = "influxdb"
)

type Plugin struct {
    qtypes.Plugin
	cli client.Client
	metricCount int
	mutex sync.Mutex

}

func New(qChan qtypes.QChan, cfg config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
		metricCount: 0,
	}
	return p, err
}

// Connect creates a connection to InfluxDB
func (p *Plugin) Connect() {
	host := p.CfgStringOr("host", "localhost")
	port, _ := p.Cfg.StringOr(fmt.Sprintf("handler.%s.port", p.Name), "8086")
	username, _ := p.Cfg.StringOr(fmt.Sprintf("handler.%s.username", p.Name), "root")
	password, _ := p.Cfg.StringOr(fmt.Sprintf("handler.%s.password", p.Name), "root")
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

func (p *Plugin) MetricsToBatchPoint(m qtypes.Metric) (pt *client.Point, err error) {
	fields := map[string]interface{}{
		"value": m.Value,
	}
	pt, err = client.NewPoint(m.Name, m.Dimensions, fields, m.Time)
	return
}
// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	p.Log("info", fmt.Sprintf("Start log handler %sv%s", p.Name, version))
	batchSize := p.CfgIntOr("batch-size", 100)
	tick := p.CfgIntOr("ticker-msec", 1000)
	p.Connect()
	bg := p.QChan.Data.Join()
	tc := p.QChan.Tick.Join()
	inputs := p.GetInputs()
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
		case val := <-bg.Read:
			switch val.(type) {
			case qtypes.Metric:
				m := val.(qtypes.Metric)
				pt, err := p.MetricsToBatchPoint(m)
				if err != nil {
					p.Log("error", fmt.Sprintf("%v", err))
					continue
				}
				bp.AddPoint(pt)
				if ! m.InputsMatch(inputs) {
					continue
				}

				if len(bp.Points()) >= batchSize {
					now := time.Now()
					bLen := len(bp.Points())
					p.Log("debug", fmt.Sprintf("%d >= %d: Write batch",bLen, batchSize))
					p.metricCount += bLen
					p.QChan.Data.Send(qtypes.NewExt(p.Name, "influxdb.batch.size", qtypes.Gauge, float64(bLen), dims, time.Now(), false))
					bp = p.WriteBatch(bp)
					took := time.Now().Sub(now)
					p.QChan.Data.Send(qtypes.NewExt(p.Name, "influxdb.batch.duration_ns", qtypes.Gauge, float64(took.Nanoseconds()), dims, time.Now(), false))
					lastTick = now
				}
			}
		case val := <-tc.Read:
			switch val.(type) {
			case qtypes.Ticker:
				tick := val.(qtypes.Ticker)
				tickDiff, skipTick := tick.SkipTick(lastTick)
				if skipTick {
					msg := fmt.Sprintf("tick '%s' | Last tick %s ago (< %s)", tick.Name, tickDiff.String(), tick.Duration.String())
					p.Log("debug", msg)
					continue
				}
				now := time.Now()
				lastTick = now
				// Might take some time
				bLen := len(bp.Points())
				p.Log("debug", fmt.Sprintf("tick '%s' | Last tick %s ago ([some wiggel room] >= %s) - Write batch of %d", tick.Name, tickDiff.String(), tick.Duration.String(), bLen))
				p.metricCount += bLen
				bp = p.WriteBatch(bp)
				took := time.Now().Sub(now)
				p.QChan.Data.Send(qtypes.NewExt(p.Name, "influxdb.batch.size", qtypes.Gauge, float64(bLen), dims, time.Now(), false))
				p.QChan.Data.Send(qtypes.NewExt(p.Name, "influxdb.batch.count", qtypes.Counter, float64(p.metricCount), dims, time.Now(), false))
				p.QChan.Data.Send(qtypes.NewExt(p.Name, "influxdb.batch.duration_ns", qtypes.Gauge, float64(took.Nanoseconds()), dims, time.Now(), false))
			default:
				p.Log("debug", fmt.Sprintf("Received Tick of type %s", reflect.TypeOf(val)))
			}
		}
	}
}
