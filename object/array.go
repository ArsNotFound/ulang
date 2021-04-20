package object

import "strings"

type Array struct {
	Elements []Object
}

func (a *Array) Bool() bool {
	return len(a.Elements) > 0
}

func (a *Array) PopLeft() Object {
	if len(a.Elements) > 0 {
		e := a.Elements[0]
		a.Elements = a.Elements[1:]
		return e
	}
	return &Null{}
}

func (a *Array) PopRight() Object {
	if len(a.Elements) > 0 {
		e := a.Elements[(len(a.Elements) - 1)]
		a.Elements = a.Elements[:(len(a.Elements) - 1)]
		return e
	}
	return &Null{}
}

func (a *Array) Prepend(obj Object) {
	a.Elements = append([]Object{obj}, a.Elements...)
}

func (a *Array) Append(obj Object) {
	a.Elements = append(a.Elements, obj)
}

func (a *Array) Copy() *Array {
	elements := make([]Object, len(a.Elements))
	for i, e := range a.Elements {
		elements[i] = e
	}
	return &Array{Elements: elements}
}

func (a *Array) Reverse() {
	for i, j := 0, len(a.Elements)-1; i < j; i, j = i+1, j-1 {
		a.Elements[i], a.Elements[j] = a.Elements[j], a.Elements[i]
	}
}

func (a *Array) Len() int {
	return len(a.Elements)
}

func (a *Array) Swap(i, j int) {
	a.Elements[i], a.Elements[j] = a.Elements[j], a.Elements[i]
}

func (a *Array) Less(i, j int) bool {
	if cmp, ok := a.Elements[i].(Comparable); ok {
		return cmp.Compare(a.Elements[j]) == -1
	}
	return false
}

func (a *Array) Compare(other Object) int {
	if obj, ok := other.(*Array); ok {
		if len(a.Elements) != len(obj.Elements) {
			return -1
		}
		for i, el := range a.Elements {
			cmp, ok := el.(Comparable)
			if !ok {
				return -1
			}
			if cmp.Compare(obj.Elements[i]) != 0 {
				return cmp.Compare(obj.Elements[i])
			}
		}

		return 0
	}

	return -1
}

func (a *Array) String() string {
	return a.Inspect()
}

func (a *Array) Inspect() string {
	var out strings.Builder

	var elems []string
	for _, e := range a.Elements {
		elems = append(elems, e.Inspect())
	}

	out.WriteRune('[')
	out.WriteString(strings.Join(elems, ", "))
	out.WriteRune(']')

	return out.String()
}

func (a *Array) Type() Type {
	return ArrayType
}
