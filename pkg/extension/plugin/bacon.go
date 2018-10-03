// Example plugin. Implements the classic bacon cinch plugin.
package main

import (
	"fmt"
	"regexp"

	"github.com/hegemone/kore/pkg/comm"
	log "github.com/sirupsen/logrus"
)

type plugin struct{}

func (p *plugin) Name() string {
	return "bacon.plugins.kore.nsk.io"
}

func (p *plugin) Help() string {
	return "Usage: !bacon [user]"
}

func (p *plugin) CmdManifest() []comm.CmdLink {
	return []comm.CmdLink{
		comm.CmdLink{
			Regexp: regexp.MustCompile(`(?P<cmd>bacon)$`),
			CmdFn:  p.CmdBacon,
		},
		comm.CmdLink{
			Regexp: regexp.MustCompile(`(?P<cmd>bacon)\s+(?P<recipient>\S+)`),
			CmdFn:  p.CmdBaconGift,
		},
	}
}

func (p *plugin) CmdBacon(c *comm.CmdDelegate) {
	log.Infof("bacon.plugins::CmdBacon, IngressMessage: %+v", c.IngressMessage)

	msg := c.IngressMessage
	identity := msg.GetIdentity()

	response := fmt.Sprintf(
		"gives %s a strip of delicious bacon.", identity,
	)

	c.SendResponse(response)
}

func (p *plugin) CmdBaconGift(c *comm.CmdDelegate) {
	log.Infof("bacon.plugins::CmdBaconGift, IngressMessage: %+v", c.IngressMessage)

	msg := c.IngressMessage
	identity := msg.GetIdentity()
	toUser := c.Submatches["recipient"]

	response := fmt.Sprintf(
		"gives %s a strip of delicious bacon as a gift from %v",
		toUser, identity,
	)

	c.SendResponse(response)
}

// Plugin is the exported type picked up by the engine
var Plugin plugin
