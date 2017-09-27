package qtypes_qchannel


import (
	"log"

	"github.com/grafov/bcast"
	"github.com/zpatrick/go-config"
	"fmt"
	"strings"
	"github.com/qframe/types/helper"
)

// QChan holds the broadcast channels to communicate
type QChan struct {
	Data *bcast.Group
	Done *bcast.Group
	Tick *bcast.Group
	Cfg *config.Config
}

func (qc *QChan) Log(logLevel, msg string) {
	// TODO: Setup in each Log() invocation seems rude
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	dL, _ := qc.Cfg.StringOr("log.level", "info")
	dI := qtypes_helper.LogStrToInt(dL)
	lI := qtypes_helper.LogStrToInt(logLevel)
	lMsg := fmt.Sprintf("[%+6s] %s", strings.ToUpper(logLevel), msg)
	if lI == 0 {
		log.Panic(lMsg)
	} else if dI >= lI {
		log.Println(lMsg)
	}
}

// NewQChan create an instance of QChan
func NewQChan() QChan {
	kv := map[string]string{"log.level": "info"}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(kv)})
	return NewCfgQChan(cfg)
}

// NewQChan create an instance of QChan
func NewCfgQChan(cfg *config.Config) QChan {
	return QChan{
		Data: bcast.NewGroup(), // create broadcast group
		Done: bcast.NewGroup(), // create broadcast group
		Tick: bcast.NewGroup(), // create broadcast group
		Cfg: cfg,
	}
}


func (qc *QChan) Broadcast() {
	qc.Log("info", "Dispatch broadcast for Data, Done and Tick")
	go qc.Data.Broadcast(0)
	go qc.Done.Broadcast(0)
	go qc.Tick.Broadcast(0)
}

func (qc *QChan) SendData(val interface{}) {
	qc.Data.Send(val)
}
