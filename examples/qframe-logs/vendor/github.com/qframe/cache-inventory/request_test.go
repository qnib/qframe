package qcache_inventory

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

func NewContainer(id, name string, ips map[string]string) types.ContainerJSON {
	cbase :=  &types.ContainerJSONBase{
		ID: id,
		Name: name,
	}

	netConfig := &types.NetworkSettings{}
	netConfig.Networks = map[string]*network.EndpointSettings{}
	for iface, ip := range ips {
		endpoint := &network.EndpointSettings{
			IPAddress: ip,
		}
		netConfig.Networks[iface] =  endpoint
	}
	cnt := types.ContainerJSON{
		ContainerJSONBase: cbase,
		NetworkSettings: netConfig,
	}
	return cnt
}

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


func TestContainer_NonEqual(t *testing.T) {
	cnt := NewContainer("CntID1", "CntName1", map[string]string{"eth0": "172.17.0.2"})
	resp := NewOKResponse(&cnt, []string{"10.0.0.1"})
	cntB := NewBridgedOnlyContainer("CntID2", "CntName2", "192.168.0.1")
	respB := NewOKResponse(&cntB, []string{"10.0.0.2"})
	checkIP := NewIPContainerRequest("src1", "172.17.0.1")
	assert.Error(t, checkIP.Equal(resp))
	assert.Error(t, checkIP.Equal(respB))
	checkName := ContainerRequest{Name: "CntNameFail"}
	assert.Error(t, checkName.Equal(resp))
	assert.Error(t, checkName.Equal(respB))
	checkID := ContainerRequest{ID: "CntIDFail"}
	assert.Error(t, checkID.Equal(resp))
	assert.Error(t, checkID.Equal(respB))
}



func TestContainer_EqualIPS(t *testing.T) {
	checkIP := NewIPContainerRequest("src1", "10.0.0.1")
	assert.NoError(t, checkIP.EqualIPS([]string{"10.0.0.1"}))
}

func TestContainer_EqualCnt(t *testing.T) {
	cnt := NewContainer("CntID1", "CntName1", map[string]string{"eth0": "172.17.0.2"})
	checkIP := NewNameContainerRequest("src1", "CntName1")
	assert.NoError(t, checkIP.EqualCnt(&cnt))
}



func TestContainer_Equal(t *testing.T) {
	cnt := NewContainer("CntID1", "CntName1", map[string]string{"eth0": "172.17.0.2"})
	resp := NewOKResponse(&cnt, []string{"10.0.0.1"})
	checkIP := NewIPContainerRequest("src1", "172.17.0.2")
	assert.NoError(t, checkIP.EqualIPS([]string{"172.17.0.2"}))
	assert.NoError(t, checkIP.EqualCnt(&cnt))
	assert.NoError(t, checkIP.Equal(resp))
	checkName := ContainerRequest{Name: "CntName1"}
	assert.NoError(t, checkName.Equal(resp))
	checkID := ContainerRequest{ID: "CntID1"}
	assert.NoError(t, checkID.Equal(resp))
	checkIP2 := NewIPContainerRequest("src1", "10.0.0.1")
	assert.NoError(t, checkIP2.Equal(resp))

}

func TestContainer_BridgeEqual(t *testing.T) {
	cnt := NewBridgedOnlyContainer("CntID2", "CntName2","172.17.0.2")
	resp := NewOKResponse(&cnt, []string{"10.0.0.1"})
	checkIP := NewIPContainerRequest("src1", "172.17.0.2")
	assert.NoError(t, checkIP.Equal(resp))
	checkName := ContainerRequest{Name: "CntName2"}
	assert.NoError(t, checkName.Equal(resp))
	checkID := ContainerRequest{ID: "CntID2"}
	assert.NoError(t, checkID.Equal(resp))
}

func TestContainerRequest_TimedOut(t *testing.T) {
	req := ContainerRequest{
		Source: "src1",
		Name: "CntName1",
	}
	req.IssuedAt = time.Now().AddDate(0,0,-1)
	assert.True(t, req.TimedOut(), "Should be timed out long ago")
}