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
	parser.write("<parameterList>\n")

	for nextToken := parser.peekNextToken(); !isSymbol(nextToken, ")"); nextToken = parser.peekNextToken() {
		argType := parser.getNextToken()
		argName := parser.getNextToken()

		if !isValidType(argType) {
			panic(fmt.Sprintf("error: invalid token: %s", argType.lexeme))
		}

		if !isIdentifier(argName) {
			panic(fmt.Sprintf("error: invalid token: %s", argName.lexeme))
		}

		parser.write(fmt.Sprintf(
			"<%s>%s</%s>\n<identifier>%s</identifier>\n",
			argType.tokenType,
			argType.lexeme,
			argType.tokenType,
			argName.lexeme,
		))

		if isSymbol(parser.peekNextToken(), ",") {
			parser.write("<symbol>,</symbol>\n")
			// discard "," token.
			parser.getNextToken()
		}
	}

	parser.write("</parameterList>\n<symbol>)</symbol>\n")

	// discard closing brace ")" token.
	parser.getNextToken()
}

func (parser *Parser) parseVarDec() {
	varDec := parser.getNextToken()
	varType := parser.getNextToken()
	varName := parser.getNextToken()

	if !isKeyword(varDec, "var") {
		panic(fmt.Sprintf("error: invalid token: %s", varDec.lexeme))
	}

	if !isValidType(varType) {
		panic(fmt.Sprintf("error: invalid token: %s", varType.lexeme))
	}

	if !isIdentifier(varName) {
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

	for nextToken := parser.peekNextToken(); !isSymbol(nextToken, ";"); nextToken = parser.peekNextToken() {
		commaSymbol := parser.getNextToken()
		nextVarName := parser.getNextToken()

		if !isSymbol(commaSymbol, ",") {
			panic(fmt.Sprintf("error: invalid token: %s", commaSymbol.lexeme))
		}

		if !isIdentifier(nextVarName) {
			panic(fmt.Sprintf("error: invalid token: %s", nextVarName.lexeme))
		}

		parser.write(fmt.Sprintf(
			"\t\t\t\t\t<symbol>,</symbol>\n"+
				"\t\t\t\t\t<identifier>%s</identifier>\n",
			nextVarName.lexeme,
		))
	}

	parser.write(
		"\t\t\t\t\t<symbol>;</symbol>\n" +
			"\t\t\t\t</varDec>\n",
	)
	// discard ";" token.
	parser.getNextToken()
}

func (parser *Parser) parseExpression() {
}

func (parser *Parser) parseDoStatement() {
	doKeyword := parser.getNextToken()
	subroutetineOrClassName := parser.getNextToken()

	if !isKeyword(doKeyword, "do") {
		panic(fmt.Sprintf("error: invalid token: %s", doKeyword.lexeme))
	}

	if !isIdentifier(subroutetineOrClassName) {
		panic(fmt.Sprintf("error: invalid token: %s", subroutetineOrClassName.lexeme))
	}

	parser.write("<keyword>do</keyword>\n")
	periodOrLparen := parser.peekNextToken()

	if isSymbol(periodOrLparen, ".") {
		className := subroutetineOrClassName
		// discard "." token.
		parser.getNextToken()
		subroutineName := parser.getNextToken()

		if !isIdentifier(subroutineName) {
			panic(fmt.Sprintf("error: invalid token: %s", subroutineName.lexeme))
		}

		if nextToken := parser.peekNextToken(); !isSymbol(nextToken, "(") {
			panic(fmt.Sprintf("error: invalid token: %s", nextToken.lexeme))
		}

		parser.write(fmt.Sprintf(
			"<identifier>%s</identifier>\n"+
				"<symbol>.</symbol>\n"+
				"<identifier>%s</identifier>\n",
			className.lexeme,
			subroutineName.lexeme,
		))
	} else if isSymbol(periodOrLparen, "(") {
		subroutineName := subroutetineOrClassName
		parser.write(fmt.Sprintf(
			"<identifier>%s</identifier>\n"+
				"<symbol>(</symbol>\n",
			subroutineName.lexeme,
		))
	} else {
		panic(fmt.Sprintf("error: invalid token: %s", periodOrLparen.lexeme))
	}

	// discard the opening paren char "(" and parse the first expression in the expression list.
	parser.getNextToken()

	parser.write("<symbol>(</symbol>\n<expressionList>\n")

	// Parse all expressions in list until we hit the closing paren ")"char.
	for end := parser.peekNextToken(); !isSymbol(end, ")"); end = parser.peekNextToken() {
		parser.parseExpression()
		nextToken := parser.peekNextToken()

		if !isSymbol(nextToken, ",") && !isSymbol(nextToken, ")") {
			panic(fmt.Sprintf("error: invalid token: %s", nextToken.lexeme))
		}

		if isSymbol(nextToken, ",") {
			// A comma symbol means there are more expressions in the list.
			// Discard the "comma" token and let the "parseExpresssion" function handle the expression
			// in the next iteration of this loop.
			parser.getNextToken()
		}
	}
	parser.write("</expressionList>\n<symbol>)</symbol>")
	// discard ")" and ";" tokens
	parser.getNextToken()
	parser.getNextToken()
}

func (parser *Parser) parseStatements() {
	token := parser.peekNextToken()

	if !isStatementDec(token) {
		panic(fmt.Sprintf("error: invalid token: %s", token.lexeme))
	}

	parser.write("<statements>\n")

	for !isSymbol(token, "}") {
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
			panic(fmt.Sprintf("error: invalid token: %s", token.lexeme))
		}
		token = parser.peekNextToken()
	}
	parser.write("</statements>\n")
}

