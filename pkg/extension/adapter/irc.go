// Example irc adapter. Expected to be built as a standalone .so.
package main

import (
	irc "github.com/fluffle/goirc/client"
	"github.com/hegemone/kore/pkg/comm"
	log "github.com/sirupsen/logrus"
)

type adapter struct {
	client      *irc.Conn
	ingressChan chan<- comm.RawIngressMessage
}

func (a *adapter) Name() string {
	return "ex-irc.adapters.kore.nsk.io"
}

func (a *adapter) Listen(ingressCh chan<- comm.RawIngressMessage) {
	log.Debug("ex-irc.adapters::Listen")
	a.ingressChan = ingressCh

	cfg := irc.NewConfig("kore")
	cfg.Server = "irc.geekshed.net:6667"

	a.client = irc.Client(cfg)

	a.client.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		conn.Join("#jbot-test")
	})

	a.client.HandleFunc(irc.PRIVMSG, func(conn *irc.Conn, line *irc.Line) {
		a.ingressChan <- comm.RawIngressMessage{
			Identity:   line.Nick,
			RawContent: line.Text(),
			ChannelID:  "#jbot-test",
		}
	})

	if err := a.client.Connect(); err != nil {
		log.Printf("Connection error: %s\n", err.Error())
	}
}

func (a *adapter) SendMessage(m comm.EgressMessage) {
	a.client.Privmsg(m.ChannelID, m.Serialize())
}

// Adapter is the exported plugin symbol picked up by engine
var Adapter adapter
