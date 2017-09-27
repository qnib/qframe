package qcache_inventory

/******************** Inventory Request
 Sends a query for a key or an IP and provides a back-channel, so that the requesting partner can block on the request
 until it arrives - honouring a timeout...
*/

import (
	"strings"
	"time"
	"github.com/docker/docker/api/types"
	"fmt"
)

type ContainerRequest struct {
	IssuedAt 	time.Time
	Source 		string
	Timeout	 	time.Duration
	Name 		string
	ID 			string
	IP 			string
	Back 		chan Response
}

func NewContainerRequest(src string, to time.Duration) ContainerRequest {
	cr := ContainerRequest{
		IssuedAt: 	time.Now(),
		Source: 	src,
		Timeout:  	to,
		Back: 		make(chan Response, 5),
	}
	return cr
}


func NewIDContainerRequest(src, id string) ContainerRequest {
	cr := NewContainerRequest(src, time.Second)
	cr.ID = id
	return cr
}

func NewNameContainerRequest(src, name string) ContainerRequest {
	cr := NewContainerRequest(src, time.Second)
	cr.Name =  name
	return cr
}

func NewIPContainerRequest(src, ip string) ContainerRequest {
	cr := NewContainerRequest(src, time.Duration(2)*time.Second)
	cr.IP =  ip
	return cr
}

func (this ContainerRequest) Equal(other Response) (err error) {
	err = this.EqualCnt(other.Container)
	if err != nil {
		return
	}
	err = this.EqualIPS(other.Ips)
	return
}

func (this ContainerRequest) EqualCnt(other *types.ContainerJSON) (err error) {
	if this.ID == "" && this.Name == "" {
		return
	}
	if this.ID == other.ID || this.Name == strings.Trim(other.Name, "/") {
		return
	}
	return fmt.Errorf("this.ID:%s != %s:ID.other || this.Name:%s != %s:Name.other", this.ID, other.ID, this.Name, other.Name)
}


func (this ContainerRequest) EqualIPS(ips []string) (err error) {
	if this.IP == "" {
		return
	}
	for _, ip := range ips {
		if this.IP == ip {
			return
		}
	}
	return fmt.Errorf("this.IP:%s not in %v", this.IP, ips)
}

func (cr *ContainerRequest) TimedOut()  bool {
	tDiff := time.Now().Sub(cr.IssuedAt)
	return tDiff > cr.Timeout
}