package main

import (
	"log"
	"time"

	"github.com/zpatrick/go-config"
	"github.com/qframe/filter-grok"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/types/messages"
	"os"
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
)

func Run(qChan qtypes_qchannel.QChan, cfg *config.Config, name string) {
	p, _ := qfilter_grok.New(qChan, cfg, name)
	p.Run()
}

func main() {
	qChan := qtypes_qchannel.NewQChan()
	qChan.Broadcast()
	cfgMap := map[string]string{
		"log.level": "debug",
		"filter.grok.pattern": openTSDBCase.Pattern,
		"filter.grok.pattern-dir": "./patterns/",
		"filter.grok.inputs": "test",
	}
	cfg := config.NewConfig(
		[]config.Provider{
			config.NewStatic(cfgMap),
		},
	)
	p, err := qfilter_grok.New(qChan, cfg, "grok")
	if err != nil {
		log.Printf("[EE] Failed to create filter: %v", err)
		return
	}
	go p.Run()
	time.Sleep(time.Duration(100)*time.Millisecond)
	ticker := time.NewTicker(time.Millisecond*time.Duration(2000)).C
	bg := qChan.Data.Join()
	res := []string{}
	for _, c := range openTSDBCase.Tests {
		b := qtypes_messages.NewBase("test")
		qm := qtypes_messages.NewMessage(b, c.Str)
		log.Printf("Send message '%s", qm.Message)
		qChan.Data.Send(qm)
	}
	for {
		select {
		case val := <- bg.Read:
			switch val.(type) {
			case qtypes_messages.Message:
				qm := val.(qtypes_messages.Message)
				if ! qm.InputsMatch([]string{"grok"}) {
					continue
				}
				res = append(res, qm.GetLastSource())
				log.Printf("#### Received result from grok (pattern:%s) filter for input: %s || Tags:%v\n", p.GetPattern(), qm.Message, qm.Tags)
			}
			if len(res) == len(openTSDBCase.Tests) {
				break
			}
		case <- ticker:
			log.Println("Ticker came along, time's up...")
			os.Exit(1)
		}
		if len(res) == len(openTSDBCase.Tests) {
			break
		}
	}
	os.Exit(0)
}
