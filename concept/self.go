package concept

// Self is a Value that returns the context.
var Self self

type self struct{}

func (self) Value(ctx interface{}) []interface{} {
	return []interface{}{ctx}
}
