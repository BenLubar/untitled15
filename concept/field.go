package concept

import "reflect"

// Field is a Value that returns a named field of a Parent struct.
type Field struct {
	Parent Value
	Name   string
}

func (f Field) Value(ctx interface{}) []interface{} {
	values := f.Parent.Value(ctx)

	for i, v := range values {
		values[i] = reflect.Indirect(reflect.ValueOf(v)).FieldByName(f.Name).Interface()
	}

	return values
}
