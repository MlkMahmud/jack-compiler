package analyzer

import (
	"github.com/MlkMahmud/jack-compiler/lexer"
	"github.com/MlkMahmud/jack-compiler/parser"
)

type JackAnalyzer struct {
	lexer  *lexer.Lexer
	parser *parser.Parser
}

func NewAnalyzer() *JackAnalyzer {
	return &JackAnalyzer{
		lexer:  lexer.NewLexer(),
		parser: parser.NewParser(),
	}
}

func (analyzer *JackAnalyzer) Run(src string) {
	tokens := analyzer.lexer.Tokenize(src)
	analyzer.parser.Parse(tokens)
}