func (parser *Parser) parseSubroutineBody() {
	parser.write(
		"<subroutineBody>\n" +
			"<symbol>{</symbol>\n",
	)

	for end := parser.peekNextToken(); !isSymbol(end, "}"); end = parser.peekNextToken() {
		if isKeyword(end, "var") {
			parser.parseVarDec()
		} else {
			parser.parseStatements()
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

	if !isSubroutineDec(subroutineType) {
		panic(fmt.Sprintf("error: invalid token: %s", subroutineType.lexeme))
	}

	if !isValidType(subroutineReturnType) && subroutineReturnType.lexeme != "void" {
		panic(fmt.Sprintf("error: invalid token: %s", subroutineReturnType.lexeme))
	}

	if !isIdentifier(subroutineName) {
		panic(fmt.Sprintf("error: invalid token: %s", subroutineName.lexeme))
	}

	if !isSymbol(leftParen, "(") {
		panic(fmt.Sprintf("error: invalid token: %s", leftParen.lexeme))
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

	if !isSymbol(nextToken, "{") {
		panic(fmt.Sprintf("error: invalid token: %s", nextToken.lexeme))
	}
	// discard left brace token "{"
	parser.getNextToken()

	if !isSymbol(parser.peekNextToken(), "}") {
		parser.parseSubroutineBody()
	}

	parser.write("</subroutineDec>\n")
}

func (parser *Parser) parseClassVarDec() {
	varType := parser.getNextToken()
	dataType := parser.getNextToken()
	varName := parser.getNextToken()

	if !isClassVarDec(varType) {
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

	if !isIdentifier(varName) {
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

	for !isSymbol(nextToken, ";") {
		// Check if it's a multi var declaration.
		commaSymbol := parser.getNextToken()
		nextVarName := parser.getNextToken()

		if !isSymbol(commaSymbol, ",") {
			panic(fmt.Sprintf("error: invalid token: %s", commaSymbol.lexeme))
		}

		if !isIdentifier(nextVarName) {
			panic(fmt.Sprintf("error: invalid token: %s", nextVarName.lexeme))
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

	if !isKeyword(classKeyword, "class") {
		panic(fmt.Sprintf("error: invalid token: %s", classKeyword.lexeme))
	}

	if !isIdentifier(className) {
		panic(fmt.Sprintf("error: invalid token: %s", className.lexeme))
	}

	if !isSymbol(leftBrace, "{") {
		panic(fmt.Sprintf("error: invalid token: %s", leftBrace.lexeme))
	}

	parser.write(fmt.Sprintf(
		"<class>\n"+
			"<keyword>class</keyword>\n"+
			"<identifier>%s</identifier>\n"+
			"<symbol> { </symbol>\n",
		className.lexeme,
	))

	for end := parser.peekNextToken(); !isSymbol(end, "}"); end = parser.peekNextToken() {
		if isClassVarDec(end) {
			parser.parseClassVarDec()
		} else if isSubroutineDec(end) {
			parser.parseSubroutineDec()
		} else {
			panic(fmt.Sprintf("error: invalid token: %s", end.lexeme))
		}
	}
	// Signifies the end of the jack program, so it must be the last token in the list.
	if parser.tokens.Len() > 1 {
		panic("error: invalid Jack program.")
	}

	// Discard the token
	parser.getNextToken()

	parser.write(
		"<symbol>}</symbol>\n" +
			"</class>\n",
	)
}
