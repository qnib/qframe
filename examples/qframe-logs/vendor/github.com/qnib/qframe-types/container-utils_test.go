package qtypes

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/docker/docker/api/types"
)

func TestAssembleServiceSlot(t *testing.T) {
	cnt := &types.Container{
		Labels: map[string]string{},
	}
	got := AssembleServiceSlot(cnt)
	assert.Equal(t, "<nil>", got)
	cnt.Labels = map[string]string{
		"com.docker.swarm.task.name": "service.1.id",
	}
	got = AssembleServiceSlot(cnt)
	assert.Equal(t, "service.1", got)

}
