package qtypes_ticker


import (
	"time"
	"github.com/qframe/types/qchannel"

)

type Ticker struct {
	Name     			string
	Duration 			time.Duration
	DurationMs			int
	DurationTolerance 	float64
	ForceTick			bool
	Tick				time.Time
	LastTick			time.Time
}

func NewTicker(name string, durMs int) Ticker {
	return Ticker{
		Name: 		name,
		Duration: 	time.Duration(durMs)*time.Millisecond,
		DurationMs: durMs,
		DurationTolerance: 0.025,
		ForceTick:  false,
		LastTick: 	time.Now().AddDate(-1, 0, 0),

	}
}
func NewForceTicker(name string, durMs int) Ticker {
	t := NewTicker(name, durMs)
	t.ForceTick = true
	return t
}

func (t Ticker) DispatchTicker(qchan qtypes_qchannel.QChan) {
	ticker := time.NewTicker(t.Duration)
	for tick := range ticker.C {
		t.Tick = tick
		qchan.Tick.Send(t)
		t.LastTick = tick
	}
}

func (t *Ticker) SkipTick(lastTick time.Time) (time.Duration, bool) {
	diff := t.Tick.Sub(lastTick)
	msSince := diff.Seconds() + float64(diff.Nanoseconds()/1e3)
	// TODO: Substract DurationTolerance from Duration to tolerate slightly smaller differences
	msCompare := float64(t.DurationMs) * (1 - t.DurationTolerance)
	if msSince < msCompare {
		return diff, true
	}
	return diff, false
}

