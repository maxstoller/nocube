package rapa102

import (
	"context"
	"time"

	"github.com/coral/nocube/pkg"
	"github.com/grandcat/zeroconf"
	log "github.com/sirupsen/logrus"
)

type RMan struct {
	registeredDevices []RAPA102
	connectedDevices  []*RAPA102
}

func New() *RMan {
	return &RMan{}
}

func (rm *RMan) Init() {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.WithFields(log.Fields{
				"Name": entry.ServiceRecord.Instance,
				"IP":   entry.AddrIPv4[0],
				"Port": entry.Port,
			}).Info("Found RAPA102 controller")

			d := RAPA102{
				Name:        entry.ServiceRecord.Instance,
				IP:          entry.AddrIPv4[0],
				Port:        entry.Port,
				PixelStream: make(chan []pkg.ColorLookupResult),
			}

			rm.registeredDevices = append(rm.registeredDevices, d)

			go func() {
				err := d.Connect()
				if err != nil {
					//TODO handle error here ffs
					return
				}
				rm.connectedDevices = append(rm.connectedDevices, &d)
			}()
		}
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err = resolver.Browse(ctx, "_apabridge._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
}

func (rm *RMan) ModuleName() string {
	return "RAPA102"
}

func (rm *RMan) Write(d []pkg.ColorLookupResult) {
	for _, conectedRAPA102 := range rm.connectedDevices {
		conectedRAPA102.PixelStream <- d
	}
}
