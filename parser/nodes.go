package parser

import "fmt"

type Class struct {
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

type SubroutineDecl struct {
	Name       string
	Params     []Expr
	Kind       SubroutineKind
	ReturnType string
	Statements []Stmt
	Vars       []VarDecl
}

type DoStmt struct {
	Arguments      []Expr
	ObjectName     string
	SubroutineName string
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

type Bool int

const (
	False Bool = iota
	True
)

func (b Bool) String() string {
	return []string{"false", "true"}[b]
}

type BoolLiteral struct {
	Value Bool
}

type IntLiteral struct {
	Value int16
}

type NullLiteral struct {}

func (n NullLiteral) String() string {
	return "null"
}

type StringLiteral struct {
	Value string
}

type ThisLiteral struct {}

func (t ThisLiteral) String() string {
	return "this"
}