package qfilter_grok

import (
	"fmt"
	"testing"
	"time"
	"log"
	"github.com/zpatrick/go-config"

	"github.com/stretchr/testify/assert"
	"reflect"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/types/messages"
)

type TestCase struct {
	Pattern string
	Tests []Test
}

type Test struct {
	Str string
	Success bool
	result map[string]string
}

var (
	ip4_1 = Test{"192.168.0.1", true, map[string]string{"ip":"192.168.0.1"}}
	ip4_2 = Test{"127.0.0.1", true, map[string]string{"ip":"127.0.0.1"}}
	ip4_FAIL2 = Test{Str: "abc.0.0.1", Success: false}
	ip6_1 = Test{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", true, map[string]string{"ip":"2001:0db8:85a3:0000:0000:8a2e:0370:7334"}}
	ipCase = TestCase{
		"%{IP:ip}",
		[]Test{ip4_1, ip4_2,	ip6_1,ip4_FAIL2},
	}
	openTSDB1 = Test{"put sys.cpu.load 1505290630 40", true, map[string]string{"key":"sys.cpu.load", "time":"1505290630", "value":"40"}}
	openTSDB2 = Test{
		"put sys.cpu.load 1505290630 40 host=webserver1 cpu=1", true,
		map[string]string{"key":"sys.cpu.load", "time":"1505290630", "value":"40", "tags":"host=webserver1 cpu=1"}}
	openTSDB3 = Test{
		"sys.cpu.load 1505290630 40 host=webserver1 cpu=1", true,
		map[string]string{"key":"sys.cpu.load", "time":"1505290630", "value":"40", "tags":"host=webserver1 cpu=1"}}
	openTSDBCase = TestCase{
		"%{OPENTSDB}",
		[]Test{openTSDB1,openTSDB2,openTSDB3},
	}
	syslog1 = Test{"<134>1 2017-09-13T18:35:47.443Z host - - event - Hello World", true, map[string]string{"syslog5424_msg":`Hello World`}}
	ceeRes1 = map[string]string{
		"syslog5424_pri": "134",
		"syslog5424_ver": "1",
		"syslog5424_ts": "2017-09-13T18:35:47.443Z",
		"syslog5424_host": "host",
		"syslog5424_app": "app",
		"syslog5424_proc": "proc",
		"syslog5424_sd": "",
		"syslog5424_msg":`@cee:{"time":"2017-09-13T18:35:47.443Z"}`,

	}
	ceeLog1 = Test{`<134>1 2017-09-13T18:35:47.443Z host app proc msgid - @cee:{"time":"2017-09-13T18:35:47.443Z"}`, true, ceeRes1}
	ceeResWithSD = map[string]string{
		"syslog5424_pri": "134",
		"syslog5424_ver": "1",
		"syslog5424_ts": "2017-09-13T18:35:47.443Z",
		"syslog5424_host": "host",
		"syslog5424_app": "app",
		"syslog5424_proc": "proc",
		"syslog5424_sd": "[key=val]",
		"syslog5424_msg":`Hello World`,
	}
	ceeLogWithSD = Test{`<134>1 2017-09-13T18:35:47.443Z host app proc msgid [key=val] Hello World`, true, ceeResWithSD}
	syslogCase = TestCase{
		"%{QNIB_SYSLOG5424LINE}",
		[]Test{syslog1,ceeLog1,ceeLogWithSD},
	}
	cfgMap = map[string]string{
		"log.level": "trace",
		"filter.grok.inputs": "test",
		"filter.grok.pattern-files": "./resources/patterns/opentsdb,./resources/patterns/linux-syslog",
	}
)

func NewCfgMap(key, val string) (map[string]string) {
	res := map[string]string{}
	for k, v := range cfgMap {
		res[k] = v
	}
	res[key] = val
	return res
}

/******* Helper */
func RunTest(t *testing.T,tCase TestCase) {
	cfgMap1 := NewCfgMap("filter.grok.pattern", tCase.Pattern)
	cfg := config.NewConfig([]config.Provider{config.NewStatic(cfgMap1)})
	qChan := qtypes_qchannel.NewCfgQChan(cfg)
	// simple
	p, err := New(qChan, cfg, "grok")
	assert.NoError(t, err)
	p.InitGrok()
	for _, c := range tCase.Tests {
		got, ok := p.Match(c.Str)
		assert.Equal(t, ok, c.Success)
		if !ok {
			continue
		}
		for k,v := range c.result {
			val, ok := got[k]
			assert.True(t, ok, fmt.Sprintf("key %s could not be found in result", k))
			m := v == val
			assert.True(t, m, fmt.Sprintf("'%s'!='%s'", v, val))
		}
	}
}

func Receive(qchan qtypes_qchannel.QChan, source string, endCnt int) {
	bg := qchan.Data.Join()
	allCnt := 1
	cnt := 1
	for {
		select {
		case val := <-bg.Read:
			allCnt++
			switch val.(type) {
			case qtypes_messages.Message:
				qm := val.(qtypes_messages.Message)
				if qm.IsLastSource(source) {
					cnt++
				}
			default:
				fmt.Printf("Dunno received msg %d: type=%s\n", allCnt, reflect.TypeOf(val))

			}
		}
		if endCnt == cnt {
			qchan.Data.Send(cnt)
			break
		}
	}
}

/******* Tests */
func TestPlugin_MatchIP(t *testing.T) {
	RunTest(t, ipCase)
}

func TestPlugin_MatchOpenTSDB(t *testing.T) {
	RunTest(t, openTSDBCase)
}

func TestPlugin_MatchSyslog5424(t *testing.T) {
	RunTest(t, syslogCase)
}

func TestPlugin_MatchInt(t *testing.T) {
	cfgMap := map[string]string{
		"log.level": "error",
		"filter.grok.pattern": "%{INT:number}",
		"filter.grok.inputs": "test",
		"filter.grok.pattern-dir": "./resources/patterns/",
	}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(cfgMap)})
	qChan := qtypes_qchannel.NewCfgQChan(cfg)
	// simple
	p, err := New(qChan, cfg, "grok")
	assert.NoError(t, err)
	p.InitGrok()
	got, ok := p.Match("test1")
	assert.True(t, ok, "test1 should match pattern")
	assert.Equal(t, map[string]string{"number": "1"}, got)
	cfgMap["filter.grok.pattern"] = "test%{INT:number} %{WORD:str}"
	cfg = config.NewConfig([]config.Provider{config.NewStatic(cfgMap)})
	p1, err := New(qChan, cfg, "grok")
	assert.NoError(t, err)
	p1.InitGrok()
	g1, ok1 := p1.Match("test1 sometext")
	assert.True(t, ok1, "test1 should match pattern")
	assert.Equal(t, map[string]string{"number": "1", "str": "sometext"}, g1)
	g2, ok2 := p1.Match("testsometext")
	assert.False(t, ok2, "should NOT match pattern")
	assert.Equal(t, map[string]string{}, g2)
}

