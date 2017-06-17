package statsq

import (
	"bytes"
	"fmt"
	"github.com/qnib/qframe-types"
	"github.com/zpatrick/go-config"
	"io"
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	MAX_UNPROCESSED_PACKETS = 1000
	TCP_READ_SIZE           = 4096
	version = "0.1.0"
)

type StatsQ struct {
	Name            string
	Version			string
	Parser          MsgParser
	Signalchan      chan os.Signal
	Cfg             *config.Config
	In              chan *qtypes.StatsdPacket
	Counters        map[string]float64
	Gauges          map[string]float64
	Timers          map[string]Float64Slice
	CountInactivity map[string]int64
	Sets            map[string][]string
	ReceiveCounter  string
	QChan           qtypes.QChan
	Percentiles     Percentiles
	BucketMapping   map[string]BucketID
}

func NewStatsQ(cfg *config.Config) StatsQ {
	return NewNamedStatsQ("", cfg, qtypes.NewQChan())
}

func NewNamedStatsQ(name string, cfg *config.Config, qchan qtypes.QChan) StatsQ {
	sd := StatsQ{
		Name:            name,
		Version:  		 version,
		Parser:          MsgParser{debug: true},
		Signalchan:      make(chan os.Signal, 1),
		Cfg:             cfg,
		In:              make(chan *qtypes.StatsdPacket, MAX_UNPROCESSED_PACKETS),
		Counters:        make(map[string]float64),
		Gauges:          make(map[string]float64),
		Timers:          make(map[string]Float64Slice),
		CountInactivity: make(map[string]int64),
		Sets:            make(map[string][]string),
		Percentiles:     Percentiles{},
		QChan:           qchan,
		BucketMapping:   map[string]BucketID{},
	}
	sd.ReceiveCounter = sd.StringOr("receive-counter", "")
	sd.Log("info", fmt.Sprintf("Pctls: %s", sd.StringOr("percentiles", "")))
	for _, pctl := range strings.Split(sd.StringOr("percentiles", ""), ",") {
		sd.Percentiles.Set(pctl)
	}
	return sd
}

func (sd *StatsQ) Log(logLevel, msg string) {
	// TODO: Setup in each Log() invocation seems rude
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	dL, _ := sd.Cfg.StringOr("log.level", "info")
	debug := sd.BoolOr("debug", false)
	dI := qtypes.LogStrToInt(dL)
	lI := qtypes.LogStrToInt(logLevel)
	if dI >= lI || debug {
		log.Printf("[%+6s] %15s Name:%-10s >> %s", strings.ToUpper(logLevel), "statsq v"+sd.Version, sd.Name, msg)
	}
}

func (sd *StatsQ) StringOr(path, alt string) string {
	if sd.Name != "" {
		path = fmt.Sprintf("%s.%s", sd.Name, path)
	}
	res, err := sd.Cfg.String(path)
	if err != nil {
		res = alt
	}
	return res
}

func (sd *StatsQ) String(path string) string {
	return sd.StringOr(path, "")
}

func (sd *StatsQ) Bool(path string) bool {
	return sd.BoolOr(path, false)
}

func (sd *StatsQ) BoolOr(path string, alt bool) bool {
	if sd.Name != "" {
		path = fmt.Sprintf("%s.%s", sd.Name, path)
	}
	res, err := sd.Cfg.Bool(path)
	if err != nil {
		res = alt
	}
	return res
}

func (sd *StatsQ) IntOr(path string, alt int) int {
	if sd.Name != "" {
		path = fmt.Sprintf("%s.%s", sd.Name, path)
	}
	res, err := sd.Cfg.Int(path)
	if err != nil {
		res = alt
	}
	return res
}

func (sd *StatsQ) Int(path string) int {
	return sd.IntOr(path, 0)
}

func (sd *StatsQ) Run() {
	signal.Notify(sd.Signalchan, syscall.SIGTERM)
	go sd.startUDPListener()
	go sd.startTCPListener()
	sd.LoopChannel()
}

// RelayMetrics listens for metrics on the QChan.Data channel and sends the metrics to the backends
func (sd *StatsQ) RelayMetrics() {

}

func (sd *StatsQ) startUDPListener() {
	serviceAddress := sd.StringOr("address", ":8125")
	address, _ := net.ResolveUDPAddr("udp", serviceAddress)
	sd.Log("info", fmt.Sprintf("listening on %s", address))
	listener, err := net.ListenUDP("udp", address)
	if err != nil {
		log.Fatalf("ERROR: ListenUDP - %s", err)
	}
	sd.ParseTo(listener, false)
}

