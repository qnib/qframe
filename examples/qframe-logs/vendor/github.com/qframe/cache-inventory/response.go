package qcache_inventory

import (
	"github.com/docker/docker/api/types"
)


type Response struct {
	Container *types.ContainerJSON
	Engine    *types.Info
	Ips       []string
	Error     error
}

func NewOKResponse(cnt *types.ContainerJSON, info *types.Info, ips []string) Response {
	var err error
	return Response{
		Container: cnt,
		Engine: info,
		Ips: ips,
		Error: err,
	}
}

func NewFAILResponse(err error) Response {
	var cnt *types.ContainerJSON
	var info *types.Info
	return Response{
		Container: cnt,
		Ips: []string{},
		Engine: info,
		Error: err,
	}
}