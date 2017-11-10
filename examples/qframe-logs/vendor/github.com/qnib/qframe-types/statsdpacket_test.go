package qtypes

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewStatsdPacket(t *testing.T) {
	sp := NewStatsdPacket("testCounter", "1.2", "c")
	assert.Equal(t,"testCounter", sp.Bucket)
}

func TestStatsdPacket_String(t *testing.T) {
	sp := NewStatsdPacket("testCounter", "1.2", "c")
	assert.Equal(t,"", sp.DimensionString())
	sp.AddDimension("key1", "val1")
	assert.Equal(t,"key1=val1", sp.DimensionString())
	sp.AddDimension("key2", "val2")
	assert.Equal(t,"key1=val1,key2=val2", sp.DimensionString())
}

func TestStatsdPacket_GetBucketKey(t *testing.T) {
	sp := NewStatsdPacket("testCounter", "1.2", "c")
	got := sp.GetBucketKey()
	assert.Equal(t, "testCounter", got)
}

func TestNewStatsdPacketDims(t *testing.T) {
	dims := NewDimensions()
	dims.Add("key1", "val1")
	sp := NewStatsdPacketDims("testCounter", "1.2", "c", dims)
	assert.Equal(t,"testCounter", sp.Bucket)
	assert.Equal(t,"key1=val1", sp.DimensionString())
}

func TestStatsdPacket_GenerateID(t *testing.T) {
	dims := NewDimensionsPre(map[string]string{"key1": "val1"})
	sp := NewStatsdPacketDims("bucketName", "1.2", "c", dims)
	assert.Equal(t, "2af96db5523ec73ecccef75192990635df2067b5", sp.GenerateID())
	dims.Add("key2", "val2")
	sp = NewStatsdPacketDims("bucketName", "1.2", "c", dims)
	assert.Equal(t, "23fe036dd9a06100c8056bc26a27f07ce9b601d4", sp.GenerateID())
}

func TestParseValStr(t *testing.T) {
	str, val, err := ParseValStr("10")
	assert.NoError(t, err, "Should work")
	assert.Equal(t, "", str)
	assert.Equal(t, val, float64(10))
	str, val, err = ParseValStr("-12")
	assert.NoError(t, err, "Should work")
	assert.Equal(t, "-", str)
	assert.Equal(t, val, float64(12))
	str, val, err = ParseValStr("+11")
	assert.NoError(t, err, "Should work")
	assert.Equal(t, "+", str)
	assert.Equal(t, val, float64(11))
	str, val, err = ParseValStr("-g")
	assert.Error(t, err, "Should not work")
}