// config is responsible for managing the config of the system and presenting
// that to comm in a controlled way. For the purposes of the POC, it is
// mostly hardcoded in its implementation. In practice, this pkg should handle
// things like validation, and file loading/parsing, in addition to env var
// loading for things like secret credentials.
package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const (
	ConfigEnvVar     = "KORE_CONFIG"
	PluginDirEnvVar  = "KORE_PLUGIN_DIR"
	AdapterDirEnvVar = "KORE_ADAPTER_DIR"
)

type Engine struct {
	BufferSize uint
}

type Plugin struct {
	Dir     string
	Enabled []string
}

type Adapter struct {
	Dir     string
	Enabled []string
}

type Config struct {
	Engine   `yaml:"engine"`
	Plugins  Plugin  `yaml:"plugins"`
	Adapters Adapter `yaml:"adapters"`
}

func New() (*Config, error) {
	c := Config{
		Engine: Engine{
			BufferSize: 8,
		},
		Plugins:  Plugin{},
		Adapters: Adapter{},
	}

	var f []byte
	var err error

	if env, exists := os.LookupEnv(ConfigEnvVar); exists {
		f, err = ioutil.ReadFile(env)
		if err != nil {
			return nil, err
		}
	} else {
		f, err = ioutil.ReadFile("config.yaml")
		if os.IsNotExist(err) {
			return &c, nil
		}
		if err != nil {
			return nil, err
		}
	}

	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal yaml: %s", err)
	}

	return &c, nil
}

func (c *Config) GetEngine() Engine {
	return c.Engine
}

func (c *Config) GetPlugin() Plugin {
	return c.Plugins
}

func (c *Config) GetAdapter() Adapter {
	return c.Adapters
}
