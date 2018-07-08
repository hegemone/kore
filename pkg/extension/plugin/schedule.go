package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/hegemone/kore/pkg/comm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type plugin struct {
	calendar *calendar.Service
}

func (p *plugin) Name() string {
	return "schedule.plugins.kore.300.io"
}

func (p *plugin) Help() string {
	return "Usage: !next OR !next <show> OR !schedule"
}

func (p *plugin) CmdManifest() []comm.CmdLink {
	return []comm.CmdLink{
		comm.CmdLink{
			Regexp: regexp.MustCompile(`(?P<cmd>next)$`),
			CmdFn:  p.Next,
		},
		comm.CmdLink{
			Regexp: regexp.MustCompile(`(?P<cmd>next)\s+(?P<show>.+)`),
			CmdFn:  p.ShowNext,
		},
		comm.CmdLink{
			Regexp: regexp.MustCompile(`(?P<cmd>schedule)$`),
			CmdFn:  p.Schedule,
		},
	}
}

func (p *plugin) Next(c *comm.CmdDelegate) {
	if p.calendar == nil {
		p.auth()
	}
	t := time.Now().Format(time.RFC3339)
	events, err := p.calendar.Events.List("jalb5frk4cunnaedbfemuqbhv4@group.calendar.google.com").ShowDeleted(false).SingleEvents(true).TimeMin(t).MaxResults(1).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events. %v", err)
	}
	c.SendResponse(fmt.Sprintf("%s (%s)\n", events.Items[0].Summary, events.Items[0].Start.DateTime))
}

func (p *plugin) ShowNext(c *comm.CmdDelegate) {
	c.SendResponse("Not implemented yet")
}
func (p *plugin) Schedule(c *comm.CmdDelegate) {
	if p.calendar == nil {
		p.auth()
	}
	t := time.Now().Format(time.RFC3339)
	events, err := p.calendar.Events.List("jalb5frk4cunnaedbfemuqbhv4@group.calendar.google.com").ShowDeleted(false).SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events. %v", err)
	}

	if len(events.Items) > 0 {
		var response string
		for _, i := range events.Items {
			var when string
			// If the DateTime is an empty string the Event is an all-day Event.
			// So only Date is available.
			if i.Start.DateTime != "" {
				when = i.Start.DateTime
			} else {
				when = i.Start.Date
			}
			response += fmt.Sprintf("%s (%s)\n", i.Summary, when)
		}
		c.SendResponse(response)
	} else {
		c.SendResponse("No upcoming events found.")
	}
}

func (p *plugin) auth() {
	data, err := ioutil.ReadFile(os.Getenv("GOOGLE_SERVICE_ACCOUNT"))
	if err != nil {
		log.Fatal(err)
	}
	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/calendar.readonly")
	if err != nil {
		log.Fatal(err)
	}
	// Initiate an http.Client. The following GET request will be
	// authorized and authenticated on the behalf of
	// your service account.
	client := conf.Client(oauth2.NoContext)

	p.calendar, err = calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}
}

// Plugin is the type picked up by the engine
var Plugin plugin
