package qtypes_health

/*****
Health messages are send by collectors and filters as a heartbeat, beating every ticker event within qframe
******/


import (
	"github.com/qframe/types/messages"
)

type HealthBeat struct {
	qtypes_messages.Base
	Type string
	Actor string
	Action string
}

func NewHealthBeat(b qtypes_messages.Base, t,actor,action string) HealthBeat {
	return HealthBeat{
		Base: b,
		Type: t,
		Actor: actor,
		Action: action,
	}
}
