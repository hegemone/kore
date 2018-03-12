package config

import (
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	want := &Config{
		Engine: Engine{
			BufferSize: 10,
		},
		Adapters: Adapter{
			Dir:     "adapts",
			Enabled: []string{"yolos.io", "rolos.io"},
		},
		Plugins: Plugin{
			Dir:     "plugs",
			Enabled: []string{"bacon.io", "jam.io"},
		},
	}
	os.Setenv("KORE_CONFIG", "testdata/config.yaml")

	c, err := New()
	if err != nil {
		t.Errorf("unable to create config: %s", err)
	}
	if !cmp.Equal(want, c) {
		t.Errorf("have: %v; want %v", c, want)
	}

	os.Unsetenv("KORE_CONFIG")

	d, err := New()
	if err != nil {
		t.Errorf("unable to create confgi: %s", err)
	}

	want = &Config{
		Engine: Engine{
			BufferSize: 8,
		},
		Adapters: Adapter{},
		Plugins:  Plugin{},
	}

	if !cmp.Equal(want, d) {
		t.Errorf("have: %v; want %v", d, want)
	}
}

func TestGetEngine(t *testing.T) {
	want := Engine{
		BufferSize: 9,
	}

	have := &Config{
		Engine: Engine{
			BufferSize: 9,
		},
	}

	if !cmp.Equal(have.GetEngine(), want) {
		t.Errorf("have %v; want %v", have.GetEngine(), want)
	}
}

func TestGetAdapter(t *testing.T) {
	want := Adapter{
		Dir:     "adapters",
		Enabled: []string{"roloyolo.io"},
	}

	have := &Config{
		Adapters: Adapter{
			Dir:     "adapters",
			Enabled: []string{"roloyolo.io"},
		},
	}

	if !cmp.Equal(have.GetAdapter(), want) {
		t.Errorf("have %v; want %v", have.GetAdapter(), want)
	}
}

func TestGetPlugin(t *testing.T) {
	want := Plugin{
		Dir:     "plugins",
		Enabled: []string{"baconbaconbacon.io"},
	}

	have := &Config{
		Plugins: Plugin{
			Dir:     "plugins",
			Enabled: []string{"baconbaconbacon.io"},
		},
	}

	if !cmp.Equal(have.GetPlugin(), want) {
		t.Errorf("have %v; want %v", have.GetPlugin(), want)
	}
}
