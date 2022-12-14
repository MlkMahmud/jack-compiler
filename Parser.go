package main

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

type ParserErrorType int

const (
	UNEXPECTED_TOKEN ParserErrorType = iota
	UNEXPECTED_END_OF_INPUT
)

type ParserError struct {
	errorMessage string
}

func (e *ParserError) Error() string {
	return e.errorMessage
}

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
			"<%s> %s </%s>\n",
			tag,
			writeSymbol(lexeme.(string)),
			tag,
		)
	}
	_, err := parser.outFile.WriteString(str)
	if err != nil {
		panic(err)
	}
}

func (parser *Parser) emitError(errorType ParserErrorType, token any) {
	srcFile := replaceFileExt(parser.outFile.Name(), ".jack")
	var errorMessage string

	switch errorType {
	case UNEXPECTED_TOKEN:
		errorMessage = fmt.Sprintf(
			"(%s):[%d:%d]: Syntax error: unexpected token '%s'\n",
			srcFile,
			token.(Token).lineNum,
			token.(Token).colNum,
			token.(Token).lexeme,
		)
		
	case UNEXPECTED_END_OF_INPUT:
		errorMessage = fmt.Sprintf(
			"(%s): Syntax error: unexpected end of input\n",
			srcFile,
		)

	default: 
		panic(fmt.Sprintf("Error Type: [%d] is not a valid parser error", errorType))
	}

	panic(&ParserError{errorMessage})
}

func (parser *Parser) getNextToken() Token {
	item := parser.tokens.Front()

	if item == nil {
		parser.emitError(UNEXPECTED_END_OF_INPUT, nil)
	}

	parser.tokens.Remove(item)
	return item.Value.(Token)
}

func (parser *Parser) peekNextToken() Token {
	item := parser.tokens.Front()

	if item == nil {
		parser.emitError(UNEXPECTED_END_OF_INPUT, nil)
	}

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
	parser.emitError(UNEXPECTED_TOKEN, token)
}

func (parser *Parser) parseParameterList() {
	// GRAMMAR: ((type varName), (',' type varName)*)?
	parser.write("parameterList", nil)

	for nextToken := parser.peekNextToken(); !isSymbol(nextToken, []string{")"}); nextToken = parser.peekNextToken() {
		parser.parseToken([]string{"boolean", "char", "className", "int"})
		parser.parseToken([]string{"varName"})
		if isSymbol(parser.peekNextToken(), []string{","}) {
			parser.parseToken([]string{","})
		}
	}

	parser.write("/parameterList", nil)
}

func (parser *Parser) parseVarDec() {
	// GRAMMAR: 'var' type varName (',' varName)* ';'
	parser.write("varDec", nil)
	parser.parseToken([]string{"var"})
	parser.parseToken([]string{"boolean", "char", "className", "int"})
	parser.parseToken([]string{"varName"})

	for nextToken := parser.peekNextToken(); !isSymbol(nextToken, []string{";"}); nextToken = parser.peekNextToken() {
		parser.parseToken([]string{","})
		parser.parseToken([]string{"varName"})
	}

	parser.parseToken([]string{";"})
	parser.write("/varDec", nil)
}

