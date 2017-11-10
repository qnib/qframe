package qtypes

import (
	"testing"
	//"github.com/stretchr/testify/assert"
	"time"
)

func BenchmarkQChan_SendData(t *testing.B) {
	qc := NewQChan()
	qc.Broadcast()
	now := time.Unix(1495270731, 0)
	//dc := qc.Data.Join()
	for i := 0; i < t.N; i++ {
		m := NewExt("source", "name", Gauge, float64(1.0), map[string]string{}, now, false)
		qc.SendData(m)
		//got := <- dc.Read
		//assert.Equal(t, m, got)

	}
}
