package qcollector_docker_logs

import (
	"fmt"
	"strings"
	"github.com/docker/docker/api/types"
	"github.com/qframe/types/docker-events"
	"github.com/qframe/types/health"
	"github.com/qframe/types/messages"
)

func SkipContainer(cjson *types.ContainerJSON, logEnv string) (skip bool, err error) {
	for _, v := range cjson.Config.Env {
		s := strings.Split(v,"=")
		if len(s) != 2 {
			err = fmt.Errorf("Could not parse environment variable '%s'", v)
			continue
		}
		if s[0] == logEnv && s[1] == "true" {
			skip = true
			break
		}
	}
	return
}


func createHealthhbeats(pName, rName, action string, ce qtypes_docker_events.ContainerEvent) (hbs []qtypes_health.HealthBeat) {
	b := qtypes_messages.NewTimedBase(pName, ce.Time)
	hbs = append(hbs, qtypes_health.NewHealthBeat(b, rName, ce.Container.ID[:12], action))
	hbs = append(hbs, qtypes_health.NewHealthBeat(b, "vitals", pName, fmt.Sprintf("%s.%s", ce.Container.ID[:12], action)))
	return
}
