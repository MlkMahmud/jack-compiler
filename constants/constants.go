package constants

type TokenType int

func (tokenType TokenType) ToString() string {
	return []string{
		"identifier", "integerConstant", "keyword", "stringConstant", "symbol",
	}[tokenType]
}
const (
	IDENTIFIER TokenType = iota
	INTEGER_CONSTANT
	KEYWORD
	STRING_CONSTANT
	SYMBOL
)

type Token struct {
	ColNum    int
	Lexeme    string
	LineNum   int
	TokenType TokenType
}

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