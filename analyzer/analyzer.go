package analyzer

import (
	"github.com/MlkMahmud/jack-compiler/lexer"
	"github.com/MlkMahmud/jack-compiler/parser"
	"github.com/MlkMahmud/jack-compiler/helpers"
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
	outFile := helpers.ReplaceFileExt(src, ".xml")
	tokens := analyzer.lexer.Tokenize(src)
	analyzer.parser.Parse(tokens, outFile)
}