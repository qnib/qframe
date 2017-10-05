package qhandler_kafka

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/zpatrick/go-config"
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/qnib/qframe-types"
	"github.com/qframe/types/docker-events"
	"strings"
	"github.com/qframe/types/plugin"
	"github.com/qframe/types/metrics"
	"github.com/qframe/types/qchannel"
)

const (
	version = "0.1.0"
	pluginTyp = qtypes.HANDLER
	pluginPkg = "kafka"
)

type Plugin struct {
    *qtypes_plugin.Plugin
	producer *kafka.Producer
	mutex sync.Mutex
	deliveryChan chan kafka.Event
}

func New(qChan qtypes_qchannel.QChan, cfg *config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes_plugin.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
		deliveryChan: make(chan kafka.Event),
	}
	return p, err
}

// Connect creates a connection to InfluxDB
func (p *Plugin) Connect() (err error) {
	bPort := p.CfgStringOr("broker.port", "9092")
	bList := []string{}
	for _, b := range strings.Split(p.CfgStringOr("broker.hosts", "tasks.broker"), ",") {
		bList = append(bList, fmt.Sprintf("%s:%s", b, bPort))
	}
	brokers := strings.Join(bList, ",")
	p.Log("info", fmt.Sprintf("Connect to broker: %s", brokers))
	p.producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})

	if err != nil {
		return
	} else {
		msg := fmt.Sprintf("Created Producer ID:%d Name:%s\n", p.MyID, p.Name)
		p.Log("info", msg)
	}
	return
}


// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start handler %s, v%s", p.Name, version))
	err := p.Connect()
	if err != nil {
		msg := fmt.Sprintf("Failed to create producer: %s\n", err)
		p.Log("error", msg)
		return
	}
	bg := p.QChan.Data.Join()
	/*dims := map[string]string{
		"version": version,
		"plugin": p.Name,
	}*/
	for {
		select {
		case val := <-bg.Read:
			p.Log("trace", fmt.Sprintf("received event: %s | %v", reflect.TypeOf(val), val))
			switch val.(type) {
			case qtypes_metrics.Metric:
				m := val.(qtypes_metrics.Metric)
				/*if m.StopProcessing(p.Plugin, false) {
					continue
				}*/
				p.PushToKafka(m)
			case qtypes_docker_events.ServiceEvent:
				se := val.(qtypes_docker_events.ServiceEvent)
				//se.StopProcessing(p.Plugin, false)
				p.PushToKafka(se)
			case qtypes_docker_events.ContainerEvent:
				ce := val.(qtypes_docker_events.ContainerEvent)
				//ce.StopProcessing(p.Plugin, false)
				p.PushToKafka(ce)
			}
		}
	}
}

type Payload struct {
	Topic string
	Data map[string]interface{}
}

// ToJSON creates JSON payload depending on the type of Event.
// The Payload is placed in a map, in which the key defines the topic to push the payload to.
func (p *Plugin) ToPayload(e interface{}) (payloads []Payload, err error) {
	switch e.(type) {
	case qtypes_docker_events.ContainerEvent:
		ce := e.(qtypes_docker_events.ContainerEvent)
		switch ce.Event.Action {
		case "start","create":
			// In case the container starts, the information about the start is passed
			payloads = append(payloads, Payload{Topic: "cnt_details", Data: ce.ContainerToFlatJSON()})
		case "exec_create":
			return
		}
		// Add normal DockerEvent
		payloads = append(payloads, Payload{Topic: "cnt_event", Data: ce.EventToFlatJSON()})
	case qtypes_docker_events.ServiceEvent:
		se := e.(qtypes_docker_events.ServiceEvent)
		switch se.Event.Action {
		case "create":
			// In case the container starts, the information about the start is passed
			payloads = append(payloads, Payload{Topic: "srv_details", Data: se.ServiceToJSON()})
		case "exec_create":
			return
		}
		payloads = append(payloads, Payload{Topic: "srv_event", Data: se.EventToJSON()})

	default:
		p.Log("info", fmt.Sprintf("Skip sending to kafka: %s", reflect.TypeOf(e)))

	}
	return
}

func (p *Plugin) PushToKafka(e interface{}) (err error) {
	payloads, err := p.ToPayload(e)
	for _, payload := range payloads {
		topic := payload.Topic
		data := payload.Data
		val, err := json.Marshal(data)
		if err != nil {
			p.Log("error", fmt.Sprintf("Marshaling failed: %v", err.Error()))
			return err
		}
		err = p.producer.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Value: val}, p.deliveryChan)
		e := <-p.deliveryChan
		m := e.(*kafka.Message)

		if m.TopicPartition.Error != nil {
			p.Log("error", fmt.Sprintf("Delivery failed: %v", m.TopicPartition.Error))
			return m.TopicPartition.Error
		} else {
			msg := fmt.Sprintf("Delivered message to topic %s [%d] at offset %v",
				*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			p.Log("tracing", msg)
		}
	}
	return
}
