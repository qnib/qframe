package qtypes

import (
	"log"
	"fmt"
	"strings"
	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-utils"
	"time"
	"github.com/pkg/errors"
	"strconv"
)

const (
	FILTER = "filter"
	COLLECTOR = "collector"
	HANDLER = "handler"
)

type Plugin struct {
	QChan 			QChan
	Cfg 			*config.Config
	MyID			int
	Typ				string
	Pkg				string
	Version 		string
	Name 			string
	LogOnlyPlugs 	[]string
	MsgCount		map[string]float64
	LocalCfg 		map[string]string

}

func NewPlugin(qChan QChan, cfg *config.Config) Plugin {
	return Plugin{
		QChan: qChan,
		Cfg: cfg,
	}
}


func NewNamedPlugin(qChan QChan, cfg *config.Config, typ, pkg, name, version string) Plugin {
	p := Plugin{
		QChan: 			qChan,
		Cfg:   			cfg,
		Typ:   			typ,
		Pkg:  			pkg,
		Version:		version,
		Name: 			name,
		LogOnlyPlugs:   []string{},
		MsgCount:       map[string]float64{
			"received": 0.0,
			"loopDrop": 0.0,
			"inputDrop": 0.0,
			"successDrop": 0.0,
		},
	}
	p.LocalCfg, _  = cfg.Settings()
	logPlugs, err := cfg.String("log.only-plugins")
	if err == nil {
		p.LogOnlyPlugs = strings.Split(logPlugs, ",")
	}
	return p
}


func LogStrToInt(level string) int {
	def := 6
	switch level {
	case "panic":
		return 0
	case "error":
		return 3
	case "warn":
		return 4
	case "notice":
		return 5
	case "info":
		return 6
	case "debug":
		return 7
	case "trace":
		return 8
	default:
		return def
	}
}

func (p *Plugin) CfgString(path string) (string, error) {
	key := fmt.Sprintf("%s.%s.%s", p.Typ, p.Name, path)
	if res, ok := p.LocalCfg[key]; ok {
		return res, nil
	}
	if res, ok := p.LocalCfg[path]; ok {
		return res, nil
	}
	return "", errors.New("Could not find "+key)
}

func (p *Plugin) CfgStringOr(path, alt string) string {
	res, err := p.CfgString(path)
	if err != nil {
		return alt
	}
	return res
}

func (p *Plugin) CfgInt(path string) (int, error) {
	key := fmt.Sprintf("%s.%s.%s", p.Typ, p.Name, path)
	if res, ok := p.LocalCfg[key]; ok {
		return strconv.Atoi(res)
	}
	return 0, errors.New("Could not find "+key)
}

func (p *Plugin) CfgIntOr(path string, alt int) int {
	res, err := p.CfgInt(path)
	if err != nil {
		return alt
	}
	return res
}

func (p *Plugin) CfgBool(path string) (bool, error) {
	key := fmt.Sprintf("%s.%s.%s", p.Typ, p.Name, path)
	if res, ok := p.LocalCfg[key]; ok {
		switch res {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return false, errors.New(fmt.Sprintf("Key '%s' neither false not true, but %s: ", key, res))

		}
	}
	return false, errors.New("Could not find "+key)
}

func (p *Plugin) CfgBoolOr(path string, alt bool) bool {
	res, err := p.CfgBool(path)
	if err != nil {
		return alt
	}
	return res
}

func (p *Plugin) GetInputs() []string {
	inStr, err := p.CfgString("inputs")
	if err != nil {
		inStr = ""
	}
	return strings.Split(inStr, ",")
}

func (p *Plugin) GetCfgItems(key string) []string {
	inStr, err := p.CfgString(key)
	if err != nil {
		inStr = ""
	}
	return strings.Split(inStr, ",")
}

func (p *Plugin) Log(logLevel, msg string) {
	if len(p.LogOnlyPlugs) != 0 && ! qutils.IsItem(p.LogOnlyPlugs, p.Name) {
		return
	}
	// TODO: Setup in each Log() invocation seems rude
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	dL, _ := p.Cfg.StringOr("log.level", "info")
	dI := LogStrToInt(dL)
	lI := LogStrToInt(logLevel)
	lMsg := fmt.Sprintf("[%+6s] %15s Name:%-10s >> %s", strings.ToUpper(logLevel), p.Pkg, p.Name, msg)
	if lI == 0 {
		log.Panic(lMsg)
	} else if dI >= lI {
		log.Println(lMsg)
	}
}

func (p *Plugin) StartTicker(name string, durMs int) Ticker {
	p.Log("info", fmt.Sprintf("Start ticker '%s' with duration of %dms", name, durMs))
	ticker := NewTicker(name, durMs)
	go ticker.DispatchTicker(p.QChan)
	return ticker
}

