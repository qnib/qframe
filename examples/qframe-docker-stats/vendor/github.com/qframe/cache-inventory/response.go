package qcache_inventory

import (
	"github.com/docker/docker/api/types"
)


type Response struct {
	Container types.ContainerJSON
	Error     error
}

func NewOKResponse(cnt types.ContainerJSON) Response {
	var err error
	return Response{
		Container: cnt,
		Error: err,
	}
}

func NewFAILResponse(err error) Response {
	var cnt types.ContainerJSON
	return Response{
		Container: cnt,
		Error: err,
	}
}