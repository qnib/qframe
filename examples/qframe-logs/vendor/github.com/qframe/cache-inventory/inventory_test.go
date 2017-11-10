package qcache_inventory


import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

var (
	cnt1 = NewContainer("CntID1", "CntName1", map[string]string{"eth0": "172.17.0.1"})
	resp1 = NewOKResponse(&cnt1, []string{"192.168.0.1"})
	cnt2 = NewContainer("CntID2", "CntName2", map[string]string{"eth0": "172.17.0.2"})
	resp2 = NewOKResponse(&cnt2, []string{"192.168.0.2"})
)

func TestInventory_SetItem(t *testing.T) {
	i := NewInventory()
	assert.IsType(t, Inventory{}, i)
	i.Data[cnt1.ID] = resp1
	assert.Len(t, i.Data, 1)
	i.SetItem(cnt2.ID, &cnt2, []string{"192.168.0.2"})
	assert.Len(t, i.Data, 2)
}

func TestInventory_GetItem(t *testing.T) {
	i := NewInventory()
	i.SetItem(cnt1.ID, &cnt1, []string{"192.168.0.1"})
	i.SetItem(cnt2.ID, &cnt2, []string{"192.168.0.2"})
	got, err := i.GetItem(cnt1.ID)
	assert.NoError(t, err)
	assert.Equal(t, resp1, got)
}

func TestInventory_filterItem(t *testing.T) {
	req1 := NewContainerRequest("src", time.Second)
	req1.Name = "CntName1"
	got, err := filterItem(req1, resp1)
	assert.NoError(t, err)
	assert.Equal(t, resp1, got)
	got, err = filterItem(req1, resp2)
	assert.Error(t, err, err.Error())
}


func TestInventory_HandleRequest(t *testing.T) {
	i := NewInventory()
	i.SetItem(cnt1.ID, &cnt1, []string{"192.168.0.1"})
	req := NewNameContainerRequest("src", cnt1.Name)
	err := i.HandleRequest(req)
	assert.NoError(t, err)
	resp := <-req.Back
	assert.Equal(t, resp1.Container, resp.Container)
	reqID := NewIDContainerRequest("src", "FakeID")
	err = i.HandleRequest(reqID)
	assert.Error(t, err)

}

func TestInventory_ServeRequest(t *testing.T) {
	i := NewInventory()
	i.SetItem(cnt1.ID, &cnt1, []string{"192.168.0.1"})
	req := NewNameContainerRequest("src", cnt1.Name)
	i.ServeRequest(req)
	assert.Equal(t, len(i.PendingRequests), 0)
	req2 := NewNameContainerRequest("src", cnt2.Name)
	i.ServeRequest(req2)
	assert.Equal(t, len(i.PendingRequests), 1)
}


func TestInventory_CheckRequest(t *testing.T) {
	i := NewInventory()
	req := NewNameContainerRequest("src", cnt1.Name)
	i.ServeRequest(req)
	assert.Equal(t, len(i.PendingRequests), 1)
	i.SetItem(cnt1.ID, &cnt1, []string{"192.168.0.1"})
	i.CheckRequests()
	resp := <-req.Back
	assert.Equal(t, resp1.Container, resp.Container)
}

func TestInventory_CheckMultipleRequest(t *testing.T) {
	i := NewInventory()
	req := NewNameContainerRequest("src", cnt1.Name)
	i.ServeRequest(req)
	req2 := NewNameContainerRequest("src", cnt2.Name)
	req2.Timeout = time.Duration(5)*time.Second
	i.ServeRequest(req2)
	assert.Equal(t, len(i.PendingRequests), 2)
	i.SetItem(cnt1.ID, &cnt1, []string{"192.168.0.1"})
	i.CheckRequests()
	r1 := <-req.Back
	assert.Equal(t, resp1.Container, r1.Container)
	assert.Equal(t, 1, len(i.PendingRequests))
	i.SetItem(cnt2.ID, &cnt2, []string{"192.168.0.2"})
	i.CheckRequests()
	r2 := <-req2.Back
	assert.NoError(t, r2.Error, "Should be fine")
	assert.Equal(t, resp2.Container, r2.Container)
}
