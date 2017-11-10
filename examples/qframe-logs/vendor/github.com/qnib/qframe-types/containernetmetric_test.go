package qtypes

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/docker/docker/api/types"
)


func TestNetMetric_Aggregate(t *testing.T) {
	m1 := NetStats{
		Base: 	Base{
			BaseVersion: 	"0.0.0",
			Time: 			time.Now(),
			SourceID: 		0,
			SourcePath:		[]string{"src1"},
			SourceSuccess: 	false,
		},
		Container: 		&types.Container{},
		NameInterface: 	"eth0",
		RxBytes:    	float64(1.0),
		RxDropped:  	float64(1.0),
		RxErrors:  	 	float64(1.0),
		RxPackets:  	float64(1.0),
		TxBytes:    	float64(1.0),
		TxDropped:  	float64(1.0),
		TxErrors:   	float64(1.0),
		TxPackets:  	float64(1.0),
	}
	m2 := NetStats{
		Base: 	Base{
			BaseVersion: 	"0.0.0",
			Time: 			time.Now(),
			SourceID: 		0,
			SourcePath:		[]string{"src1"},
			SourceSuccess: 	false,
		},
		Container: 		&types.Container{},
		NameInterface: 	"eth1",
		RxBytes:    	float64(2.0),
		RxDropped:  	float64(2.0),
		RxErrors:  	 	float64(2.0),
		RxPackets:  	float64(2.0),
		TxBytes:    	float64(2.0),
		TxDropped:  	float64(2.0),
		TxErrors:   	float64(2.0),
		TxPackets:  	float64(2.0),
	}
	a := AggregateNetStats("total", m1, m2)
	assert.Equal(t, "total", a.NameInterface)
}
