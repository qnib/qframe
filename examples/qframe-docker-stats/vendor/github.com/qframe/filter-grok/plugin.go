package qfilter_grok

import (
	"fmt"
	"path/filepath"
	"os"
	"reflect"
	"strings"
	"sync"
	"github.com/vjeantet/grok"
	"github.com/zpatrick/go-config"

	"github.com/qframe/types/plugin"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/types/messages"
	"github.com/deckarep/golang-set"
)

const (
	version = "0.1.14"
	pluginTyp = "filter"
	pluginPkg = "grok"
)

type Plugin struct {
	*qtypes_plugin.Plugin
	mu 				sync.Mutex
	grok    		*grok.Grok
	pattern 		string
}

func (p *Plugin) GetOverwriteKeys() []string {
	inStr, err := p.CfgString("overwrite-keys")
	if err != nil {
		inStr = ""
	}
	return strings.Split(inStr, ",")
}

// New creates a new instance of the Plugin.
func New(qChan qtypes_qchannel.QChan, cfg *config.Config, name string) (p Plugin, err error) {
	p = Plugin{
		Plugin: qtypes_plugin.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg,  name, version),
	}
	return p, err
}

// Match matches the pattern against a string and returns the values extracted.
func (p *Plugin) Match(str string) (map[string]string, bool) {
	match := true
	val, _ := p.grok.Parse(p.pattern, str)
	keys := reflect.ValueOf(val).MapKeys()
	if len(keys) == 0 {
		match = false
	}
	if match {
		p.Log("trace", fmt.Sprintf("Pattern '%s' matched '%s'", p.pattern, str))
	} else {
		p.Log("trace", fmt.Sprintf("Pattern '%s' DID NOT match '%s'", p.pattern, str))
	}
	return val, match
}

// GetPattern returns the pattern used.
func (p *Plugin) GetPattern() string {
	return p.pattern
}

// InitGrok() kicks of grok and adds patterns to the grok-instance.
func (p *Plugin) InitGrok() {
	p.grok, _ = grok.New()
	var err error
	p.pattern, err = p.CfgString("pattern")
	if err != nil {
		p.Log("fatal", "Could not find pattern in config")
	}
	pFileSet := mapset.NewSet()
	pDir, err := p.CfgString("pattern-dir")
	if err == nil && pDir != "" {
		err := filepath.Walk(pDir, func(path string, f os.FileInfo, err error) error {
			if ! f.IsDir() {
				pFileSet.Add(path)
			}
			return nil
		})
		if err != nil {
			p.Log("error", err.Error())
		}
	}
	pFiles, err := p.CfgString("pattern-files")
	for _, f := range strings.Split(pFiles, ",") {
		if f == "" {
			continue
		}
		pFileSet.Add(f)
	}
	for f := range pFileSet.Iterator().C {
		p.Log("trace", fmt.Sprintf("Iterate %s", f))
		err := p.grok.AddPatternsFromPath(f.(string))
		if err != nil {
			p.Log("error", err.Error())
		} else {
			p.Log("info", fmt.Sprintf("Added pattern-file '%s'", f))
		}
	}
}

// Lock locks the plugins' mutex.
func (p *Plugin) Lock() {
	p.mu.Lock()
}

// Unlock unlocks the plugins' mutex.
func (p *Plugin) Unlock() {
	p.mu.Unlock()
}

// Run fetches everything from the Data channel and flushes it to stdout.
func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start grok filter v%s", p.Version))
	p.InitGrok()
	ignoreContainerEvents := p.CfgBoolOr("ignore-container-events", true)
	bg := p.QChan.Data.Join()
	msgKey := p.CfgStringOr("overwrite-message-key", "")
	for {
		val := bg.Recv()
		p.Log("trace", fmt.Sprintf("received %s", reflect.TypeOf(val)))
		switch val.(type) {
		case qtypes_messages.ContainerMessage:
			cm := val.(qtypes_messages.ContainerMessage)
			if cm.StopProcessing(p.Plugin, false) {
				continue
			}
			cm.AppendSource(p.Name)
			var kv map[string]string
			kv, cm.SourceSuccess = p.Match(cm.Message.Message)
			if cm.SourceSuccess {
				p.Log("debug", fmt.Sprintf("Matched pattern '%s': ", p.pattern))
				for k,v := range kv {
					p.Log("debug", fmt.Sprintf("    %15s: %s", k, v))
					cm.Tags[k] = v
					if msgKey == k {
						cm.Message.Message = v
					}
				}
			} else {
				p.Log("debug", fmt.Sprintf("No match of '%s' for message '%s'", p.pattern, cm.Message.Message))
			}
			p.QChan.Data.Send(cm)
		case qtypes_messages.Message:
			qm := val.(qtypes_messages.Message)
			if qm.StopProcessing(p.Plugin, false) {
				continue
			}
			qm.AppendSource(p.Name)
			kv, SourceSuccess := p.Match(qm.Message)
			qm.SourceSuccess = SourceSuccess
			if qm.SourceSuccess {
				p.Log("debug", fmt.Sprintf("Matched pattern '%s'", p.pattern))
				for k,v := range kv {
					p.Log("trace", fmt.Sprintf("    %15s: %s", k,v ))
					qm.Tags[k] = v
					if msgKey == k {
						qm.Message = v
					}
				}
				p.QChan.Data.Send(qm)
				continue
			} else {
				p.Log("debug", fmt.Sprintf("No match of '%s' for message '%s'", p.pattern, qm.Message))
			}
			p.QChan.Data.Send(qm)
		default:
			if ignoreContainerEvents {
				continue
			}
			p.Log("trace", fmt.Sprintf("No match for type '%s'", reflect.TypeOf(val)))
		}
	}
}
