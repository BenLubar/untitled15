package concept

// Join is a Value that appends any number of other Value.
type Join []Value

func (j Join) Value(ctx interface{}) []interface{} {
	var values []interface{}

	for _, v := range j {
		values = append(values, v.Value(ctx)...)
	}

	return values
}
