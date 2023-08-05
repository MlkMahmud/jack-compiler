package parser

import "fmt"

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

type LogicalOperator int

const (
	And LogicalOperator = iota
	Or
)

func (op LogicalOperator) String() string {
	return []string{"&", "|"}[op]
}

type LogicalExpr struct {
	Operator LogicalOperator
	Left     Expr
	Right    Expr
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

type NullLiteral struct{}

func (n NullLiteral) String() string {
	return "null"
}

type StringLiteral struct {
	Value string
}

type ThisLiteral struct{}

func (t ThisLiteral) String() string {
	return "this"
}

type Indentifier struct {
	Name string
}

func (i Indentifier) String() string {
	return i.Name
}
