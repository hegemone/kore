package msg

// Egress is the structured outgoing message Adapters should implement
// how to handle in their `SendMessage` function.
type Egress struct {
	ChannelID string
	Content   string
}

// EgressBuffer is the message type that wraps the Egress message type
// with metadata necessary for routing to the correct adapter.
type EgressBuffer struct {
	Originator    Originator
	EgressMessage Egress
}

// Serialize simply serializes an `Egress`.
func (e *Egress) Serialize() string {
	// NOTE: Might want to expand on this in the future
	return e.Content
}
