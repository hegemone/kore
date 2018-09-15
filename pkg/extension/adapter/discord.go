// Example discord adapter. Expected to be built as a standalone .so.
package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type Message struct {
	RawContent string
	Identity   string
	ChannelID  string
	Response   string
}

func (im Message) GetAdapterName() string {
	return Adapter.Name()
}

func (im Message) GetIdentity() string {
	return im.Identity
}

func (im Message) GetRawMessage() string {
	return im.RawContent
}

func (im Message) GetParsedMessage() string {
	return im.RawContent[0:len(im.RawContent)]
}

func (im Message) SetPluginResponse(response string) {
	im.Response = response
}

func (im Message) GetPluginResponse() string {
	return im.Response
}

type adapter struct {
	ingressChan chan<- Message
	client      *discordgo.Session
}

func (a *adapter) Name() string {
	return "ex-discord.adapters.kore.nsk.io"
}

func (a *adapter) Listen(ingressCh chan<- Message) {
	log.Debug("ex-discord.adapters::Listen")

	a.ingressChan = ingressCh

	var err error
	a.client, err = discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("DISCORD_TOKEN")))
	if err != nil {
		log.Errorf("unable to establish Discord session: %s", err)
		panic(err)
	}

	a.client.AddHandler(a.messageCreate)

	err = a.client.Open()
	if err != nil {
		panic(err)
	}
}

func (a *adapter) SendMessage(m Message) {
	a.client.ChannelMessageSend(m.ChannelID, m.GetPluginResponse())
	//a.client.SendMessage(m.Serialize())
}

func (a *adapter) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	a.ingressChan <- Message{Identity: m.Author.Username, ChannelID: m.ChannelID, RawContent: m.Content}
}

// Adapter is the exported plugin symbol picked up by the engine
var Adapter adapter
