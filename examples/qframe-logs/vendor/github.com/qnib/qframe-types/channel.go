package qtypes

import (
	"os"
	"log"

	"github.com/grafov/bcast"
)

// QChan holds the broadcast channels to communicate
type QChan struct {
	Data *bcast.Group
	Back *bcast.Group
	Tick *bcast.Group
	Done chan os.Signal
}

// NewQChan create an instance of QChan
func NewQChan() QChan {
	return QChan{
		Data: bcast.NewGroup(), // create broadcast group
		Back: bcast.NewGroup(), // create broadcast group
		Tick: bcast.NewGroup(), // create broadcast group
		Done: make(chan os.Signal, 1),
	}
}

func (qc *QChan) Broadcast() {
	log.Println("[II] Dispatch broadcast for Back, Data and Tick")
	go qc.Data.Broadcast(0)
	go qc.Back.Broadcast(0)
	go qc.Tick.Broadcast(0)
}

func (qc *QChan) SendData(val interface{}) {
	qc.Data.Send(val)
}
