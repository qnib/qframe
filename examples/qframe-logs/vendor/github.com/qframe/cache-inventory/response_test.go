package qcache_inventory

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
)

var (
	cnt = NewContainer("CntID1", "CntName1", map[string]string{"eth0": "172.17.0.1"})
)

func TestNewOKResponse(t *testing.T) {
	r := NewOKResponse(&cnt, []string{})
	assert.NoError(t, r.Error, "Should create a response with a nil error")
}

func TestNewFAILResponse(t *testing.T) {
	r := NewFAILResponse(errors.New("FAIL!"))
	assert.Error(t, r.Error, "Should create a response with an error")
}