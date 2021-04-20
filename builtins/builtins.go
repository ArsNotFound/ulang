package builtins

import (
	"fmt"
	"sort"

	. "github.com/Ars2014/ulang/object"
)

var Builtins = map[string]*Builtin{
	"abs":      {Name: "abs", Fn: Abs},
	"args":     {Name: "args", Fn: Args},
	"assert":   {Name: "assert", Fn: Assert},
	"bin":      {Name: "bin", Fn: Bin},
	"bool":     {Name: "bool", Fn: Bool},
	"chr":      {Name: "chr", Fn: Chr},
	"divmod":   {Name: "divmod", Fn: Divmod},
	"exit":     {Name: "exit", Fn: Exit},
	"find":     {Name: "find", Fn: Find},
	"first":    {Name: "first", Fn: First},
	"float":    {Name: "float", Fn: ToFloat},
	"hash":     {Name: "hash", Fn: HashOf},
	"hex":      {Name: "hex", Fn: Hex},
	"id":       {Name: "id", Fn: IdOf},
	"input":    {Name: "input", Fn: Input},
	"int":      {Name: "int", Fn: Int},
	"join":     {Name: "join", Fn: Join},
	"last":     {Name: "last", Fn: Last},
	"len":      {Name: "len", Fn: Len},
	"lower":    {Name: "lower", Fn: Lower},
	"max":      {Name: "max", Fn: Max},
	"min":      {Name: "min", Fn: Min},
	"oct":      {Name: "oct", Fn: Oct},
	"ord":      {Name: "ord", Fn: Ord},
	"pop":      {Name: "pop", Fn: Pop},
	"pow":      {Name: "pow", Fn: Pow},
	"print":    {Name: "print", Fn: Print},
	"push":     {Name: "push", Fn: Push},
	"rest":     {Name: "rest", Fn: Rest},
	"reversed": {Name: "reversed", Fn: Reversed},
	"sorted":   {Name: "sorted", Fn: Sorted},
	"split":    {Name: "split", Fn: Split},
	"str":      {Name: "str", Fn: Str},
	"typeof":   {Name: "typeof", Fn: TypeOf},
	"upper":    {Name: "upper", Fn: Upper},
}

var BuiltinsIndex []*Builtin

func init() {
	var keys []string
	for k := range Builtins {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		BuiltinsIndex = append(BuiltinsIndex, Builtins[k])
	}
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
