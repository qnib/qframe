package qframe_filter_inventory

import (
	"C"
	"fmt"
	"time"
	"github.com/zpatrick/go-config"

	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-utils"
	"github.com/qnib/qframe-inventory/lib"
)

const (
	version = "0.1.2"
	pluginTyp = qtypes.FILTER
	pluginPkg = "inventory"
)

type Plugin struct {
	qtypes.Plugin
	Inventory qframe_inventory.Inventory
}

func New(qChan qtypes.QChan, cfg config.Config, name string) Plugin {
	return Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
		Inventory: qframe_inventory.NewInventory(),
	}
}

// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start inventory v%s", p.Version))
	myId := qutils.GetGID()
	dc := p.QChan.Data.Join()
	tickerTime := p.CfgIntOr("ticker-ms", 2500)
	ticker := time.NewTicker(time.Millisecond * time.Duration(tickerTime)).C
	for {
		select {
		case val := <-dc.Read:
			switch val.(type) {
			case qtypes.ContainerEvent:
				ce := val.(qtypes.ContainerEvent)
				if ce.SourceID == int(myId) {
					continue
				}
				p.Log("debug", fmt.Sprintf("Received Event: %s.%s",ce.Event.Type, ce.Event.Action))
				if ce.Event.Type == "container" && ce.Event.Action == "start" {
					p.Inventory.SetItem(ce.Container.ID, ce.Container)
				}
			case qframe_inventory.ContainerRequest:
				req := val.(qframe_inventory.ContainerRequest)
				p.Log("info", fmt.Sprintf("Received InventoryRequest for %v", req))
				p.Inventory.ServeRequest(req)
			}
		case <- ticker:
			p.Log("debug", "Ticker came along: p.Inventory.CheckRequests()")
			p.Inventory.CheckRequests()
			continue
		}
	}
}
