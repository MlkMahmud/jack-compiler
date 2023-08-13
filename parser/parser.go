package parser

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"

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
				if val == "className" || val == "subroutineName" || val == "varName" {
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

func (parser *Parser) parseParameterList() (params []Parameter) {
	// GRAMMAR: ((type varName), (',' type varName)*)?
	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{")"}); nextToken = parser.peekNextToken() {
		paramTypeToken := parser.getNextToken()
		paramNameToken := parser.getNextToken()
		parser.assertToken(paramTypeToken, []string{"boolean", "char", "className", "int"})
		parser.assertToken(paramNameToken, []string{"varName"})
		params = append(params, Parameter{Name: paramNameToken.Lexeme, Type: paramTypeToken.Lexeme})

		if helpers.IsSymbol(parser.peekNextToken(), []string{","}) {
			parser.getNextToken()
		}
	}

	return params
}

func (parser *Parser) parseVarDec() (vars []VarDecl) {
	// GRAMMAR: 'var' type varName (',' varName)* ';'
	parser.assertToken(parser.getNextToken(), []string{"var"})

	varTypeToken := parser.getNextToken()
	varNameToken := parser.getNextToken()
	parser.assertToken(varTypeToken, []string{"boolean", "char", "className", "int"})
	parser.assertToken(varNameToken, []string{"varName"})

	vars = append(vars, VarDecl{Name: varNameToken.Lexeme, Type: varTypeToken.Lexeme, Kind: Var})

	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{";"}); nextToken = parser.peekNextToken() {
		parser.assertToken(parser.getNextToken(), []string{","})
		nextVarNameToken := parser.getNextToken()
		parser.assertToken(nextVarNameToken, []string{"varName"})

		vars = append(vars, VarDecl{Name: nextVarNameToken.Lexeme, Type: varTypeToken.Lexeme, Kind: Var})
	}

	parser.assertToken(parser.getNextToken(), []string{";"})
	return vars
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

