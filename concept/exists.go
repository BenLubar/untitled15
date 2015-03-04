package concept

// Exists is a Predicate that is true if the Value has a result.
type Exists struct {
	Value Value
}

func (e Exists) Test(ctx interface{}) bool {
	return len(e.Value.Value(ctx)) != 0
}
