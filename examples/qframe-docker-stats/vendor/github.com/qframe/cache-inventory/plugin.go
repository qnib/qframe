package qcache_inventory

import (
	"C"
	"fmt"
	"time"
	"github.com/zpatrick/go-config"

	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-utils"
)

const (
	version = "0.3.0"
	pluginTyp = qtypes.CACHE
	pluginPkg = "inventory"
)

type Plugin struct {
	qtypes.Plugin
	Inventory Inventory
}

func New(qChan qtypes.QChan, cfg *config.Config, name string) (Plugin, error) {
	return Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
		Inventory: NewInventory(),
	}, nil
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
				if ce.SourceID == myId {
					continue
				}
				p.Log("trace", fmt.Sprintf("Received Event: %s.%s",ce.Event.Type, ce.Event.Action))
				if ce.Event.Type == "container" && ce.Event.Action == "start" {
					p.Inventory.SetItem(ce.Container.ID, ce.Container)
				}
			case ContainerRequest:
				req := val.(ContainerRequest)
				p.Log("trace", fmt.Sprintf("Received InventoryRequest for %v", req))
				p.Inventory.ServeRequest(req)
			}
		case <- ticker:
			p.Log("trace", "Ticker came along: p.Inventory.CheckRequests()")
			p.Inventory.CheckRequests()
			continue
		}
	}
}
