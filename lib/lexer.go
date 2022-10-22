package lib

import (
	"fmt"
	"io"
	"jack-compiler/utils"
	"log"
	"os"
	"regexp"
	"strings"
)

var keywords = utils.Set([]string{
	"class", "constructor", "method", "function", "int",
	"boolean", "char", "void", "var", "static", "field",
	"let", "do", "if", "else", "while", "return", "true",
	"false", "null", "this",
})

var symbols = utils.Set([]string{
	"(", ")", "{", "}", "[", "]", "/", "-", "+", "*",
	",", ".", "=", ";", "&", "|", "<", ">", "~",
})

type lexer struct {
	colNum  int
	lineNum int
	source  *os.File
}

type token struct {
	colNum    int
	lineNum   int
	tokenType string
	value     string
}

func NewLexer() *lexer {
	return new(lexer)
}

func (lx *lexer) appendToken(tokens *[]token, entry token) {
	if len(entry.value) > 1 {
		// Set the current token's colNum to the position of its first character.
		entry.colNum = lx.colNum - len(entry.value)
	} else {
		entry.colNum = lx.colNum
	}
	entry.lineNum = lx.lineNum
	*tokens = append(*tokens, entry)
}

func (lx *lexer) read() interface{} {
	var buffer = make([]byte, 1)
	var bytes, err = lx.source.Read(buffer)

	if err != nil {
		if err == io.EOF {
			return nil
		}
		panic(err)
	}

	var char = string(buffer[:bytes])

	if char == "\n" {
		lx.colNum = 0
		lx.lineNum++
	} else {
		lx.colNum++
	}

	return char
}

func (lx *lexer) Tokenize(src string) []token {
	defer func() {
		var err = lx.source.Close()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}()

	var file, err = os.Open(src)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	lx.source = file
	lx.colNum = 0
	lx.lineNum = 1

	var tokens = make([]token, 0)
	var char = lx.read()

	for char != nil {
		if char == "/" {
			var nextChar = lx.read()
			if nextChar == "/" {
				// This is a single line comment
				// Advance until we hit the next newline char or EOF

				var newlineChar = lx.read()
				for {
					if newlineChar == "\n" || newlineChar == nil {
						break
					}
					newlineChar = lx.read()
				}
				char = newlineChar
			} else if nextChar == "*" {
				// This is a multiline comment
				// Advance until we hit the "*/" terminator

				/* 
					Save a reference to the starting column and line number of this comment. 
					Deduct one from the "colNum" to account for the read operation on "nextChar".
				*/
				var startColNum = lx.colNum - 1
				var startLineNum = lx.lineNum

				var asteriskChar = lx.read()
				var forwardSlashChar = lx.read()

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
					forwardSlashChar = lx.read()
				}
				char = lx.read()
			} else {
				// This is a division symbol
				lx.appendToken(&tokens, token{
					tokenType: "symbol",
					value:     "/",
				})
				char = nextChar
			}
		} else if symbols.Has(char.(string)) {
			lx.appendToken(&tokens, token{
				tokenType: "symbol",
				value:     char.(string),
			})
			char = lx.read()
		} else if char == `"` {
			char = lx.read()
			var chars = []string{}

			var startColNum = lx.colNum
			var startLineNum = lx.lineNum

			for {
				if char == "\n" || char == "nil" {
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
				char = lx.read()
			}

			var word = strings.Join(chars, "")
	
			lx.appendToken(&tokens, token{
				tokenType: "string",
				value:     fmt.Sprintf(`"%s"`, word),
			})

			char = lx.read()
		} else if regexp.MustCompile(`(?i)[_a-z]`).MatchString(char.(string)) {
			var chars = []string{char.(string)}
			char = lx.read()

			for regexp.MustCompile(`\w`).MatchString(char.(string)) {
				chars = append(chars, char.(string))
				char = lx.read()
			}

			var word = strings.Join(chars, "")
			var token = token{value: word}
			if keywords.Has(word) {
				token.tokenType = "keyword"
			} else {
				token.tokenType = "identifier"
			}
			lx.appendToken(&tokens, token)
		} else if regexp.MustCompile(`\d`).MatchString(char.(string)) {
			var chars = []string{char.(string)}
			char = lx.read()
			for regexp.MustCompile(`\d`).MatchString(char.(string)) {
				chars = append(chars, char.(string))
				char = lx.read()
			}
			var word = strings.Join(chars, "")
			lx.appendToken(&tokens, token{
				tokenType: "integer",
				value:     word,
			})
		} else {
			char = lx.read()
		}
	}
	return tokens
}
