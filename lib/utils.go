package lib

import (
	"fmt"
	"path/filepath"
)

type CompilerError struct {
	errorMessage string
}

func (e *CompilerError) Error() string {
	return e.errorMessage
}

func isValidType(token Token) bool {
	if token.tokenType == IDENTIFIER {
		return true
	}

	for _, val := range []string{"boolean", "char", "int"} {
		if token.tokenType == KEYWORD && token.lexeme == val {
			return true
		}
	}
	return false
}

func isClassVarDec(token Token) bool {
	for _, val := range []string{"field", "static"} {
		if token.lexeme == val && token.tokenType == KEYWORD {
			return true
		}
	}
	return false
}

func isSubroutineDec(token Token) bool {
	for _, val := range []string{"constructor", "function", "method"} {
		if token.lexeme == val && token.tokenType == KEYWORD {
			return true
		}
	}
	return false
}

func isIdentifier(token Token) bool {
	return token.tokenType == IDENTIFIER
}

func isKeyword(token Token, lexeme string) bool {
	return token.tokenType == KEYWORD && token.lexeme == lexeme
}

func isSymbol(token Token, lexeme string) bool {
	return token.tokenType == SYMBOL && token.lexeme == lexeme
}

func isStatementDec(token Token) bool {
	for _, lexeme := range []string{"do", "if", "let", "return", "while"} {
		if token.lexeme == lexeme && token.tokenType == KEYWORD {
			return true
		}
	}
	return false
}

func isKeywordConstant(token Token) bool {
	for _, val := range []string{"false", "null", "this", "true"} {
		if isKeyword(token, val) {
			return true
		}
	}
	return false
}

func isOperator(token Token) bool {
	for _, val := range []string{"+", "-", "*", "/", "&", "|", "<", ">", "="} {
		if isSymbol(token, val) {
			return true
		}
	}
	return false
}

func isUnaryOperator(token Token) bool {
	return isSymbol(token, "-") || isSymbol(token, "~")
}

func printErr(token Token) {
	panic(fmt.Sprintf("error: invalid token: %s", token.lexeme))
}

func writeSymbol(token Token) string {
	switch token.lexeme {
	case "<":
		return "&lt;"
	case ">":
		return "&gt;"
	case "&":
		return "&amp;"
	default:
		return token.lexeme
	}
}

func isEndOfExpression(token Token) bool {
	return isSymbol(token, "]") || isSymbol(token, ")") || isSymbol(token, ";")
}

func updateFileExt(filename, ext string) string {
	extension := filepath.Ext(filename)
	basename := filename[0: len(filename) - len(extension)]
	return basename + ext
}