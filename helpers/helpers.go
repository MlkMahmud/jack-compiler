package helpers

import (
	"path/filepath"

	"github.com/MlkMahmud/jack-compiler/constants"
)

func IsKeyword(token constants.Token, lexemes []string) bool {
	for _, val := range lexemes {
		if token.TokenType == constants.KEYWORD && token.Lexeme == val {
			return true
		}
	}
	return false
}

func IsSymbol(token constants.Token, lexemes []string) bool {
	for _, val := range lexemes {
		if token.TokenType == constants.SYMBOL && token.Lexeme == val {
			return true
		}
	}
	return false
}

func IsBinaryOperator(token constants.Token) bool {
	return token.TokenType == constants.SYMBOL && Contains([]string{"+", "-", "*", "/", "<", ">"}, token.Lexeme)
}

func IsLogicalOperator(token constants.Token) bool {
	return token.TokenType == constants.SYMBOL && Contains([]string{"&", "|"}, token.Lexeme)
}

func IsLiteralType(token constants.Token) bool {
	return token.TokenType == constants.INTEGER_CONSTANT ||
		token.TokenType == constants.STRING_CONSTANT ||
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
