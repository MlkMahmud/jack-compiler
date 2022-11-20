package lib

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