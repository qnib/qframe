package qtypes_messages

import (
	"encoding/json"
	"github.com/deckarep/golang-set"
	"github.com/qframe/types/plugin"
	"fmt"
)

type Message struct {
	Base
	Message string
}

func NewMessage(b Base, msg string) Message {
	m := Message{
		Base: b,
		Message: msg,
	}
	m.GenDefaultID()
	return m
}

// ParseJSONMap iterates over set of potential keys and if unmarshalls the string value of all keys into a new map.
func (m *Message) ParseJsonMap(p *qtypes_plugin.Plugin, keys mapset.Set, kv map[string]string) map[string]string {
	res := map[string]string{}
	it := keys.Iterator()
	for val := range it.C {
		key := val.(string)
		v, ok := kv[key]
		if !ok {
			p.Log("debug", fmt.Sprintf("Could not find key '%s' in Tags: %v", key, kv))
			continue
		}
		p.Log("debug", fmt.Sprintf("unmarshall: %s", v))
		byt := []byte(v)
		var dat map[string]interface{}
		json.Unmarshal(byt, &dat)
		for k, v := range dat {
			res[k] = fmt.Sprintf("%v", v)
		}
	}
	return res
}

// ToStringRFC54242 returns a string in RFC5424 format
func (m *Message) ToStringRFC54242() (res string, err error) {

	return
}
