package di

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"

	. "github.com/wk8/go-ordered-map/v2"
)

// Any retrieves the value for the given key.
// It returns an error if the key is missing.
func (c Config) Any(key string) (any, error) {
	v, ok := c.m[key]
	if !ok {
		return nil, fmt.Errorf(`di: "%s.%s" is required`, c.name, key)
	}
	return v, nil
}

// EnvString retrieves the value for the given key as a string,
// treating values starting with "$" as environment variables.
// It returns an error if the environment variable is empty.
func (c Config) EnvString(key string) (string, error) {
	s, err := c.String(key)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(s, "$") {
		return s, nil
	}
	s = os.Getenv(s[1:])
	if s == "" {
		return "", fmt.Errorf(`di: "%s.%s" env is empty`, c.name, key)
	}
	return s, nil
}

// String retrieves the value for the given key as a string.
func (c Config) String(key string) (string, error) {
	v, err := c.Any(key)
	return cast.ToString(v), err
}

// Bool retrieves the value for the given key as a boolean.
func (c Config) Bool(key string) (bool, error) {
	v, err := c.Any(key)
	return cast.ToBool(v), err
}

// Int retrieves the value for the given key as an integer.
func (c Config) Int(key string) (int, error) {
	v, err := c.Any(key)
	return cast.ToInt(v), err
}

// Float retrieves the value for the given key as a float64.
func (c Config) Float(key string) (float64, error) {
	v, err := c.Any(key)
	return cast.ToFloat64(v), err
}

// Duration retrieves the value for the given key as a `time.Duration`.
func (c Config) Duration(key string) (time.Duration, error) {
	v, err := c.Any(key)
	return cast.ToDuration(v), err
}

var diGlobal struct {
	Deps    OrderedMap[string, M] `yaml:"di"`
	Imports []string              `yaml:"imports"`
	Version string                `yaml:"version"`
	parsed  bool
}

func parse(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("di/parse: %w", err)
	}
	defer f.Close()

	d := yaml.NewDecoder(f)
	if err := d.Decode(&diGlobal); err != nil {
		return fmt.Errorf("di/parse: %w", err)
	}

	diGlobal.parsed = true
	return nil
}

var emptyConfig = Config{m: make(map[string]any)}

func pickConfig(name string) Config {
	if !diGlobal.parsed {
		return emptyConfig
	}

	m, ok := diGlobal.Deps.Get(name)
	if !ok {
		return emptyConfig
	}

	return Config{name: name, m: m}
}
