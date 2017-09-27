package qhandler_elasticsearch

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/OwnLocal/goes"
	"github.com/zpatrick/go-config"

	"github.com/qframe/types/plugin"
	"github.com/qframe/types/constants"
	"github.com/qframe/types/docker-events"
	"reflect"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/types/messages"
)

const (
	version   = "0.1.14"
	pluginTyp = qtypes_constants.HANDLER
	pluginPkg = "elasticsearch"
)

// Plugin holds a buffer and the initial information from the server
type Plugin struct {
	*qtypes_plugin.Plugin
	buffer      chan interface{}
	indexPrefix string
	indexName   string
	KVtoFields  map[string]string
	SkipKV		[]string
	last        time.Time
	cli        	*goes.Connection
}

// New returns an initial instance
func New(qChan qtypes_qchannel.QChan, cfg *config.Config, name string) (Plugin, error) {
	p := qtypes_plugin.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version)
	esp := Plugin{
		Plugin: p,
		buffer: make(chan interface{}, 1000),
	}
	nameSplit := strings.Split(p.Name, "_")
	idxDef := esp.Name
	if len(nameSplit) != 0 {
		idxDef = nameSplit[len(nameSplit)-1]
	}
	idx := p.CfgStringOr("index-prefix", idxDef)
	esp.ParseKVtoFields()
	esp.ParseSkipKV()
	esp.indexPrefix = idx
	esp.last = time.Now().Add(-24 * time.Hour)
	return esp, nil
}

func (p *Plugin) ParseSkipKV() {
	kvCfg := p.CfgStringOr("kv-skip", "")
	p.SkipKV = strings.Split(kvCfg, ",")
}

func (p *Plugin) ParseKVtoFields() {
	p.KVtoFields = map[string]string{}
	kvCfg := p.CfgStringOr("kv-to-field", "")
	for _, tuple := range strings.Split(kvCfg, ",") {
		slice := strings.Split(tuple, ":")
		if len(slice) != 2 {
			p.Log("error", fmt.Sprintf("Could not split kv-to-field by ':': %s", tuple))
			return
		}
		p.Log("info", fmt.Sprintf("KV key '%s' will replace '%s' when indexing", slice[0], slice[1]))
		p.KVtoFields[slice[0]] = slice[1]
	}
}

// Takes log from framework and buffers it in elasticsearch buffer
func (p *Plugin) pushToBuffer() {
	bg := p.QChan.Data.Join()
	for {
		val := bg.Recv()
		switch val.(type) {
		case qtypes_messages.ContainerMessage:
			msg := val.(qtypes_messages.ContainerMessage)
			if msg.StopProcessing(p.Plugin, false) {
				continue
			}
			p.buffer <- msg
		case qtypes_docker_events.ContainerEvent:
			msg := val.(qtypes_docker_events.ContainerEvent)
			if msg.StopProcessing(p.Plugin, false) {
				continue
			}
			p.buffer <- msg
		default:
			p.Log("trace", fmt.Sprintf("No case for type %s", reflect.TypeOf(val)))

		}
	}
}

func (p *Plugin) createESClient() (err error) {
	host := p.CfgStringOr("host", "localhost")
	port := p.CfgStringOr("port", "9200")
	now := time.Now()
	p.indexName = fmt.Sprintf("%s-%04d-%02d-%02d", p.indexPrefix, now.Year(), now.Month(), now.Day())
	p.Log("info", fmt.Sprintf("Connecting to %s:%s", host, port))
	p.cli = goes.NewConnection(host, port)
	return
}

func (p *Plugin) createIndex() (err error) {
	idxCfg := map[string]interface{}{
		"settings": map[string]interface{}{
			"index.number_of_shards":   1,
			"index.number_of_replicas": 0,
		},
		"mappings": map[string]interface{}{
			"_default_": map[string]interface{}{
				"_source": map[string]interface{}{
					"enabled": true,
				},
				"_all": map[string]interface{}{
					"enabled": false,
				},
			},
		},
	}
	indices := []string{p.indexName}
	idxExist, _ := p.cli.IndicesExist(indices)
	if idxExist {
		//p.Log("debug", fmt.Sprintf("Index '%s' already exists", p.indexName)
		return err
	}
	//log.Printf("[DD] Index '%v' does not exists", indices)
	_, err = p.cli.CreateIndex(p.indexName, idxCfg)
	if err != nil {
		//p.Log("warn", fmt.Sprintf("Index '%s' could not be created", p.indexName)
		return err
	}
	p.Log("debug", fmt.Sprintf("Created index '%s'.", p.indexName))
	return err
}

