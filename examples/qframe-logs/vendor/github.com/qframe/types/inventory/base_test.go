package qtypes_inventory

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/qframe/types/messages"
)


const (
	inv1 = `{"time": 0, "subject":"node1", "object": "node2", "action": "connected", "tags": {}}`
)

func TestNewBaseFromJson(t *testing.T) {
	ts := time.Unix(1499156134, 0)
	qb := qtypes_messages.NewTimedBase("src1", ts)
	b, err := NewBaseFromJson(qb, inv1)
	assert.NoError(t, err)
	assert.Equal(t, "node1", b.Subject)
}

func TestSplitUnixNano(t *testing.T) {
	now := 1257894000000000011
	s, n := SplitUnixNano(int64(now))
	assert.Equal(t, int64(1257894000), s)
	assert.Equal(t, int64(11), n)
}
