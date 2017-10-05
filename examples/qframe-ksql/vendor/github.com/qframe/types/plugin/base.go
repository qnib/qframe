package qtypes_plugin

import (
	"github.com/zpatrick/go-config"
	"github.com/qframe/types/qchannel"
)

const (
	version = "0.1.8"
)

type Base struct {
	BaseVersion 	string
	Version 		string
	QChan 			qtypes_qchannel.QChan
	ErrChan			chan error
	Cfg 			*config.Config
	MsgCount		map[string]float64


}

func NewBase(qChan qtypes_qchannel.QChan, cfg *config.Config) Base {
	b := Base{
		BaseVersion: version,
		QChan: qChan,
		ErrChan: make(chan error),
		Cfg: cfg,
		MsgCount:       map[string]float64{
			"received": 0.0,
			"loopDrop": 0.0,
			"inputDrop": 0.0,
			"successDrop": 0.0,
		},
	}
	return b
}


/*
func (p *Base) DispatchMsgCount() {

	tickMs := p.CfgIntOr("count-ticker-ms", 5000)
	p.Log("info", fmt.Sprintf("Dispatch goroutine to send MsgCount every %dms", tickMs))
	ticker := time.NewTicker(time.Duration(tickMs)*time.Millisecond).C
	pre := map[string]float64{}
	for {
		tick := <-ticker
		pre = p.SendMsgCount(tick, pre)
	}
}

func (p *Base) SendMsgCount(tick time.Time, pre map[string]float64) map[string]float64 {
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
*/