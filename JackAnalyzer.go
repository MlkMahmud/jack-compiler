package main

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
	analyzer.parser.Parse(tokens, outFile)
}
