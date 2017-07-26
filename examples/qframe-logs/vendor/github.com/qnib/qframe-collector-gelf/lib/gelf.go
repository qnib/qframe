package qframe_collector_gelf

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/qnib/qframe-types"
	"github.com/zpatrick/go-config"
)

const (
	version   = "0.1.2"
	pluginTyp = "collector"
	pluginPkg = "gelf"
)

type Plugin struct {
	qtypes.Plugin
}

func NewPlugin(qChan qtypes.QChan, cfg config.Config, name string) Plugin {
	return Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
	}
}

func (p *Plugin) Run() {
	p.Log("info", fmt.Sprintf("Start GELF collector %s v%s", p.Name, version))
	port := p.CfgStringOr("port", "12201")
	/* Lets prepare a address at any address at port 12201*/
	ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", port))
	if err != nil {
		p.Log("error", fmt.Sprintf("%v", err))
	}
	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		p.Log("error", fmt.Sprintf("%v", err))
	}
	defer ServerConn.Close()
	p.Log("info", fmt.Sprintf("Start GELF server on '%s'", ServerAddr))
	buf := make([]byte, 1024)

	p.Log("info", fmt.Sprintf("Wait for incomming GELF message"))
	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			p.Log("error", fmt.Sprintf("%v", err))
			continue
		}
		qm := qtypes.NewQMsg("collector", p.Name)
		gmsg := GelfMsg{}
		json.Unmarshal(buf[0:n], &gmsg)
		gmsg.SourceAddr = addr.String()
		qm.Msg = gmsg.ShortMsg
		qm.Data = gmsg
		p.QChan.Data.Send(qm)
	}
}
