package lib

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

/*
	The Jack parser receives an array of tokens.
	It attemps to parse the tokens and generate a
	parse tree according to the rules of Jack's grammar.
*/

type Parser struct {
	outFile *os.File
	tokens  *list.List
}

func NewParser() *Parser {
	return new(Parser)
}

func (parser *Parser) write(tag string, lexeme any) {
	var str string
	if lexeme == nil {
		str = fmt.Sprintf("<%s>\n", tag)
	} else {
		str = fmt.Sprintf(
			"<%s>%s</%s>\n",
			tag,
			lexeme,
			tag,
		)
	}
	_, err := parser.outFile.WriteString(str)
	if err != nil {
		panic(err)
	}
}

func (parser *Parser) throwSyntaxError(token Token) {
	srcFile := replaceFileExt(parser.outFile.Name(), ".jack")
	errorMessage := fmt.Sprintf(
		"(%s):[%d:%d]: syntax error: unexpected token '%s'\n",
		srcFile,
		token.lineNum,
		token.colNum,
		token.lexeme,
	)
	panic(&CompilerError{errorMessage})
}

func (parser *Parser) getNextToken() Token {
	item := parser.tokens.Front()
	parser.tokens.Remove(item)
	return item.Value.(Token)
}

func (parser *Parser) peekNextToken() Token {
	item := parser.tokens.Front()
	return item.Value.(Token)
}

func (parser *Parser) parseToken(terminals []string) {
	token := parser.getNextToken()
	
	switch token.tokenType {
	case KEYWORD, SYMBOL:
		{
			for _, val := range terminals {
				if val == token.lexeme {
					parser.write(token.tokenType.toString(), token.lexeme)
					return
				}
			}
		}

	case INTEGER_CONSTANT, STRING_CONSTANT:
		{
			for _, val := range terminals {
				if val == token.tokenType.toString() {
					parser.write(token.tokenType.toString(), token.lexeme)
					return
				}
			}
		}

	case IDENTIFIER:
		{
			for _, val := range terminals {
				if val == "className" || val == "subroutineName" || val == "varName" {
					parser.write(token.tokenType.toString(), token.lexeme)
					return
				}
			}
		}
	}
	parser.throwSyntaxError(token)
}

func (parser *Parser) parseParameterList() {
	parser.write("parameterList", nil)

	for nextToken := parser.peekNextToken(); !isSymbol(nextToken, ")"); nextToken = parser.peekNextToken() {
		parser.parseToken([]string{"boolean", "char", "className", "int"})
		parser.parseToken([]string{"varName"})
		if isSymbol(parser.peekNextToken(), ",") {
			parser.parseToken([]string{","})
		}
	}

	parser.write("/parameterList", nil)
}

func (parser *Parser) parseVarDec() {
	parser.write("varDec", nil)
	parser.parseToken([]string{"var"})
	parser.parseToken([]string{"boolean", "char", "className", "int"})
	parser.parseToken([]string{"varName"})

	for nextToken := parser.peekNextToken(); !isSymbol(nextToken, ";"); nextToken = parser.peekNextToken() {
		parser.parseToken([]string{","})
		parser.parseToken([]string{"varName"})
	}

	parser.parseToken([]string{";"})
	parser.write("/varDec", nil)
}

func (parser *Parser) parseDoStatement() {
	parser.write("doStatement", nil)
	parser.parseToken([]string{"do"})
	parser.parseToken([]string{"className", "subroutineName"})

	if isSymbol(parser.peekNextToken(), ".") {
		parser.parseToken([]string{"."})
		parser.parseToken([]string{"subroutineName"})
	}

	parser.parseToken([]string{"("})

	// parser.parseExpressionList()
	parser.parseToken([]string{")"})
	parser.parseToken([]string{";"})
	parser.write("/doStatement", nil)
}

func (parser *Parser) parseStatements() {
	parser.write("statements", nil)

	for token := parser.peekNextToken(); !isSymbol(token, "}"); token = parser.peekNextToken() {
		switch token.lexeme {
		case "do":
			parser.parseDoStatement()
		case "if":
			// parseIfStatement
		case "let":
			// parseLetStatement
		case "return":
			// parseReturnStatement
		case "while":
			// parseWhileStatement
		default:
			parser.throwSyntaxError(token)
		}
	}
	parser.write("/statements", nil)
}

func (parser *Parser) parseSubroutineBody() {
	parser.write("subroutineBody", nil)
	parser.parseToken([]string{"{"})

	for end := parser.peekNextToken(); !isSymbol(end, "}"); end = parser.peekNextToken() {
		if isKeyword(end, "var") {
			parser.parseVarDec()
		} else {
			parser.parseStatements()
		}
	}
	parser.parseToken([]string{"}"})
	parser.write("/subroutineBody", nil)
}

func (parser *Parser) parseSubroutineDec() {
	parser.write("subroutineDec", nil)
	parser.parseToken([]string{"constructor", "function", "method"})
	parser.parseToken([]string{"boolean", "char", "className", "int", "void"})
	parser.parseToken([]string{"subroutineName"})
	parser.parseToken([]string{"("})

	parser.parseParameterList()
	parser.parseToken([]string{")"})
	parser.parseSubroutineBody()
	parser.write("/subroutineDec", nil)
}

func (parser *Parser) parseClassVarDec() {
	parser.write("classVarDec", nil)

	parser.parseToken([]string{"field", "static"})
	parser.parseToken([]string{"boolean", "char", "className", "int"})
	parser.parseToken([]string{"varName"})

	nextToken := parser.peekNextToken()

	// Check if it's a multi var declaration.
	for !isSymbol(nextToken, ";") {
		parser.parseToken([]string{","})
		parser.parseToken([]string{"varName"})

		nextToken = parser.peekNextToken()
	}

	parser.parseToken([]string{";"})
	parser.write("/classVarDec", nil)
}

func (parser *Parser) Parse(tokens *list.List, outFile string) {
	file, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}

		if r := recover(); r != nil {
			var compilerError *CompilerError
			if errors.As(r.(error), &compilerError) {
				fmt.Println(r)
			} else {
				debug.PrintStack()
			}
			os.Exit(1)
		}
	}()

	parser.outFile = file
	parser.tokens = tokens

	parser.write("class", nil)
	parser.parseToken([]string{"class"})
	parser.parseToken([]string{"className"})
	parser.parseToken([]string{"{"})

	for end := parser.peekNextToken(); !isSymbol(end, "}"); end = parser.peekNextToken() {
		if isClassVarDec(end) {
			parser.parseClassVarDec()
		} else if isSubroutineDec(end) {
			parser.parseSubroutineDec()
		} else {
			parser.throwSyntaxError(end)
		}
	}

	parser.parseToken([]string{"}"})
	parser.write("/class", nil)
}
