package qtypes_helper


import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestIsInput(t *testing.T) {
	assert.True(t, IsInput([]string{"one","two"}, "one"))
	assert.True(t, IsInput([]string{"two"}, "one-two"))
}

func TestIsItem(t *testing.T) {
	assert.True(t, IsInput([]string{"one","two"}, "*"))
	assert.False(t, IsInput([]string{"one","two"}, "three"))

}

func TestIsLastSource(t *testing.T) {
	assert.True(t, IsLastSource([]string{"one","two"}, "two"))
	assert.False(t, IsLastSource([]string{"one","two"}, "one"))
}

