package di

import (
	"fmt"
	"reflect"
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

		vf := v.Field(i)
		// Set dependency into the target field.
		uf := reflect.NewAt(vf.Type(), unsafe.Pointer(vf.UnsafeAddr())).Elem()
		uf.Set(reflect.ValueOf(dep))
	}

	return err
}
