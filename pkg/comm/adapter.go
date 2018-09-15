package comm

import (
	"fmt"
	goplugin "plugin"
	"regexp"

	"github.com/hegemone/kore/pkg/msg"
)

// The prefix trigger used to denote a cmd.
// Example: !droplet start skynet
const adapterCmdTriggerPrefix = "!"

// Regexp applied to check isCmd.
var adapterCmdRegexp, _ = regexp.Compile(fmt.Sprintf("^%s\\S*($| )",
	adapterCmdTriggerPrefix))

// Adapter is an abstraction that should be implemented to present a standard
// interface to the comm server for communicating to, and from external platforms.
// Similar to the `Plugin`, it is a facade type that delegates  actions like
// sending and receiving messages to concrete implementations dynamically loaded
// from shared libraries.
type Adapter interface {
	// SendMessage is the public trigger indicating a dynamically loaded adapter
	// should transmit an `EgressMessage` to its platform. Dynamically loaded
	// adapters must define how that is done.
	SendMessage(msg.MessageInterface)
	Name() string
	// Listen is the public trigger that initiates an adapter to start listening
	// to external platform events. It should be implemented as non-blocking and
	// push `RawIngressMessage`s to the inChan on the receipt of raw messages
	// from the external platform.
	Listen(chan<- msg.MessageInterface)
}

// LoadAdapter loads adapter behavior from a given .so adapter file
func LoadAdapter(adapterFile string) (Adapter, error) {
	// TODO: Need a *lot* of validation here to make sure a bad adapter doesn't
	// just crash the server.
	// -> Actually confirm the casts are valid and these functions look like they should?
	// TODO: Can the hardcoded pattern of $PROPERTY Lookup -> Cast be made more elegant?
	rawGoPlugin, err := goplugin.Open(adapterFile)
	if err != nil {
		return nil, err
	}

	aSym, err := rawGoPlugin.Lookup("Adapter")

	if err != nil {
		return nil, err
	}

	a := aSym.(Adapter)
	return a, nil
}

func isCmd(rawContent string) bool {
	// isCmd is where the adapter defines whether or not raw content is indeed a Cmd.
	return adapterCmdRegexp.MatchString(rawContent)
}
