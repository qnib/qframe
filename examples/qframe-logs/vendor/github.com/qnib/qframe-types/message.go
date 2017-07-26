package qtypes

import (
	"strings"
	"github.com/docker/docker/api/types"
)

const (
	MsgCEE = "cee"
	MsgTCP = "tcp"
	MsgDLOG = "docker-log"
)


type Message struct {
	Base
	Container   types.ContainerJSON
	Name       	string            	`json:"name"`
	LogLevel       string				`json:"loglevel"`
	MessageType	string            	`json:"type"`
	Message     string            	`json:"value"`
	KV			map[string]string 	`json:"data"`
}

func NewMessage(base Base, name, mType, msg string) Message {
	return Message{
		Base: base,
		Name: name,
		Container: types.ContainerJSON{},
		LogLevel: "INFO",
		MessageType: mType,
		Message: msg,
		KV: map[string]string{},
	}
}

func NewContainerMessage(base Base, cnt types.ContainerJSON, name, mType, msg string) Message {
	m := NewMessage(base, name, mType, msg)
	m.Container = cnt
	return m
}

func (m *Message) GetContainerName() string {
	if m.Container.Name != "" {
		return strings.Trim(m.Container.Name, "/")
	} else {
		return "<none>"
	}
}