func (p *Plugin) indexContainerEvent(msg qtypes_docker_events.ContainerEvent) (err error) {
	data := map[string]interface{}{
		"msg_version": 	msg.BaseVersion,
		"Timestamp":   	msg.Time.Format("2006-01-02T15:04:05.999999-07:00"),
		"msg":         	msg.Message,
		"source_path": 	strings.Join(msg.SourcePath,","),
	}
	if msg.GetContainerName() != "" {
		data["container_id"] = msg.Container.ID
		data["container_name"] = msg.GetContainerName()
		data["container_cmd"] = strings.Join(msg.Container.Config.Cmd, " ")
		data["image"] = msg.Container.Image
		data["image_name"] = msg.Container.Config.Image
		if msg.Container.Node != nil {
			p.Log("debug", "Set msg.Container.Node")
			//data["node_name"] = msg.Container.Node.Name
			//data["node_ip"] = msg.Container.Node.IPAddress
		}

	}
	for k,v := range data {
		p.Log("trace", fmt.Sprintf("%30s: %s", k, v))
	}
	d := goes.Document{
		Index:  p.indexName,
		Type: "container-event",
		Fields: data,
	}
	extraArgs := make(url.Values, 1)
	//extraArgs.Set("ttl", "86400000")
	response, err := p.cli.Index(d, extraArgs)
	_ = response
	return
}

func (p *Plugin) indexContainerMessage(msg qtypes_messages.ContainerMessage) (err error) {
	data := map[string]interface{}{
		"msg_version": 	msg.BaseVersion,
		"Timestamp":   	msg.Time.Format("2006-01-02T15:04:05.999999-07:00"),
		"msg":         	msg.Message.ToJSON(),
		"docker_engine":  map[string]interface{}{
			"name": msg.Engine.Name,
			"id": msg.Engine.ID,
			"kernel": msg.Engine.KernelVersion,
			"server_version": msg.Engine.ServerVersion,
			"labels": msg.Engine.Labels,
		},
		"swarm": map[string]interface{}{
			"node_id": msg.Engine.Swarm.NodeID,
		},
		"source_path": 	strings.Join(msg.SourcePath,","),
	}
	if msg.GetContainerName() != "" {
		data["container_id"] = msg.Container.ID
		data["container_name"] = msg.GetContainerName()
		data["container_cmd"] = strings.Join(msg.Container.Config.Cmd, " ")
		data["image"] = msg.Container.Image
		data["image_name"] = msg.Container.Config.Image
		if msg.Container.Node != nil {
			p.Log("debug", "Set msg.Container.Node")
			//data["node_name"] = msg.Container.Node.Name
			//data["node_ip"] = msg.Container.Node.IPAddress
		}

	}
	for k,v := range data {
		p.Log("debug", fmt.Sprintf("%30s: %s", k, v))
	}
	d := goes.Document{
		Index:  p.indexName,
		Type: 	"container-message",
		Fields: data,
	}
	extraArgs := make(url.Values, 1)
	//extraArgs.Set("ttl", "86400000")
	_, err = p.cli.Index(d, extraArgs)
	return
}

func (p *Plugin) indexDoc(doc interface{}) (err error) {
	now := time.Now()
	if p.last.Day() != now.Day() {
		p.indexName = fmt.Sprintf("%s-%04d-%02d-%02d", p.indexPrefix, now.Year(), now.Month(), now.Day())
		p.createIndex()
		p.last = now
	}
	switch doc.(type) {
	case qtypes_docker_events.ContainerEvent:
		msg := doc.(qtypes_docker_events.ContainerEvent)
		return p.indexContainerEvent(msg)
	case qtypes_messages.ContainerMessage:
		msg := doc.(qtypes_messages.ContainerMessage)
		return p.indexContainerMessage(msg)
	}
	return
}

// Run pushes the logs to elasticsearch
func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start elasticsearch handler: %sv%s", p.Name, version))
	go p.pushToBuffer()
	err := p.createESClient()
	p.createIndex()
	_ = err
	for {
		qm := <-p.buffer
		err := p.indexDoc(qm)
		if err != nil {
			p.Log("error", fmt.Sprintf("Failed to index msg: %s || %v", qm, err))
		} else {
			p.Log("trace", "Indexed message...")
		}
	}
}
