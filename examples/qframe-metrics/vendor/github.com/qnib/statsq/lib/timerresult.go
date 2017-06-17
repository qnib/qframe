package statsq

import (
	"fmt"
	"strings"
)

type TimerResult struct {
	TestsGot map[string]string
	TestsExp map[string]float64
}

func NewTimerResult(exp map[string]float64) TimerResult {
	tr := TimerResult{
		TestsExp: exp,
		TestsGot: map[string]string{},
	}
	for k := range exp {
		tr.TestsGot[k] = "missed"
	}
	return tr
}

func (tr *TimerResult) Input(name string, value float64) {
	if v, ok := tr.TestsExp[name]; ok {
		if value == v {
			tr.TestsGot[name] = "ok"
		} else {
			tr.TestsGot[name] = fmt.Sprintf("%v!=%v", v, value)
		}
	}
}

func (tr *TimerResult) Check() bool {
	nok := false
	for _, v := range tr.TestsGot {
		if v != "ok" {
			nok = true
		}
	}
	return !nok
}

func (tr *TimerResult) Result() string {
	var res []string
	for k, v := range tr.TestsGot {
		res = append(res, fmt.Sprintf("%-20s: %f : %s", k, tr.TestsExp[k], v))
	}
	return strings.Join(res, "\n")
}
