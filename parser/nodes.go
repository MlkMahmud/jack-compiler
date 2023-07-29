package parser

type ClassDeclaration struct {
	Identifier             string
	ClassVarDeclarations   []ClassVarDeclaration
	SubroutineDeclarations []SubroutineDeclaration
}

type ClassVarDeclaration struct {
	Identifier string
	Kind       string
	Type       string
}

type SubroutineDeclaration struct {
	Identifier      string
	Kind            string
	Type            string
	VarDeclarations []VarDeclaration
	Statements      []any
}

type VarDeclaration struct {
	Identifier string
	Type       string
}
