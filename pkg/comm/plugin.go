package comm

import (
	goplugin "plugin"
	"regexp"

	log "github.com/sirupsen/logrus"
)

// CmdFn is the required function signature of a plugin command. It should accept
// a `*CmdDelegate` for interfacing back with the `Engine`.
type CmdFn func(*CmdDelegate)

// CmdLink pairs a regexp to a CmdFn. If the engine matches the regexp against
// an `IngressMessage`, it will route the command to the `CmdLink`'s `CmdFn`.
type CmdLink struct {
	Regexp *regexp.Regexp
	CmdFn  CmdFn
}

// Plugin is the primary abstraction representing dynamic, extensible behavioral features.
// In reality, it's a facade that presents a controller public interface for
// consumers, while delegating much of its functionality to dynamically loaded
// functions sourced from a shared library.
type Plugin interface {
	Name() string
	Help() string
	CmdManifest() []CmdLink
}

// LoadPlugin loads dynamic plugin behavior from a given .so plugin file
func LoadPlugin(pluginFile string) (Plugin, error) {
	// TODO: Need a *lot* of validation here to make sure a bad plugin doesn't
	// just crash the server.
	// -> Actually confirm the casts are valid and these functions look like they should?
	// TODO: Can the hardcoded pattern of $PROPERTY Lookup -> Cast be made more elegant?
	rawGoPlugin, err := goplugin.Open(pluginFile)
	if err != nil {
		return nil, err
	}

	pSym, err := rawGoPlugin.Lookup("Plugin")
	if err != nil {
		log.Error("Error occurred while loading plugin %s: %s", pluginFile, err.Error())
	}

	p := pSym.(Plugin)

	return p, nil
}
