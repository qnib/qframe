package qtypes_docker_events

import (
	"testing"
	"time"
	"github.com/qframe/types/messages"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types"

	"github.com/stretchr/testify/assert"
	"github.com/qframe/types/helper"
)

func TestNewDockerEvent(t *testing.T) {
	ts := time.Unix(1499156134, 123124)
	b := qtypes_messages.NewTimedBase("src1", ts)
	info := types.Info{ID: "EngineID1"}
	event := events.Message{
		Actor: events.Actor{ID: "123"},
		Action: "start",
		Type: "container",
	}
	exp := DockerEvent{
		Base: b,
		Message: "container.start",
		Event: event,
		Engine: info,
	}
	got := NewDockerEvent(b, info, event)
	assert.Equal(t, exp, got)
}


func TestDockerEvent_ToJSON(t *testing.T) {
	ts := time.Unix(1499156134, 123124)
	b := qtypes_messages.NewTimedBase("src1", ts)
	event := events.Message{
		Actor: events.Actor{ID: "123"},
		Action: "start",
		Type: "container",
	}
	info := types.Info{ID: "EngineID1"}

	de := NewDockerEvent(b, info, event)
	exp := map[string]interface{}{
		"base_version": b.BaseVersion,
		"id": "",
		"time": ts.String(),
		"time_unix_nano": ts.UnixNano(),
		"source_id": 0,
		"source_path": []string{"src1"},
		"source_success": true,
		"tags": map[string]string{},
		"message": "container.start",
		"event": event,
	}
	got := de.EventToJSON()
	assert.Equal(t, exp["time"], got["time"])
	res := qtypes_helper.CompareMap(exp, got)
	assert.True(t, res, "Not deeply equal")
}

func TestDockerEvent_ToFlatJSON(t *testing.T) {
	ts := time.Unix(1499156134, 123124)
	b := qtypes_messages.NewTimedBase("src1", ts)
	event := events.Message{
		Actor: events.Actor{ID: "ContainerID1"},
		Action: "start",
		Type: "container",
	}
	info := types.Info{ID: "EngineID1"}

	de := NewDockerEvent(b, info, event)
	got := de.EventToFlatJSON()
	assert.Equal(t, "src1", got["msg_source_path"])
	assert.Equal(t, "1499156134000123124", got["msg_time_unix_nano"])
	assert.Equal(t, "EngineID1", got["engine_id"])
	assert.Equal(t, "container", got["event_type"])
	assert.Equal(t, "start", got["event_action"])
	assert.Equal(t, "ContainerID1", got["container_id"])
}