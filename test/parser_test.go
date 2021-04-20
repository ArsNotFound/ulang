package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Ars2014/ulang/ast"
	"github.com/Ars2014/ulang/lexer"
	"github.com/Ars2014/ulang/parser"
)

func TestBlockStatement(t *testing.T) {
	input := []byte("{ a / b + k ; c+ (-5) ; !d ; a = d ; fn(a){}(a) }")
	expected := []string{"((a / b) + k)", "(c + (-5))", "(!d)", "(a = d)", "fn (a) {}(a)"}

	l := lexer.NewLexer(input)
	p := parser.NewParser()
	res, err := p.Parse(l)
	assert.NoError(t, err)
	assert.IsType(t, &ast.Program{}, res)

	prog := res.(*ast.Program)
	assert.Len(t, prog.Statements, 1)
	assert.IsType(t, &ast.BlockStatement{}, prog.Statements[0])

	block := prog.Statements[0].(*ast.BlockStatement)
	assert.Len(t, block.Statements, len(expected))
	for i, stmt := range block.Statements {
		assert.IsType(t, &ast.ExpressionStatement{}, stmt)
		assert.Equal(t, expected[i], stmt.String())
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"return test", "return test"},
		{"return", "return"},
	}

	p := parser.NewParser()
	for _, tt := range tests {
		l := lexer.NewLexer([]byte(tt.input))
		res, err := p.Parse(l)
		assert.NoError(t, err)
		assert.IsType(t, &ast.Program{}, res)

		prog := res.(*ast.Program)
		assert.Len(t, prog.Statements, 1)
		assert.IsType(t, &ast.ReturnStatement{}, prog.Statements[0])

		ret := prog.Statements[0].(*ast.ReturnStatement)
		assert.Equal(t, tt.expected, ret.String())
	}
}

func TestIndexExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedLeft  string
		expectedIndex string
	}{
		{"a[b]", "a", "b"},
		{"(a + b)[c - d]", "(a + b)", "(c - d)"},
		{"[a, b, c][0]", "[a, b, c]", "0"},
		{"(if a != b { 0 } else { 1 })[0]", "if (a != b) {\n0\n} else {\n1\n}", "0"},
	}

	p := parser.NewParser()
	for _, tt := range tests {
		l := lexer.NewLexer([]byte(tt.input))
		res, err := p.Parse(l)
		assert.NoError(t, err)
		assert.IsType(t, &ast.Program{}, res)

		prog := res.(*ast.Program)
		assert.Len(t, prog.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, prog.Statements[0])

		expr := prog.Statements[0].(*ast.ExpressionStatement).Expression
		assert.IsType(t, &ast.IndexExpression{}, expr)

		index := expr.(*ast.IndexExpression)
		assert.Equal(t, tt.expectedLeft, index.Left.String())
		assert.Equal(t, tt.expectedIndex, index.Index.String())
	}
}

func TestCallExpresion(t *testing.T) {
	tests := []struct {
		input            string
		expectedFunction string
		expectedArgs     []string
	}{
		{"a(b)", "a", []string{"b"}},
		{"f()", "f", []string{}},
		{"a(a,b)", "a", []string{"a", "b"}},
		{"fn (a,b,c) {}(1, 2, 3)", "fn (a, b, c) {}", []string{"1", "2", "3"}},
	}

	p := parser.NewParser()
	for _, tt := range tests {
		l := lexer.NewLexer([]byte(tt.input))
		res, err := p.Parse(l)
		assert.NoError(t, err)
		assert.IsType(t, &ast.Program{}, res)

		prog := res.(*ast.Program)
		assert.Len(t, prog.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, prog.Statements[0])

		expr := prog.Statements[0].(*ast.ExpressionStatement).Expression
		assert.IsType(t, &ast.CallExpression{}, expr)

		call := expr.(*ast.CallExpression)
		assert.Equal(t, tt.expectedFunction, call.Function.String())
		assert.Len(t, call.Arguments, len(tt.expectedArgs))
		for i, arg := range call.Arguments {
			assert.Equal(t, tt.expectedArgs[i], arg.String())
		}
	}
}

func TestSelectorExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedLeft  string
		expectedRight string
	}{
		{"a.b", "a", "b"},
		{"(a + c).d", "(a + c)", "d"},
	}

	p := parser.NewParser()
	for _, tt := range tests {
		l := lexer.NewLexer([]byte(tt.input))
		res, err := p.Parse(l)
		assert.NoError(t, err)
		assert.IsType(t, &ast.Program{}, res)

		prog := res.(*ast.Program)
		assert.Len(t, prog.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, prog.Statements[0])

		expr := prog.Statements[0].(*ast.ExpressionStatement).Expression
		assert.IsType(t, &ast.SelectorExpression{}, expr)

		sel := expr.(*ast.SelectorExpression)
		assert.Equal(t, tt.expectedLeft, sel.Left.String())
		assert.Equal(t, tt.expectedRight, sel.Right.String())
	}
}

func TestAssignExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedLeft  string
		expectedRight string
	}{
		{"a = b", "a", "b"},
		{"a[0] = b", "(a[0])", "b"},
		{"a.c = b", "a.c", "b"},
	}

	p := parser.NewParser()
	for _, tt := range tests {
		l := lexer.NewLexer([]byte(tt.input))
		res, err := p.Parse(l)
		assert.NoError(t, err)
		assert.IsType(t, &ast.Program{}, res)

		prog := res.(*ast.Program)
		assert.Len(t, prog.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, prog.Statements[0])

		expr := prog.Statements[0].(*ast.ExpressionStatement).Expression
		assert.IsType(t, &ast.AssignExpression{}, expr)

		assign := expr.(*ast.AssignExpression)
		assert.Equal(t, tt.expectedLeft, assign.Left.String())
		assert.Equal(t, tt.expectedRight, assign.Right.String())
	}
}

func TestForExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"for a != b { 0 }", "for (a != b) {\n0\n}"},
		{"for a = 0; b != c; c = c + 1 { }", "for (a = 0); (b != c); (c = (c + 1)) {}"},
		{"for ; b != c && c == 0; c = c + 1 { }", "for ; ((b != c) && (c == 0)); (c = (c + 1)) {}"},
		{"for {}", "for {}"},
	}

	p := parser.NewParser()
	for _, tt := range tests {
		l := lexer.NewLexer([]byte(tt.input))
		res, err := p.Parse(l)
		assert.NoError(t, err)
		assert.IsType(t, &ast.Program{}, res)

		prog := res.(*ast.Program)
		assert.Len(t, prog.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, prog.Statements[0])

		expr := prog.Statements[0].(*ast.ExpressionStatement).Expression
		assert.IsType(t, &ast.ForExpression{}, expr)

		f := expr.(*ast.ForExpression)
		assert.Equal(t, tt.expected, f.String())
	}
}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"if a != b {}", "if (a != b) {}"},
		{"if a != b {} else {}", "if (a != b) {} else {}"},
		{"if a != b {} else if b != c && a != c {} else {}", "if (a != b) {} else if ((b != c) && (a != c)) {} else {}"},
	}

	p := parser.NewParser()
	for _, tt := range tests {
		l := lexer.NewLexer([]byte(tt.input))
		res, err := p.Parse(l)
		assert.NoError(t, err)
		assert.IsType(t, &ast.Program{}, res)

		prog := res.(*ast.Program)
		assert.Len(t, prog.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, prog.Statements[0])

		expr := prog.Statements[0].(*ast.ExpressionStatement).Expression
		assert.IsType(t, &ast.IfExpression{}, expr)

		f := expr.(*ast.IfExpression)
		assert.Equal(t, tt.expected, f.String())
	}
}

func TestIdentifiers(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test", "test"},
		{"_test", "_test"},
		{"_00333", "_00333"},
	}

	p := parser.NewParser()
	for _, tt := range tests {
		l := lexer.NewLexer([]byte(tt.input))
		res, err := p.Parse(l)
		assert.NoError(t, err)
		assert.IsType(t, &ast.Program{}, res)

		prog := res.(*ast.Program)
		assert.Len(t, prog.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, prog.Statements[0])

		expr := prog.Statements[0].(*ast.ExpressionStatement).Expression
		assert.IsType(t, &ast.Identifier{}, expr)

		ident := expr.(*ast.Identifier)
		assert.Equal(t, tt.expected, ident.String())
	}
}

func TestNull(t *testing.T) {
	input := []byte("null")
	l := lexer.NewLexer(input)
	p := parser.NewParser()
	res, err := p.Parse(l)
	assert.NoError(t, err)
	assert.IsType(t, &ast.Program{}, res)

	prog := res.(*ast.Program)
	assert.Len(t, prog.Statements, 1)
	assert.IsType(t, &ast.ExpressionStatement{}, prog.Statements[0])

	expr := prog.Statements[0].(*ast.ExpressionStatement).Expression
	assert.IsType(t, &ast.Null{}, expr)
}

func TestBooleanLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	p := parser.NewParser()
	for _, tt := range tests {
		l := lexer.NewLexer([]byte(tt.input))
		res, err := p.Parse(l)
		assert.NoError(t, err)
		assert.IsType(t, &ast.Program{}, res)

		prog := res.(*ast.Program)
		assert.Len(t, prog.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, prog.Statements[0])

		expr := prog.Statements[0].(*ast.ExpressionStatement).Expression
		assert.IsType(t, &ast.BooleanLiteral{}, expr)

		boolean := expr.(*ast.BooleanLiteral)
		assert.Equal(t, tt.expected, boolean.Value)
	}
}

func TestIntegerLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"10", 10},
		{"014", 014},
		{"0xDEADBEEF", 0xDEADBEEF},
	}

	p := parser.NewParser()
	for _, tt := range tests {
		l := lexer.NewLexer([]byte(tt.input))
		res, err := p.Parse(l)
		assert.NoError(t, err)
		assert.IsType(t, &ast.Program{}, res)

		prog := res.(*ast.Program)
		assert.Len(t, prog.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, prog.Statements[0])

		expr := prog.Statements[0].(*ast.ExpressionStatement).Expression
		assert.IsType(t, &ast.IntegerLiteral{}, expr)

		integer := expr.(*ast.IntegerLiteral)
		assert.Equal(t, tt.expected, integer.Value)
	}
}

func TestFloatLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"10.10", 10.10},
		{"1.16e+10", 1.16e+10},
		{"1.16e-10", 1.16e-10},
		{"15e10", 15e10},
		{".2", .2},
		{".1e5", .1e5},
	}

	p := parser.NewParser()
	for _, tt := range tests {
		l := lexer.NewLexer([]byte(tt.input))
		res, err := p.Parse(l)
		assert.NoError(t, err)
		assert.IsType(t, &ast.Program{}, res)

		prog := res.(*ast.Program)
		assert.Len(t, prog.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, prog.Statements[0])

		expr := prog.Statements[0].(*ast.ExpressionStatement).Expression
		assert.IsType(t, &ast.FloatLiteral{}, expr)

		float := expr.(*ast.FloatLiteral)
		assert.Equal(t, tt.expected, float.Value)
	}
}

func TestStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"test\n"`, "test\n"},
		{"`test\\n`", "test\\n"},
	}

	p := parser.NewParser()
	for _, tt := range tests {
		l := lexer.NewLexer([]byte(tt.input))
		res, err := p.Parse(l)
		assert.NoError(t, err)
		assert.IsType(t, &ast.Program{}, res)

		prog := res.(*ast.Program)
		assert.Len(t, prog.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, prog.Statements[0])

		expr := prog.Statements[0].(*ast.ExpressionStatement).Expression
		assert.IsType(t, &ast.StringLiteral{}, expr)

		str := expr.(*ast.StringLiteral)
		assert.Equal(t, tt.expected, str.Value)
	}
}

func TestFunctionLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"fn (test, test1) {}", "fn (test, test1) {}"},
		{"fn () { test }", "fn () {\ntest\n}"},
	}

	p := parser.NewParser()
	for _, tt := range tests {
		l := lexer.NewLexer([]byte(tt.input))
		res, err := p.Parse(l)
		assert.NoError(t, err)
		assert.IsType(t, &ast.Program{}, res)

		prog := res.(*ast.Program)
		assert.Len(t, prog.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, prog.Statements[0])

		expr := prog.Statements[0].(*ast.ExpressionStatement).Expression
		assert.IsType(t, &ast.FunctionLiteral{}, expr)
		assert.Equal(t, tt.expected, expr.String())
	}
}

func TestArrayLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"[]", []string{}},
		{"[a, b, c]", []string{"a", "b", "c"}},
		{"[0.1, 1]", []string{"0.1", "1"}},
		{`["test", true, 1]`, []string{"\"test\"", "true", "1"}},
	}

	p := parser.NewParser()
	for _, tt := range tests {
		l := lexer.NewLexer([]byte(tt.input))
		res, err := p.Parse(l)
		assert.NoError(t, err)
		assert.IsType(t, &ast.Program{}, res)

		prog := res.(*ast.Program)
		assert.Len(t, prog.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, prog.Statements[0])

		expr := prog.Statements[0].(*ast.ExpressionStatement).Expression
		assert.IsType(t, &ast.ArrayLiteral{}, expr)

		array := expr.(*ast.ArrayLiteral)

		for i, el := range array.Elements {
			assert.Equal(t, tt.expected[i], el.String())
		}
	}
}

func TestHashLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]string
	}{
		{"{,}", map[string]string{}},
		{"{a + b : b}", map[string]string{"(a + b)": "b"}},
		{"{0.1 : t + a}", map[string]string{"0.1": "(t + a)"}},
	}

	p := parser.NewParser()
	for _, tt := range tests {
		l := lexer.NewLexer([]byte(tt.input))
		res, err := p.Parse(l)
		assert.NoError(t, err)
		assert.IsType(t, &ast.Program{}, res)

		prog := res.(*ast.Program)
		assert.Len(t, prog.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, prog.Statements[0])

		expr := prog.Statements[0].(*ast.ExpressionStatement).Expression
		assert.IsType(t, &ast.HashLiteral{}, expr)

		hash := expr.(*ast.HashLiteral)
		for k, v := range hash.Pairs {
			assert.Contains(t, tt.expected, k.String())
			assert.Equal(t, tt.expected[k.String()], v.String())
		}
	}
}
