package parser

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/MlkMahmud/jack-compiler/helpers"
	"github.com/MlkMahmud/jack-compiler/types"
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

ifStatement      'if' '(' expression ')' '{' statements '}' ('else' '{' statements '}')?

whileStatement:  'while' '(' expression ')' '{' statements '}'

doStatement:     'do' (identifier '.')? identifier '(' expressionList ')' ';'

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
	filename string
	tokens   []types.Token
}

func NewParser() *Parser {
	return new(Parser)
}

func (parser *Parser) emitError(errorType ParserErrorType, token any) {
	var errorMessage string

	switch errorType {
	case UNEXPECTED_TOKEN:
		errorMessage = fmt.Sprintf(
			"(%s):[%d:%d]: Syntax error: unexpected token '%s'\n",
			parser.filename,
			token.(types.Token).LineNum,
			token.(types.Token).ColNum,
			token.(types.Token).Lexeme,
		)

	case UNEXPECTED_END_OF_INPUT:
		errorMessage = fmt.Sprintf(
			"(%s): Syntax error: unexpected end of input\n",
			parser.filename,
		)

	default:
		panic(fmt.Sprintf("Error Type: [%d] is not a valid parser error", errorType))
	}
	panic(&ParserError{errorMessage})
}

func (parser *Parser) getNextToken() types.Token {
	if len(parser.tokens) == 0 {
		parser.emitError(UNEXPECTED_END_OF_INPUT, nil)
	}

	token := parser.tokens[0]
	parser.tokens = parser.tokens[1:]
	return token
}

func (parser *Parser) peekNextToken() types.Token {
	if len(parser.tokens) == 0 {
		parser.emitError(UNEXPECTED_END_OF_INPUT, nil)
	}

	return parser.tokens[0]
}

func (parser *Parser) peekNthToken(index int) types.Token {
	if !(index >= 0 && index < len(parser.tokens)) {
		parser.emitError(UNEXPECTED_END_OF_INPUT, nil)
	}

	return parser.tokens[index]
}

func (parser *Parser) assertToken(token types.Token, terminals []string) {
	switch token.TokenType {
	case types.KEYWORD, types.SYMBOL:
		{
			for _, val := range terminals {
				if val == token.Lexeme {
					return
				}
			}
		}

	case types.INTEGER_CONSTANT, types.STRING_CONSTANT:
		{
			for _, val := range terminals {
				if val == token.TokenType.String() {
					return
				}
			}
		}

	case types.IDENTIFIER:
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

func (parser *Parser) parseParameterList() (params []types.Parameter) {
	// GRAMMAR: ((type varName), (',' type varName)*)?
	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{")"}); nextToken = parser.peekNextToken() {
		paramTypeToken := parser.getNextToken()
		paramNameToken := parser.getNextToken()
		parser.assertToken(paramTypeToken, []string{"boolean", "char", "className", "int"})
		parser.assertToken(paramNameToken, []string{"varName"})
		params = append(params, types.Parameter{Name: paramNameToken.Lexeme, Type: paramTypeToken.Lexeme})

		if helpers.IsSymbol(parser.peekNextToken(), []string{","}) {
			parser.getNextToken()
		}
	}

	return params
}

func (parser *Parser) parseVarDec() (vars []types.VarDecl) {
	// GRAMMAR: 'var' type varName (',' varName)* ';'
	parser.assertToken(parser.getNextToken(), []string{"var"})

	varTypeToken := parser.getNextToken()
	varNameToken := parser.getNextToken()
	parser.assertToken(varTypeToken, []string{"boolean", "char", "className", "int"})
	parser.assertToken(varNameToken, []string{"varName"})

	vars = append(vars, types.VarDecl{Name: varNameToken.Lexeme, Type: varTypeToken.Lexeme, Kind: types.Var})

	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{";"}); nextToken = parser.peekNextToken() {
		parser.assertToken(parser.getNextToken(), []string{","})
		nextVarNameToken := parser.getNextToken()
		parser.assertToken(nextVarNameToken, []string{"varName"})

		vars = append(vars, types.VarDecl{Name: nextVarNameToken.Lexeme, Type: varTypeToken.Lexeme, Kind: types.Var})
	}

	parser.assertToken(parser.getNextToken(), []string{";"})
	return vars
}

