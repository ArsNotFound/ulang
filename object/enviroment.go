package object

import "unicode"

type Environment struct {
	store  map[string]Object
	parent *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func (e *Environment) ExportedHash() *Hash {
	pairs := make(map[HashKey]HashPair)
	for k, v := range e.store {
		if unicode.IsUpper(rune(k[0])) {
			s := &String{Value: k}
			pairs[s.HashKey()] = HashPair{Key: s, Value: v}
		}
	}
	return &Hash{Pairs: pairs}
}

func (e *Environment) NewChild() *Environment {
	env := NewEnvironment()
	env.parent = e
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.parent != nil {
		obj, ok = e.parent.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
