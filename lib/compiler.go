package lib

import (
	"strings"
)

type Compiler struct {
	lexer  *Lexer
	parser *Parser
}

func NewCompiler() *Compiler {
	return &Compiler{
		lexer:  NewLexer(),
		parser: NewParser(),
	}
}

func (compiler *Compiler) Compile(src string) {
	outFile := strings.Replace(src, ".jack", ".xml", -1)
	tokens := compiler.lexer.Tokenize(src)
	compiler.parser.Parse(tokens, outFile)
}
