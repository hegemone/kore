package comm

import (
	"fmt"
	"path/filepath"

	"github.com/hegemone/kore/pkg/config"
	"github.com/hegemone/kore/pkg/msg"
	log "github.com/sirupsen/logrus"
)

// Engine is the heart of korecomm. It's responsible for routing traffic amongst
// buffers in a concurrent way, as well as the loading and execution of extensions.
type Engine struct {
	// Messaging buffers
	ingressBuffer chan msg.MessageInterface
	egressBuffer  chan msg.MessageInterface

	// Extensions
	plugins  map[string]Plugin
	adapters map[string]Adapter
}

// NewEngine creates a new Engine.
func NewEngine(c *config.Config) *Engine {
	// Configurable size of the internal message buffers
	bufferSize := c.GetEngine().BufferSize

	return &Engine{
		ingressBuffer: make(chan msg.MessageInterface, bufferSize),
		egressBuffer:  make(chan msg.MessageInterface, bufferSize),
		plugins:       make(map[string]Plugin),
		adapters:      make(map[string]Adapter),
	}
}

// LoadExtensions will attempt to load enabled plugins and extensions. Includes
// extension init (used for things like establishing connections with platforms).
func (e *Engine) LoadExtensions() error {
	log.Info("Loading extensions")
	if err := e.loadPlugins(); err != nil {
		return err
	}
	return e.loadAdapters()
}

// These are all helper functions to allow for testing. I'll need to look
// into if there isn't a better way to structure the code to make testing
// easier but for now, this'll do.
var (
	eHandleIngress = (*Engine).handleIngress
	eHandleEgress  = (*Engine).handleEgress
	funcDone       = func() {}
)

// Start will cause the engine to start listening on all successfully loaded
// adapters. On the receipt of any new message from an adapter, it will parse
// the message and determine if the contents are a command. If the message does
// contain a command, it will be transformed to an `IngressMessage` and routed
// to matching plugin commands. If the plugin sends back a message to the
// originator, it will be transformed to an `EgressMessage` and routed to the
// originating adapter for transmission via the client.
func (e *Engine) Start() {
	log.Debug("Engine::Start")

	// Spawn listening routines for each adapter
	for _, adapter := range e.adapters {
		adapterIngressChannel := make(chan msg.MessageInterface, 2)

		go func(adapter Adapter, adapterIngressChannel chan msg.MessageInterface) {
			// Tell the adapter to start listening and sending messages back via
			// their own ingress channel. Listen should be non-blocking!
			adapter.Listen(adapterIngressChannel)

			// Engine listens to the N channels the adapters are transmitting on
			// for RawIngressMessages. Adapter channels are fanned-in to the
			// rawIngressBuffer for parsing.
			for message := range adapterIngressChannel {
				e.ingressBuffer <- message
				funcDone()
			}
		}(adapter, adapterIngressChannel)
	}

	// Wire up messaging events to their handlers
	for {
		select {
		case m := <-e.ingressBuffer:
			eHandleIngress(e, m)
			funcDone()
		case m := <-e.egressBuffer:
			eHandleEgress(e, m)
			funcDone()
		}
	}
}

// handleRawIngress main function is to filter commands from raw messages.
// If a message is determined to be a command, it is parsed and structured as
// an `IngressMessage`, then passed to the ingressBuffer for further handling.
func (e *Engine) handleIngress(message msg.MessageInterface) {
	go func() {
		if !isCmd(message.GetRawMessage()) {
			return
		}

		if string(message.GetRawMessage()[0]) != adapterCmdTriggerPrefix {
			log.Warningf(
				"raw content was flagged as a command, but does not contain trigger prefix, skipping...",
			)
			log.Warning(message.GetRawMessage)
			return
		}

		cmdMatches := e.applyCmdManifests(message.GetParsedMessage())

		for _, cmdMatch := range cmdMatches {
			delegate := NewCmdDelegate(message, cmdMatch.Submatches)

			// Execute plugin command and pass delegate as an intermediary
			cmdMatch.CmdFn(&delegate)

			// If the plugin has sent a response to the delegate, let's build
			// an `EgressMessage` and push that onto the outgoing buffer for dispatch
			if delegate.response != "" {
				log.Debugf("Engine::handleIngress: plugin response is %+v", delegate.response)
				message.SetPluginResponse(delegate.response)
				log.Debugf("Engine::handleIngress: sending message %+v to egressBuffer", message)
				e.egressBuffer <- message
			}
		}
	}()
}

