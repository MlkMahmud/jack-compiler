package lib

import (
	"fmt"
	"jack-compiler/util"
	"log"
	"os"
	"strings"
)

// /*
// 	The Jack parser receives an array of tokens.
// 	It attemps to parse the tokens and generate a
// 	parse tree according to the rules of Jack's grammar.
// */

type Parser struct {
	outFile *os.File
	tokens  util.Queue
}

func NewParser() *Parser {
	return new(Parser)
}

func (parser *Parser) getNextToken() Token {
	return parser.tokens.Dequeue().(Token)
}

func (parser *Parser) peekNextToken() Token {
	return parser.tokens.Peek().(Token)
}

func (parser *Parser) write(bytes string) {
	var _, err = parser.outFile.Write([]byte(bytes))
	if err != nil {
		panic(err)
	}
}

func (parser *Parser) parseClassDec() {
	var classKeyword = parser.getNextToken()
	var className = parser.getNextToken()
	var leftBrace = parser.getNextToken()

	if classKeyword.lexeme != "class" {
		panic("class declaration must begin with the 'class' keyword")
	}

	if className.tokenType != IDENTIFIER {
		panic("class name must be a valid identifier")
	}

	if leftBrace.lexeme != "{" {
		panic("class name must be followed by an opening brace '{' character")
	}

	parser.write(fmt.Sprintf(
		`
		<class>
		<keyword>class</keyword>
		<identifier>%s</identifier>
		<symbol> { </symbol>
		`,
		classKeyword.lexeme,
	))

	var nextToken = parser.peekNextToken()

	if nextToken.lexeme == "}" {
		// This is a class with an empty body.
		parser.write(("</class>"))
		// Discard the token
		parser.getNextToken()
		return
	}

	/*
		If the class body is not empty, the next token must be a class var declaration
		or a subroutine declaration.
	*/
	for _, val := range []string{"constructor", "field", "function", "method", "static"} {
		if val == nextToken.lexeme {
			return
		}
	}

	panic(fmt.Sprintf("invalid token %s", nextToken.lexeme))
}

func (parser *Parser) parseSubroutineDec() {}

func (parser *Parser) parseClassVarDec() {
	var varType = parser.getNextToken()
	var dataType = parser.getNextToken()
	var varName = parser.getNextToken()

	if varType.lexeme != "field" && varType.lexeme != "static" {
		panic(
			fmt.Sprintf(
				"invalid variable type: %s.\nclass variables must be decalred with either the 'field' or 'static' keyword.",
				varName.lexeme,
			),
		)
	}

	if dataType.tokenType != IDENTIFIER &&
		dataType.tokenType != KEYWORDS["boolean"] &&
		dataType.tokenType != KEYWORDS["char"] &&
		dataType.tokenType != KEYWORDS["int"] {
		panic(
			fmt.Sprintf(
				"invalid variable dataType: %s",
				dataType.lexeme,
			),
		)
	}

	if varName.tokenType != IDENTIFIER {
		panic(
			fmt.Sprintf(
				"invalid variable name: %s",
				varName.lexeme,
			),
		)
	}

	parser.write(
		fmt.Sprintf(
			`
			<classVarDec>
			<keyword>%s</keyword>
			<keyword>%s</keyword>
			<identifier>%s</identifier>
			`,
			varType.lexeme,
			dataType.lexeme,
			varName.lexeme,
		),
	)

	var nextToken = parser.peekNextToken()

	for nextToken.lexeme != ";" {
		// Check if it's a multi var declaration.
		var commaSymbol = parser.getNextToken()
		var varName = parser.getNextToken()

		if commaSymbol.lexeme != "," && varName.tokenType != IDENTIFIER {
			panic("invalid class variable declaration.")
		}

		parser.write(
			fmt.Sprintf(
			`
				<symbol>,</symbol>
				<identifier>%s</identifier>
			`,
			varName.lexeme,
			),
		)

		nextToken = parser.peekNextToken()
	}

	parser.write("</classVarDec>")
	// discard final ";" token.
	parser.tokens.Dequeue()
}

func (parser *Parser) Parse(tokens util.Queue, outFile string) {
	if tokens.Size() < 4 {
		//	The minimum number of tokens for a valid Jack program is 4. "class className {}"
		var srcFile = strings.Replace(outFile, ".xml", ".jack", -1)
		log.Fatalf(fmt.Sprintf("File: %s is not a valid Jack program", srcFile))
	}

	var file, err = os.Create(outFile)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		var err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	parser.outFile = file
	parser.tokens = tokens

	for tokens.Size() > 0 {
		var nextToken = parser.peekNextToken()

		switch nextToken.lexeme {
		case "class":
			parser.parseClassDec()
		case "field", "static":
			parser.parseClassVarDec()
		case "constructor", "function", "method":
			parser.parseSubroutineDec()
		case "}":
			parser.write("</class>")
			parser.tokens.Dequeue()
		}
	}
}
