package concept

// Predicate returns true or false depending on its context.
type Predicate interface {
	Test(context interface{}) bool
}
