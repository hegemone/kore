package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/hegemone/kore/pkg/comm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type plugin struct {
	client *godo.Client
}

// TokenSource is a type used to create the oauth2 client.
type TokenSource struct {
	AccessToken string
}

// Token is a function to implement the oauth2 TokenSource interface.
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func (p *plugin) Name() string {
	return "digitalocean.plugins.kore.300.io"
}

func (p *plugin) Help() string {
	return "Usage: !droplet list OR !droplet start OR !droplet stop"
}

func (p *plugin) CmdManifest() []comm.CmdLink {
	return []comm.CmdLink{
		comm.CmdLink{
			Regexp: regexp.MustCompile(`^(?P<cmd>droplet)\s(?P<action>list)$`),
			CmdFn:  p.List,
		},
		comm.CmdLink{
			Regexp: regexp.MustCompile(`^(?P<cmd>droplet)\s(?P<action>start)\s(?P<id>[0-9]+)`),
			CmdFn:  p.Start,
		},
		comm.CmdLink{
			Regexp: regexp.MustCompile(`^(?P<cmd>droplet)\s(?P<action>stop)\s(?P<id>[0-9]+)`),
			CmdFn:  p.Stop,
		},
	}
}

func (p *plugin) List(c *comm.CmdDelegate) {
	if p.client == nil {
		log.Infof("Authenticating to Digital Ocean API")
		p.auth()
	}
	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	ctx := context.TODO()
	var names string
	for {
		droplets, resp, err := p.client.Droplets.List(ctx, opt)
		if err != nil {
			panic(err)
		}

		// append the current page's droplets to our list
		for _, d := range droplets {
			names += fmt.Sprintf("%s %s\n", d.Name, d.ID)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			panic(err)
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}
	c.SendResponse(names)
}

func (p *plugin) Start(c *comm.CmdDelegate) {
	if p.client == nil {
		log.Infof("Authenticating to Digital Ocean API")
		p.auth()
	}

	ctx := context.TODO()
	id, err := strconv.Atoi(c.Submatches["id"])

	if err != nil {
		log.Fatalf("Plugin::DigitalOcean - unable to convert ID string to int")
	}

	log.Infof("Plugin::DigitalOcean - Starting droplet with ID %d", id)
	_, _, err = p.client.DropletActions.PowerOn(ctx, id)

	if err != nil {
		c.SendResponse(fmt.Sprintf("Plugin::DigitalOcean - starting Droplet failed: %s", err))
	}

	c.SendResponse(fmt.Sprintf("Successfully started droplet %s", id))
}

func (p *plugin) Stop(c *comm.CmdDelegate) {
	if p.client == nil {
		log.Infof("Authenticating to Digital Ocean API")
		p.auth()
	}

	ctx := context.TODO()
	id, err := strconv.Atoi(c.Submatches["id"])

	if err != nil {
		log.Fatalf("Plugin::DigitalOcean - unable to convert ID string to int")
	}

	log.Infof("Plugin::DigitalOcean - Shutting down droplet with ID %d", id)
	_, _, err = p.client.DropletActions.Shutdown(ctx, id)

	if err != nil {
		c.SendResponse(fmt.Sprintf("Plugin::DigitalOcean - shutting down Droplet failed: %s", err))
	}

	c.SendResponse(fmt.Sprintf("Successfully started droplet %s", id))
}

func (p *plugin) auth() {
	tokenSource := &TokenSource{
		AccessToken: os.Getenv("DO_TOKEN"),
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	p.client = godo.NewClient(oauthClient)

}

// Plugin is the type picked up by the engine
var Plugin plugin
