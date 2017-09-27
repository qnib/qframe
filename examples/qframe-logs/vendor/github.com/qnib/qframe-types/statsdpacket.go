package qtypes

import (
	"strings"
	"crypto/sha1"
	"encoding/hex"
	"strconv"
)

type StatsdPacket struct {
	Bucket   	string
	ValFlt   	float64
	ValStr   	string
	Modifier 	string
	Sampling 	float32
	Dimensions	Dimensions
}

type Packet struct {
	Bucket   string
	ValFlt   float64
	ValStr   string
	Modifier string
	Sampling float32
}

func (sd *StatsdPacket) GetBucketKey() string {
	return sd.Bucket
}

func (sd *StatsdPacket) GetDims() map[string]string {
	return sd.Dimensions.GetDims()
}

func (sd *StatsdPacket) GenerateID() string {
	idRaw := []string{sd.Bucket}
	for key, val := range sd.Dimensions.Map {
		idRaw = append(idRaw, key+"="+val)
	}
	s := strings.Join(idRaw, "_")
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}


func NewStatsdPacket(bucket string, val string, modifier string) *StatsdPacket {
	return NewStatsdPacketDims(bucket, val, modifier, NewDimensions())
}

func NewStatsdPacketDims(bucket string, val string, modifier string, dims Dimensions) *StatsdPacket {
	vStr, vFlt, _ := ParseValStr(val)
	return &StatsdPacket{
		Bucket: bucket,
		ValFlt: vFlt,
		ValStr: vStr,
		Modifier: modifier,
		Sampling: float32(1),
		Dimensions: dims,
	}
}

func ParseValStr(s string) (str string, flt float64, err error) {
	if strings.HasPrefix(s, "+") {
		str = "+"
		s = s[1:]
	} else if strings.HasPrefix(s, "-") {
		str = "-"
		s = s[1:]
	}
	flt, err = strconv.ParseFloat(s, 64)
	if err != nil {
		return
	}
	return
}

func (sd *StatsdPacket) AddDimension(key, val string) {
	sd.Dimensions.Add(key, val)
}

func (sd *StatsdPacket) DimensionString() string {
	return sd.Dimensions.String()
}
