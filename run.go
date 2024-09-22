package di

func run[R Runnable](r R) (Runnable, error) {
	injected := make(map[string]Dependency)
	for i := range globalDeps {
		var (
			dep  = globalDeps[i]
			name = dep.Name()
		)

		// It is expected by the implementation that dependencies are
		// already injected for the New constructor. Thus, inject them.
		if err := inject(dep, injected); err != nil {
			return nil, err
		}

		idep, err := dep.New(globalConfig(name))
		if err != nil {
			return nil, err
		}

		if dep != idep {
			// It is allowed for implementation to return a newly allocated
			// instance of the dependency, which means it won't have previously
			// injected fields. In this case, we need to inject them again.
			if err := inject(idep, injected); err != nil {
				return nil, err
			}
		}

		injected[name] = idep

		// Update the original dependency from global state
		// with initialized one.
		globalDeps[i] = idep
	}
	return r, inject(&r, injected)
}
