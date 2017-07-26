package statsq

import (
	"github.com/qnib/qframe-types"
)

type Packet struct {
	Bucket   string
	ValFlt   float64
	ValStr   string
	Modifier string
	Sampling float32
	Args     string
}

func NewStatsdPacketFromPacket(p *Packet) *qtypes.StatsdPacket {
	return &qtypes.StatsdPacket{
		Bucket:     p.Bucket,
		ValFlt:     p.ValFlt,
		ValStr:     p.ValStr,
		Modifier:   p.Modifier,
		Sampling:   p.Sampling,
		Dimensions: qtypes.NewDimensions(),
	}
}
