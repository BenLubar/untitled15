package concept

import "reflect"

// Slice is a Value that expands a slice or array to its elements.
type Slice struct {
	Parent Value
}

func (s Slice) Value(ctx interface{}) []interface{} {
	var values []interface{}

	for _, slice := range s.Parent.Value(ctx) {
		v := reflect.ValueOf(slice)
		for i, l := 0, v.Len(); i < l; i++ {
			values = append(values, v.Index(i).Interface())
		}
	}

	return values
}
