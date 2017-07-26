package qframe_collector_tcp

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
	"net"
	"reflect"

	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-inventory/lib"
)

const (
	version = "0.2.3"
	pluginTyp = "collector"
	pluginPkg = "tcp"
)

type Plugin struct {
	qtypes.Plugin
	buffer chan interface{}
}

func New(qChan qtypes.QChan, cfg *config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
		buffer: make(chan interface{}, 1000),
	}
	return p, err
}

func (p *Plugin) HandleInventoryRequest(qm qtypes.Message) {
	if _, ok := qm.KV["host"]; !ok {
		p.Log("error", fmt.Sprintf("Got msg '%s' without host key!", qm))
		qm.SourceSuccess = false
		p.QChan.Data.Send(qm)
		return
	}
	p.Log("trace", fmt.Sprintf("Got msg from %s: %s", qm.KV["host"], qm.Message))
	req := qframe_inventory.NewIPContainerRequest(strings.Join(qm.SourcePath,","), qm.KV["host"])
	tout := p.CfgIntOr("inventory-timeout-ms", 2000)
	timeout := time.NewTimer(time.Duration(tout)*time.Millisecond).C
	p.QChan.Data.Send(req)
	select {
	case resp := <- req.Back:
		if resp.Error != nil {
			p.Log("error", resp.Error.Error())
			qm.SourceSuccess = false
		} else {
			qm.Container = resp.Container
			p.Log("trace", fmt.Sprintf("Got InventoryResponse: ContainerName:%s | Image:%s", qm.GetContainerName(), qm.Container.Config.Image))
		}
	case <- timeout:
		p.Log("debug", fmt.Sprintf("Experience timeout for IP %s... continue w/o Container info (SourcePath: %s)", qm.KV["host"], strings.Join(qm.SourcePath, ",")))
	}
	p.QChan.Data.Send(qm)
}

func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start collector v%s", p.Version))
	host := p.CfgStringOr("bind-host", "0.0.0.0")
	port := p.CfgStringOr("bind-port", "11001")
	// Listen for incoming connections.
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		p.Log("error", fmt.Sprintln("Error listening:", err.Error()))
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	p.Log("info", fmt.Sprintln("Listening on " + host + ":" + port))
	go p.handleRequests(l)
	for {
		select {
		case msg := <- p.buffer:
			switch msg.(type) {
			case IncommingMsg:
				im := msg.(IncommingMsg)
				base := qtypes.NewTimedBase("tcp", time.Now())
				qm := qtypes.NewMessage(base, p.Name, qtypes.MsgTCP, im.Msg)
				qm.KV["host"] = im.Host
				go p.HandleInventoryRequest(qm)
			default:
				p.Log("warn", fmt.Sprintf("Unkown data type: %s", reflect.TypeOf(msg)))
			}
		}
	}

}

func (p *Plugin) handleRequests(l net.Listener) {
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go p.handleRequest(conn)
	}
}

type IncommingMsg struct {
	Msg string
	Host string
}

func (p *Plugin) handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1048576)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	} else {
		n := bytes.Index(buf, []byte{0})
		addrTuple := strings.Split(conn.RemoteAddr().String(), ":")
		im := IncommingMsg{
			Msg: string(buf[:n-1]),
			Host: addrTuple[0],
		}
		p.Log("trace", fmt.Sprintf("Received Raw TCP message '%s' from '%s'", im.Msg, im.Host))
		p.buffer <- im
	}
	// Close the connection when you're done with it.
	conn.Close()
}
