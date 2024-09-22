package di

type composer struct{}

func (c composer) Run() error { return nil }

func compose(deps ...Dependency) error {
	for _, dep := range deps {
		globalDeps = append(globalDeps, dep)
	}
	_, err := run(composer{})
	return err
}
