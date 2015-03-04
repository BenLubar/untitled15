package concept

import (
	"fmt"

	"github.com/BenLubar/untitled15/enum"
)

// Enum is a Value that returns a named enumeration value.
type Enum struct {
	Name string
	Key  string
}

func (e Enum) Value(ctx interface{}) []interface{} {
	v, ok := enum.Lookup(e.Name).New().Parse(e.Key)
	if !ok {
		panic(fmt.Errorf("concept: enumeration %q does not have a value named %q", e.Name, e.Key))
	}

	return []interface{}{v}
}
