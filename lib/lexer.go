package lib

import (
	"errors"
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
	colNum int
	lineNum int
	source *os.File
}

type token struct {
	colNum int
	lineNum int
	tokenType  string
	value string
}

func NewLexer() *lexer {
	var lx = lexer{ colNum: 0, lineNum: 1 }
	return &lx
}

func (lx *lexer) appendToken(tokens *[]token, entry token) {
	entry.lineNum = lx.lineNum
	if len(entry.value) > 1 {
		entry.colNum = lx.colNum - len(entry.value)
	} else {
		entry.colNum = lx.colNum
	}
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
		err := lx.source.Close()
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
				var asteriskChar = lx.read()
				// Advance until we hit the "*/" terminator
				var forwardSlashChar = lx.read()

				for fmt.Sprintf("%s%s", asteriskChar, forwardSlashChar) != "*/" {
					if forwardSlashChar == nil {
						panic(errors.New("SyntaxError: invalid multiline comment"))
					}
					asteriskChar = forwardSlashChar
					forwardSlashChar = lx.read()
				}
				char = lx.read()
			} else {
				// This is a division symbol
				lx.appendToken(&tokens, token{
					tokenType: "symbol",
					value: "/",
				})
				char = nextChar
			}
		} else if symbols.Has(char.(string)) {
			lx.appendToken(&tokens, token{
				tokenType: "symbol",
				value: char.(string),
			})
			char = lx.read()
		} else if char == `"` {
			char = lx.read()
			var chars = []string{}

			for {
				if char == `"` || char == nil {
					break
				}
				chars = append(chars, char.(string))
				char = lx.read()
			}
			var word = strings.Join(chars, "")

			if len(word) > 0 {
				lx.appendToken(&tokens, token{
					tokenType: "string",
					value: word,
				})
			}
	
			if char == `"` {
				char = lx.read()
			}
		} else if regexp.MustCompile(`(?i)[_a-z]`).MatchString(char.(string)) {
			var chars = []string{char.(string)}
			char = lx.read().(string)

			for regexp.MustCompile(`\w`).MatchString(char.(string)) {
				chars = append(chars, char.(string))
				char = lx.read().(string)
			}

			var word = strings.Join(chars, "")
			var token = token{ value: word }
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
				value: word,
			})
		} else {
			char = lx.read()
		}
	}
	return tokens
}
