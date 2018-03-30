// Example irc adapter. Expected to be built as a standalone .so.
package main

import (
	irc "github.com/fluffle/goirc/client"
	"github.com/hegemone/kore/pkg/msg"
	log "github.com/sirupsen/logrus"
	"strings"
)

type adapter struct {
	client      *irc.Conn
	ingressChan chan<- msg.RawIngress
}

func (a *adapter) Name() string {
	return "ex-irc.adapters.kore.nsk.io"
}

func (a *adapter) Listen(ingressCh chan<- msg.RawIngress) {
	log.Debug("ex-irc.adapters::Listen")
	a.ingressChan = ingressCh

	cfg := irc.NewConfig("kore")
	cfg.Server = "irc.geekshed.net:6667"

	a.client = irc.Client(cfg)

	a.client.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		conn.Join("#jbot-test")
	})

	a.client.HandleFunc(irc.PRIVMSG, func(conn *irc.Conn, line *irc.Line) {
		a.ingressChan <- msg.RawIngress{
			Identity:   line.Nick,
			RawContent: line.Text(),
			ChannelID:  "#jbot-test",
		}
	})

	if err := a.client.Connect(); err != nil {
		log.Printf("Connection error: %s\n", err.Error())
	}
}

func (a *adapter) SendMessage(m msg.Egress) {
	// The irc library we are using truncates messages with \n characters to
	// the first line. As a workaround, split the message on the newline and
	// send each line individually.
	for _, i := range strings.Split(m.Serialize(), "\n") {
		a.client.Privmsg(m.ChannelID, i)
	}
}

// Adapter is the exported plugin symbol picked up by engine
var Adapter adapter