func (sd *StatsQ) startTCPListener() {
	serviceAddress := sd.StringOr("tcpaddr", "")
	if serviceAddress == "" {
		return
	}
	address, _ := net.ResolveTCPAddr("tcp", serviceAddress)
	log.Printf("listening on %s", address)
	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		log.Fatalf("ERROR: ListenTCP - %s", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatalf("ERROR: AcceptTCP - %s", err)
		}
		go sd.ParseTo(conn, true)
	}
}

func (sd *StatsQ) ParseTo(conn io.ReadCloser, partialReads bool) {
	defer conn.Close()
	maxUdpPacketSize := sd.Int("max-udp-packet-size")
	prefix := sd.String("prefix")
	postfix := sd.String("postfix")
	debug := sd.Bool("debug")
	parser := NewParser(conn, partialReads, debug, maxUdpPacketSize, prefix, postfix)
	sd.Log("debug", "Start ParseTo Loop")
	for {
		p, more := parser.Next()
		sd.Log("debug", fmt.Sprintf("Received: %v", p))
		if p != nil {
			sd.In <- p
		}
		if !more {
			break
		}
	}
}

func (sd *StatsQ) LoopChannel() {
	tickMs := sd.IntOr("send-metric-ms", 1000)
	sd.Log("info", fmt.Sprintf("StatsQ ticker: %dms", tickMs))
	ticker := time.NewTicker(time.Duration(tickMs) * time.Millisecond).C
	for {
		select {
		case s := <-sd.In:
			sd.HandlerStatsdPacket(s)
		case <-ticker:
			sd.FanOutMetrics()
		}
	}
}

func (sd *StatsQ) HandlerStatsdPacket(sp *qtypes.StatsdPacket) {
	if sd.ReceiveCounter != "" {
		v, ok := sd.Counters[sd.ReceiveCounter]
		if !ok || v < 0 {
			sd.Counters[sd.ReceiveCounter] = 0
		}
		sd.Counters[sd.ReceiveCounter] += 1
	}
	bid := NewBucketID(sp.Bucket, sp.Dimensions)
	bkey := bid.ID
	_, ok := sd.BucketMapping[bkey]
	if !ok {
		log.Printf("Include bid '%s' w/ key '%s' in BucketMapping", bid.BucketName, bkey)
		sd.BucketMapping[bkey] = bid
	}
	switch sp.Modifier {
	case "ms":
		_, ok := sd.Timers[bkey]
		if !ok {
			var t Float64Slice
			sd.Timers[bkey] = t
		}
		sd.Timers[bkey] = append(sd.Timers[bkey], sp.ValFlt)
	case "g":
		gaugeValue, _ := sd.Gauges[bkey]
		if sp.ValStr == "" {
			gaugeValue = sp.ValFlt
		} else if sp.ValStr == "+" {
			// watch out for overflows
			if sp.ValFlt > (math.MaxFloat64 - gaugeValue) {
				gaugeValue = math.MaxFloat64
			} else {
				gaugeValue += sp.ValFlt
			}
		} else if sp.ValStr == "-" {
			// subtract checking for negative numbers
			if sp.ValFlt > gaugeValue {
				gaugeValue = 0
			} else {
				gaugeValue -= sp.ValFlt
			}
		}
		sd.Gauges[bkey] = gaugeValue
	case "c":
		_, ok := sd.Counters[bkey]
		if !ok {
			sd.Counters[bkey] = 0
		}
		sd.Counters[bkey] += sp.ValFlt * float64(1/sp.Sampling)
	case "s":
		_, ok := sd.Sets[bkey]
		if !ok {
			sd.Sets[bkey] = make([]string, 0)
		}
		sd.Sets[bkey] = append(sd.Sets[bkey], sp.ValStr)
	}
}

func (sd *StatsQ) FanOutMetrics() {
	now := time.Now()
	sd.FanOutCounters(now)
	sd.FanOutGauges(now)
	sd.FanOutSets(now)
	sd.FanOutTimers(now)

}

func (sd *StatsQ) ParseLine(msg string) (err error) {
	sp := sd.Parser.parseLine([]byte(msg))
	sd.HandlerStatsdPacket(sp)
	return
}

func (sd *StatsQ) FanOutCounters(now time.Time) int64 {
	var num int64
	// continue sending zeros for counters for a short period of time even if we have no new data
	for id, value := range sd.Counters {
		bid, ok := sd.BucketMapping[id]
		if !ok {
			sd.Log("error", fmt.Sprintf("Could not find BucketID for key '%s'", id))
			return num
		}
		m := qtypes.NewExt(sd.Name, bid.BucketName, qtypes.Counter, value, bid.Dimensions.Map, now, false)
		sd.sendMetric(m)
		delete(sd.Counters, id)
		sd.CountInactivity[id] = 0
		num++
	}
	for id, purgeCount := range sd.CountInactivity {
		bid, ok := sd.BucketMapping[id]
		if !ok {
			sd.Log("error", fmt.Sprintf("Could not find BucketID for key '%s'", id))
			return num
		}
		if purgeCount > 0 {
			m := qtypes.NewExt(sd.Name, bid.BucketName, qtypes.Counter, 0.0, bid.Dimensions.Map, now, false)
			sd.sendMetric(m)
			num++
		}
		sd.CountInactivity[id] += 1
		if sd.CountInactivity[id] > int64(sd.Int("persist-count-keys")) {
			delete(sd.CountInactivity, id)
		}
	}
	return num
}