func TestPlugin_GetOverwriteKeys(t *testing.T) {
	cfgMap := map[string]string{}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(cfgMap)})
	qChan := qtypes_qchannel.NewCfgQChan(cfg)
	p, err := New(qChan, cfg, "grok")
	assert.NoError(t, err)
	assert.Equal(t, []string{""}, p.GetOverwriteKeys())
	cfgMap["filter.grok.overwrite-keys"] = "msg"
	cfg = config.NewConfig([]config.Provider{config.NewStatic(cfgMap)})
	p, err = New(qChan, cfg, "grok")
	assert.NoError(t, err)
	assert.Equal(t, []string{"msg"}, p.GetOverwriteKeys())

}

func TestPlugin_GetPattern(t *testing.T) {
	cfgMap := map[string]string{
		"log.level": "error",
		"filter.grok.pattern-dir": "./resources/patterns/",
		"filter.grok.pattern": "test%{INT:number}",
	}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(cfgMap)})
	qChan := qtypes_qchannel.NewCfgQChan(cfg)
	p, err := New(qChan, cfg, "grok")
	assert.NoError(t, err)
	p.InitGrok()
	assert.Equal(t, "test%{INT:number}", p.GetPattern())
}

func TestPlugin_Run(t *testing.T) {
	endCnt := 2
	cfgMap := map[string]string{
		"log.level": "error",
		"filter.grok.pattern": "test%{INT:number}",
		"filter.grok.inputs": "test",
	}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(cfgMap)})
	qChan := qtypes_qchannel.NewCfgQChan(cfg)
	qChan.Broadcast()
	go Receive(qChan, "grok", endCnt)
	p, err := New(qChan, cfg, "grok")
	if err != nil {
		log.Printf("[EE] Failed to create filter: %v", err)
		return
	}
	dc := qChan.Data.Join()
	go p.Run()
	time.Sleep(time.Duration(50)*time.Millisecond)
	p.Log("info", fmt.Sprintf("Benchmark sends %d messages to grok", endCnt))
	bMsg := qtypes_messages.NewBase("test")
	qm := qtypes_messages.NewMessage(bMsg, "testMsg")
	for i := 1; i <= endCnt; i++ {
		msg := fmt.Sprintf("test%d", i)
		qm.Message = msg
		qChan.Data.Send(qm)
	}
	done := false
	for {
		select {
		case val := <- dc.Read:
			switch val.(type) {
			case int:
				vali := val.(int)
				assert.Equal(t, endCnt, vali)
				done = true
			}
		case <-time.After(5 * time.Second):
			t.Fatal("metrics receive timeout")
		}
		if done {
			break
		}
	}
}

