package object

type Null struct{}

func (n *Null) Bool() bool {
	return false
}

func (n *Null) Compare(other Object) int {
	if _, ok := other.(*Null); ok {
		return 0
	}

	return 1
}

func (n *Null) String() string {
	return n.Inspect()
}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) Type() Type {
	return NullType
}
