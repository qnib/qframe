package qframe_filter_grok

import (
	"C"
	"fmt"
	"os"
	"reflect"
	"strings"
	"github.com/vjeantet/grok"
	"github.com/zpatrick/go-config"

	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-utils"
)

const (
	version = "0.1.9"
	pluginTyp = "filter"
	pluginPkg = "grok"
	defPatternDir = "/etc/grok-patterns"
)

type Plugin struct {
	qtypes.Plugin
	grok    *grok.Grok
	pattern string
}

func (p *Plugin) GetOverwriteKeys() []string {
	inStr, err := p.CfgString("overwrite-keys")
	if err != nil {
		inStr = ""
	}
	return strings.Split(inStr, ",")
}

func New(qChan qtypes.QChan, cfg config.Config, name string) (p Plugin, err error) {
	p = Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg,  name, version),
	}
	p.grok, _ = grok.New()
	p.pattern, err = p.CfgString("pattern")
	if err != nil {
		p.Log("error", "Could not find pattern in config")
		return p, err
	}
	pDir, err := p.CfgString("pattern-dir")
	if err != nil {
		if _, err := os.Stat(defPatternDir); err == nil {
			pDir = defPatternDir
			p.Log("info", fmt.Sprintf("Add patterns from DEFAULT directory '%s'", pDir))
		}
	} else {
		p.Log("info", fmt.Sprintf("Add patterns from directory '%s'", pDir))
	}
	if _, err := os.Stat(pDir); err != nil {
		p.Log("error", fmt.Sprintf("Patterns directory does not exist '%s'", pDir))
	} else {
		p.grok.AddPatternsFromPath(pDir)
	}
	return p, err
}

func (p *Plugin) Match(str string) (map[string]string, bool) {
	match := true
	val, _ := p.grok.Parse(p.pattern, str)
	keys := reflect.ValueOf(val).MapKeys()
	if len(keys) == 0 {
		match = false
	}
	return val, match
}

func (p *Plugin) GetPattern() string {
	return p.pattern
}

// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start grok filter v%s", p.Version))
	myId := qutils.GetGID()
	bg := p.QChan.Data.Join()
	inputs := p.GetInputs()
	srcSuccess := p.CfgBoolOr("source-success", true)
	msgKey := p.CfgStringOr("overwrite-message-key", "")
	for {
		val := bg.Recv()
		switch val.(type) {
		case qtypes.QMsg:
			qm := val.(qtypes.QMsg)
			if qm.SourceID == myId {
				continue
			}
			if len(inputs) != 0 && !qutils.IsInput(inputs, qm.Source) {
				continue
			}
			if qm.SourceSuccess != srcSuccess {
				continue
			}
			qm.Type = "filter"
			qm.Source = p.Name
			qm.SourceID = myId
			qm.SourcePath = append(qm.SourcePath, p.Name)
			qm.KV, qm.SourceSuccess = p.Match(qm.Msg)
			p.QChan.Data.Send(qm)
		case qtypes.Message:
			qm := val.(qtypes.Message)
			if ! qm.InputsMatch(inputs) {
				continue
			}
			if qm.SourceSuccess != srcSuccess {
				continue
			}
			if qm.IsLastSource(p.Name) {
				p.Log("debug", "IsLastSource() = true")
				continue
			}
			if len(inputs) != 0 && ! qm.InputsMatch(inputs) {
				p.Log("debug", fmt.Sprintf("InputsMatch(%v) = false", inputs))
				continue
			}
			if qm.SourceSuccess != srcSuccess {
				p.Log("debug", "qcs.SourceSuccess != srcSuccess")
				continue
			}
			qm.AppendSource(p.Name)
			var kv map[string]string
			kv, qm.SourceSuccess = p.Match(qm.Message)
			if qm.SourceSuccess {
				p.Log("debug", fmt.Sprintf("Matched pattern '%s'", p.pattern))
				for k,v := range kv {
					p.Log("debug", fmt.Sprintf("    %15s: %s", k,v ))
					qm.KV[k] = v
					if msgKey == k {
						qm.Message = v
					}
				}
			} else {
				p.Log("debug", fmt.Sprintf("No match of '%s' for message '%s'", p.pattern, qm.Message))
			}
			p.QChan.Data.Send(qm)
		}
	}
}
