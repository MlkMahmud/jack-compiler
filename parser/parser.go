package parser

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/MlkMahmud/jack-compiler/constants"
	"github.com/MlkMahmud/jack-compiler/helpers"
)

/*
LEXICAL ELEMENTS

keyword:         'class' | 'constructor' | 'function' | 'method' |
								 'field' | 'static' | 'var' | 'int' | 'char' | 'boolean' |
								 'void' | 'true' | 'false' |'null' | 'this' | 'let' | 'do' |
								 'if' | 'else' | 'while' | 'return'

symbol:          '{' | '}' | '(' | ')' | '[' | ']' | '.' | ',' | ';' |
                 '+' | '-' | '*' | '/' | '&' | '|' | '<' | '>' | '=' | '~'

integerConstant: decimal integer in range 0...32767

stringConstant:  Sequence of Unicode characters surrounded by double quotes (cannot contain newline or double quote characters)

identifier:      Sequence of letters, digits, underscores NOT starting with a digit.


PROGRAM STRUCTURE

class:           'class' className '{' classVarDec* subroutineDec* '}'

classVarDec:     ('static' | 'field') type varName (',' varName)* ';'

type:            'int' | 'char' | 'boolean' | className

subroutineDec:   ('constructor' | 'function' | 'method') ('void' | type) subroutineName '(' parameterList ')' subroutineBody

parameterList:   ((type identifier) (',' <type> identifier)*)?

subroutineBody:  '{' varDec* statements '}'

varDec:          'var' ('int' | 'char' | 'boolean' | identifier) identifier (',' identifier)* ';'

className:       identifier

subroutineName:  identifier

varName:         identifier

STATEMENTS

statements:      statement*

statement:       letStatement | ifStatement | whileStatement | doStatement | returnStatement

letStatement:    'let' identifier ('[' expression ']')? '=' expression ';'

ifStatement      'if' '(' <expression> ')' '{' statements> '}' ('else' '{' statements> '}')?

whileStatement:  'while' '(' expression ')' '{' statements '}' ('else' '{' statements '}')?

doStatement:     'do' (identifier '.')? identifier '(' <expressionList> ')' ';'
	
returnStatement: 'return' <expression>? ';'


EXPRESSIONS

expression:      term (op term)*

term:            integerConstant | stringConstant | keywordConstant
  			         varName | varName '[' expression ']' | subroutineCall |
  			         '(' expression ')' | unaryOp term

subroutineCall:  subroutineName '(' expressionList ')' | (className|varName) '.' subroutineName '(' expressionList ')'

expressionList:  (expression (',' expression)* )?

op:              '+' | '-' | '*' | '/' | '&' | '|' | '<' | '>' | '='

unaryOp:         '-' | '~' 

keywordConstant: 'true' | 'false' | 'null' | 'this'
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
	tokens  []constants.Token
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
			helpers.WriteSymbol(lexeme.(string)),
			tag,
		)
	}
	_, err := parser.outFile.WriteString(str)
	if err != nil {
		panic(err)
	}
}

func (parser *Parser) emitError(errorType ParserErrorType, token any) {
	srcFile := helpers.ReplaceFileExt(parser.outFile.Name(), ".jack")
	var errorMessage string

	switch errorType {
	case UNEXPECTED_TOKEN:
		errorMessage = fmt.Sprintf(
			"(%s):[%d:%d]: Syntax error: unexpected token '%s'\n",
			srcFile,
			token.(constants.Token).LineNum,
			token.(constants.Token).ColNum,
			token.(constants.Token).Lexeme,
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

func (parser *Parser) getNextToken() constants.Token {
	if len(parser.tokens) == 0 {
		parser.emitError(UNEXPECTED_END_OF_INPUT, nil)
	}

	token := parser.tokens[0]
	parser.tokens = parser.tokens[1:]
	return token
}

func (parser *Parser) peekNextToken() constants.Token {
	if len(parser.tokens) == 0 {
		parser.emitError(UNEXPECTED_END_OF_INPUT, nil)
	}

	return parser.tokens[0]
}

func (parser *Parser) assertToken(token constants.Token, terminals []string) {
	switch token.TokenType {
	case constants.KEYWORD, constants.SYMBOL:
		{
			for _, val := range terminals {
				if val == token.Lexeme {
					return
				}
			}
		}

	case constants.INTEGER_CONSTANT, constants.STRING_CONSTANT:
		{
			for _, val := range terminals {
				if val == token.TokenType.String() {
					return
				}
			}
		}

	case constants.IDENTIFIER:
		{
			for _, val := range terminals {
				if val == "identifier" {
					return
				}
			}
		}
	}
	parser.emitError(UNEXPECTED_TOKEN, token)
}

func (parser *Parser) parseToken(terminals []string) {
	token := parser.getNextToken()

	switch token.TokenType {
	case constants.KEYWORD, constants.SYMBOL:
		{
			for _, val := range terminals {
				if val == token.Lexeme {
					parser.write(token.TokenType.String(), token.Lexeme)
					return
				}
			}
		}

	case constants.INTEGER_CONSTANT, constants.STRING_CONSTANT:
		{
			for _, val := range terminals {
				if val == token.TokenType.String() {
					parser.write(token.TokenType.String(), token.Lexeme)
					return
				}
			}
		}

	case constants.IDENTIFIER:
		{
			for _, val := range terminals {
				if val == "className" || val == "subroutineName" || val == "varName" {
					parser.write(token.TokenType.String(), token.Lexeme)
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

	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{")"}); nextToken = parser.peekNextToken() {
		parser.parseToken([]string{"boolean", "char", "className", "int"})
		parser.parseToken([]string{"varName"})
		if helpers.IsSymbol(parser.peekNextToken(), []string{","}) {
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

	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{";"}); nextToken = parser.peekNextToken() {
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
	if helpers.IsSymbol(token, []string{"("}) {
		parser.parseToken([]string{"("})
		parser.parseExpression()
		parser.parseToken([]string{")"})
	} else if token.TokenType == constants.INTEGER_CONSTANT ||
		token.TokenType == constants.STRING_CONSTANT ||
		token.TokenType == constants.KEYWORD {
		parser.parseToken([]string{"integerConstant", "stringConstant", "true", "false", "null", "this"})
	} else if token.TokenType == constants.IDENTIFIER {
		parser.parseToken([]string{"varName"})
		// If the next token is a left bracket. We're dealing with an array indexing expression.
		if helpers.IsSymbol(parser.peekNextToken(), []string{"["}) {
			parser.parseToken([]string{"["})
			parser.parseExpression()
			parser.parseToken([]string{"]"})
			// If the next token is a left parenthesis or a period symbol. We're dealing with a subroutine call.
		} else if helpers.IsSymbol(parser.peekNextToken(), []string{"(", "."}) {
			if helpers.IsSymbol(parser.peekNextToken(), []string{"."}) {
				parser.parseToken([]string{"."})
				parser.parseToken([]string{"subroutineName"})
			}
			parser.parseToken([]string{"("})
			parser.parseExpressionList()
			parser.parseToken([]string{")"})
		}
	} else if helpers.IsSymbol(token, []string{"-", "~"}) {
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
	if helpers.IsSymbol(parser.peekNextToken(), []string{";", "]", ")"}) {
		return
	}
	parser.write("expression", nil)
	parser.parseTerm()
	for nextToken := parser.peekNextToken(); helpers.IsSymbol(nextToken, []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}); nextToken = parser.peekNextToken() {
		parser.parseToken([]string{"+", "-", "*", "/", "&", "|", "<", ">", "="})
		parser.parseTerm()
	}
	parser.write("/expression", nil)
}

func (parser *Parser) parseExpressionList() {
	// GRAMMAR: (expression (',' expression)*)?
	parser.write("expressionList", nil)

	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{")"}); nextToken = parser.peekNextToken() {
		parser.parseExpression()

		if helpers.IsSymbol(parser.peekNextToken(), []string{","}) {
			parser.parseToken([]string{","})
		}
	}
	parser.write("/expressionList", nil)
}

func (parser *Parser) parseSubroutineCall() {
	// GRAMMAR: subroutineName '(' expressionList ')'
	// GRAMMAR: (className | varName) '.' subroutineName '(' expressionList ')'
	parser.parseToken([]string{"className", "subroutineName", "varName"})
	if helpers.IsSymbol(parser.peekNextToken(), []string{"."}) {
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

	if helpers.IsKeyword(parser.peekNextToken(), []string{"else"}) {
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
	if helpers.IsSymbol(parser.peekNextToken(), []string{"["}) {
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

	for token := parser.peekNextToken(); !helpers.IsSymbol(token, []string{"}"}); token = parser.peekNextToken() {
		switch token.Lexeme {
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

	for end := parser.peekNextToken(); !helpers.IsSymbol(end, []string{"}"}); end = parser.peekNextToken() {
		if helpers.IsKeyword(end, []string{"var"}) {
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

func (parser *Parser) parseClassVarDec() (vars []VarDecl) {
	varKindToken := parser.getNextToken()
	varTypeToken := parser.getNextToken()
	varNameToken := parser.getNextToken()

	parser.assertToken(varKindToken, []string{"field", "static"})
	parser.assertToken(varTypeToken, []string{"boolean", "char", "className", "int"})
	parser.assertToken(varNameToken, []string{"varName"})

	var varKind VarKind;

	if varKindToken.Lexeme == "field" {
		varKind = Field
	} else {
		varKind = Static
	}

	vars = append(vars, VarDecl{ Name: varNameToken.Lexeme, Type: varTypeToken.Lexeme, Kind: varKind })

	// Check if it's a multi var declaration.
	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{";"}); nextToken = parser.peekNextToken() {
		parser.assertToken(parser.getNextToken(), []string{","})
		parser.assertToken(parser.peekNextToken(), []string{"varName"})

		varNameToken := parser.getNextToken()
		vars = append(vars, VarDecl{ Name: varNameToken.Lexeme, Type: varTypeToken.Lexeme, Kind: varKind })
	}

	parser.assertToken(parser.getNextToken(), []string{";"})
	return vars
}

func (parser *Parser) Parse(tokens []constants.Token) CompilationUnit {
	defer func() {
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

	if len(tokens) < 1 {
		return CompilationUnit{}
	}

	parser.tokens = tokens
	parser.assertToken(parser.getNextToken(), []string{"class"})
	classNameToken := parser.getNextToken()
	parser.assertToken(classNameToken, []string{constants.IDENTIFIER.String()})
	parser.assertToken(parser.getNextToken(), []string{"{"})

	compilationUnit := CompilationUnit{Name: classNameToken.Lexeme}

	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{"}"}); nextToken = parser.peekNextToken() {
		if helpers.IsKeyword(nextToken, []string{"field", "static"}) {
			vars := parser.parseClassVarDec()
			compilationUnit.Vars = append(compilationUnit.Vars, vars...)
		} else if helpers.IsKeyword(nextToken, []string{"constructor", "function", "method"}) {
			parser.parseSubroutineDec()
		} else {
			parser.emitError(UNEXPECTED_TOKEN, nextToken)
		}
	}

	return compilationUnit
}