func (parser *Parser) parseTerm() {
	parser.write("term", nil)
	token := parser.peekNextToken()

	// An expression enclosed in parentheses
	if isSymbol(token, []string{"("}) {
		parser.parseToken([]string{"("})
		parser.parseExpression()
		parser.parseToken([]string{")"})
	} else if token.tokenType == INTEGER_CONSTANT ||
		token.tokenType == STRING_CONSTANT ||
		token.tokenType == KEYWORD {
		parser.parseToken([]string{"integerConstant", "stringConstant", "true", "false", "null", "this"})
	} else if token.tokenType == IDENTIFIER {
		parser.parseToken([]string{"varName"})
		// If the next token is a left bracket. We're dealing with an array indexing expression.
		if isSymbol(parser.peekNextToken(), []string{"["}) {
			parser.parseToken([]string{"["})
			parser.parseExpression()
			parser.parseToken([]string{"]"})
			// If the next token is a left parenthesis or a period symbol. We're dealing with a subroutine call.
		} else if isSymbol(parser.peekNextToken(), []string{"(", "."}) {
			if isSymbol(parser.peekNextToken(), []string{"."}) {
				parser.parseToken([]string{"."})
				parser.parseToken([]string{"subroutineName"})
			}
			parser.parseToken([]string{"("})
			parser.parseExpressionList()
			parser.parseToken([]string{")"})
		}
	} else if isSymbol(token, []string{"-", "~"}) {
		// If the token is a unary symbol it must be followed by a term
		parser.parseToken([]string{"-", "~"})
		parser.parseTerm()
	} else {
		parser.emitError(UNEXPECTED_TOKEN, token)
	}
	parser.write("/term", nil)
}

func (parser *Parser) parseExpression() {
	// GRAMMAR: term (op term)*
	if isSymbol(parser.peekNextToken(), []string{";", "]", ")"}) {
		return;
	}
	parser.write("expression", nil)
	parser.parseTerm()
	for nextToken := parser.peekNextToken(); isSymbol(nextToken, []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}); nextToken = parser.peekNextToken() {
		parser.parseToken([]string{"+", "-", "*", "/", "&", "|", "<", ">", "="})
		parser.parseTerm()
	}
	parser.write("/expression", nil)
}

func (parser *Parser) parseExpressionList() {
	// GRAMMAR: (expression (',' expression)*)?
	parser.write("expressionList", nil)

	for nextToken := parser.peekNextToken(); !isSymbol(nextToken, []string{")"}); nextToken = parser.peekNextToken() {
		parser.parseExpression()

		if isSymbol(parser.peekNextToken(), []string{","}) {
			parser.parseToken([]string{","})
		}
	}
	parser.write("/expressionList", nil)
}

func (parser *Parser) parseSubroutineCall() {
	// GRAMMAR: subroutineName '(' expressionList ')'
	// GRAMMAR: (className | varName) '.' subroutineName '(' expressionList ')'
	parser.parseToken([]string{"className", "subroutineName", "varName"})
	if isSymbol(parser.peekNextToken(), []string{"."}) {
		parser.parseToken([]string{"."})
		parser.parseToken([]string{"subroutineName"})
	}

	parser.parseToken([]string{"("})
	parser.parseExpressionList()
	parser.parseToken([]string{")"})
}

func (parser *Parser) parseDoStatement() {
	// GRAMMAR: 'do' subroutineCall ';'
	parser.write("doStatement", nil)
	parser.parseToken([]string{"do"})
	parser.parseSubroutineCall()
	parser.parseToken([]string{";"})
	parser.write("/doStatement", nil)
}

func (parser *Parser) parseIfStatement() {
	// GRAMMAR: 'if' '(' expression ')' '{' statements '}' ('else' '{' statements '}')?
	parser.write("ifStatement", nil)

	parser.parseToken([]string{"if"})
	parser.parseToken([]string{"("})
	parser.parseExpression()
	parser.parseToken([]string{")"})
	parser.parseToken([]string{"{"})
	parser.parseStatements()
	parser.parseToken([]string{"}"})

	if isKeyword(parser.peekNextToken(), []string{"else"}) {
		parser.parseToken([]string{"else"})
		parser.parseToken([]string{"{"})
		parser.parseStatements()
		parser.parseToken([]string{"}"})
	}

	parser.write("/ifStatement", nil)
}

func (parser *Parser) parseLetStatement() {
	// GRAMMAR: 'let' varName ('[' expression ']')? '=' expression ';'
	parser.write("letStatement", nil)
	parser.parseToken([]string{"let"})
	parser.parseToken([]string{"varName"})
	
	// If the next token is a left bracket "[". We're dealing with an array index assignment.
	if isSymbol(parser.peekNextToken(), []string{"["}) {
		parser.parseToken([]string{"["})
		parser.parseExpression()
		parser.parseToken([]string{"]"})
	}

	parser.parseToken([]string{"="})
	parser.parseExpression()
	parser.parseToken([]string{";"})
	parser.write("/letStatement", nil)
}

