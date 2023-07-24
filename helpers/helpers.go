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