func (sd *StatsQ) FanOutGauges(now time.Time) int64 {
	var num int64
	for id, currentValue := range sd.Gauges {
		bid, ok := sd.BucketMapping[id]
		if !ok {
			sd.Log("error", fmt.Sprintf("Could not find BucketID for key '%s'", id))
			return num
		}
		m := qtypes.NewExt(sd.Name, bid.BucketName, qtypes.Gauge, currentValue, bid.Dimensions.Map, now, false)
		sd.sendMetric(m)
		num++
		if sd.Bool("delete-gauges") {
			sd.Log("info", fmt.Sprintf("Delete gauges with id '%s'", id))
			delete(sd.Gauges, id)
		}
	}
	return num
}

func (sd *StatsQ) FanOutSets(now time.Time) int64 {
	num := int64(len(sd.Sets))
	for id, set := range sd.Sets {
		bid, ok := sd.BucketMapping[id]
		if !ok {
			sd.Log("error", fmt.Sprintf("Could not find BucketID for key '%s'", id))
			return num
		}
		uniqueSet := map[string]bool{}
		for _, str := range set {
			uniqueSet[str] = true
		}
		m := qtypes.NewExt(sd.Name, bid.BucketName, qtypes.Gauge, float64(len(uniqueSet)), bid.Dimensions.Map, now, false)
		sd.sendMetric(m)
		delete(sd.Sets, id)
	}
	return num
}

func (sd *StatsQ) FanOutTimers(now time.Time) int64 {
	var num int64
	//postfix := sd.String("postfix")
	for id, timer := range sd.Timers {
		bid, ok := sd.BucketMapping[id]
		if !ok {
			sd.Log("error", fmt.Sprintf("Could not find BucketID for key '%s'", id))
			return num
		}
		//bucketWithoutPostfix := bid.BucketName[:len(bid.BucketName)-len(postfix)]
		num++

		sort.Sort(timer)
		min := timer[0]
		max := timer[len(timer)-1]
		maxAtThreshold := max
		count := len(timer)

		sum := float64(0)
		for _, value := range timer {
			sum += value
		}
		mean := sum / float64(len(timer))

		for _, pct := range sd.Percentiles {
			if len(timer) > 1 {
				var abs float64
				if pct.float >= 0 {
					abs = pct.float
				} else {
					abs = 100 + pct.float
				}
				// poor man's math.Round(x):
				// math.Floor(x + 0.5)
				indexOfPerc := int(math.Floor(((abs / 100.0) * float64(count)) + 0.5))
				if pct.float >= 0 {
					indexOfPerc -= 1 // index offset=0
				}
				maxAtThreshold = timer[indexOfPerc]
			}

			var name string
			if pct.float >= 0 {
				name = fmt.Sprintf("%s.upper_%s", bid.BucketName, pct.str)
			} else {
				name = fmt.Sprintf("%s.lower_%s", bid.BucketName, pct.str[1:])
			}
			m := qtypes.NewExt(sd.Name, name, qtypes.Gauge, maxAtThreshold, bid.GetDims(), now, false)
			sd.sendMetric(m)
		}

		name := fmt.Sprintf("%s.mean", bid.BucketName)
		m := qtypes.NewExt(sd.Name, name, qtypes.Gauge, mean, bid.GetDims(), now, false)
		sd.sendMetric(m)
		name = fmt.Sprintf("%s.upper", bid.BucketName)
		m = qtypes.NewExt(sd.Name, name, qtypes.Gauge, max, bid.GetDims(), now, false)
		sd.sendMetric(m)
		name = fmt.Sprintf("%s.lower", bid.BucketName)
		m = qtypes.NewExt(sd.Name, name, qtypes.Gauge, min, bid.GetDims(), now, false)
		sd.sendMetric(m)
		name = fmt.Sprintf("%s.count", bid.BucketName)
		m = qtypes.NewExt(sd.Name, name, qtypes.Gauge, float64(count), bid.GetDims(), now, false)
		sd.sendMetric(m)
		delete(sd.Timers, id)
	}
	return num
}

func (sd *StatsQ) sendMetric(m qtypes.Metric) {
	sd.Log("trace", m.ToOpenTSDB())
	sd.QChan.Data.Send(m)
}

