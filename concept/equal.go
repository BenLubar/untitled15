package concept

// Equal is a Predicate that is true if all of the Value results contain an
// identical element.
type Equal []Value

func (e Equal) Test(ctx interface{}) bool {
	values := make([][]interface{}, len(e))

	for i, v := range e {
		values[i] = v.Value(ctx)
	}

outer:
	for _, v1 := range values[0] {
	inner:
		for _, set := range values[1:] {
			for _, v2 := range set {
				if v1 == v2 {
					continue inner
				}
			}
			continue outer
		}

		return true
	}

	return false
}
