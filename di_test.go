package di

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testDependency struct {
	data string
}

func (testDependency) Name() string {
	return "test/d"
}

func (testDependency) New(c C) (D, error) {
	return testDependency{
		data: Must(c.String("data")),
	}, nil
}

func TestNew(t *testing.T) {
	global.Deps.Set("test/d", M{"data": "test"})

	d, err := New[testDependency]("test/d")
	require.ErrorContains(t, err, "test/d not found")

	Inject(testDependency{})

	d, err = New[testDependency]("test/d")
	require.NotNil(t, d)
	require.NoError(t, err)
	require.Equal(t, "test/d", d.Name())
	require.Equal(t, "test", d.data)
}
