package qframe_handler_elasticsearch

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/OwnLocal/goes"
	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-utils"
	"github.com/zpatrick/go-config"
)

const (
	version   = "0.1.12"
	pluginTyp = "handler"
	pluginPkg = "elasticsearch"
)

// Elasticsearch holds a buffer and the initial information from the server
type Elasticsearch struct {
	qtypes.Plugin
	buffer      chan interface{}
	indexPrefix string
	indexName   string
	KVtoFields  map[string]string
	SkipKV		[]string
	last        time.Time
	conn        *goes.Connection
}

// NewElasticsearch returns an initial instance
func New(qChan qtypes.QChan, cfg *config.Config, name string) (Elasticsearch, error) {
	p := Elasticsearch{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg,  name, version),
		buffer: make(chan interface{}, 1000),
	}
	nameSplit := strings.Split(p.Name, "_")
	idxDef := p.Name
	if len(nameSplit) != 0 {
		idxDef = nameSplit[len(nameSplit)-1]
	}
	idx := p.CfgStringOr("index-prefix", idxDef)
	p.ParseKVtoFields()
	p.ParseSkipKV()
	p.indexPrefix = idx
	p.last = time.Now().Add(-24 * time.Hour)
	return p, nil
}

func (p *Elasticsearch) ParseSkipKV() {
	kvCfg := p.CfgStringOr("kv-skip", "")
	p.SkipKV = strings.Split(kvCfg, ",")
}

