// Example discord adapter. Expected to be built as a standalone .so.
package main

import (
	"github.com/hegemone/kore/pkg/comm"
	"github.com/hegemone/kore/pkg/mock"
	log "github.com/sirupsen/logrus"
)

type adapter struct {
	client mock.PlatformClient
}

func (a *adapter) Name() string {
	return "discord"
}

func (a *adapter) Listen(ingressCh chan<- comm.RawIngressMessage) {
	log.Debug("ex-discord.adapters::Listen")

	a.client = *mock.NewPlatformClient("discord")
	a.client.Connect()

	go func() {
		for clientMsg := range a.client.Chat {
			ingressCh <- comm.RawIngressMessage{
				Identity:   clientMsg.User,
				RawContent: clientMsg.Message,
			}
		}
	}()
}

func (a *adapter) SendMessage(m comm.EgressMessage) {
	a.client.SendMessage(m.Serialize())
}

// Adapter is the exported plugin symbol picked up by the engine
var Adapter adapter
