package object

type Type string

const (
	IntegerType  = "int"
	FloatType    = "float"
	StringType   = "str"
	BooleanType  = "bool"
	NullType     = "null"
	ReturnType   = "return"
	ErrorType    = "error"
	FunctionType = "fn"
	BuiltInType  = "builtin"
	ArrayType    = "array"
	HashType     = "hash"
)

type Object interface {
	String() string
	Type() Type
	Bool() bool
	Inspect() string
}

type Comparable interface {
	Compare(other Object) int
}

type Sizeable interface {
	Len() int
}

type Immutable interface {
	Clone() Object
}

type Hashable interface {
	HashKey() HashKey
}

type BuiltinFunction func(args ...Object) Object
