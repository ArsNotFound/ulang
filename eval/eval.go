package eval

import (
	"fmt"
	"strings"

	"github.com/Ars2014/ulang/ast"
	"github.com/Ars2014/ulang/builtins"
	"github.com/Ars2014/ulang/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)

	// Statements
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.Return{Value: val}

	// Expressions
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.AssignExpression:
		return evalAssignExpression(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ForExpression:
		return evalForExpression(node, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.BooleanLiteral:
		return fromNativeBoolean(node.Value)
	case *ast.Null:
		return NULL

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Parameters, Env: env, Body: node.Body}

	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(fn, args)

	case *ast.ArrayLiteral:
		elems := evalExpressions(node.Elements, env)
		if len(elems) == 1 && isError(elems[0]) {
			return elems[0]
		}
		return &object.Array{Elements: elems}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.SelectorExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		return evalSelectorExpression(left, node.Right)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	}

	return nil
}

func evalAssignExpression(expr *ast.AssignExpression, env *object.Environment) object.Object {
	value := Eval(expr.Right, env)
	if isError(value) {
		return value
	}

	switch e := expr.Left.(type) {
	case *ast.Identifier:
		env.Set(e.Value, value)
		return value

	case *ast.IndexExpression:
		left := Eval(e.Left, env)
		if isError(left) {
			return left
		}
		switch obj := left.(type) {
		case *object.Array:
			index := Eval(e.Index, env)
			if isError(index) {
				return index
			}
			if id, ok := index.(*object.Integer); ok {
				if id.Value >= 0 && int(id.Value) < obj.Len() {
					obj.Elements[id.Value] = value
				} else {
					return newError("array[%d] index out of range: %d", obj.Len(), id.Value)
				}
			} else {
				return newError("cannot index array with %#v", index)
			}

		case *object.Hash:
			key := Eval(e.Index, env)
			if isError(key) {
				return key
			}
			if hashKey, ok := key.(object.Hashable); ok {
				hashed := hashKey.HashKey()
				obj.Pairs[hashed] = object.HashPair{Key: key, Value: value}
			} else {
				return newError("cannot index hash with %T", key)
			}

		default:
			return newError("object type %T does not support item assignment", obj)
		}

	case *ast.SelectorExpression:
		left := Eval(e.Left, env)
		if hash, ok := left.(*object.Hash); ok {
			key := &object.String{Value: e.Right.Value}
			hashed := key.HashKey()
			hash.Pairs[hashed] = object.HashPair{Key: key, Value: value}
		} else {
			return newError("object type %T does not support item assignment", left)
		}

	default:
		return newError("expected identifier, index expression or selector expression got=%T", e)
	}

	return NULL
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.Return:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.ReturnType || rt == object.ErrorType {
				return result
			}
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch r := right.(type) {
	case *object.Integer:
		return evalIntegerPrefixOperatorExpression(operator, r)
	case *object.Float:
		return evalFloatPrefixOperatorExpression(operator, r)
	case *object.Boolean:
		return evalBooleanPrefixOperatorExpression(operator, r)
	default:
		return newError("%s doesn't support prefix operators", right.Type())
	}
}

