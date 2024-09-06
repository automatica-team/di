package di

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

func inject(target any, deps map[string]D) (err error) {
	if len(deps) == 0 {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("di/inject: %v", r)
		}
	}()

	v := reflect.ValueOf(target).Elem()
	vt := v.Type()

	for i := 0; i < v.NumField(); i++ {
		name := vt.Field(i).Tag.Get("di")
		if name == "" {
			continue
		}

		dep, ok := deps[name]
		if !ok {
			return fmt.Errorf("di/inject: dependency %s not exists", name)
		}

		// TODO: compare dependency type with the target field type?

		vf := v.Field(i)
		vft := vf.Type()

		switch {
		// If the target field is an Optional type.
		case strings.HasPrefix(vft.String(), "di.Optional"):
			of := vf.FieldByName("v")
			// Set dependency into the Optional's v field.
			ua := reflect.NewAt(of.Type(), unsafe.Pointer(of.UnsafeAddr()))
			ua.Elem().Set(reflect.ValueOf(dep))
		default:
			// Set dependency into the target field.
			ua := reflect.NewAt(vft, unsafe.Pointer(vf.UnsafeAddr()))
			ua.Elem().Set(reflect.ValueOf(dep))
		}
	}

	return err
}
