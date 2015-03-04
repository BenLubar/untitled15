package concept

// Value returns zero or more values based on its context.
type Value interface {
	Value(context interface{}) []interface{}
}
