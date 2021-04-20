package object

type Boolean struct {
	Value bool
}

func (b *Boolean) Bool() bool {
	return b.Value
}

func (b *Boolean) Int() int {
	if b.Value {
		return 1
	}
	return 0
}

func (b *Boolean) Compare(other Object) int {
	if obj, ok := other.(*Boolean); ok {
		return b.Int() - obj.Int()
	}
	return 1
}

func (b *Boolean) String() string {
	return b.Inspect()
}

func (b *Boolean) Inspect() string {
	if b.Value {
		return "true"
	}
	return "false"
}

func (b *Boolean) Type() Type {
	return BooleanType
}

func (b *Boolean) Clone() Object {
	return &Boolean{Value: b.Value}
}

func (b *Boolean) HashKey() HashKey {
	return HashKey{Type: b.Type(), Value: uint64(b.Int())}
}
