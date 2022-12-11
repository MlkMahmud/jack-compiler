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
			writeSymbol(lexeme.(string)),
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
		parser.throwSyntaxError(token)
	}
	parser.write("/term", nil)
}

func (parser *Parser) parseExpression() {
	// GRAMMAR: term (op term)*
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

func (parser *Parser) parseReturnStatement() {
	// GRAMMAR: 'return' expression? ';'
	parser.write("returnStatement", nil)
	parser.parseToken([]string{"return"})
	
	if !isSymbol(parser.peekNextToken(), []string{";"}) {
		parser.parseExpression()
	}

	parser.parseToken([]string{";"})
	parser.write("/returnStatement", nil)
}

func (parser *Parser) parseStatements() {
	parser.write("statements", nil)

	for token := parser.peekNextToken(); !isSymbol(token, []string{"}"}); token = parser.peekNextToken() {
		switch token.lexeme {
		case "do":
			parser.parseDoStatement()
		case "if":
			// parseIfStatement
		case "let":
			// parseLetStatement
		case "return":
			parser.parseReturnStatement()
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
	for !isSymbol(nextToken, []string{";"}) {
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

	for nextToken := parser.peekNextToken(); !isSymbol(nextToken, []string{"}"}); nextToken = parser.peekNextToken() {
		if isKeyword(nextToken, []string{"field", "static"}) {
			parser.parseClassVarDec()
		} else if isKeyword(nextToken, []string{"constructor", "function", "method"}) {
			parser.parseSubroutineDec()
		} else {
			parser.throwSyntaxError(nextToken)
		}
	}

	parser.parseToken([]string{"}"})
	parser.write("/class", nil)
}
