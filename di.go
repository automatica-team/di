package di

import "errors"

// A set of frequently used aliases for the package.
type (
	D = Dependency
	C = Config
	R = Runnable
	M = map[string]any
)

type (
	// Dependency represents an injectable dependency.
	Dependency interface {
		// Name returns the name of the dependency.
		// The name is used to identify the dependency in the `di.yml`.
		Name() string
		// New creates a new instance of the dependency using the provided
		// config `di.Config`. Don't worry about filling the new instance with
		// other injectable dependencies, it will be done automatically.
		//
		// NOTE: Always return the pointer to the dependency.
		New(C) (D, error)
	}

	// Config represents the configuration of the dependency.
	// It exposes a bunch of useful methods to get the values.
	Config struct {
		name string
		m    M
	}

	// Runnable represents a final runnable instance, that uses all the
	// dependencies provided in the end. Using `di.Run` the dependencies of the
	// runnable instance will be injected automatically.
	Runnable interface {
		// Run is the main method of the runnable instance.
		Run() error
	}
)

var globalDeps []Dependency

// Inject adds a dependency to the global scope.
// Always inject dependencies you are going to use.
// Always pass the pointer to the dependency.
func Inject(d Dependency) {
	globalDeps = append(globalDeps, d)
}

// Parse parses the global `di.yml` configuration file.
// It also sets the `Version` of current dependencies state.
func Parse(path ...string) error {
	if len(path) == 0 {
		path = []string{"di.yml"}
	}
	return parse(path[0])
}

// Run runs the `di` flow injecting all the dependencies and running the
// provided runnable instance. Your runnable instance will be injected with
// dependencies as well.
func Run[R Runnable]() error {
	if global.Running {
		return errors.New("di: already running")
	}
	r, err := run[R](*new(R))
	if err != nil {
		return err
	}
	return r.Run()
}

// Idle is a dummy runnable that does nothing. Use it for initializing
// your dependencies without running anything.
type Idle struct{}

// Run implements `di.Runnable` interface.
func (Idle) Run() error { return nil }

// Get returns the dependency by its name.
// It should be already initialized.
func Get[T Dependency](name string) (T, error) {
	return globalGet[T](name)
}

// Must is a helper function when working with `di.Config` that returns zero
// value of the type if the error is not nil.
func Must[T any](v T, _ error) T {
	return v
}

// Version returns the version of the current dependencies state.
func Version() string {
	return global.Version
}
