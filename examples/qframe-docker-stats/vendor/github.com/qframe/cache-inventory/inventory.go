package qcache_inventory

import (
	"errors"
	"sync"
	"fmt"
	"github.com/docker/docker/api/types"

)

type Inventory struct {
	Version string
	Data   map[string]Response
	PendingRequests []ContainerRequest
	mux sync.Mutex
}

func NewInventory() Inventory {
	return Inventory{
		Version: version,
		Data: map[string]Response{},
		PendingRequests: []ContainerRequest{},
	}
}

func (i *Inventory) SetItem(key string, item *types.ContainerJSON, info types.Info, ips []string) (err error) {
	i.mux.Lock()
	defer i.mux.Unlock()
	resp := NewOKResponse(item, &info, ips)
	i.Data[key] = resp
	return
}

func (i *Inventory) GetItem(key string) (out Response, err error) {
	i.mux.Lock()
	defer i.mux.Unlock()
	if item, ok := i.Data[key];ok {
		return item, err
	}
	return out, errors.New(fmt.Sprintf("No item found with key '%s'", key))
}

func filterItem(in ContainerRequest, other Response) (out Response, err error) {
	err = in.Equal(other)
	if err == nil {
		return other, err
	}
	return out, err
}


func (i *Inventory) HandleRequest(req ContainerRequest) (err error) {
	if len(i.Data) == 0 {
		return errors.New("inventory is empty so far")
	}
	for _, resp := range i.Data {
		res, err := filterItem(req, resp)
		if err == nil {
			req.Back <- NewOKResponse(res.Container, resp.Engine, resp.Ips)
			return err
		} else if req.TimedOut() {
			err = errors.New(fmt.Sprintf("Timed out after %s", req.Timeout.String()))
			req.Back <- NewFAILResponse(err)
			return err
		}
	}
	return errors.New("could not match filter")
}


func (i *Inventory) ServeRequest(req ContainerRequest) (err error) {
	err = i.HandleRequest(req)
	if err != nil {
		i.mux.Lock()
		i.PendingRequests = append(i.PendingRequests, req)
		i.mux.Unlock()
	}
	return
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

