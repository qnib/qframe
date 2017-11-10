package qutils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestIsInput(t *testing.T) {
	assert.True(t, IsInput([]string{"one","two"}, "one"))
	assert.True(t, IsInput([]string{"two"}, "one-two"))
}

func TestIsLastSource(t *testing.T) {
	assert.True(t, IsInput([]string{"one","two"}, "one"))
	assert.True(t, IsInput([]string{"two"}, "one-two"))
}
