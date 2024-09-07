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

		vf := v.Field(i)
		vft := vf.Type()

		optional := strings.HasPrefix(vft.String(), "di.Optional")
		if optional {
			// Override target field with internal Optional `v` field.
			vf = vf.FieldByName("v")
		}

		dep, ok := deps[name]
		if !ok && !optional {
			return fmt.Errorf("di/inject: dependency %s does not exist", name)
		}

		// Set dependency into the target field.
		ua := reflect.NewAt(vf.Type(), unsafe.Pointer(vf.UnsafeAddr()))
		ua.Elem().Set(reflect.ValueOf(dep))
	}

	return err
}
