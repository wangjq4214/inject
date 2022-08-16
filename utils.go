package inject

import (
	"fmt"
	"reflect"
)

// Determine if v is a pointer to a structure
func isStructPointer(v reflect.Type) bool {
	return v.Kind() == reflect.Pointer && v.Elem().Kind() == reflect.Struct
}

// Determine if v is a constructor
func isFunction(v reflect.Type) bool {
	return v.Kind() == reflect.Func
}

// Determine if an object is a supported type.
func isCorrectType(o *Object, t reflect.Type) error {
	if isStructPointer(t) {
		o.t = structPointer
		return nil
	}

	if isFunction(t) {
		o.t = function
		return nil
	}

	return fmt.Errorf(
		"expected object value to be a pointer to a struct or constructor but got type %v with value %v",
		o.reflectType,
		o.Value,
	)
}
