package qcollector_docker_logs

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/types/messages"
)


// struct to keep info and channels to goroutine
// -> get heartbeats so that we know it's still alive
// -> allow for gracefully shutdown of the supervisor
type ContainerSupervisor struct {
	Plugin
	TailRunning string
	Action      string
	CntID 		string 			 // ContainerID
	CntName 	string			 // sanatized name of container
	Info		*types.Info
	Container 	*types.ContainerJSON
	Com 		chan interface{} // Channel to communicate with goroutine
	cli 		*client.Client
	qChan 		qtypes_qchannel.QChan
	TimeRexex	*regexp.Regexp
}


func (cs ContainerSupervisor) Run() {
	msg := fmt.Sprintf("Start listener for: '%s' [%s]", cs.CntName, cs.CntID)
	cs.Log("info", msg)

	logOpts := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Tail: "any", Timestamps: true}
	if cs.Action == "running" {
		logOpts.Tail = cs.TailRunning
	}
	reader, err := cs.cli.ContainerLogs(ctx, cs.CntID, logOpts)
	if err != nil {
		msg := fmt.Sprintf("Error when connecting to log of %s: %s", cs.CntName, err.Error())
		cs.Log("error", msg)
		return
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		sText := strings.Split(line, " ")
		shostL := strings.TrimLeft(strings.Join(sText[1:], " "), " ")
		var base qtypes_messages.Base
		lTime, err := cs.fuzzyParseTime(sText[0])
		if err != nil {
			base = qtypes_messages.NewBase(cs.Name)
		} else {
			base = qtypes_messages.NewTimedBase(cs.Name, lTime)
		}
		qm := qtypes_messages.NewContainerMessage(base, cs.Container, shostL)
		qm.AddEngineInfo(cs.Info)
		cs.Log("debug", fmt.Sprintf("MsgDigest:'%s'  | Container '%s': %s", qm.GetMessageDigest(), cs.Container.Name, shostL))
		cs.qChan.Data.Send(qm)
	}
	err = scanner.Err()
	if err != nil {
		msg := fmt.Sprintf("Something went wrong going through the log: %s", err.Error())
		cs.Log("error", msg)
		return
	}
	for {
		select {
		case msg := <-cs.Com:
			switch msg {
			case "died":
				msg := fmt.Sprintf("Container [%s]->'%s' died -> BYE!", cs.CntID, cs.CntName)
				cs.Log("debug", msg)
				return
			default:
				msg := fmt.Sprintf("Container [%s]->'%s' got message from cs.Com: %v", cs.CntID, cs.CntName, msg)
				cs.Log("debug", msg)
			}
		}
	}
}

func (cs ContainerSupervisor) fuzzyParseTime(s string) (t time.Time, err error) {
	t, err = time.Parse(time.RFC3339, s)
	if err != nil {
		match := cs.TimeRegex.FindString(s)
		if match != "" {
			t, err = time.Parse(time.RFC3339, match)
			return
		}
	}
	return
}