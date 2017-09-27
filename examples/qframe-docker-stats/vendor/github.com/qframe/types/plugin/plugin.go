package qtypes_plugin

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"github.com/pkg/errors"
	"github.com/zpatrick/go-config"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/types/helper"
	"github.com/qframe/types/ticker"
)

type Plugin struct {
	Base
	MyID		int
	Typ			string
	Pkg			string
	Version 	string
	Name 		string
}


func NewNamedPlugin(qChan qtypes_qchannel.QChan, cfg *config.Config, typ, pkg, name, version string) *Plugin {
	b := NewBase(qChan, cfg)
	return NewBasePlugin(b, typ, pkg, name, version)
}


func NewBasePlugin(b Base, typ, pkg, name, version string) *Plugin {
	p := &Plugin{
		Base: b,
		Typ:   		typ,
		Pkg:  		pkg,
		Version:	version,
		Name: 		name,
	}
	return p
}

func (p *Plugin) GetInfo() (typ,pkg,name string) {
	return p.Typ, p.Pkg, p.Name
}

func (p *Plugin) GetLogOnlyPlugs() []string {
	return p.LogOnlyPlugs
}

// Config

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

func (p *Plugin) GetInputs() (res []string) {
	inStr, err := p.CfgString("inputs")
	if err == nil {
		res = strings.Split(inStr, ",")
	}
	return res
}

func (p *Plugin) GetCfgItems(key string) []string {
	inStr, err := p.CfgString(key)
	if err != nil {
		inStr = ""
	}
	return strings.Split(inStr, ",")
}

// Log
func (p *Plugin) Log(logLevel, msg string) {
	if len(p.LogOnlyPlugs) != 0 && ! qtypes_helper.IsItem(p.LogOnlyPlugs, p.Name) {
		return
	}
	// TODO: Setup in each Log() invocation seems rude
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	dL, _ := p.Cfg.StringOr("log.level", "info")
	dI := qtypes_helper.LogStrToInt(dL)
	lI := qtypes_helper.LogStrToInt(logLevel)
	lMsg := fmt.Sprintf("[%+6s] %15s Name:%-10s >> %s", strings.ToUpper(logLevel), p.Pkg, p.Name, msg)
	if lI == 0 {
		log.Panic(lMsg)
	} else if dI >= lI {
		log.Println(lMsg)
	}
}

// Ticker
func (p *Plugin) StartTicker(name string, durMs int) qtypes_ticker.Ticker {
	p.Log("info", fmt.Sprintf("Start ticker '%s' with duration of %dms", name, durMs))
	ticker := qtypes_ticker.NewTicker(name, durMs)
	go ticker.DispatchTicker(p.QChan)
	return ticker
}
