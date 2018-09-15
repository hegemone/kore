// Example irc adapter. Expected to be built as a standalone .so.
package main

import (
	"strings"

	irc "github.com/fluffle/goirc/client"
	"github.com/hegemone/kore/pkg/msg"
	log "github.com/sirupsen/logrus"
)

type adapter struct {
	client      *irc.Conn
	ingressChan chan<- msg.MessageInterface
}

type message struct {
	RawContent string
	Identity   string
	ChannelID  string
	Response   string
}

func (im *message) GetAdapterName() string {
	return Adapter.Name()
}

func (im *message) GetMetadata() interface{} {
	return &message{
		im.RawContent,
		im.Identity,
		im.ChannelID,
		im.Response,
	}
}

func (im *message) GetIdentity() string {
	return im.Identity
}

func (im *message) GetRawMessage() string {
	return im.RawContent
}

func (im *message) GetParsedMessage() string {
	return im.RawContent[0:len(im.RawContent)]
}

func (im *message) SetPluginResponse(response string) {
	im.Response = response
}

func (im *message) GetPluginResponse() string {
	return im.Response
}

func (a *adapter) Name() string {
	return "ex-irc.adapters.kore.nsk.io"
}

func (a *adapter) Listen(ingressCh chan<- msg.MessageInterface) {
	log.Debug("ex-irc.adapters::Listen")
	a.ingressChan = ingressCh

	cfg := irc.NewConfig("kore")
	cfg.Server = "irc.geekshed.net:6667"

	a.client = irc.Client(cfg)

	a.client.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		conn.Join("#jbot-test")
	})

	a.client.HandleFunc(irc.PRIVMSG, func(conn *irc.Conn, line *irc.Line) {
		a.ingressChan <- &message{
			Identity:   line.Nick,
			RawContent: line.Text(),
			ChannelID:  "#jbot-test",
		}
	})

	if err := a.client.Connect(); err != nil {
		log.Printf("Connection error: %s\n", err.Error())
	}
}

func (a *adapter) SendMessage(m msg.MessageInterface) {
	// The irc library we are using truncates messages with \n characters to
	// the first line. As a workaround, split the message on the newline and
	// send each line individually.
	mesg := m.GetMetadata()
	aMessage := mesg.(*message)
	log.Debugf("ex-irc.adapters::SendMessage: message is %+v", aMessage)
	log.Debugf("ex-irc.adapters::SendMessage: channelID is %v", aMessage.ChannelID)
	log.Debugf("ex-irc.adapters::SendMessage: client is %v", a.client)

	for _, i := range strings.Split(m.GetPluginResponse(), "\n") {
		a.client.Privmsg(aMessage.ChannelID, i)
	}
}

// Adapter is the exported plugin symbol picked up by engine
var Adapter adapter
