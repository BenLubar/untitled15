package concept

// Negate is a Predicate that is true if its Parent is false.
type Negate struct {
	Parent Predicate
}

func (n Negate) Test(ctx interface{}) bool {
	return !n.Parent.Test(ctx)
}

// And is a Predicate that is true if all of its elements are true.
type And []Predicate

func (a And) Test(ctx interface{}) bool {
	for _, p := range a {
		if !p.Test(ctx) {
			return false
		}
	}

	return true
}

// Or is a Predicate that is true if any of its elements are true.
type Or []Predicate

func (o Or) Test(ctx interface{}) bool {
	for _, p := range o {
		if p.Test(ctx) {
			return true
		}
	}

	return false
}
