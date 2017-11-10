package qtypes_helper


import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLogStrToInt(t *testing.T) {
	exp := map[string]int{
		"panic": 0,
		"error": 3,
		"warn": 4,
		"notice": 5,
		"info": 6,
		"debug": 7,
		"trace": 8,
		"nil": 6,
	}
	for l,i := range exp {
		assert.Equal(t, i, LogStrToInt(l))
	}
}

