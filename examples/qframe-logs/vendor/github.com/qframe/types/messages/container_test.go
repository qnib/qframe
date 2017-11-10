package qtypes_messages

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

var (
	t1 = time.Unix(1504086453, 0)
)

func NewBridgedOnlyContainer(id, name, ip string) types.ContainerJSON {
	cbase :=  &types.ContainerJSONBase{
		ID: id,
		Name: name,
	}

	netConfig := &types.NetworkSettings{DefaultNetworkSettings: types.DefaultNetworkSettings{IPAddress: ip}}
	netConfig.Networks = map[string]*network.EndpointSettings{}
	cnt := types.ContainerJSON{
		ContainerJSONBase: cbase,
		NetworkSettings: netConfig,
	}
	return cnt
}

func TestUnitContainerMessage_GetContainerName(t *testing.T) {
	cnt1 := NewBridgedOnlyContainer("CntID1", "CntName1", "192.168.0.1")
	cnt2 := NewBridgedOnlyContainer("CntID1", "", "192.168.0.1")
	b := NewTimedBase("test1", t1)
	cm := NewContainerMessage(b, &cnt1, "This is a message")
	got := cm.GetContainerName()
	assert.Equal(t, "CntName1", got)
	cm = NewContainerMessage(b, &cnt2, "This is a message")
	got = cm.GetContainerName()
	assert.Equal(t, "<none>", got)
}