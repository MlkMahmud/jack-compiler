package types

import (
	"fmt"
	"strings"
)

type Expr interface {
	fmt.Stringer
}

type BinaryExpr struct {
	Operator BinaryOperator
	Left     Expr
	Right    Expr
}

func (expr BinaryExpr) String() string {
	return fmt.Sprintf(
		"%s %s %s",
		expr.Left,
		expr.Operator,
		expr.Right,
	)
}

type CallExpr struct {
	Arguments []Expr
	Callee    Expr
}

func (expr CallExpr) String() string {
	var args []string

	for _, arg := range expr.Arguments {
		args = append(args, arg.String())
	}

	return fmt.Sprintf(
		"%s(%s)",
		expr.Callee,
		strings.Join(args, ", "))
}

type IndexExpr struct {
	Indexer Expr
	Object  Ident
}

func (expr IndexExpr) String() string {
	return fmt.Sprintf(
		"%s[%s]",
		expr.Object,
		expr.Indexer,
	)
}

type LogicalExpr struct {
	Operator LogicalOperator
	Left     Expr
	Right    Expr
}

func (expr LogicalExpr) String() string {
	return fmt.Sprintf(
		"%s %s %s",
		expr.Left,
		expr.Operator,
		expr.Right,
	)
}

type MemberExpr struct {
	Object   Ident
	Property Ident
}

func (expr MemberExpr) String() string {
	return fmt.Sprintf(
		"%s.%s",
		expr.Object,
		expr.Property,
	)
}

type ParenExpr struct {
	Expression Expr
}

func (expr ParenExpr) String() string {
	return fmt.Sprintf("(%s)", expr.Expression)
}

type UnaryExpr struct {
	Operator UnaryOperator
	Operand  Expr
}

func (expr UnaryExpr) String() string {
	return fmt.Sprintf(
		"%s %s",
		expr.Operator,
		expr.Operand,
	)
}

type Literal struct {
	Type  LiteralType
	Value string
}

func (literal Literal) String() string {
	return fmt.Sprintf("%v", literal.Value)
}

type Ident struct {
	Name string
}

func (i Ident) String() string {
	return i.Name
}
