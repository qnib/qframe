package qcache_inventory

import (
	"errors"
	"sync"
	"fmt"
	"github.com/docker/docker/api/types"

)

type Inventory struct {
	Version string
	Data   map[string]types.ContainerJSON
	PendingRequests []ContainerRequest
	mux sync.Mutex
}

func NewInventory() Inventory {
	return Inventory{
		Version: version,
		Data: make(map[string]types.ContainerJSON),
		PendingRequests: []ContainerRequest{},
	}
}

func (i *Inventory) SetItem(key string, item types.ContainerJSON) (err error) {
	i.mux.Lock()
	defer i.mux.Unlock()
	i.Data[key] = item
	return
}

func (i *Inventory) GetItem(key string) (cntOut types.ContainerJSON, err error) {
	i.mux.Lock()
	defer i.mux.Unlock()
	if item, ok := i.Data[key];ok {
		return item, err
	}
	return cntOut, errors.New(fmt.Sprintf("No item found with key '%s'", key))
}

func filterItem(in ContainerRequest, other types.ContainerJSON) (out types.ContainerJSON, err error) {
	if in.Equal(other) {
		return other, err
	}
	return out, errors.New("filter does not match")
}


func (i *Inventory) HandleRequest(req ContainerRequest) (err error) {
	if len(i.Data) == 0 {
		return errors.New("Inventory is empty so far")
	}
	for _, cnt := range i.Data {
		res, err := filterItem(req, cnt)
		if err == nil {
			req.Back <- NewOKResponse(res)
			return err
		} else if req.TimedOut() {
			err = errors.New(fmt.Sprintf("Timed out after %s", req.Timeout.String()))
			req.Back <- NewFAILResponse(err)
			return err
		}
	}
	return errors.New("Could not match filter")
}


func (i *Inventory) ServeRequest(req ContainerRequest) {
	err := i.HandleRequest(req)
	if err != nil {
		i.mux.Lock()
		i.PendingRequests = append(i.PendingRequests, req)
		i.mux.Unlock()
	}
}


// CheckRequests iterates over all requests and responses if the request can be fulfilled
func (inv *Inventory) CheckRequests() {
	if len(inv.PendingRequests) == 0 {
		return
	}
	inv.mux.Lock()
	defer inv.mux.Unlock()
	remainReq := []ContainerRequest{}
	for _, req := range inv.PendingRequests {
		err := inv.HandleRequest(req)
		if err != nil {
			remainReq = append(remainReq, req)
		}
	}
	inv.PendingRequests = remainReq
}

