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

var IDENTIFIER = "identifier"

var KEYWORDS = map[string]bool{
	"class": true, "constructor": true, "method": true,
	"function": true, "int": true, "boolean": true,
	"char": true, "void": true, "var": true, "static": true,
	"field": true, "let": true, "do": true, "if": true,
	"else": true, "while": true, "return": true, "true": true,
	"false": true, "null": true, "this": true,
}

var SYMBOLS = map[string]bool{
	"(": true, ")": true, "{": true, "}": true, "[": true, "]": true, "/": true,
	"-": true, "+": true, "*": true, ",": true, ".": true, "=": true, ";": true,
	"&": true, "|": true, "<": true, ">": true, "~": true,
}

type Lexer struct {
	colNum  int
	lineNum int
	source  *os.File
	tokens  *list.List
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
	buffer := make([]byte, 1)
	bytes, err := lexer.source.Read(buffer)

	if err != nil {
		if err == io.EOF {
			return nil
		}
		panic(err)
	}

	char := string(buffer[:bytes])

	if char == "\n" {
		lexer.colNum = 0
		lexer.lineNum++
	} else {
		lexer.colNum++
	}

	return char
}

func (lexer *Lexer) Tokenize(src string) *list.List {
	defer func() {
		err := lexer.source.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	file, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}

	lexer.colNum = 0
	lexer.lineNum = 1
	lexer.source = file
	lexer.tokens = list.New()

	char := lexer.read()

	for char != nil {
		if char == "/" {
			nextChar := lexer.read()
			if nextChar == "/" {
				// This is a single line comment
				// Advance until we hit the next newline char or EOF
				newlineChar := lexer.read()

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
				startColNum := lexer.colNum - 1
				startLineNum := lexer.lineNum

				asteriskChar := lexer.read()
				forwardSlashChar := lexer.read()

				// Save all characters until the end of the comment. (Used only for error logging).
				chars := []string{"/", "*", asteriskChar.(string)}

				for fmt.Sprintf("%s%s", asteriskChar, forwardSlashChar) != "*/" {
					if forwardSlashChar == nil {
						message := fmt.Sprintf(
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
					tokenType: "symbol",
					lexeme:    "/",
				})
				char = nextChar
			}
		} else if SYMBOLS[char.(string)] {
			lexer.appendToken(Token{
				tokenType: "symbol",
				lexeme:    char.(string),
			})
			char = lexer.read()
		} else if char == `"` {
			char = lexer.read()
			chars := []string{}

			startColNum := lexer.colNum
			startLineNum := lexer.lineNum

			for {
				if char == "\n" || char == nil {
					errMessage := fmt.Sprintf(
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

			word := strings.Join(chars, "")

			lexer.appendToken(Token{
				tokenType: "stringConstant",
				lexeme:    fmt.Sprintf(`"%s"`, word),
			})

			char = lexer.read()
		} else if regexp.MustCompile(`(?i)[_a-z]`).MatchString(char.(string)) {
			chars := []string{char.(string)}
			char = lexer.read()

			for regexp.MustCompile(`\w`).MatchString(char.(string)) {
				chars = append(chars, char.(string))
				char = lexer.read()
			}

			word := strings.Join(chars, "")
			token := Token{lexeme: word}
			if KEYWORDS[word] {
				token.tokenType = "keyword"
			} else {
				token.tokenType = IDENTIFIER
			}
			lexer.appendToken(token)
		} else if regexp.MustCompile(`\d`).MatchString(char.(string)) {
			chars := []string{char.(string)}
			char = lexer.read()
			for regexp.MustCompile(`\d`).MatchString(char.(string)) {
				chars = append(chars, char.(string))
				char = lexer.read()
			}
			word := strings.Join(chars, "")
			lexer.appendToken(Token{
				tokenType: "integerConstant",
				lexeme:    word,
			})
		} else {
			char = lexer.read()
		}
	}
	return lexer.tokens
}
