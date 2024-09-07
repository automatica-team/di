package di

// Optional is a wrapping type that can be used to inject
// optional dependencies.
type Optional[T Dependency] struct {
	v *T
}

// Use calls the provided function with the value of the
// optional dependency if it was injected. See With for
// the usage example.
func (o Optional[T]) Use(f func(*T)) {
	if o.v != nil {
		f(o.v)
	}
}

// With calls the provided function with the value of the
// optional dependency if it was injected. This is an
// error-powered wrapper.
//
// Usage:
//
//	var s struct {
//		cache di.Optional[cache.Cache] `di:"x/cache"`
//	}
//	s.cache.With(func(c *cache.Cache) error {
//		return c.Set("key", "value")
//	})
func (o Optional[T]) With(f func(*T) error) (err error) {
	if o.v != nil {
		return f(o.v)
	}
	return nil
}

// Get returns the value of the optional dependency and a
// boolean flag indicating whether it was injected.
func (o Optional[T]) Get() (*T, bool) {
	if o.v == nil {
		return nil, false
	}
	return o.v, true
}