type cmdMatch struct {
	CmdFn      CmdFn
	Submatches map[string]string
}

// applyCmdManifests runs the content against all registered plugin `CmdLink`s
// to determine the set of plugin cmd's that need to be executed.
func (e *Engine) applyCmdManifests(content string) []cmdMatch {
	matches := make([]cmdMatch, 0)

	for _, plugin := range e.plugins {
		for _, cmdLink := range plugin.CmdManifest() {
			re := cmdLink.Regexp
			subm := re.FindStringSubmatch(content)
			args := re.SubexpNames()

			if len(subm) > 0 {
				log.Infof("Found matching plugin command manifest: %s matches %s", cmdLink, subm)
				parsedMatches := make(map[string]string)
				for i, v := range subm {
					parsedMatches[args[i]] = v
				}
				log.Debugf("Parsed matches %v", parsedMatches)
				matches = append(matches, cmdMatch{
					CmdFn:      cmdLink.CmdFn,
					Submatches: parsedMatches,
				})
			}
		}
	}

	return matches
}

// handleEgress simply routes an `EgressMessage` off the egressBuffer to an
// adapter for transmission.
func (e *Engine) handleEgress(ebm msg.MessageInterface) {
	log.Debugf("Engine::handleEgress: %+v", ebm)
	go func() {
		log.Debugf("Engine::handleEgress: sending message to adapter %+v", ebm.GetAdapterName())
		e.adapters[ebm.GetAdapterName()].SendMessage(ebm)
	}()
}

// TODO: load{Plugins,Adapters} are almost identical. Should make extension
// loading generic.
func (e *Engine) loadPlugins() error {
	c, err := config.New()
	if err != nil {
		return err
	}
	plugConf := c.GetPlugin()
	log.Infof("Loading plugins from: %v", plugConf.Dir)

	// TODO: Check that requested plugins are available in dir, log if not
	for _, pluginName := range plugConf.Enabled {
		log.Infof("-> %v", pluginName)
		pluginFile := filepath.Join(
			plugConf.Dir,
			fmt.Sprintf("%s.so", pluginName),
		)

		loadedPlugin, err := LoadPlugin(pluginFile)
		if err != nil {
			// TODO: Probably want this to be more resilient so the comm server can
			// skip problematic plugins while still loading valid ones.
			return err
		}

		e.plugins[loadedPlugin.Name()] = loadedPlugin
	}

	log.Info("Successfully loaded plugins:")
	for pluginName := range e.plugins {
		log.Infof("-> %s", pluginName)
	}

	return nil
}

func (e *Engine) loadAdapters() error {
	c, err := config.New()
	if err != nil {
		return err
	}
	adapterConf := c.GetAdapter()
	log.Infof("Loading adapters from: %v", adapterConf.Dir)

	// TODO: Check that requested adapters are available in dir, log if not
	for _, adapterName := range adapterConf.Enabled {
		log.Infof("-> %v", adapterName)
		adapterFile := filepath.Join(
			adapterConf.Dir,
			fmt.Sprintf("%s.so", adapterName),
		)
		log.Infof("file: %s", adapterFile)

		loadedAdapter, err := LoadAdapter(adapterFile)
		if err != nil {
			// TODO: Probably want this to be more resilient so the comm server can
			// skip problematic adapters while still loading valid ones.
			return err
		}

		e.adapters[loadedAdapter.Name()] = loadedAdapter
	}

	log.Info("Successfully loaded adapters:")
	for adapterName := range e.adapters {
		log.Infof("-> %s", adapterName)
	}

	return nil
}
