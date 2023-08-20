package types

type TokenType int

func (tokenType TokenType) String() string {
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
	Filename  string
	Lexeme    string
	LineNum   int
	TokenType TokenType
}
