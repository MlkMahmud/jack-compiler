package lib

import (
	"path/filepath"
)

func isKeyword(token Token, lexemes []string) bool {
	for _, val := range lexemes {
		if token.tokenType == KEYWORD && token.lexeme == val {
			return true
		}
	}
	return false
}

func isSymbol(token Token, lexemes []string) bool {
	for _, val := range lexemes {
		if token.tokenType == SYMBOL && token.lexeme == val {
			return true
		}
	}
	return false
}

func writeSymbol(lexeme string) string {
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

func replaceFileExt(filename, ext string) string {
	extension := filepath.Ext(filename)
	basename := filename[0 : len(filename)-len(extension)]
	return basename + ext
}
