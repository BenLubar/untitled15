package enum

import (
	"fmt"
	"reflect"
	"strconv"
)

type Value interface {
	Parse(string) (v Value, ok bool)
	fmt.Stringer
}

type Enum struct {
	name    string
	forward []string
	reverse map[string]uint64
	t       reflect.Type
}

var registered = make(map[string]*Enum)

func Lookup(name string) *Enum {
	return registered[name]
}

func Register(zero Value, name string, values []string) *Enum {
	e := &Enum{
		name:    name,
		forward: values,
		reverse: make(map[string]uint64, len(values)),
		t:       reflect.TypeOf(zero),
	}

	for i, s := range values {
		if n, ok := e.reverse[s]; ok {
			panic(fmt.Errorf("enum: duplicate value in enumeration %q: %d and %d are both %q", name, n, i, s))
		}
		e.reverse[s] = uint64(i)
	}

	if _, ok := registered[name]; ok {
		panic(fmt.Errorf("enum: duplicate registration for %q", name))
	}
	registered[name] = e

	return e
}

func (e *Enum) New() Value {
	return reflect.New(e.t).Elem().Interface().(Value)
}

func (e *Enum) Parse(s string, bits int) (n uint64, ok bool) {
	if n, ok = e.reverse[s]; ok {
		return
	}

	n, err := strconv.ParseUint(s, 10, bits)
	if err == nil && n >= uint64(len(e.forward)) {
		ok = true
	}
	return
}

func (e *Enum) String(n uint64) string {
	if uint64(len(e.forward)) > n {
		return e.forward[n]
	}
	return strconv.FormatUint(n, 10)
}
