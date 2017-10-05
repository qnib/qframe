package qutils

import (
	"strings"
	"fmt"
	"regexp"
	"errors"
)

type DockerServiceTask struct {
	Name string
	Slot string
	TaskID string
}
func SanatizeContainerName(names []string) (string) {
	// TODO: Is there any reason to care about other names?
	return strings.TrimPrefix(names[0], "/")
}

func ContainerNameExtractService(names []string) (task DockerServiceTask, err error) {
	name := SanatizeContainerName(names)
	regEx := regexp.MustCompile(`(?P<service>[a-z\-\.\_0-9]+)\.(?P<slot>[0-9]+)\.(?P<task_id>[a-z0-9]+)`)
	m := GetParams(regEx, name)
	if len(m) == 3 {
		task.TaskID = m["task_id"]
		task.Slot, _ = m["slot"]
		task.Name = m["service"]
	} else {
		err = errors.New(fmt.Sprintf("Container Name '%s' does not match the service-task nameing scheme", name))
	}
	return
}
