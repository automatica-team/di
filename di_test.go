package di

type testDependency struct {
	data string
}

func (testDependency) Name() string {
	return "test/d"
}

func (testDependency) New(c C) (D, error) {
	return &testDependency{
		data: Must(c.String("data")),
	}, nil
}
