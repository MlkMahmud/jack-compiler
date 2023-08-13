package parser

import (
	"fmt"
	"strings"
)

type CompilationUnit struct {
	Name        string
	Subroutines []SubroutineDecl
	Vars        []VarDecl
}

type SubroutineKind int

const (
	Constructor SubroutineKind = iota
	Function
	Method
)

func (s SubroutineKind) String() string {
	return []string{"constructor", "function", "method"}[s]
}

type VarKind int

const (
	Field VarKind = iota
	Static
	Var
)

func (v VarKind) String() string {
	return []string{"field", "static", "var"}[v]
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
	Statements []Stmt
	Vars       []VarDecl
}

type BlockStmt struct {
	Statements []Stmt
	Vars       []VarDecl
}

type DoStmt struct {
	Arguments      []Expr
	ObjectName     string
	SubroutineName string
}

func (d DoStmt) String() string {
	var args []string

	for _, arg := range d.Arguments {
		args = append(args, arg.String())
	}

	if d.ObjectName != "" {
		return fmt.Sprintf(
			"%s.%s(%s)",
			d.ObjectName,
			d.SubroutineName,
			strings.Join(args, ", "),
		)
	}

	return fmt.Sprintf(
		"%s(%s)",
		d.SubroutineName,
		strings.Join(args, ", "),
	)
}

type IfStmt struct {
	Condition Expr
	ThenStmt  Stmt
	ElseStmt  Stmt
}

type LetStmt struct {
	Name        string
	ArrayAccess Expr
	Value       Expr
}

type ReturnStmt struct {
	Value Expr
}

type WhileStmt struct {
	Condition  Expr
	Statements []Stmt
}

type AssignExpr struct {
	Target Expr
	Value  Expr
}

type BinaryOperator int

const (
	Addition BinaryOperator = iota
	Subraction
	Multiplication
	Division
	LessThan
	GreaterThan
)

func (op BinaryOperator) String() string {
	return []string{"+", "-", "*", "/", "<", ">"}[op]
}

type BinaryExpr struct {
	Operator BinaryOperator
	Left     Expr
	Right    Expr
}

type CallExpr struct {
	Arguments []Expr
	Callee    interface{}
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

type LogicalOperator int

const (
	And LogicalOperator = iota
	Or
)

func (op LogicalOperator) String() string {
	return []string{"&", "|"}[op]
}

type IndexExpr struct {
	Indexer Expr
	Object  Indentifier
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

type MemberExpr struct {
	Object   Indentifier
	Property Indentifier
}

type ParenExpr struct {
	Expression Expr
}

func (expr ParenExpr) String() string {
	return fmt.Sprintf("(%s)", expr.Expression)
}

type UnaryOperator int

const (
	ArithmeticNegation UnaryOperator = iota
	BooleanNegation
)

func (op UnaryOperator) String() string {
	return []string{"-", "~"}[op]
}

type UnaryExpr struct {
	Operator UnaryOperator
	Operand  string
}

type Literal struct {
	Value interface{}
}

func (literal Literal) String() string {
	return fmt.Sprintf("%v", literal.Value)
}

type Indentifier struct {
	Name string
}

func (i Indentifier) String() string {
	return i.Name
}