func evalIntegerPrefixOperatorExpression(operator string, right *object.Integer) object.Object {
	value := right.Value

	switch operator {
	case "+":
		return right
	case "!":
		return fromNativeBoolean(!right.Bool())
	case "~":
		return &object.Integer{Value: ^value}
	case "-":
		return &object.Integer{Value: -value}
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalFloatPrefixOperatorExpression(operator string, right *object.Float) object.Object {
	value := right.Value

	switch operator {
	case "+":
		return right
	case "!":
		return fromNativeBoolean(!right.Bool())
	case "-":
		return &object.Float{Value: -value}
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBooleanPrefixOperatorExpression(operator string, right *object.Boolean) object.Object {
	switch operator {
	case "!":
		return fromNativeBoolean(!right.Bool())
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case operator == "+" && left.Type() == object.HashType && right.Type() == object.HashType:
		leftVal := left.(*object.Hash).Pairs
		rightVal := right.(*object.Hash).Pairs
		pairs := make(map[object.HashKey]object.HashPair)
		for k, v := range leftVal {
			pairs[k] = v
		}
		for k, v := range rightVal {
			pairs[k] = v
		}
		return &object.Hash{Pairs: pairs}

	case operator == "+" && left.Type() == object.ArrayType && right.Type() == object.ArrayType:
		leftVal := left.(*object.Array).Elements
		rightVal := right.(*object.Array).Elements
		elements := make([]object.Object, len(leftVal)+len(rightVal))
		elements = append(leftVal, rightVal...)
		return &object.Array{Elements: elements}

	case operator == "*" && left.Type() == object.ArrayType && right.Type() == object.IntegerType:
		leftVal := left.(*object.Array).Elements
		rightVal := right.(*object.Integer).Value
		elements := leftVal
		for i := rightVal; i > 1; i-- {
			elements = append(elements, leftVal...)
		}
		return &object.Array{Elements: elements}

	case operator == "*" && left.Type() == object.IntegerType && right.Type() == object.ArrayType:
		rightVal := right.(*object.Array).Elements
		leftVal := left.(*object.Integer).Value
		elements := rightVal
		for i := leftVal; i > 1; i-- {
			elements = append(elements, rightVal...)
		}
		return &object.Array{Elements: elements}

	case operator == "*" && left.Type() == object.StringType && right.Type() == object.IntegerType:
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.Integer).Value
		return &object.String{Value: strings.Repeat(leftVal, int(rightVal))}

	case operator == "*" && left.Type() == object.IntegerType && right.Type() == object.StringType:
		rightVal := right.(*object.String).Value
		leftVal := left.(*object.Integer).Value
		return &object.String{Value: strings.Repeat(rightVal, int(leftVal))}

	case left.Type() == object.IntegerType && right.Type() == object.FloatType:
		fLeft := &object.Float{Value: float64(left.(*object.Integer).Value)}
		return evalFloatInfixExpression(operator, fLeft, right.(*object.Float))

	case left.Type() == object.FloatType && right.Type() == object.IntegerType:
		fRight := &object.Float{Value: float64(right.(*object.Integer).Value)}
		return evalFloatInfixExpression(operator, left.(*object.Float), fRight)

	case operator == "==":
		return fromNativeBoolean(left.(object.Comparable).Compare(right) == 0)
	case operator == "!=":
		return fromNativeBoolean(left.(object.Comparable).Compare(right) != 0)
	case operator == "<=":
		return fromNativeBoolean(left.(object.Comparable).Compare(right) < 1)
	case operator == ">=":
		return fromNativeBoolean(left.(object.Comparable).Compare(right) > -1)
	case operator == "<":
		return fromNativeBoolean(left.(object.Comparable).Compare(right) == -1)
	case operator == ">":
		return fromNativeBoolean(left.(object.Comparable).Compare(right) == 1)

	case left.Type() == right.Type():
		switch left.Type() {
		case object.BooleanType:
			return evalBooleanInfixExpression(operator, left.(*object.Boolean), right.(*object.Boolean))
		case object.IntegerType:
			return evalIntegerInfixExpression(operator, left.(*object.Integer), right.(*object.Integer))
		case object.FloatType:
			return evalFloatInfixExpression(operator, left.(*object.Float), right.(*object.Float))
		case object.StringType:
			return evalStringInfixExpression(operator, left.(*object.String), right.(*object.String))
		}
		fallthrough

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

}

func evalBooleanInfixExpression(operator string, left *object.Boolean, right *object.Boolean) object.Object {
	leftVal := left.Value
	rightVal := right.Value

	switch operator {
	case "&&":
		return fromNativeBoolean(leftVal && rightVal)
	case "||":
		return fromNativeBoolean(leftVal || rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left *object.Integer, right *object.Integer) object.Object {
	leftVal := left.Value
	rightVal := right.Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero: %d %s %d", leftVal, operator, rightVal)
		}
		return &object.Float{Value: float64(leftVal) / float64(rightVal)}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "|":
		return &object.Integer{Value: leftVal | rightVal}
	case "^":
		return &object.Integer{Value: leftVal ^ rightVal}
	case "&":
		return &object.Integer{Value: leftVal & rightVal}
	case "<<":
		return &object.Integer{Value: leftVal << rightVal}
	case ">>":
		return &object.Integer{Value: leftVal >> rightVal}
	case "==":
		return fromNativeBoolean(left.Compare(right) == 0)
	case "!=":
		return fromNativeBoolean(left.Compare(right) != 0)
	case "<=":
		return fromNativeBoolean(left.Compare(right) < 1)
	case ">=":
		return fromNativeBoolean(left.Compare(right) > -1)
	case "<":
		return fromNativeBoolean(left.Compare(right) == -1)
	case ">":
		return fromNativeBoolean(left.Compare(right) == 1)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left *object.Float, right *object.Float) object.Object {
	leftVal := left.Value
	rightVal := right.Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0.0 {
			return newError("division by zero: %f %s %f", leftVal, operator, rightVal)
		}
		return &object.Float{Value: leftVal / rightVal}
	case "==":
		return fromNativeBoolean(left.Compare(right) == 0)
	case "!=":
		return fromNativeBoolean(left.Compare(right) != 0)
	case "<=":
		return fromNativeBoolean(left.Compare(right) < 1)
	case ">=":
		return fromNativeBoolean(left.Compare(right) > -1)
	case "<":
		return fromNativeBoolean(left.Compare(right) == -1)
	case ">":
		return fromNativeBoolean(left.Compare(right) == 1)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left *object.String, right *object.String) object.Object {
	leftVal := left.Value
	rightVal := right.Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(expr *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(expr.Condition, env)
	if isError(condition) {
		return condition
	}

	if condition.Bool() {
		return Eval(expr.Consequence, env)
	} else if expr.Alternative != nil {
		return Eval(expr.Alternative, env)
	} else {
		return NULL
	}
}

func evalForExpression(expr *ast.ForExpression, env *object.Environment) object.Object {
	var result object.Object

	if expr.Initializer != nil {
		init := Eval(expr.Initializer, env)
		if isError(init) {
			return init
		}
	}

	var condition ast.Expression
	if expr.Condition != nil {
		condition = expr.Condition
	} else {
		condition = &ast.BooleanLiteral{Value: true}
	}

	var counter ast.Expression
	if expr.Counter != nil {
		counter = expr.Counter
	} else {
		counter = &ast.BooleanLiteral{Value: false}
	}

	for {
		cond := Eval(condition, env)
		if isError(cond) {
			return cond
		}

		if cond.Bool() {
			result = Eval(expr.Consequence, env)
		} else {
			break
		}

		Eval(counter, env)
	}

	if result != nil {
		return result
	}

	return NULL
}

func evalIdentifier(ident *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(ident.Value); ok {
		return val
	}

	if builtin, ok := builtins.Builtins[ident.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + ident.Value)
}

func evalExpressions(list ast.ExpressionList, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range list {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		env := extendFunctionEnv(fn, args)
		return unwrapReturnValue(Eval(fn.Body, env))

	case *object.Builtin:
		if result := fn.Fn(args...); result != nil {
			return result
		}
		return NULL

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := fn.Env.NewChild()

	for paramId, param := range fn.Parameters {
		env.Set(param.Value, args[paramId])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.Return); ok {
		return returnValue.Value
	}

	return obj
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch {
	case left.Type() == object.StringType && index.Type() == object.IntegerType:
		return evalStringIndexExpression(left.(*object.String), index.(*object.Integer))
	case left.Type() == object.ArrayType && index.Type() == object.IntegerType:
		return evalArrayIndexExpression(left.(*object.Array), index.(*object.Integer))
	case left.Type() == object.HashType:
		return evalHashIndexExpression(left.(*object.Hash), index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalStringIndexExpression(str *object.String, index *object.Integer) object.Object {
	idx := index.Value
	max := int64(len(str.Value) - 1)

	if idx < 0 || idx > max {
		return &object.String{Value: ""}
	}

	return &object.String{Value: string(str.Value[idx])}
}

func evalArrayIndexExpression(array *object.Array, index *object.Integer) object.Object {
	idx := index.Value
	max := int64(len(array.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return array.Elements[idx]
}

func evalHashIndexExpression(hash *object.Hash, index object.Object) object.Object {
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hash.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalSelectorExpression(left object.Object, right *ast.Identifier) object.Object {
	hash, ok := left.(*object.Hash)
	if !ok {
		return newError("%s does not support selection", left.Type())
	}

	index := &object.String{Value: right.Value}
	key := index.HashKey()

	pair, ok := hash.Pairs[key]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func fromNativeBoolean(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorType
	}
	return false
}