func (parser *Parser) parseTerm() types.Expr {
	token := parser.peekNextToken()

	if helpers.IsSymbol(token, []string{"("}) {
		return parser.parseParenExpression()
	}

	if token.TokenType == types.IDENTIFIER {
		var expr types.Expr
		nextToken := parser.peekNthToken(1)

		if helpers.IsSymbol(nextToken, []string{"(", "."}) {
			expr = parser.parseSubroutineCall()
		} else if helpers.IsSymbol(nextToken, []string{"["}) {
			expr = parser.parseIndexExpression()
		} else {
			expr = types.Ident{Name: parser.getNextToken().Lexeme}
		}
		return expr
	}

	if helpers.IsLiteralType(token) {
		return parser.parseLiteralExpression()
	}

	return parser.parseUnaryExpression()
}

func (parse *Parser) parseLiteralExpression() types.Literal {
	// GRAMMAR: 'true' | 'false' | 'null' | 'this' | integerConstant | stringConstant
	token := parse.getNextToken()

	if helpers.IsKeyword(token, []string{"true", "false"}) {
		return types.Literal{Type: types.BooleanLiteral, Value: token.Lexeme}
	}

	if helpers.IsKeyword(token, []string{"null"}) {
		return types.Literal{Type: types.ThisLiteral, Value: token.Lexeme}
	}

	if helpers.IsKeyword(token, []string{"this"}) {
		return types.Literal{Type: types.ThisLiteral, Value: token.Lexeme}
	}

	if token.TokenType == types.INTEGER_CONSTANT {
		return types.Literal{Type: types.IntegerLiteral, Value: token.Lexeme}
	}

	parse.assertToken(token, []string{"stringConstant"})
	return types.Literal{Type: types.StringLiteral, Value: token.Lexeme}
}

func (parser *Parser) parseIndexExpression() types.IndexExpr {
	// GRAMMAR: varName '[' expression ']'
	var expr types.IndexExpr
	identToken := parser.getNextToken()
	parser.assertToken(identToken, []string{"varName"})
	parser.assertToken(parser.getNextToken(), []string{"["})
	expr = types.IndexExpr{
		Object:  types.Ident{Name: identToken.Lexeme},
		Indexer: parser.parseExpression(),
	}
	parser.assertToken(parser.getNextToken(), []string{"]"})
	return expr
}

func (parser *Parser) parseParenExpression() types.ParenExpr {
	// GRAMMAR: '(' expression ')'
	var expr types.ParenExpr
	parser.assertToken(parser.getNextToken(), []string{"("})
	expr.Expression = parser.parseExpression()
	parser.assertToken(parser.getNextToken(), []string{")"})
	return expr
}

func (parser *Parser) parseSubroutineCall() types.CallExpr {
	// GRAMMAR: subroutineName '(' expressionList ')' | (className | varName) '.' subroutineName '(' expressionList ')
	var expr types.CallExpr

	token := parser.getNextToken()
	parser.assertToken(token, []string{"className", "subroutineName", "varName"})

	if helpers.IsSymbol(parser.peekNextToken(), []string{"("}) {
		expr.Callee = types.Ident{Name: token.Lexeme}
	} else {
		parser.assertToken(parser.getNextToken(), []string{"."})
		subroutineNameToken := parser.getNextToken()
		parser.assertToken(subroutineNameToken, []string{"subroutineName"})

		expr.Callee = types.MemberExpr{
			Object:   types.Ident{Name: token.Lexeme},
			Property: types.Ident{Name: subroutineNameToken.Lexeme},
		}
	}
	expr.Arguments = parser.parseExpressionList()
	return expr
}

func (parser *Parser) parseUnaryExpression() types.UnaryExpr {
	// GRAMMAR: ('-' | '~') term
	opToken := parser.getNextToken()
	parser.assertToken(opToken, []string{"-", "~"})
	operator, _ := types.GetUnaryOperator(opToken.Lexeme)

	return types.UnaryExpr{
		Operator: operator,
		Operand:  parser.parseExpression(),
	}
}

