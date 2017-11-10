package qtypes_docker_events



import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types"

	"github.com/qframe/types/messages"
	"github.com/qframe/types/helper"
	"github.com/docker/docker/api/types/swarm"
)

func TestServiceEvent_ToJSON(t *testing.T) {
	ts := time.Unix(1499156134, 123124)
	b := qtypes_messages.NewTimedBase("src1", ts)
	event := events.Message{
		Actor: events.Actor{ID: "123"},
		Action: "create",
		Type: "service",
	}
	srv := swarm.Service{
		ID: "12345",
	}
	info := types.Info{ID: "EngineID1"}
	de := NewDockerEvent(b, info, event)
	se := NewServiceEvent(de, srv)
	exp := map[string]interface{}{
		"base_version": b.BaseVersion,
		"id": "",
		"time": ts.String(),
		"time_unix_nano": ts.UnixNano(),
		"source_id": 0,
		"source_path": []string{"src1"},
		"source_success": true,
		"tags": map[string]string{},
		"message": "service.create",
		"service": srv,
	}
	got := se.ServiceToJSON()
	assert.Equal(t, exp["time"], got["time"])
	res := qtypes_helper.CompareMap(exp, got)
	assert.True(t, res, "Not deeply equal")
}
