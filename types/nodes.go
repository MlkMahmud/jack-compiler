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

type Stmt interface {
	fmt.Stringer
}

type VarDecl struct {
	Name string
	Kind SymbolKind
	Type string
}

type Parameter struct {
	Name string
	Type string
}

type SubroutineDecl struct {
	Name   Ident
	Params []Parameter
	Kind   SymbolKind
	Type   string
	Body   SubroutineBody
}

type SubroutineBody struct {
	Statements []Stmt
	Vars       []VarDecl
}

type BlockStmt struct {
	Statements []Stmt
}

func (s BlockStmt) String() string {
	stmts := []string{}

	for _, stmt := range s.Statements {
		stmts = append(stmts, stmt.String())
	}

	return fmt.Sprintf(
		"{ %s }\n",
		strings.Join(stmts, "\n"),
	)
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

func (stmt ReturnStmt) String() string {
	return fmt.Sprintf("return %s;\n", stmt.Expression)
}

type WhileStmt struct {
	Body      BlockStmt
	Condition Expr
}

func (stmt WhileStmt) String() string {
	return fmt.Sprintf(
		"while (%s) %s",
		stmt.Condition,
		stmt.Body,
	)
}