func (parser *Parser) parseExpression_() Expr {
	// GRAMMAR: term (op term)*
	token := parser.getNextToken()
	var expr Expr

	// Parenthesized expression
	if helpers.IsSymbol(token, []string{"("}) {
		expression := parser.parseExpression_()
		expr = ParenExpr{Expression: expression}
		parser.assertToken(parser.getNextToken(), []string{")"})
	} else if helpers.IsLiteralType(token) {
		if token.TokenType == constants.INTEGER_CONSTANT {
			value, err := strconv.ParseInt(token.Lexeme, 10, 16)
			if err != nil {
				parser.emitError(UNEXPECTED_TOKEN, token)
			}
			expr = Literal{Value: value}
		} else if helpers.Contains([]string{"true", "false"}, token.Lexeme) {
			value, err := strconv.ParseBool(token.Lexeme)
			if err != nil {
				parser.emitError(UNEXPECTED_TOKEN, token)
			}
			expr = Literal{Value: value}
		} else {
			expr = Literal{Value: token.Lexeme}
		}
	} else if token.TokenType == constants.IDENTIFIER {
		nextToken := parser.peekNextToken()
		// If the next token is a left parenthesis '(' or a period symbol '.' we're dealing with a call expression
		if helpers.IsSymbol(nextToken, []string{"("}) {
			expr = CallExpr{
				Arguments: parser.parseExpressionList_(),
				Callee: Indentifier{ Name: token.Lexeme },
			}
		} else if helpers.IsSymbol(nextToken, []string{"."}) {
			varNameToken := parser.getNextToken()
			parser.assertToken(varNameToken, []string{"varName"})
			expr = CallExpr{
				Arguments: parser.parseExpressionList_(),
				Callee: MemberExpr{ Object: Indentifier{ Name: token.Lexeme }, Property: Indentifier{ varNameToken.Lexeme } },
			}
			// If the next token is a left bracket. We're dealing with an array indexing expression.
		} else if helpers.IsSymbol(nextToken, []string{"["}) {
			
		}
	}
	return expr
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

func (parser *Parser) parseExpressionList_() (args []Expr) {
	// GRAMMAR: (expression (',' expression)*)?
	parser.assertToken(parser.getNextToken(), []string{"("})

	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{")"}); nextToken = parser.peekNextToken() {
		args = append(args, parser.parseExpression_())
		if helpers.IsKeyword(parser.peekNextToken(), []string{","}) {
			// discard "," token before next expresssion in list
			parser.getNextToken()
		}
	}

	parser.assertToken(parser.getNextToken(), []string{")"})
	return args
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

func (parser *Parser) parseDoStatement() (stmt DoStmt) {
	// GRAMMAR: 'do' subroutineName '(' expressionList ')' ';' | 'do' (className | varName) '.' subroutineName '(' expressionList ') ';'
	parser.assertToken(parser.getNextToken(), []string{"do"})
	token := parser.getNextToken()
	parser.assertToken(token, []string{"className", "subroutineName", "varName"})

	if helpers.IsSymbol(parser.peekNextToken(), []string{"."}) {
		// discard '.' token
		parser.getNextToken()
		subroutineNameToken := parser.getNextToken()
		parser.assertToken(subroutineNameToken, []string{"subroutineName"})

		stmt.ObjectName = token.Lexeme
		stmt.SubroutineName = subroutineNameToken.Lexeme
	} else {
		stmt.SubroutineName = token.Lexeme
	}

	stmt.Arguments = parser.parseExpressionList_()
	return stmt
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

func (parser *Parser) parseStatement() Stmt {
	token := parser.peekNextToken()
	var stmt Stmt

	switch token.Lexeme {
	case "do":
		stmt = parser.parseDoStatement()
	default:
		parser.emitError(UNEXPECTED_TOKEN, token)
	}
	return stmt
}

func (parser *Parser) parseSubroutineBody() (body BlockStmt) {
	// GRAMMAR: '{' varDec* statements '}'
	parser.assertToken(parser.getNextToken(), []string{"{"})

	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{"}"}); nextToken = parser.peekNextToken() {
		if helpers.IsKeyword(nextToken, []string{"var"}) {
			vars := parser.parseVarDec()
			body.Vars = append(body.Vars, vars...)
		} else {
			body.Statements = append(body.Statements, parser.parseStatement())
		}
	}

	parser.assertToken(parser.getNextToken(), []string{"}"})
	return body
}

func (parser *Parser) parseSubroutineDec() (subroutine SubroutineDecl) {
	// GRAMMAR: ('constructor' | 'function' | 'method') ('void' | type) subroutineName '(' parameterList ')' subroutineBody
	subroutineKindToken := parser.getNextToken()
	subroutineTypeToken := parser.getNextToken()
	subroutineNameToken := parser.getNextToken()

	parser.assertToken(subroutineKindToken, []string{"constructor", "function", "method"})
	parser.assertToken(subroutineTypeToken, []string{"boolean", "char", "className", "int", "void"})
	parser.assertToken(subroutineNameToken, []string{"subroutineName"})

	var subroutineKind SubroutineKind

	switch subroutineKindToken.Lexeme {
	case "constructor":
		subroutineKind = Constructor
	case "function":
		subroutineKind = Function
	default:
		subroutineKind = Method
	}

	subroutine.Name = subroutineNameToken.Lexeme
	subroutine.Kind = subroutineKind
	subroutine.Type = subroutineTypeToken.Lexeme

	parser.assertToken(parser.getNextToken(), []string{"("})
	subroutine.Params = append(subroutine.Params, parser.parseParameterList()...)
	parser.assertToken(parser.getNextToken(), []string{")"})

	body := parser.parseSubroutineBody()
	subroutine.Statements = body.Statements
	subroutine.Vars = body.Vars

	return subroutine
}

func (parser *Parser) parseClassVarDec() (vars []VarDecl) {
	// GRAMMAR: ('static' | 'field') type varName (',' varName)* ';'
	varKindToken := parser.getNextToken()
	varTypeToken := parser.getNextToken()
	varNameToken := parser.getNextToken()

	parser.assertToken(varKindToken, []string{"field", "static"})
	parser.assertToken(varTypeToken, []string{"boolean", "char", "className", "int"})
	parser.assertToken(varNameToken, []string{"varName"})

	var varKind VarKind

	if varKindToken.Lexeme == "field" {
		varKind = Field
	} else {
		varKind = Static
	}

	vars = append(vars, VarDecl{Name: varNameToken.Lexeme, Type: varTypeToken.Lexeme, Kind: varKind})

	// Check if it's a multi var declaration.
	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{";"}); nextToken = parser.peekNextToken() {
		parser.assertToken(parser.getNextToken(), []string{","})
		parser.assertToken(parser.peekNextToken(), []string{"varName"})

		nextVarNameToken := parser.getNextToken()
		vars = append(vars, VarDecl{Name: nextVarNameToken.Lexeme, Type: varTypeToken.Lexeme, Kind: varKind})
	}

	parser.assertToken(parser.getNextToken(), []string{";"})
	return vars
}

func (parser *Parser) Parse(tokens []constants.Token) (compilationUnit CompilationUnit) {
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
		return compilationUnit
	}

	parser.tokens = tokens
	parser.assertToken(parser.getNextToken(), []string{"class"})
	classNameToken := parser.getNextToken()
	parser.assertToken(classNameToken, []string{"className"})
	parser.assertToken(parser.getNextToken(), []string{"{"})

	compilationUnit.Name = classNameToken.Lexeme

	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{"}"}); nextToken = parser.peekNextToken() {
		if helpers.IsKeyword(nextToken, []string{"field", "static"}) {
			vars := parser.parseClassVarDec()
			compilationUnit.Vars = append(compilationUnit.Vars, vars...)
		} else if helpers.IsKeyword(nextToken, []string{"constructor", "function", "method"}) {
			subroutine := parser.parseSubroutineDec()
			compilationUnit.Subroutines = append(compilationUnit.Subroutines, subroutine)
		} else {
			parser.emitError(UNEXPECTED_TOKEN, nextToken)
		}
	}

	return compilationUnit
}
