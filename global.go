package di

import (
	"fmt"
	"os"

	orderedmap "github.com/wk8/go-ordered-map/v2"
	"gopkg.in/yaml.v3"
)

var global = struct {
	// Running disables the user to Run something more than once.
	Running bool

	// Version is the version of the dependencies state.
	Version string `yaml:"version"`

	// Deps is a map of dependencies and their configs.
	Deps *orderedmap.OrderedMap[string, M] `yaml:"di"`

	// Imports is a list of import paths to found injected
	// external dependencies from.
	Imports []string `yaml:"imports"`
}{
	Version: "0",
	Deps:    orderedmap.New[string, M](),
}

func parse(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("di/parse: %w", err)
	}
	defer f.Close()

	d := yaml.NewDecoder(f)
	if err := d.Decode(&global); err != nil {
		return fmt.Errorf("di/parse: %w", err)
	}

	return nil
}

func globalConfig(name string) Config {
	m, ok := global.Deps.Get(name)
	if !ok {
		return emptyConfig
	}
	return Config{name: name, m: m}
}

func globalGet[T Dependency](name string) (zero T, _ error) {
	deps := make(map[string]D)
	for _, dep := range globalDeps {
		deps[dep.Name()] = dep
	}

	d, ok := deps[name]
	if !ok {
		return zero, fmt.Errorf("di: dependency %s not found", name)
	}

	return d.(T), nil
}
