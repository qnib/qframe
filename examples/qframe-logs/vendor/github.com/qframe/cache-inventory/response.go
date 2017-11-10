package qcache_inventory

import (
	"github.com/docker/docker/api/types"
)


type Response struct {
	Container *types.ContainerJSON
	Ips       []string
	Error     error
}

func NewOKResponse(cnt *types.ContainerJSON, ips []string) Response {
	var err error
	return Response{
		Container: cnt,
		Ips: ips,
		Error: err,
	}
}

func NewFAILResponse(err error) Response {
	var cnt *types.ContainerJSON
	return Response{
		Container: cnt,
		Ips: []string{},
		Error: err,
	}
}