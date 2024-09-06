package di

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptional(t *testing.T) {
	type Target struct {
		Opt Optional[testDependency] `di:"test/d"`
	}
	deps := map[string]D{
		"test/d": &testDependency{data: "optional"},
	}

	t.Run("Optional is injected", func(t *testing.T) {
		var target Target
		err := inject(&target, deps)
		require.NoError(t, err)

		v, ok := target.Opt.Get()
		require.True(t, ok)
		const expected = "optional"
		require.Equal(t, expected, v.data)

		err = target.Opt.With(func(d testDependency) {
			require.Equal(t, expected, d.data)
		})
		require.NoError(t, err)

		err = target.Opt.With(func(d testDependency) error {
			require.Equal(t, expected, d.data)
			return errors.New("test")
		})
		require.ErrorContains(t, err, "test")
	})

	t.Run("Optional is not injected", func(t *testing.T) {
		var target Target
		err := inject(&target, map[string]D{})
		require.NoError(t, err)

		_, ok := target.Opt.Get()
		require.False(t, ok)

		err = target.Opt.With(func(d testDependency) {
			t.Error("should not be called")
		})
		require.NoError(t, err)

		err = target.Opt.With(func(d testDependency) error {
			t.Error("should not be called")
			return errors.New("test")
		})
		require.NoError(t, err)
	})
}