/******* Benchmarks */
func BenchmarkGrokINT(b *testing.B) {
	endCnt := b.N
	cfgMap := map[string]string{
		"log.level": "error",
		"filter.grok.pattern": "test%{INT:number}",
		"filter.grok.inputs": "test",
	}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(cfgMap)})
	qChan := qtypes_qchannel.NewCfgQChan(cfg)
	qChan.Broadcast()
	go Receive(qChan, "grok", endCnt)
	p, err := New(qChan, cfg, "grok")
	if err != nil {
		log.Printf("[EE] Failed to create filter: %v", err)
		return
	}
	dc := qChan.Data.Join()
	go p.Run()
	time.Sleep(time.Duration(50)*time.Millisecond)
	p.Log("info", fmt.Sprintf("Benchmark sends %d messages to grok", endCnt))
	bMsg := qtypes_messages.NewBase("test")
	qm := qtypes_messages.NewMessage(bMsg, "testMsg")
	for i := 1; i <= endCnt; i++ {
		msg := fmt.Sprintf("test%d", i)
		qm.Message = msg
		qChan.Data.Send(qm)
	}
	done := false
	for {
		select {
		case val := <- dc.Read:
			switch val.(type) {
			case int:
				vali := val.(int)
				assert.Equal(b, endCnt, vali)
				done = true
			}
		case <-time.After(5 * time.Second):
				b.Fatal("metrics receive timeout")
		}
		if done {
			break
		}
	}
}

func BenchmarkGrokOpenTSDB(b *testing.B) {
	endCnt := b.N
	cfgMap := map[string]string{
		"log.level": "error",
		"filter.grok.pattern": "%{OPENTSDB}",
		"filter.grok.inputs": "test",
		"filter.grok.pattern-files": "./resources/patterns/opentsdb",
	}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(cfgMap)})
	qChan := qtypes_qchannel.NewCfgQChan(cfg)
	qChan.Broadcast()
	go Receive(qChan, "grok", endCnt)
	p, err := New(qChan, cfg, "grok")
	if err != nil {
		log.Printf("[EE] Failed to create filter: %v", err)
		return
	}
	dc := qChan.Data.Join()
	go p.Run()
	time.Sleep(time.Duration(50)*time.Millisecond)
	p.Log("info", fmt.Sprintf("Benchmark sends %d messages to grok", endCnt))
	bMsg := qtypes_messages.NewBase("test")
	qm := qtypes_messages.NewMessage(bMsg, "testMsg")
	for i := 1; i <= endCnt; i++ {
		msg := fmt.Sprintf("put sys.cpu.load 1505290630 %d host=webserver1 iteration=%d", i, i)
		qm.Message = msg
		qChan.Data.Send(qm)
	}
	done := false
	for {
		select {
		case val := <- dc.Read:
			switch val.(type) {
			case int:
				vali := val.(int)
				assert.Equal(b, endCnt, vali)
				done = true
			}
		case <-time.After(5 * time.Second):
			b.Fatal("metrics receive timeout")
		}
		if done {
			break
		}
	}
}