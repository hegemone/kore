package msg

type MessageInterface interface {
	GetIdentity() string
	GetMetadata() interface{}
	GetAdapterName() string
	GetParsedMessage() string
	GetRawMessage() string
	GetPluginResponse() string
	SetPluginResponse(string)
}
