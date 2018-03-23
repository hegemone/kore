// Example discord adapter. Expected to be built as a standalone .so.
package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/hegemone/kore/pkg/msg"
	log "github.com/sirupsen/logrus"
	"os"
)

type adapter struct {
	ingressChan chan<- msg.RawIngress
	client      *discordgo.Session
}

func (a *adapter) Name() string {
	return "ex-discord.adapters.kore.nsk.io"
}

func (a *adapter) Listen(ingressCh chan<- msg.RawIngress) {
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

func (a *adapter) SendMessage(m msg.Egress) {
	a.client.ChannelMessageSend(m.ChannelID, m.Serialize())
	//a.client.SendMessage(m.Serialize())
}

func (a *adapter) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	a.ingressChan <- msg.RawIngress{m.Author.Username, m.ChannelID, m.Content}
}

// Adapter is the exported plugin symbol picked up by the engine
var Adapter adapter