func (parser *Parser) parseReturnStatement() {
	// GRAMMAR: 'return' expression? ';'
	parser.write("returnStatement", nil)
	parser.parseToken([]string{"return"})
	parser.parseExpression()
	parser.parseToken([]string{";"})
	parser.write("/returnStatement", nil)
}

func (parser *Parser) parseWhileStatement() {
	// GRAMMAR: 'while' '(' expression ')' '{' statements '}'
	parser.write("whileStatement", nil)
	
	parser.parseToken([]string{"while"})
	parser.parseToken([]string{"("})
	parser.parseExpression()
	parser.parseToken([]string{")"})
	parser.parseToken([]string{"{"})
	parser.parseStatements()
	parser.parseToken([]string{"}"})

	parser.write("/whileStatement", nil)
}

func (parser *Parser) parseStatements() {
	// GRAMMAR: statement*
	parser.write("statements", nil)

	for token := parser.peekNextToken(); !isSymbol(token, []string{"}"}); token = parser.peekNextToken() {
		switch token.lexeme {
		case "do":
			parser.parseDoStatement()
		case "if":
			parser.parseIfStatement()
		case "let":
			parser.parseLetStatement()
		case "return":
			parser.parseReturnStatement()
		case "while":
			parser.parseWhileStatement()
		default:
			parser.emitError(UNEXPECTED_TOKEN, token)
		}
	}
	parser.write("/statements", nil)
}

func (parser *Parser) parseSubroutineBody() {
	// GRAMMAR: '{' varDec* statements '}'
	parser.write("subroutineBody", nil)
	parser.parseToken([]string{"{"})

	for end := parser.peekNextToken(); !isSymbol(end, []string{"}"}); end = parser.peekNextToken() {
		if isKeyword(end, []string{"var"}) {
			parser.parseVarDec()
		} else {
			parser.parseStatements()
		}
	}
	parser.parseToken([]string{"}"})
	parser.write("/subroutineBody", nil)
}

func (parser *Parser) parseSubroutineDec() {
	// GRAMMAR: ('constructor' | 'function' | 'method') ('void' | type) subroutineName '(' parameterList ')' subroutineBody
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
	// GRAMMAR: ('static' | 'field') ('int' | 'boolean' | 'char' | className) varName (',' varName)* ';'
	parser.write("classVarDec", nil)

	parser.parseToken([]string{"field", "static"})
	parser.parseToken([]string{"boolean", "char", "className", "int"})
	parser.parseToken([]string{"varName"})

	nextToken := parser.peekNextToken()

	// Check if it's a multi var declaration.
	for !isSymbol(nextToken, []string{";"}) {
		parser.parseToken([]string{","})
		parser.parseToken([]string{"varName"})

		nextToken = parser.peekNextToken()
	}

	parser.parseToken([]string{";"})
	parser.write("/classVarDec", nil)
}

func (parser *Parser) Parse(tokens *list.List, outFile string) {
	// If the tokens list is empty, we're dealing with an empty file.
	if (tokens.Len() < 1) {
		return
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

		if r := recover(); r != nil {
			var parserError *ParserError
			if errors.As(r.(error), &parserError) {
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

	for nextToken := parser.peekNextToken(); !isSymbol(nextToken, []string{"}"}); nextToken = parser.peekNextToken() {
		if isKeyword(nextToken, []string{"field", "static"}) {
			parser.parseClassVarDec()
		} else if isKeyword(nextToken, []string{"constructor", "function", "method"}) {
			parser.parseSubroutineDec()
		} else {
			parser.emitError(UNEXPECTED_TOKEN, nextToken)
		}
	}

	parser.parseToken([]string{"}"})
	parser.write("/class", nil)
}
