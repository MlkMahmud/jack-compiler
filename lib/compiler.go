package lib

import (
	"strings"
)

type Compiler struct {
	lexer *Lexer
	parser *Parser
}

func NewCompiler() *Compiler {
	var compiler = new(Compiler)
	compiler.lexer = NewLexer()
	compiler.parser = NewParser()
	return compiler
}

func (compiler *Compiler) Compile(src string) {
	var outFile = strings.Replace(src, ".jack", ".xml", -1)
	var tokens = compiler.lexer.Tokenize(src)
	compiler.parser.Parse(tokens, outFile)
}