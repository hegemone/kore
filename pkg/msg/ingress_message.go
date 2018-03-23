package msg

// Originator contains originating information for an `IngressMessage`. It
// contains the "identity" of the person/system that triggered the incoming
// message, in addition to the adapter name that procuded the event.
type Originator struct {
	Identity    string
	AdapterName string
	ChannelID   string
}

// Ingress is the structured, parsed message representing an incoming
// command to be processed. An `IngressMessage` has been determined by the
// system to have content containing a command.
type Ingress struct {
	Content    string
	Originator Originator
}

// RawIngress is an unprocessed message passed from the adapter to the
// engine. Has not yet been parsed to determine if the message is a cmd or not.
type RawIngress struct {
	Identity   string
	ChannelID  string
	RawContent string
}

// IngressBuffer is the buffer messasge type that carries the parsed ingress
// message along with any necessary metadata for proper routing. Currently,
// there is no metadata.
type IngressBuffer struct {
	IngressMessage Ingress
}

// RawIngressBuffer messages are internal messaging types usually containing a public
// payload + some kind of metadata, ex: to facilitate routing
type RawIngressBuffer struct {
	AdapterName string // e.g. Discord
	// the raw message, i.e. `!cmdTrigger cmd `
	RawIngressMessage RawIngress
}
