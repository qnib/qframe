package qtypes

import (
	"strings"
	"github.com/docker/docker/api/types"
	"github.com/qnib/qframe-utils"
	"fmt"
)

const (
	MsgCEE = "cee"
	MsgTCP = "tcp"
	MsgDLOG = "docker-log"
	MsgMetric = "metric" //needs to have name,time and value field ; optional tags (key1=val1,key2=val2)
)


type Message struct {
	Base
	Container   types.ContainerJSON
	Name       	string            	`json:"name"`
	LogLevel    string				`json:"loglevel"`
	MessageType	string            	`json:"type"`
	Message     string            	`json:"value"`
	KV			map[string]string 	`json:"data"`
}

func NewMessage(base Base, name, mType, msg string) Message {
	m := Message{
		Base: base,
		Name: name,
		Container: types.ContainerJSON{},
		LogLevel: "INFO",
		MessageType: mType,
		Message: msg,
		KV: map[string]string{},
	}
	m.SourceID = int(qutils.GetGID())
	return m
}

func NewContainerMessage(base Base, cnt types.ContainerJSON, name, mType, msg string) Message {
	m := NewMessage(base, name, mType, msg)
	m.Container = cnt
	m.ID = m.GenContainerMsgID()
	return m
}

// GenContainerMsgID uses "<container_id>-<time.UnixNano()>-<MSG>" and does a sha1 hash.
func (m *Message) GenContainerMsgID() string {
	s := fmt.Sprintf("%s-%d-%s", m.Container.ID, m.Time.UnixNano(), m.Message)
	return Sha1HashString(s)
}

func (m *Message) GetContainerName() string {
	if m.Container.Name != "" {
		return strings.Trim(m.Container.Name, "/")
	} else {
		return "<none>"
	}
}
