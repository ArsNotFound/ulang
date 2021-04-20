package ast

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/Ars2014/ulang/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type StatementList []Statement

func NewStatementList(stmt Statement) (StatementList, error) {
	return StatementList{stmt}, nil
}

func AppendStatement(sl StatementList, s Statement) (StatementList, error) {
	return append(sl, s), nil
}

type Expression interface {
	Node
	expressionNode()
}

type ExpressionList []Expression

func NewExpressionList(expr Expression) (ExpressionList, error) {
	if expr == nil {
		return ExpressionList{}, nil
	}
	return ExpressionList{expr}, nil
}

func AppendExpression(exprList ExpressionList, expr Expression) (ExpressionList, error) {
	return append(exprList, expr), nil
}

type ExpressionMap map[Expression]Expression

func NewExpressionMap(key, value Expression) (ExpressionMap, error) {
	pairs := make(ExpressionMap)
	if key != nil {
		pairs[key] = value
	}

	return pairs, nil
}

func AppendExpressionPair(pairs ExpressionMap, key, value Expression) (ExpressionMap, error) {
	pairs[key] = value
	return pairs, nil
}

type Program struct {
	Statements StatementList
}

func NewProgram(stmtList StatementList) (*Program, error) {
	return &Program{Statements: stmtList}, nil
}

func (p *Program) TokenLiteral() string {
	return "PROGRAM"
}

