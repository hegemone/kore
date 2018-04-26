package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/digitalocean/godo"
	"github.com/hegemone/kore/pkg/comm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	pat = "89f4d579f3053052680477cc3788caef5447f2917bb84b9e9df856745a40faff"
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
			Regexp: regexp.MustCompile(`droplet\slist$`),
			CmdFn:  p.List,
		},
		comm.CmdLink{
			Regexp: regexp.MustCompile(`droplet\sstart$`),
			CmdFn:  p.Start,
		},
		comm.CmdLink{
			Regexp: regexp.MustCompile(`droplet\sstop$`),
			CmdFn:  p.Stop,
		},
	}
}

func (p *plugin) List(c *comm.CmdDelegate) {
	if p.client == nil {
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

		log.Infof("Received page of droplets: %v", droplets)

		// append the current page's droplets to our list
		for _, d := range droplets {
			names += fmt.Sprintf("%s\n", d.Name)
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
	c.SendResponse("Not implemented yet")
}
func (p *plugin) Stop(c *comm.CmdDelegate) {
	c.SendResponse("Not implemented yet")
}

func (p *plugin) auth() {
	tokenSource := &TokenSource{
		AccessToken: pat,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	p.client = godo.NewClient(oauthClient)

}

// Plugin is the type picked up by the engine
var Plugin plugin
