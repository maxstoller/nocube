package rapa102

import (
	"net"
	"net/url"
	"strconv"

	"github.com/coral/nocube/pkg"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type RAPA102 struct {
	Name string
	IP   net.IP
	Port int
	conn *websocket.Conn

	PixelStream chan []pkg.ColorLookupResult
}

func (r *RAPA102) Connect() error {

	addr := r.IP.String() + ":" + strconv.Itoa(r.Port)
	u := url.URL{Scheme: "ws", Host: addr, Path: "/data"}
	log.WithFields(log.Fields{
		"Name": r.Name,
		"IP":   r.IP,
		"Port": r.Port,
		"URL":  u.String(),
	}).Info("Connecting to discovered RAPA102")

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.WithError(err).Error("Could not connect to RAPA102 on " + u.String())
		return err
	}

	r.conn = c

	go r.handlePixels()

	return nil
}
func (r *RAPA102) handlePixels() {
	for {
		select {
		case p := <-r.PixelStream:
			var bytes = []byte{}
			for _, _ = range p {
				bytes = append(bytes, []byte{
					/* 					utils.Clamp255(color.Color[0] * 255),
					   					utils.Clamp255(color.Color[1] * 255),
					   					utils.Clamp255(color.Color[2] * 255), */
					255, 0, 0,
				}...)
			}

			r.conn.WriteMessage(websocket.BinaryMessage, bytes)
		}
	}
}