func (p *Program) String() string {
	var out strings.Builder

	var statements []string
	for _, s := range p.Statements {
		statements = append(statements, s.String())
	}

	if len(statements) > 0 {
		out.WriteString(strings.Join(statements, ";\n"))
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements StatementList
}

func NewBlockStatement(t *token.Token, stmtList StatementList) (*BlockStatement, error) {
	return &BlockStatement{Token: *t, Statements: stmtList}, nil
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return string(bs.Token.Lit) }
func (bs *BlockStatement) String() string {
	var out strings.Builder

	var statements []string
	for _, s := range bs.Statements {
		statements = append(statements, s.String())
	}

	out.WriteRune('{')
	if len(statements) > 0 {
		out.WriteRune('\n')
		out.WriteString(strings.Join(statements, ";\n"))
		out.WriteRune('\n')
	}
	out.WriteRune('}')

	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func NewReturnStatement(t *token.Token, rv Expression) (*ReturnStatement, error) {
	return &ReturnStatement{Token: *t, ReturnValue: rv}, nil
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return string(rs.Token.Lit) }
func (rs *ReturnStatement) String() string {
	var out strings.Builder

	out.WriteString(rs.TokenLiteral())

	if rs.ReturnValue != nil {
		out.WriteString(" " + rs.ReturnValue.String())
	}

	return out.String()
}

type ExpressionStatement struct {
	Expression Expression
}

func NewExpressionStatement(expr Expression) (*ExpressionStatement, error) {
	return &ExpressionStatement{Expression: expr}, nil
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	if es.Expression != nil {
		return es.Expression.TokenLiteral()
	}

	return ""
}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func NewPrefixExpression(t *token.Token, right Expression) (*PrefixExpression, error) {
	return &PrefixExpression{Token: *t, Operator: string(t.Lit), Right: right}, nil
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return string(pe.Token.Lit) }
func (pe *PrefixExpression) String() string {
	var out strings.Builder

	out.WriteRune('(')
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteRune(')')

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func NewInfixExpression(left Expression, t *token.Token, right Expression) (*InfixExpression, error) {
	return &InfixExpression{Left: left, Token: *t, Operator: string(t.Lit), Right: right}, nil
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return string(ie.Token.Lit) }
func (ie *InfixExpression) String() string {
	var out strings.Builder

	out.WriteRune('(')
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteRune(')')

	return out.String()
}

type AssignExpression struct {
	Token token.Token
	Left  Expression // Identifier, IndexExpression or SelectorExpression
	Right Expression
}

func NewAssignExpression(left Expression, t *token.Token, right Expression) (*AssignExpression, error) {
	return &AssignExpression{Token: *t, Left: left, Right: right}, nil
}

func (ae *AssignExpression) expressionNode()      {}
func (ae *AssignExpression) TokenLiteral() string { return string(ae.Token.Lit) }
func (ae *AssignExpression) String() string {
	var out strings.Builder

	out.WriteRune('(')
	out.WriteString(ae.Left.String())
	out.WriteString(" " + ae.TokenLiteral() + " ")
	out.WriteString(ae.Right.String())
	out.WriteRune(')')

	return out.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func NewIfExpression(t *token.Token, cond Expression, conseq *BlockStatement, alter *BlockStatement) (*IfExpression, error) {
	return &IfExpression{Token: *t, Condition: cond, Consequence: conseq, Alternative: alter}, nil
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return string(ie.Token.Lit) }
func (ie *IfExpression) String() string {
	var out strings.Builder

	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteRune(' ')
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString(" else ")
		if len(ie.Alternative.Statements) == 1 {
			if expr, ok := ie.Alternative.Statements[0].(*ExpressionStatement); ok {
				if ifExpr, ok := expr.Expression.(*IfExpression); ok {
					out.WriteString(ifExpr.String())
					goto exit
				}
			}
		}
		out.WriteString(ie.Alternative.String())
	}

exit:
	return out.String()
}

type ForExpression struct {
	Token       token.Token
	Initializer Expression
	Condition   Expression
	Counter     Expression
	Consequence *BlockStatement
}

func NewForExpression(t *token.Token, init Expression, cond Expression, count Expression, conseq *BlockStatement) (*ForExpression, error) {
	return &ForExpression{Token: *t, Initializer: init, Condition: cond, Counter: count, Consequence: conseq}, nil
}

func (fe *ForExpression) expressionNode()      {}
func (fe *ForExpression) TokenLiteral() string { return string(fe.Token.Lit) }
func (fe *ForExpression) String() string {
	var out strings.Builder

	out.WriteString("for ")

	if fe.Initializer == nil && fe.Counter == nil { // while-like syntax
		if fe.Condition != nil {
			out.WriteString(fe.Condition.String() + " ")
		}
	} else {
		if fe.Initializer != nil {
			out.WriteString(fe.Initializer.String())
		}
		out.WriteRune(';')

		if fe.Condition != nil {
			out.WriteString(" " + fe.Condition.String())
		}
		out.WriteRune(';')

		if fe.Counter != nil {
			out.WriteString(" " + fe.Counter.String())
		}
		out.WriteRune(' ')
	}

	out.WriteString(fe.Consequence.String())

	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func NewIdentifier(t *token.Token) (*Identifier, error) {
	return &Identifier{Token: *t, Value: string(t.Lit)}, nil
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return string(i.Token.Lit) }
func (i *Identifier) String() string       { return i.Value }

type IdentifierList []*Identifier

func NewIdentifierList(i *Identifier) (IdentifierList, error) {
	if i == nil {
		return IdentifierList{}, nil
	}

	return IdentifierList{i}, nil
}

func AppendIdentifier(il IdentifierList, i *Identifier) (IdentifierList, error) {
	return append(il, i), nil
}

type Null struct {
	Token token.Token
}

func NewNull(t *token.Token) (*Null, error) {
	return &Null{Token: *t}, nil
}

func (n *Null) expressionNode()      {}
func (n *Null) TokenLiteral() string { return string(n.Token.Lit) }
func (n *Null) String() string       { return n.TokenLiteral() }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func NewBooleanLiteral(t *token.Token) (*BooleanLiteral, error) {
	return &BooleanLiteral{Token: *t, Value: bytes.Equal(t.Lit, []byte("true"))}, nil
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return string(bl.Token.Lit) }
func (bl *BooleanLiteral) String() string       { return bl.TokenLiteral() }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func NewIntegerLiteral(t *token.Token) (*IntegerLiteral, error) {
	v, err := strconv.ParseInt(string(t.Lit), 0, 64)
	if err != nil {
		return nil, err
	}

	return &IntegerLiteral{Token: *t, Value: v}, nil
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return string(il.Token.Lit) }
func (il *IntegerLiteral) String() string       { return il.TokenLiteral() }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func NewFloatLiteral(t *token.Token) (*FloatLiteral, error) {
	v, err := strconv.ParseFloat(string(t.Lit), 64)
	if err != nil {
		return nil, err
	}

	return &FloatLiteral{Token: *t, Value: v}, nil
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return string(fl.Token.Lit) }
func (fl *FloatLiteral) String() string       { return fl.TokenLiteral() }

type StringLiteral struct {
	Token token.Token
	Value string
}

func NewStringLiteral(t *token.Token) (*StringLiteral, error) {
	val, err := strconv.Unquote(string(t.Lit))
	if err != nil {
		return nil, err
	}
	return &StringLiteral{Token: *t, Value: val}, nil
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return string(sl.Token.Lit) }
func (sl *StringLiteral) String() string       { return sl.TokenLiteral() }

type FunctionLiteral struct {
	Token      token.Token
	Parameters IdentifierList
	Body       *BlockStatement
}

func NewFunctionLiteral(t *token.Token, params IdentifierList, body *BlockStatement) (*FunctionLiteral, error) {
	if params == nil {
		params = IdentifierList{}
	}
	return &FunctionLiteral{Token: *t, Parameters: params, Body: body}, nil
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return string(fl.Token.Lit) }
func (fl *FunctionLiteral) String() string {
	var out strings.Builder

	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString(" (")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments ExpressionList
}

func NewCallExpression(function Expression, t *token.Token, args ExpressionList) (*CallExpression, error) {
	return &CallExpression{Function: function, Token: *t, Arguments: args}, nil
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return string(ce.Token.Lit) }
func (ce *CallExpression) String() string {
	var out strings.Builder

	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteRune('(')
	out.WriteString(strings.Join(args, ", "))
	out.WriteRune(')')

	return out.String()
}

type ArrayLiteral struct {
	Token    token.Token
	Elements ExpressionList
}

func NewArrayLiteral(t *token.Token, elems ExpressionList) (*ArrayLiteral, error) {
	if elems == nil {
		elems = ExpressionList{}
	}
	return &ArrayLiteral{Token: *t, Elements: elems}, nil
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return string(al.Token.Lit) }
func (al *ArrayLiteral) String() string {
	var out strings.Builder

	var elems []string
	for _, el := range al.Elements {
		elems = append(elems, el.String())
	}

	out.WriteRune('[')
	out.WriteString(strings.Join(elems, ", "))
	out.WriteRune(']')

	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs ExpressionMap
}

func NewHashLiteral(t *token.Token, pairs ExpressionMap) (*HashLiteral, error) {
	if pairs == nil {
		pairs = ExpressionMap{}
	}
	return &HashLiteral{Token: *t, Pairs: pairs}, nil
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return string(hl.Token.Lit) }
func (hl *HashLiteral) String() string {
	var out strings.Builder

	var pairs []string
	for k, v := range hl.Pairs {
		pairs = append(pairs, k.String()+":"+v.String())
	}

	out.WriteRune('{')
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteRune('}')

	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func NewIndexExpression(left Expression, t *token.Token, index Expression) (*IndexExpression, error) {
	return &IndexExpression{Left: left, Token: *t, Index: index}, nil
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return string(ie.Token.Lit) }
func (ie *IndexExpression) String() string {
	var out strings.Builder

	out.WriteRune('(')
	out.WriteString(ie.Left.String())
	out.WriteRune('[')
	out.WriteString(ie.Index.String())
	out.WriteRune(']')
	out.WriteRune(')')

	return out.String()
}

type SelectorExpression struct {
	Token token.Token
	Left  Expression
	Right *Identifier
}

func NewSelectorExpression(left Expression, t *token.Token, right *Identifier) (*SelectorExpression, error) {
	return &SelectorExpression{Left: left, Token: *t, Right: right}, nil
}

func (se *SelectorExpression) expressionNode()      {}
func (se *SelectorExpression) TokenLiteral() string { return string(se.Token.Lit) }
func (se *SelectorExpression) String() string {
	var out strings.Builder

	out.WriteString(se.Left.String())
	out.WriteRune('.')
	out.WriteString(se.Right.String())

	return out.String()
}