func (p *Elasticsearch) ParseKVtoFields() {
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
func (p *Elasticsearch) pushToBuffer() {
	bg := p.QChan.Data.Join()
	inputs := p.GetInputs()
	srcSuccess := p.CfgBoolOr("source-success", true)
	for {
		val := bg.Recv()
		switch val.(type) {
		case qtypes.QMsg:
			qm := val.(qtypes.QMsg)
			if len(inputs) != 0 && !qutils.IsInput(inputs, qm.Source) {
				//p.Log("debug", fmt.Sprintf("(%s) %v - skip_input src:%s not in %v", qm.Source, qm.Msg, qm.Source, inputs))
				continue
			}
			if qm.SourceSuccess != srcSuccess {
				//p.Log("debug", fmt.Sprintf("(%s) %v - skip_success %v", qm.Source, qm.Msg, qm.SourceSuccess))
				continue
			}
			p.Log("trace", fmt.Sprintf("qtypes.QMsg from '%v' || %s", qm.SourcePath, qm.Msg))
			p.buffer <- qm
		case qtypes.Message:
			msg := val.(qtypes.Message)
			if ! msg.InputsMatch(inputs) {
				continue
			}
			if msg.SourceSuccess != srcSuccess {
				continue
			}
			p.buffer <- msg
		case qtypes.ContainerEvent:
			msg := val.(qtypes.ContainerEvent)
			if ! msg.InputsMatch(inputs) {
				continue
			}
			if msg.SourceSuccess != srcSuccess {
				continue
			}
			p.buffer <- msg


		}
	}
}

func (p *Elasticsearch) createESClient() (err error) {
	host := p.CfgStringOr("host", "localhost")
	port := p.CfgStringOr("port", "9200")
	now := time.Now()
	p.indexName = fmt.Sprintf("%s-%04d-%02d-%02d", p.indexPrefix, now.Year(), now.Month(), now.Day())
	p.Log("info", fmt.Sprintf("Connecting to %s:%s", host, port))
	p.conn = goes.NewConnection(host, port)
	return
}

func (p *Elasticsearch) createIndex() (err error) {
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
	idxExist, _ := p.conn.IndicesExist(indices)
	if idxExist {
		//p.Log("debug", fmt.Sprintf("Index '%s' already exists", p.indexName)
		return err
	}
	//log.Printf("[DD] Index '%v' does not exists", indices)
	_, err = p.conn.CreateIndex(p.indexName, idxCfg)
	if err != nil {
		//p.Log("warn", fmt.Sprintf("Index '%s' could not be created", p.indexName)
		return err
	}
	p.Log("debug", fmt.Sprintf("Created index '%s'.", p.indexName))
	return err
}

func (p *Elasticsearch) indexQMsg(msg qtypes.QMsg) error {
	data := map[string]interface{}{
		"msg_version": msg.QmsgVersion,
		"Timestamp":   msg.Time.Format("2006-01-02T15:04:05.999999-07:00"),
		"msg":         msg.Msg,
		"source":      msg.Source,
		"source_path": msg.SourcePath,
		"type":        msg.Type,
		"host":        msg.Host,
		"Level":       msg.Level,
	}
	if len(msg.KV) != 0 {
		data[msg.Source] = msg.KV
	}
	switch msg.Data.(type) {
	case qtypes.GelfMsg:
		//p.Log("debug", "msg-data is GELF msg...")
		gmsg := msg.Data.(qtypes.GelfMsg)
		data["container_id"] = gmsg.ContainerID
		data["container_name"] = gmsg.ContainerName
		data["container_cmd"] = gmsg.Command
		data["container_host"] = gmsg.Host
		data["image_id"] = gmsg.ImageID
		data["image_name"] = gmsg.ImageName
	}
	d := goes.Document{
		Index:  p.indexName,
		Type:   "log",
		Fields: data,
	}
	extraArgs := make(url.Values, 1)
	//extraArgs.Set("ttl", "86400000")
	response, err := p.conn.Index(d, extraArgs)
	_ = response
	//fmt.Printf("%s | %s\n", d, response.Error)
	return err
}

func (p *Elasticsearch) indexContainerEvent(msg qtypes.ContainerEvent) (err error) {
	data := map[string]interface{}{
		"msg_version": 	msg.BaseVersion,
		"Timestamp":   	msg.Time.Format("2006-01-02T15:04:05.999999-07:00"),
		"msg":         	msg.Message,
		"source_path": 	strings.Join(msg.SourcePath,","),
	}
	for k,v := range msg.Data {
		data[k] = v
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
		Type:   "container-event",
		Fields: data,
	}
	extraArgs := make(url.Values, 1)
	//extraArgs.Set("ttl", "86400000")
	response, err := p.conn.Index(d, extraArgs)
	_ = response
	return
}

func (p *Elasticsearch) indexMessage(msg qtypes.Message) (err error) {
	data := map[string]interface{}{
		"msg_version": 	msg.BaseVersion,
		"Timestamp":   	msg.Time.Format("2006-01-02T15:04:05.999999-07:00"),
		"msg":         	msg.Message,
		"source_path": 	strings.Join(msg.SourcePath,","),
		"Level":   		msg.LogLevel,
	}

	if host, ok := msg.KV["host"]; ok {
		data["host"] = host
	}
	for k,v := range msg.Data {
		data[k] = v
	}
	for k,v := range msg.KV {
		key := fmt.Sprintf("%s.%s", msg.Name, k)
		if qutils.IsItem(p.SkipKV, key) {
			p.Log("debug", fmt.Sprintf("Skip key %s in qm.KV", key))
		} else {
			if nKey, ok := p.KVtoFields[key]; ok {
				p.Log("debug", fmt.Sprintf("Overwrite field '%s' with %s", nKey, v))
				data[nKey] = v
			} else {
				data[key] = v
			}
		}
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
	p.Log("trace", fmt.Sprintf("%30s: %s", "_id", msg.ID))
	p.Log("trace", fmt.Sprintf("%30s: %s", "_type", msg.MessageType))
	for k,v := range data {
		p.Log("trace", fmt.Sprintf("%30s: %s", k, v))
	}
	d := goes.Document{
		Index:  p.indexName,
		Type:   msg.MessageType,
		Fields: data,
	}
	if msg.ID != "" {
		d.Id = msg.ID
	}
	extraArgs := make(url.Values, 1)
	//extraArgs.Set("ttl", "86400000")
	response, err := p.conn.Index(d, extraArgs)
	_ = response
	return
}

func (p *Elasticsearch) indexDoc(doc interface{}) (err error) {
	now := time.Now()
	if p.last.Day() != now.Day() {
		p.indexName = fmt.Sprintf("%s-%04d-%02d-%02d", p.indexPrefix, now.Year(), now.Month(), now.Day())
		p.createIndex()
		p.last = now
	}
	switch doc.(type) {
	case qtypes.QMsg:
		msg := doc.(qtypes.QMsg)
		return p.indexQMsg(msg)
	case qtypes.Message:
		msg := doc.(qtypes.Message)
		return p.indexMessage(msg)
	case qtypes.ContainerEvent:
		msg := doc.(qtypes.ContainerEvent)
		return p.indexContainerEvent(msg)
	}
	return
}

// Run pushes the logs to elasticsearch
func (p *Elasticsearch) Run() {
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
