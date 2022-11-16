package lib

import (
	"container/list"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var IDENTIFIER = "IDENTIFIER"

var KEYWORDS = map[string]string{
	"class": "KEYWORD_CLASS", "constructor": "KEYWORD_CONSTRUCTOR", "method": "KEYWORD_METHOD",
	"function": "KEYWORD_FUNCTION", "int": "KEYWORD_INT", "boolean": "KEYWORD_BOOLEAN",
	"char": "KEYWORD_CHAR", "void": "KEYWORD_VOID", "var": "KEYWORD_VAR", "static": "KEYWORD_STATIC",
	"field": "KEYWORD_FIELD", "let": "KEYWORD_LET", "do": "KEYWORD_DO", "if": "KEYWORD_IF",
	"else": "KEYWORD_ELSE", "while": "KEYWORD_WHILE", "return": "KEYWORD_RETURN", "true": "KEYWORD_TRUE",
	"false": "KEYWORD_FALSE", "null": "KEYWORD_NULL", "this": "KEYWORD_THIS",
}

var SYMBOLS = map[string]string{
	"(": "LPAREN", ")": "RPAREN", "{": "LBRACE", "}": "RBRACE", "[": "LBRACKET", "]": "RBRACKET", "/": "DIV",
	"-": "MINUS", "+": "PLUS", "*": "MUL", ",": "COMMA", ".": "PERIOD", "=": "EQUALS", ";": "SEMICOLON",
	"&": "AND", "|": "OR", "<": "LESS", ">": "GREATER", "~": "TILDE",
}


type Lexer struct {
	colNum  int
	lineNum int
	source  *os.File
	tokens 	*list.List
}

type Token struct {
	tokenType string
	lexeme    string
	colNum    int
	lineNum   int
}

func NewLexer() *Lexer {
	return new(Lexer)
}

func (lexer *Lexer) appendToken(entry Token) {
	if len(entry.lexeme) > 1 {
		// Set the current token's colNum to the position of its first character.
		entry.colNum = lexer.colNum - len(entry.lexeme)
	} else {
		entry.colNum = lexer.colNum
	}
	entry.lineNum = lexer.lineNum
	lexer.tokens.PushBack(entry)
}

func (lexer *Lexer) read() interface{} {
	var buffer = make([]byte, 1)
	var bytes, err = lexer.source.Read(buffer)

	if err != nil {
		if err == io.EOF {
			return nil
		}
		panic(err)
	}

	var char = string(buffer[:bytes])

	if char == "\n" {
		lexer.colNum = 0
		lexer.lineNum++
	} else {
		lexer.colNum++
	}

	return char
}

func (lexer *Lexer) Tokenize(src string) list.List {
	defer func() {
		var err = lexer.source.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var file, err = os.Open(src)
	if err != nil {
		log.Fatal(err)
	}

	lexer.colNum = 0
	lexer.lineNum = 1
	lexer.source = file
	lexer.tokens = list.New()

	var char = lexer.read()

	for char != nil {
		if char == "/" {
			var nextChar = lexer.read()
			if nextChar == "/" {
				// This is a single line comment
				// Advance until we hit the next newline char or EOF
				var newlineChar = lexer.read()

				for {
					if newlineChar == "\n" || newlineChar == nil {
						break
					}
					newlineChar = lexer.read()
				}

				char = newlineChar
			} else if nextChar == "*" {
				// This is a multiline comment
				// Advance until we hit the "*/" terminator

				/*
				  Save a reference to the starting column and line number of this comment.
					Deduct one from the "colNum" to account for the read operation on "nextChar".
				*/
				var startColNum = lexer.colNum - 1
				var startLineNum = lexer.lineNum

				var asteriskChar = lexer.read()
				var forwardSlashChar = lexer.read()

				// Save all characters until the end of the comment. (Used only for error logging).
				var chars = []string{"/", "*", asteriskChar.(string)}

				for fmt.Sprintf("%s%s", asteriskChar, forwardSlashChar) != "*/" {
					if forwardSlashChar == nil {
						var message = fmt.Sprintf(
							"%s:%s:%s\n\n%s\nSyntaxError: invalid or unexpected token",
							src,
							fmt.Sprint(startLineNum),
							fmt.Sprint(startColNum),
							strings.Join(chars, ""),
						)
						panic(message)
					}
					chars = append(chars, forwardSlashChar.(string))
					asteriskChar = forwardSlashChar
					forwardSlashChar = lexer.read()
				}
				char = lexer.read()
			} else {
				// This is a division symbol
				lexer.appendToken(Token{
					tokenType: SYMBOLS["/"],
					lexeme:    "/",
				})
				char = nextChar
			}
		} else if SYMBOLS[char.(string)] != "" {
			lexer.appendToken(Token{
				tokenType: SYMBOLS[char.(string)],
				lexeme:    char.(string),
			})
			char = lexer.read()
		} else if char == `"` {
			char = lexer.read()
			var chars = []string{}

			var startColNum = lexer.colNum
			var startLineNum = lexer.lineNum

			for {
				if char == "\n" || char == nil {
					var errMessage = fmt.Sprintf(
						"%s:%s:%s\n\n%s\n\nSyntaxError: invalid or unexpected token",
						src,
						fmt.Sprint(startLineNum),
						fmt.Sprint(startColNum),
						fmt.Sprintf(`"%s`, strings.Join(chars, "")),
					)
					panic(errMessage)
				}

				if char == `"` {
					break
				}

				chars = append(chars, char.(string))
				char = lexer.read()
			}

			var word = strings.Join(chars, "")

			lexer.appendToken(Token{
				tokenType: "STRING",
				lexeme:    fmt.Sprintf(`"%s"`, word),
			})

			char = lexer.read()
		} else if regexp.MustCompile(`(?i)[_a-z]`).MatchString(char.(string)) {
			var chars = []string{char.(string)}
			char = lexer.read()

			for regexp.MustCompile(`\w`).MatchString(char.(string)) {
				chars = append(chars, char.(string))
				char = lexer.read()
			}

			var word = strings.Join(chars, "")
			var token = Token{lexeme: word}
			if KEYWORDS[word] != "" {
				token.tokenType = KEYWORDS[word]
			} else {
				token.tokenType = IDENTIFIER
			}
			lexer.appendToken(token)
		} else if regexp.MustCompile(`\d`).MatchString(char.(string)) {
			var chars = []string{char.(string)}
			char = lexer.read()
			for regexp.MustCompile(`\d`).MatchString(char.(string)) {
				chars = append(chars, char.(string))
				char = lexer.read()
			}
			var word = strings.Join(chars, "")
			lexer.appendToken(Token{
				tokenType: "NUMBER",
				lexeme:    word,
			})
		} else {
			char = lexer.read()
		}
	}
	return *lexer.tokens
}
