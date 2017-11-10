package qtypes

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDimensions_String(t *testing.T) {
	d := NewDimensions()
	assert.Equal(t,"", d.String())
	d.Add("key1", "val1")
	assert.Equal(t,"key1=val1", d.String())
	d.Add("key2", "val2")
	assert.Equal(t,"key1=val1,key2=val2", d.String())
}

func TestNewDimensionsPre(t *testing.T) {
	d := NewDimensionsPre(map[string]string{"key1": "val1"})
	assert.Equal(t,"key1=val1", d.String())
}

func TestNewDimensionsFromString(t *testing.T) {
	d := NewDimensionsFromString("")
	assert.Equal(t,"", d.String())
	d = NewDimensionsFromString("key1=val1")
	assert.Equal(t,"key1=val1", d.String())
	d = NewDimensionsFromString("key1=val1,val2=val2")
	assert.Equal(t,"key1=val1,val2=val2", d.String())
}

func TestNewDimensionsFromBytes(t *testing.T) {
	d := NewDimensionsFromBytes([]byte(""))
	assert.Equal(t,"", d.String())
	d = NewDimensionsFromBytes([]byte("key1=val1"))
	assert.Equal(t,"key1=val1", d.String())
	d = NewDimensionsFromBytes([]byte("key1=val1,val2=val2"))
	assert.Equal(t,"key1=val1,val2=val2", d.String())
}