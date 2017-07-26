package qframe_handler_log

import (
	"fmt"
	"github.com/zpatrick/go-config"

	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-utils"
)

const (
	version   = "0.1.4"
	pluginTyp = "handler"
	pluginPkg = "log"
)

type Plugin struct {
	qtypes.Plugin
}

func New(qChan qtypes.QChan, cfg *config.Config, name string) (Plugin, error) {
	p := Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
	}
	p.Version = version
	p.Name = name
	return p, nil
}

// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	p.Log("info", fmt.Sprintf("Start log handler v%s", p.Version))
	bg := p.QChan.Data.Join()
	inputs := p.GetInputs()
	for {
		select {
		case val := <-bg.Read:
			switch val.(type) {
			case qtypes.QMsg:
				qm := val.(qtypes.QMsg)
				if len(inputs) != 0 && !qutils.IsInput(inputs, qm.Source) {
					continue
				}
				p.Log("info" , fmt.Sprintf("%-7s sType:%-6s sName:[%d]%-10s %s\n", qm.LogString(), qm.Type, qm.SourceID, qm.Source, qm.Msg))
			}
		}
	}
}