func (parser *Parser) parseExpression() types.Expr {
	// GRAMMAR: term (op term)*
	term := parser.parseTerm()
	nexToken := parser.peekNextToken()

	if helpers.IsBinaryOperator(nexToken) {
		parser.assertToken(parser.getNextToken(), []string{"+", "-", "*", "/", "<", ">", "="})
		operator, _ := types.GetBinaryOperator(nexToken.Lexeme)

		return types.BinaryExpr{
			Left:     term,
			Operator: operator,
			Right:    parser.parseExpression(),
		}
	}

	if helpers.IsLogicalOperator(nexToken) {
		parser.assertToken(parser.getNextToken(), []string{"&", "|"})
		operator, _ := types.GetLogicalOperator(nexToken.Lexeme)
		return types.LogicalExpr{
			Left:     term,
			Operator: operator,
			Right:    parser.parseExpression(),
		}
	}

	return term
}

func (parser *Parser) parseExpressionList() (args []types.Expr) {
	// GRAMMAR: (expression (',' expression)*)?
	parser.assertToken(parser.getNextToken(), []string{"("})

	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{")"}); nextToken = parser.peekNextToken() {
		args = append(args, parser.parseExpression())
		if helpers.IsSymbol(parser.peekNextToken(), []string{","}) {
			// discard "," token before next expresssion in list
			parser.getNextToken()
		}
	}

	parser.assertToken(parser.getNextToken(), []string{")"})
	return args
}

func (parser *Parser) parseDoStatement() (stmt types.DoStmt) {
	// GRAMMAR: 'do' subroutineName '(' expressionList ')' ';' | 'do' (className | varName) '.' subroutineName '(' expressionList ') ';'
	parser.assertToken(parser.getNextToken(), []string{"do"})
	stmt.Expression = parser.parseSubroutineCall()
	parser.assertToken(parser.getNextToken(), []string{";"})
	return stmt
}

func (parser *Parser) parseIfStatement() (stmt types.IfStmt) {
	// GRAMMAR: 'if' '(' expression ')' '{' statements '}' ('else' '{' statements '}')?
	parser.assertToken(parser.getNextToken(), []string{"if"})
	parser.assertToken(parser.getNextToken(), []string{"("})
	stmt.Condition = parser.parseExpression()
	parser.assertToken(parser.getNextToken(), []string{")"})

	stmt.ThenStmt = parser.parseBlockStatement()

	if helpers.IsKeyword(parser.peekNextToken(), []string{"else"}) {
		// discard 'else' keyword and parse else statement block
		parser.getNextToken()
		stmt.ElseStmt = parser.parseBlockStatement()
	}
	return stmt
}

func (parser *Parser) parseLetStatement() (stmt types.LetStmt) {
	// GRAMMAR: 'let' varName ('[' expression ']')? '=' expression ';'
	parser.assertToken(parser.getNextToken(), []string{"let"})

	// If the token ahead of the next token is a '[' we're dealing with an index expression.
	if helpers.IsSymbol(parser.peekNthToken(1), []string{"["}) {
		stmt.Target = parser.parseIndexExpression()
	} else {
		identToken := parser.getNextToken()
		parser.assertToken(identToken, []string{"varName"})
		stmt.Target = types.Ident{Name: identToken.Lexeme}
	}

	parser.assertToken(parser.getNextToken(), []string{"="})
	stmt.Value = parser.parseExpression()
	parser.assertToken(parser.getNextToken(), []string{";"})

	return stmt
}

func (parser *Parser) parseReturnStatement() (stmt types.ReturnStmt) {
	// GRAMMAR: 'return' expression? ';'
	parser.assertToken(parser.getNextToken(), []string{"return"})

	if !helpers.IsSymbol(parser.peekNextToken(), []string{";"}) {
		stmt.Expression = parser.parseExpression()
	}

	parser.assertToken(parser.getNextToken(), []string{";"})
	return stmt
}

func (parser *Parser) parseWhileStatement() (stmt types.WhileStmt) {
	// GRAMMAR: 'while' '(' expression ')' '{' statements '}'
	parser.assertToken(parser.getNextToken(), []string{"while"})
	parser.assertToken(parser.getNextToken(), []string{"("})
	stmt.Condition = parser.parseExpression()
	parser.assertToken(parser.getNextToken(), []string{")"})
	stmt.Body = parser.parseBlockStatement()

	return stmt
}

func (parser *Parser) parseBlockStatement() (block types.BlockStmt) {
	// GRAMMAR: '{' statements '}'
	parser.assertToken(parser.getNextToken(), []string{"{"})

	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{"}"}); nextToken = parser.peekNextToken() {
		block.Statements = append(block.Statements, parser.parseStatement())
	}

	parser.assertToken(parser.getNextToken(), []string{"}"})
	return block
}

