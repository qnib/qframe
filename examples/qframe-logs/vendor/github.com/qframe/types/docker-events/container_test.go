package qtypes_docker_events

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"

	"github.com/docker/docker/api/types/events"

	"github.com/docker/docker/api/types"
	"github.com/qframe/types/messages"
	"github.com/docker/docker/api/types/container"
)

func TestContainerEvent_ContainerToFlatJSON(t *testing.T) {
	ts := time.Unix(1499156134, 123124)
	b := qtypes_messages.NewTimedBase("src1", ts)
	cbase := types.ContainerJSONBase{
		ID: "ContainerID1",
		Image: "qnib/image",
		Name: "ContainerName",
		Created: ts.String(),
		Args: []string{"tail","-f","/dev/null"},
	}
	info := types.Info{ID: "EngineID1"}
	cnt := types.ContainerJSON{
		ContainerJSONBase: &cbase,
		Config: &container.Config{Image: "configImage"},
	}
	event := events.Message{
		Actor: events.Actor{ID: "ContainerID1"},
		Action: "start",
		Type: "container",
	}
	de := NewDockerEvent(b, info, event)
	ce := NewContainerEvent(de, cnt)
	got := ce.ContainerToFlatJSON()
	assert.Equal(t, "ContainerID1", got["container_id"])
	assert.Equal(t, "EngineID1", got["engine_id"])
	assert.Equal(t, "qnib/image", got["container_image"])
	assert.Equal(t, "ContainerName", got["container_name"])
	assert.Equal(t, "2017-07-04 10:15:34.000123124 +0200 CEST", got["container_created"])
	assert.Equal(t, "tail -f /dev/null", got["container_args"])
	assert.Equal(t, "configImage", got["container_cfg_image"])

}


