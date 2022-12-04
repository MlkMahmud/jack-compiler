package lib

import (
	"fmt"
	"log"
)

type JackAnalyzer struct {
	lexer  *Lexer
	parser *Parser
}

func NewAnalyzer() *JackAnalyzer {
	return &JackAnalyzer{
		lexer:  NewLexer(),
		parser: NewParser(),
	}
}

func (analyzer *JackAnalyzer) Run(src string) {
	outFile := replaceFileExt(src, ".xml")
	tokens := analyzer.lexer.Tokenize(src)
	if tokens.Len() < 4 {
		//	The minimum number of tokens for a valid Jack program is 4. "class className {}"
		log.SetFlags(0)
		log.Fatalf(fmt.Sprintf("ERROR: %s is not a valid Jack program\n", src))
	}
	analyzer.parser.Parse(tokens, outFile)
}
	