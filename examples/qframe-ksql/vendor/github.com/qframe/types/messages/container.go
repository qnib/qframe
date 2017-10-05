package qtypes_messages


import (
	"fmt"
	"strings"
	"github.com/docker/docker/api/types"

)

type ContainerMessage struct {
	Message
	Container types.ContainerJSON
	Engine types.Info	// TODO:
}


func NewContainerMessage(base Base, cnt *types.ContainerJSON, msg string) ContainerMessage {
	m := NewMessage(base, msg)
	c := ContainerMessage{
		Message: m,
		Container: *cnt,
		Engine: types.Info{},
	}
	c.ID = c.GenContainerMsgID()
	return c
}


func (cm *ContainerMessage) AddEngineInfo(e *types.Info) {
	cm.Engine = *e
}

// GenContainerMsgID uses "<container_id>-<time.UnixNano()>-<MSG>" and does a sha1 hash.
func (c *ContainerMessage) GenContainerMsgID() string {
	s := fmt.Sprintf("%s-%d-%s", c.Container.ID, c.Time.UnixNano(), c.Message)
	return Sha1HashString(s)
}

func (c *ContainerMessage) GetContainerName() string {
	if c.Container.Name != "" {
		return strings.Trim(c.Container.Name, "/")
	} else {
		return "<none>"
	}
}
