package types

import (
	"fmt"
	"strings"
)

type Class struct {
	Name        Ident
	Subroutines []SubroutineDecl
	Vars        []VarDecl
}

type Expr interface {
	fmt.Stringer
}

type Stmt interface {
	fmt.Stringer
}

type VarDecl struct {
	Name string
	Kind VarKind
	Type string
}

type Parameter struct {
	Name string
	Type string
}

type SubroutineDecl struct {
	Name       string
	Params     []Parameter
	Kind       SubroutineKind
	Type       string
	Body       SubroutineBody
}

type SubroutineBody struct {
	Statements []Stmt
	Vars       []VarDecl
}

type BlockStmt struct {
	Statements []Stmt
}

type DoStmt struct {
	Expression CallExpr
}

func (d DoStmt) String() string {
	return d.Expression.String()
}

type IfStmt struct {
	Condition Expr
	ThenStmt  BlockStmt
	ElseStmt  BlockStmt
}

func (stmt IfStmt) String() string {
	if len(stmt.ElseStmt.Statements) < 1 {
		return fmt.Sprintf(
			"if (%s) %s",
			stmt.Condition,
			stmt.ThenStmt,
		)
	}

	return fmt.Sprintf(
		"if (%s) %s\nelse %s",
		stmt.Condition,
		stmt.ThenStmt,
		stmt.ElseStmt,
	)
}

type LetStmt struct {
	Target Expr
	Value  Expr
}

func (stmt LetStmt) String() string {
	return fmt.Sprintf(
		"let %s = %s\n",
		stmt.Target,
		stmt.Value,
	)
}

type ReturnStmt struct {
	Expression Expr
}

type WhileStmt struct {
	Condition  Expr
	Statements []Stmt
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
