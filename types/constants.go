package types

type SymbolKind string

const (
	Argument    SymbolKind = "argument"
	Constructor SymbolKind = "constructor"
	Field       SymbolKind = "field"
	Function    SymbolKind = "function"
	Method      SymbolKind = "method"
	Static      SymbolKind = "static"
	Var         SymbolKind = "var"
)

type BinaryOperator string

const (
	Addition       BinaryOperator = "+"
	Subraction     BinaryOperator = "-"
	Multiplication BinaryOperator = "*"
	Division       BinaryOperator = "/"
	LessThan       BinaryOperator = "<"
	GreaterThan    BinaryOperator = ">"
	Equals         BinaryOperator = "="
)

type LogicalOperator string

const (
	And LogicalOperator = "&"
	Or  LogicalOperator = "|"
)

type UnaryOperator string

const (
	ArithmeticNegation UnaryOperator = "-"
	BooleanNegation    UnaryOperator = "~"
)

type LiteralType string

const (
	BooleanLiteral LiteralType = "boolean"
	IntegerLiteral LiteralType = "int"
	NullLiteral    LiteralType = "null"
	StringLiteral  LiteralType = "string"
	ThisLiteral    LiteralType = "this"
)

var KEYWORDS = map[string]bool{
	"class": true, "constructor": true, "method": true,
	"function": true, "int": true, "boolean": true,
	"char": true, "void": true, "var": true, "static": true,
	"field": true, "let": true, "do": true, "if": true,
	"else": true, "while": true, "return": true, "true": true,
	"false": true, "null": true, "this": true,
}

var SYMBOLS = map[string]bool{
	"(": true, ")": true, "{": true, "}": true, "[": true, "]": true, "/": true,
	"-": true, "+": true, "*": true, ",": true, ".": true, "=": true, ";": true,
	"&": true, "|": true, "<": true, ">": true, "~": true,
}
