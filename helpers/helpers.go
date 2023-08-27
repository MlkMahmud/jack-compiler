package helpers

import (
	"github.com/MlkMahmud/jack-compiler/types"
)

func IsOneOfKeywords(token types.Token, lexemes []string) bool {
	for _, val := range lexemes {
		if token.TokenType == types.KEYWORD && token.Lexeme == val {
			return true
		}
	}
	return false
}

func IsOneOfSymbols(token types.Token, lexemes []string) bool {
	for _, val := range lexemes {
		if token.TokenType == types.SYMBOL && token.Lexeme == val {
			return true
		}
	}
	return false
}

func IsBinaryOperator(token types.Token) bool {
	return token.TokenType == types.SYMBOL && Contains([]string{"+", "-", "*", "/", "<", ">", "="}, token.Lexeme)
}

func IsLogicalOperator(token types.Token) bool {
	return token.TokenType == types.SYMBOL && Contains([]string{"&", "|"}, token.Lexeme)
}

func IsLiteralType(token types.Token) bool {
	return token.TokenType == types.INTEGER_CONSTANT ||
		token.TokenType == types.STRING_CONSTANT ||
		Contains([]string{"true", "false", "null", "this"}, token.Lexeme)
}

func Contains[T comparable](arr []T, elem T) bool {
	for _, value := range arr {
		if value == elem {
			return true
		}
	}
	return false
}
