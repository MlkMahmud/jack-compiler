package helpers

import (
	"path/filepath"

	"github.com/MlkMahmud/jack-compiler/types"
)

func IsKeyword(token types.Token, lexemes []string) bool {
	for _, val := range lexemes {
		if token.TokenType == types.KEYWORD && token.Lexeme == val {
			return true
		}
	}
	return false
}

func IsSymbol(token types.Token, lexemes []string) bool {
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

func WriteSymbol(lexeme string) string {
	switch lexeme {
	case "<":
		return "&lt;"
	case ">":
		return "&gt;"
	case "&":
		return "&amp;"
	default:
		return lexeme
	}
}

func ReplaceFileExt(filename, ext string) string {
	extension := filepath.Ext(filename)
	basename := filename[0 : len(filename)-len(extension)]
	return basename + ext
}
