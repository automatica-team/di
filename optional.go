package di

// Optional is a wrapping type that can be used to inject
// optional dependencies.
type Optional[T any] struct {
	v *T
}

// With calls the provided function with the value of the
// optional dependency if it was injected. If the function
// returns an error, it is returned by the method.
//
// Usage:
//
//	s.cache.With(func(c cache.Cache) {}) // OR:
//	err := s.cache.With(func(c cache.Cache) error {})
func (o Optional[T]) With(f any) (err error) {
	if o.v == nil {
		return nil
	}
	switch f := f.(type) {
	case func(T):
		f(*o.v)
	case func(T) error:
		return f(*o.v)
	}
	return err
}

// Get returns the value of the optional dependency and a
// boolean flag indicating whether it was injected.
func (o Optional[T]) Get() (T, bool) {
	if o.v == nil {
		var zero T
		return zero, false
	}
	return *o.v, true
}