func sanitizeBucket(bucket string) string {
	b := make([]byte, len(bucket))
	var bl int

	for i := 0; i < len(bucket); i++ {
		c := bucket[i]
		switch {
		case (c >= byte('a') && c <= byte('z')) || (c >= byte('A') && c <= byte('Z')) || (c >= byte('0') && c <= byte('9')) || c == byte('-') || c == byte('.') || c == byte('_'):
			b[bl] = c
			bl++
		case c == byte(' '):
			b[bl] = byte('_')
			bl++
		case c == byte('/'):
			b[bl] = byte('-')
			bl++
		}
	}
	return string(b[:bl])
}

func (sd *StatsQ) ProcessCounters(buffer *bytes.Buffer, now int64) int64 {
	var num int64
	// continue sending zeros for counters for a short period of time even if we have no new data
	for bucket, value := range sd.Counters {
		fmt.Fprintf(buffer, "%s %s %d\n", bucket, strconv.FormatFloat(value, 'f', -1, 64), now)
		delete(sd.Counters, bucket)
		sd.CountInactivity[bucket] = 0
		num++
	}
	for bucket, purgeCount := range sd.CountInactivity {
		if purgeCount > 0 {
			fmt.Fprintf(buffer, "%s 0 %d\n", bucket, now)
			num++
		}
		sd.CountInactivity[bucket] += 1
		if sd.CountInactivity[bucket] > int64(sd.Int("persist-count-keys")) {
			delete(sd.CountInactivity, bucket)
		}
	}
	return num
}

func (sd *StatsQ) ProcessGauges(buffer *bytes.Buffer, now int64) int64 {
	var num int64

	for bucket, currentValue := range sd.Gauges {
		fmt.Fprintf(buffer, "%s %s %d\n", bucket, strconv.FormatFloat(currentValue, 'f', -1, 64), now)
		num++
		if sd.Bool("delete-gauges") {
			delete(sd.Gauges, bucket)
		}
	}
	return num
}

func (sd *StatsQ) ProcessSets(buffer *bytes.Buffer, now int64) int64 {
	num := int64(len(sd.Sets))
	for bucket, set := range sd.Sets {

		uniqueSet := map[string]bool{}
		for _, str := range set {
			uniqueSet[str] = true
		}

		fmt.Fprintf(buffer, "%s %d %d\n", bucket, len(uniqueSet), now)
		delete(sd.Sets, bucket)
	}
	return num
}

func (sd *StatsQ) ProcessTimers(buffer *bytes.Buffer, now int64) int64 {
	var num int64
	postfix := sd.String("postfix")
	for bucket, timer := range sd.Timers {
		bucketWithoutPostfix := bucket[:len(bucket)-len(postfix)]
		num++

		sort.Sort(timer)
		min := timer[0]
		max := timer[len(timer)-1]
		maxAtThreshold := max
		count := len(timer)

		sum := float64(0)
		for _, value := range timer {
			sum += value
		}
		mean := sum / float64(len(timer))

		for _, pct := range sd.Percentiles {
			if len(timer) > 1 {
				var abs float64
				if pct.float >= 0 {
					abs = pct.float
				} else {
					abs = 100 + pct.float
				}
				// poor man's math.Round(x):
				// math.Floor(x + 0.5)
				indexOfPerc := int(math.Floor(((abs / 100.0) * float64(count)) + 0.5))
				if pct.float >= 0 {
					indexOfPerc -= 1 // index offset=0
				}
				maxAtThreshold = timer[indexOfPerc]
			}

			var tmpl string
			var pctstr string
			if pct.float >= 0 {
				tmpl = "%s.upper_%s%s %s %d\n"
				pctstr = pct.str
			} else {
				tmpl = "%s.lower_%s%s %s %d\n"
				pctstr = pct.str[1:]
			}
			threshold_s := strconv.FormatFloat(maxAtThreshold, 'f', -1, 64)
			fmt.Fprintf(buffer, tmpl, bucketWithoutPostfix, pctstr, postfix, threshold_s, now)
		}

		mean_s := strconv.FormatFloat(mean, 'f', -1, 64)
		max_s := strconv.FormatFloat(max, 'f', -1, 64)
		min_s := strconv.FormatFloat(min, 'f', -1, 64)

		fmt.Fprintf(buffer, "%s.mean%s %s %d\n", bucketWithoutPostfix, postfix, mean_s, now)
		fmt.Fprintf(buffer, "%s.upper%s %s %d\n", bucketWithoutPostfix, postfix, max_s, now)
		fmt.Fprintf(buffer, "%s.lower%s %s %d\n", bucketWithoutPostfix, postfix, min_s, now)
		fmt.Fprintf(buffer, "%s.count%s %d %d\n", bucketWithoutPostfix, postfix, count, now)

		delete(sd.Timers, bucket)
	}
	return num
}