func (parser *Parser) parseStatement() types.Stmt {
	token := parser.peekNextToken()
	var stmt types.Stmt

	switch token.Lexeme {
	case "do":
		stmt = parser.parseDoStatement()
	case "if":
		stmt = parser.parseIfStatement()
	case "let":
		stmt = parser.parseLetStatement()
	case "return":
		stmt = parser.parseReturnStatement()
	case "while":
		stmt = parser.parseWhileStatement()
	default:
		parser.emitError(UNEXPECTED_TOKEN, token)
	}
	return stmt
}

func (parser *Parser) parseSubroutineBody() (body types.SubroutineBody) {
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

func (parser *Parser) parseSubroutineDec() (subroutine types.SubroutineDecl) {
	// GRAMMAR: ('constructor' | 'function' | 'method') ('void' | type) subroutineName '(' parameterList ')' subroutineBody
	subroutineKindToken := parser.getNextToken()
	subroutineTypeToken := parser.getNextToken()
	subroutineNameToken := parser.getNextToken()

	parser.assertToken(subroutineKindToken, []string{"constructor", "function", "method"})
	parser.assertToken(subroutineTypeToken, []string{"boolean", "char", "className", "int", "void"})
	parser.assertToken(subroutineNameToken, []string{"subroutineName"})

	subroutineKind, _ := types.GetSubroutineKind(subroutineKindToken.Lexeme)

	subroutine.Name = subroutineNameToken.Lexeme
	subroutine.Kind = subroutineKind
	subroutine.Type = subroutineTypeToken.Lexeme

	parser.assertToken(parser.getNextToken(), []string{"("})
	subroutine.Params = append(subroutine.Params, parser.parseParameterList()...)
	parser.assertToken(parser.getNextToken(), []string{")"})

	subroutine.Body = parser.parseSubroutineBody()
	return subroutine
}

func (parser *Parser) parseClassVarDec() (vars []types.VarDecl) {
	// GRAMMAR: ('static' | 'field') type varName (',' varName)* ';'
	varKindToken := parser.getNextToken()
	varTypeToken := parser.getNextToken()
	varNameToken := parser.getNextToken()

	parser.assertToken(varKindToken, []string{"field", "static"})
	parser.assertToken(varTypeToken, []string{"boolean", "char", "className", "int"})
	parser.assertToken(varNameToken, []string{"varName"})

	varKind, _ := types.GetVarKind(varKindToken.Lexeme)
	vars = append(vars, types.VarDecl{Name: varNameToken.Lexeme, Type: varTypeToken.Lexeme, Kind: varKind})

	// Check if it's a multi var declaration.
	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{";"}); nextToken = parser.peekNextToken() {
		parser.assertToken(parser.getNextToken(), []string{","})
		parser.assertToken(parser.peekNextToken(), []string{"varName"})

		nextVarNameToken := parser.getNextToken()
		vars = append(vars, types.VarDecl{Name: nextVarNameToken.Lexeme, Type: varTypeToken.Lexeme, Kind: varKind})
	}

	parser.assertToken(parser.getNextToken(), []string{";"})
	return vars
}

func (parser *Parser) Parse(tokens []types.Token) (class types.Class) {
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
		return class
	}

	parser.filename = tokens[0].Filename
	parser.tokens = tokens
	parser.assertToken(parser.getNextToken(), []string{"class"})
	classNameToken := parser.getNextToken()
	parser.assertToken(classNameToken, []string{"className"})
	parser.assertToken(parser.getNextToken(), []string{"{"})

	class.Name = types.Ident{Name: classNameToken.Lexeme}

	for nextToken := parser.peekNextToken(); !helpers.IsSymbol(nextToken, []string{"}"}); nextToken = parser.peekNextToken() {
		if helpers.IsKeyword(nextToken, []string{"field", "static"}) {
			vars := parser.parseClassVarDec()
			class.Vars = append(class.Vars, vars...)
		} else if helpers.IsKeyword(nextToken, []string{"constructor", "function", "method"}) {
			subroutine := parser.parseSubroutineDec()
			class.Subroutines = append(class.Subroutines, subroutine)
		} else {
			parser.emitError(UNEXPECTED_TOKEN, nextToken)
		}
	}

	return class
}
