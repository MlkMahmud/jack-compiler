package lib

import (
	"container/list"
	"fmt"
	"log"
	"os"
	"strings"
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

func isValidType(token Token) bool {
	if token.tokenType != IDENTIFIER &&
		token.lexeme != "boolean" &&
		token.lexeme != "char" &&
		token.lexeme != "int" {
		return false
	}
	return true
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

func (parser *Parser) write(str string) {
	_, err := parser.outFile.WriteString(str)
	if err != nil {
		panic(err)
	}
}

func (parser *Parser) parseParameterList() {
	argType := parser.getNextToken()
	argName := parser.getNextToken()

	if !isValidType(argType) {
		panic(fmt.Sprintf("error: invalid token: %s", argType.lexeme))
	}

	if argName.tokenType != IDENTIFIER {
		panic(fmt.Sprintf("error: invalid token: %s", argName.lexeme))
	}

	parser.write(fmt.Sprintf(
		"\t\t<parameterList>\n"+
			"\t\t\t<%s>%s</%s>\n"+
			"\t\t\t<identifier>%s</identifier>\n",
		argType.tokenType,
		argType.lexeme,
		argType.tokenType,
		argName.lexeme,
	))

	nextToken := parser.peekNextToken()

	for nextToken.lexeme != ")" {
		commaLiteral := parser.getNextToken()
		nextArgType := parser.getNextToken()
		nextArgName := parser.getNextToken()

		if commaLiteral.lexeme != "," {
			panic(fmt.Sprintf("error: invalid token: %s", commaLiteral.lexeme))
		}

		if !isValidType(nextArgType) {
			panic(fmt.Sprintf("error: invalid token: %s", argName.lexeme))
		}

		if nextArgName.tokenType != IDENTIFIER {
			panic(fmt.Sprintf("error: invalid token: %s", nextArgName.lexeme))
		}

		parser.write(fmt.Sprintf(
			"\t\t\t<symbol>,</symbol>\n"+
				"\t\t\t<%s>%s</%s>\n"+
				"\t\t\t<identifier>%s</identifier>\n",
			nextArgType.tokenType,
			nextArgType.lexeme,
			nextArgType.tokenType,
			nextArgName.lexeme,
		))
		nextToken = parser.peekNextToken()
	}

	parser.write(
		"\t\t</parameterList>\n" +
			"\t\t<symbol>)</symbol>\n",
	)
	// discard closing brace ")" token.
	parser.getNextToken()
}

func (parser *Parser) parseVarDec() {
	varDec := parser.getNextToken()
	varType := parser.getNextToken()
	varName := parser.getNextToken()

	if varDec.lexeme != "var" {
		panic(fmt.Sprintf("error: invalid token: %s", varDec.lexeme))
	}

	if !isValidType(varType) {
		panic(fmt.Sprintf("error: invalid token: %s", varType.lexeme))
	}

	if varName.tokenType != IDENTIFIER {
		panic(fmt.Sprintf("error: invalid token: %s", varName.lexeme))
	}

	parser.write(fmt.Sprintf(
		"\t\t\t\t<varDec>\n"+
			"\t\t\t\t\t<keyword>var<keyword>\n"+
			"\t\t\t\t\t<%s>%s</%s>\n"+
			"\t\t\t\t\t<identifier>%s</identifier>\n",
		varType.tokenType,
		varType.lexeme,
		varType.tokenType,
		varName.lexeme,
	))

	nextToken := parser.peekNextToken()

	for nextToken.lexeme != ";" {
		commaSymbol := parser.getNextToken()
		nextVarName := parser.getNextToken()

		if commaSymbol.lexeme != "," {
			panic(fmt.Sprintf("error: invalid token: %s", commaSymbol.lexeme))
		}

		if nextVarName.tokenType != IDENTIFIER {
			panic(fmt.Sprintf("error: invalid token: %s", nextVarName.lexeme))
		}

		parser.write(fmt.Sprintf(
			"\t\t\t\t\t<symbol>,</symbol>\n"+
				"\t\t\t\t\t<identifier>%s</identifier>\n",
			nextVarName.lexeme,
		))

		nextToken = parser.peekNextToken()
	}

	parser.write(
		"\t\t\t\t\t<symbol>;</symbol>\n" +
			"\t\t\t\t</varDec>\n",
	)
	// discard ";" token.
	parser.getNextToken()
}

func (parser *Parser) parseExpression() {}

func (parser *Parser) parseDoStatement() {
	doKeyword := parser.getNextToken()
	subroutetineOrClassName := parser.getNextToken()

	if doKeyword.lexeme != "do" {
		panic(fmt.Sprintf("error: invalid token: %s", doKeyword.lexeme))
	}

	if subroutetineOrClassName.tokenType != IDENTIFIER {
		panic(fmt.Sprintf("error: invalid token: %s", subroutetineOrClassName.lexeme))
	}

	parser.write("<keyword>do</keyword>\n")
	periodOrLparen := parser.getNextToken()

	if periodOrLparen.lexeme == "." {
		className := subroutetineOrClassName
		subroutineName := parser.getNextToken()

		if subroutineName.tokenType != IDENTIFIER {
			panic(fmt.Sprintf("error: invalid token: %s", subroutineName.lexeme))
		}

		if nextToken := parser.peekNextToken(); nextToken.lexeme != "(" {
			panic(fmt.Sprintf("error: invalid token: %s", nextToken.lexeme))
		}

		parser.write(fmt.Sprintf(
			"<identifier>%s<identifier>\n"+
				"<symbol>.</symbol>\n"+
				"<identifier>%s</identifier>",
			className.lexeme,
			subroutineName.lexeme,
		))
	} else if periodOrLparen.lexeme == "(" {
		subroutineName := subroutetineOrClassName
		parser.write(fmt.Sprintf(
			"<identifier>%s</identifier>\n"+
				"<symbol>(</symbol>\n",
			subroutineName.lexeme,
		))
	} else {
		panic(fmt.Sprintf("error: invalid token: %s", periodOrLparen.lexeme))
	}

	// handle expression List
}

func (parser *Parser) parseStatements() {
	statementType := parser.peekNextToken()

	parser.write("<statements>\n")

	switch statementType.lexeme {
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
		panic(fmt.Sprintf("error: invalid token: %s", statementType.lexeme))
	}
	parser.write("</statements>\n")
}

func (parser *Parser) parseSubroutineBody() {
	parser.write(
		"<subroutineBody>\n" +
			"<symbol>{</symbol>\n",
	)

	rBrace := parser.peekNextToken()

	for rBrace.lexeme != "}" {
		switch rBrace.lexeme {
		case "var":
			parser.parseVarDec()
			rBrace = parser.peekNextToken()
		case "do", "if", "let", "return", "while":
			parser.parseStatements()
			rBrace = parser.peekNextToken()
		default:
			panic(fmt.Sprintf("error: invalid token: %s", rBrace.lexeme))
		}
	}

	parser.write(
		"<symbol>}</symbol>\n" +
			"</subroutineBody>\n",
	)
	// discard right brace "}" token
	parser.getNextToken()
}

func (parser *Parser) parseSubroutineDec() {
	subroutineType := parser.getNextToken()
	subroutineReturnType := parser.getNextToken()
	subroutineName := parser.getNextToken()
	leftParen := parser.getNextToken()

	if subroutineType.lexeme != "constructor" &&
		subroutineType.lexeme != "function" &&
		subroutineType.lexeme != "method" {
		panic(fmt.Sprintf("error: invalid token: %s", subroutineType.lexeme))
	}

	if !isValidType(subroutineReturnType) && subroutineReturnType.lexeme != "void" {
		panic(fmt.Sprintf("error: invalid token: %s", subroutineType.lexeme))
	}

	if subroutineName.tokenType != IDENTIFIER {
		panic(fmt.Sprintf("error: invalid token: %s", subroutineType.lexeme))
	}

	if leftParen.lexeme != "(" {
		panic(fmt.Sprintf("error: invalid token: %s", subroutineType.lexeme))
	}

	parser.write(fmt.Sprintf(
		"<subroutineDec>\n"+
			"<keyword>%s</keyword>\n"+
			"<%s>%s</%s>\n"+
			"<identifier>%s</identifier>\n"+
			"<symbol>(</symbol>\n",
		subroutineType.lexeme,
		subroutineReturnType.tokenType,
		subroutineReturnType.lexeme,
		subroutineReturnType.tokenType,
		subroutineName.lexeme,
	))

	parser.parseParameterList()
	nextToken := parser.peekNextToken()

	if nextToken.lexeme != "{" {
		panic(fmt.Sprintf("error: invalid token: %s", nextToken.lexeme))
	}
	// discard left brace token "{"
	parser.getNextToken()

	if parser.peekNextToken().lexeme != "}" {
		parser.parseSubroutineBody()
	}

	parser.write("</subroutineDec>\n")
}

func (parser *Parser) parseClassVarDec() {
	varType := parser.getNextToken()
	dataType := parser.getNextToken()
	varName := parser.getNextToken()

	if varType.lexeme != "field" && varType.lexeme != "static" {
		panic(
			fmt.Sprintf(
				"invalid variable type: %s.\nclass variables must be decalred with either the 'field' or 'static' keyword.",
				varName.lexeme,
			),
		)
	}

	if !isValidType(dataType) {
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
			"<classVarDec>\n"+
				"<keyword>%s</keyword>\n"+
				"<keyword>%s</keyword>\n"+
				"<identifier>%s</identifier>\n",
			varType.lexeme,
			dataType.lexeme,
			varName.lexeme,
		),
	)

	nextToken := parser.peekNextToken()

	for nextToken.lexeme != ";" {
		// Check if it's a multi var declaration.
		commaSymbol := parser.getNextToken()
		varName := parser.getNextToken()

		if commaSymbol.lexeme != "," && varName.tokenType != IDENTIFIER {
			panic("invalid class variable declaration.")
		}

		parser.write(
			fmt.Sprintf(
				"<symbol>,</symbol>\n"+
					"<identifier>%s</identifier>\n",
				varName.lexeme,
			),
		)

		nextToken = parser.peekNextToken()
	}

	parser.write(
		"<symbol>;</symbol>\n" +
			"</classVarDec>\n",
	)
	// discard final ";" token.
	parser.getNextToken()
}

func (parser *Parser) Parse(tokens *list.List, outFile string) {
	if tokens.Len() < 4 {
		//	The minimum number of tokens for a valid Jack program is 4. "class className {}"
		srcFile := strings.Replace(outFile, ".xml", ".jack", -1)
		log.Fatalf(fmt.Sprintf("File: %s is not a valid Jack program", srcFile))
	}

	file, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	parser.outFile = file
	parser.tokens = tokens

	classKeyword := parser.getNextToken()
	className := parser.getNextToken()
	leftBrace := parser.getNextToken()

	if classKeyword.lexeme != "class" {
		panic("All Jack programs must begin with a class declaration")
	}

	if className.tokenType != IDENTIFIER {
		panic("class name must be a valid identifier")
	}

	if leftBrace.lexeme != "{" {
		panic("class name must be followed by an opening brace '{' character")
	}

	parser.write(fmt.Sprintf(
		"<class>\n"+
			"<keyword>class</keyword>\n"+
			"<identifier>%s</identifier>\n"+
			"<symbol> { </symbol>\n",
		className.lexeme,
	))

	nextToken := parser.peekNextToken()

	for parser.tokens.Len() > 0 {
		switch nextToken.lexeme {
		case "field", "static":
			parser.parseClassVarDec()
			nextToken = parser.peekNextToken()
		case "constructor", "function", "method":
			parser.parseSubroutineDec()
			nextToken = parser.peekNextToken()
		case "}":
			// Signifies the end of the jack program, so it must be the last token in the list.
			if parser.tokens.Len() > 1 {
				panic("error: invalid Jack program.")
			}
			parser.write(
				"<symbol>}</symbol>\n" +
					"</class>\n",
			)
			// Discard the token
			parser.getNextToken()

		default:
			panic(fmt.Sprintf("error: invalid token: %s", nextToken.lexeme))
		}
	}
}
