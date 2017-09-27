package qtypes_helper

import (
	"github.com/docker/docker/api/types"
	"fmt"
	"strings"
)

// AssembleServiceSlot create {{.Service.Name}}.{{.Task.Slot}}
func AssembleServiceSlot(cnt *types.Container) string {
	if tn, tnok := cnt.Labels["com.docker.swarm.task.name"]; tnok {
		arr := strings.Split(tn, ".")
		if len(arr) != 3 {
			return "<nil>"
		}
		return fmt.Sprintf("%s.%s", arr[0], arr[1])
	}
	return "<nil>"
}

// AssembleServiceSlot create {{.Service.Name}}.{{.Task.Slot}}
func AssembleTaskSlot(cnt *types.Container) string {
	if tn, tnok := cnt.Labels["com.docker.swarm.task.name"]; tnok {
		arr := strings.Split(tn, ".")
		if len(arr) != 3 {
			return "<nil>"
		}
		return arr[1]
	}
	return "<nil>"
}

// AssembleServiceSlot create {{.Service.Name}}.{{.Task.Slot}}
func AssembleJSONServiceSlot(cnt *types.ContainerJSON) string {
	if tn, tnok := cnt.Config.Labels["com.docker.swarm.task.name"]; tnok {
		arr := strings.Split(tn, ".")
		if len(arr) != 3 {
			return "<nil>"
		}
		return fmt.Sprintf("%s.%s", arr[0], arr[1])
	}
	return "<nil>"
}

// AssembleServiceSlot create {{.Service.Name}}.{{.Task.Slot}}
func AssembleJSONTaskSlot(cnt *types.ContainerJSON) string {
	if tn, tnok := cnt.Config.Labels["com.docker.swarm.task.name"]; tnok {
		arr := strings.Split(tn, ".")
		if len(arr) != 3 {
			return "<nil>"
		}
		return arr[1]
	}
	return "<nil>"
}

func AssembleDefaultDimensions(cnt *types.Container) map[string]string {
	dims := map[string]string{
		"container_id":   cnt.ID,
		"container_name": strings.Trim(cnt.Names[0], "/"),
		"image_name":     cnt.Image,
		"service_slot":   AssembleServiceSlot(cnt),
		"task_slot":   	  AssembleTaskSlot(cnt),
		"command":        strings.Replace(cnt.Command, " ", "#", -1),
		"created":        fmt.Sprintf("%d", cnt.Created),
	}
	for k, v := range cnt.Labels {
		dv := strings.Replace(v, " ", "#", -1)
		dv = strings.Replace(v, ".", "_", -1)
		dims[k] = dv
	}
	return dims
}

func AssembleJSONDefaultDimensions(cnt *types.ContainerJSON) map[string]string {
	dims := map[string]string{
		"container_id":   cnt.ID,
		"container_name": strings.Trim(cnt.Name, "/"),
		"image_name":     cnt.Image,
		"service_slot":   AssembleJSONServiceSlot(cnt),
		"task_slot":   	  AssembleJSONTaskSlot(cnt),
		"command":        strings.Replace(strings.Join(cnt.Config.Cmd, "#"), " ", "#", -1),
		"created":        fmt.Sprintf("%d", cnt.Created),
	}
	for k, v := range cnt.Config.Labels {
		dv := strings.Replace(v, " ", "#", -1)
		dv = strings.Replace(v, ".", "_", -1)
		dims[k] = dv
	}
	return dims
}