func (p *Plugin) StopProcessingMessage(qm Message, allowEmptyInput bool) bool {
	p.MsgCount["received"]++
	if p.MyID == qm.SourceID {
		p.Log("debug", "Msg came from the same GID")
		p.MsgCount["loopDrop"]++
		return true
	}
	// TODO: Most likely invoked often, so check if performant enough
	inputs := p.GetInputs()
	if ! allowEmptyInput && len(inputs) == 0 {
		msg := fmt.Sprintf("Plugin '%s' does not allow empty imputs, please set '%s.%s.inputs'", p.Name, p.Typ, p.Name)
		log.Fatal(msg)
	}
	srcSuccess := p.CfgBoolOr("source-success", true)
	if ! qm.InputsMatch(inputs) {
		p.Log("debug", fmt.Sprintf("InputsMatch(%v) != %s | Msg: %s", inputs, qm.GetLastSource(), qm.Message))
		p.MsgCount["inputDrop"]++
		return true
	}
	if qm.SourceSuccess != srcSuccess {
		p.Log("debug", fmt.Sprintf("qm.SourceSuccess (%v) != (%v) srcSuccess", qm.SourceSuccess, srcSuccess))
		p.MsgCount["successDrop"]++
		return true
	}
	return false
}

func (p *Plugin) StopProcessingMetric(qm Metric, allowEmptyInput bool) bool {
	p.MsgCount["received"]++
	// TODO: Most likely invoked often, so check if performant enough
	inputs := p.GetInputs()
	if ! allowEmptyInput && len(inputs) == 0 {
		msg := fmt.Sprintf("Plugin '%s' does not allow empty imputs, please set '%s.%s.inputs'", p.Name, p.Typ, p.Name)
		log.Fatal(msg)
	}
	srcSuccess := p.CfgBoolOr("source-success", true)
	if ! qm.InputsMatch(inputs) {
		p.Log("debug", fmt.Sprintf("InputsMatch(%v) != %s", inputs, qm.GetLastSource()))
		p.MsgCount["inputDrop"]++
		return true
	}
	if qm.SourceSuccess != srcSuccess {
		p.Log("debug", fmt.Sprintf("qcs.SourceSuccess != srcSuccess (%v)", srcSuccess))
		p.MsgCount["successDrop"]++
		return true
	}
	return false
}

func (p *Plugin) StopProcessingCntEvent(ce ContainerEvent, allowEmptyInput bool) bool {
	p.MsgCount["received"]++
	// TODO: Most likely invoked often, so check if performant enough
	inputs := p.GetInputs()
	if ! allowEmptyInput && len(inputs) == 0 {
		msg := fmt.Sprintf("Plugin '%s' does not allow empty imputs, please set '%s.%s.inputs'", p.Name, p.Typ, p.Name)
		log.Fatal(msg)
	}
	if ! ce.InputsMatch(inputs) {
		p.Log("debug", fmt.Sprintf("InputsMatch(%v) = false", inputs))
		p.MsgCount["inputDrop"]++
		return true
	}
	srcSuccess := p.CfgBoolOr("source-success", true)
	if ce.SourceSuccess != srcSuccess {
		p.Log("debug", "ce.SourceSuccess != srcSuccess")
		p.MsgCount["successDrop"]++
		return true
	}
	return false
}

func (p *Plugin) StopProcessingCntStats(qcs ContainerStats, allowEmptyInput bool) bool {
	p.MsgCount["received"]++
	inputs := p.GetInputs()
	if ! allowEmptyInput && len(inputs) == 0 {
		msg := fmt.Sprintf("Plugin '%s' does not allow empty imputs, please set '%s.%s.inputs'", p.Name, p.Typ, p.Name)
		log.Fatal(msg)
	}
	if qcs.IsLastSource(p.Name) {
		p.Log("debug", "IsLastSource() = true")
		return true

	}
	if ! qcs.InputsMatch(inputs) {
		p.Log("debug", fmt.Sprintf("InputsMatch(%v) = false", inputs))
		return true

	}
	srcSuccess := p.CfgBoolOr("source-success", true)
	if qcs.SourceSuccess != srcSuccess {
		p.Log("debug", "qcs.SourceSuccess != srcSuccess")
		return true
	}
	return false
}

func (p *Plugin) DispatchMsgCount() {
	tickMs := p.CfgIntOr("count-ticker-ms", 5000)
	p.Log("info", fmt.Sprintf("Dispatch goroutine to send MsgCount every %dms", tickMs))
	ticker := time.NewTicker(time.Duration(tickMs)*time.Millisecond).C
	pre := map[string]float64{}
	for {
		tick := <-ticker
		pre = p.SendMsgCount(tick, pre)
	}
}

func (p *Plugin) SendMsgCount(tick time.Time, pre map[string]float64) map[string]float64 {
	dims := map[string]string{
		"plugin_name": p.Name,
		"plugin_version": p.Version,
		"plugin_type": p.Typ,
	}
	qm := NewExt(p.Name, "none", Counter, 0.0, dims, tick, false)
	for k,v := range p.MsgCount {
		if _, ok := pre[k]; !ok {
			pre[k] = v
		} else if pre[k] == v {
			continue
		}
		qm.Name = fmt.Sprintf("msg.%s", k)
		p.Log("debug", fmt.Sprintf("Send MsgCount %s=%f", qm.Name,v))
		qm.Value = float64(v)
		p.QChan.SendData(qm)
	}
	return pre
}
