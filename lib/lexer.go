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

type Lexer struct {
	source *os.File
}

type Token struct {
	Type  string
	Value string
}

func NewLexer() *Lexer {
	return new(Lexer)
}


/*
This method returns a single byte (as a string char) from the Lexer's source file.
If the source file is at EOF it returns nil
*/
func (lexer *Lexer) Read() interface{} {
	var buffer = make([]byte, 1)
	var bytes, err = lexer.source.Read(buffer)

	if err != nil {
		if err == io.EOF {
			return nil
		}
		log.Fatal(err)
		os.Exit(1)
	}
	return string(buffer[:bytes])
}

func (lexer *Lexer) Tokenize(src string) []Token {
	defer func() {
		err := lexer.source.Close()
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

	lexer.source = file
	var tokens = make([]Token, 0)
	var char = lexer.Read()

	for char != nil {
		if char == "/" {
			var nextChar = lexer.Read()
			if nextChar == "/" {
				// This is a single line comment
				// Advance until we hit the next newline char or EOF
				var newlineChar = lexer.Read()
				for {
					if newlineChar == "\n" || newlineChar == nil {
						break
					}
					newlineChar = lexer.Read()
				}
			} else if nextChar == "*" {
				// This is a multiline comment
				// Advance until we hit the "*/" terminator
				var asteriskChar = lexer.Read()
				var forwardSlashChar = lexer.Read()

				for fmt.Sprintf("%s%s", asteriskChar, forwardSlashChar) != "*/" {
					if forwardSlashChar == nil {
						log.Fatal(errors.New("SyntaxError: invalid multiline comment"))
						os.Exit(1)
					}
					asteriskChar = forwardSlashChar
					forwardSlashChar = lexer.Read()
				}
			} else {
				// This is a division symbol
				tokens = append(tokens, Token{Type: "symbol", Value: "/"})
				char = nextChar
			}
		} else if symbols.Has(char.(string)) {
			tokens = append(tokens, Token{Type: "symbol", Value: char.(string)})
			char = lexer.Read()
		} else if char == `"` {
			char = lexer.Read()
			var chars = []string{}

			for {
				if char == `"` || char == nil {
					break
				}
				chars = append(chars, char.(string))
				char = lexer.Read()
			}
			var word = strings.Join(chars, "")

			if len(word) > 0 {
				tokens = append(tokens, Token{Type: "String", Value: word})
			}
	
			if char == `"` {
				char = lexer.Read()
			}
		} else if regexp.MustCompile(`(?i)[_a-z]`).MatchString(char.(string)) {
			var chars = []string{char.(string)}
			char = lexer.Read().(string)

			for regexp.MustCompile(`\w`).MatchString(char.(string)) {
				chars = append(chars, char.(string))
				char = lexer.Read().(string)
			}

			var word = strings.Join(chars, "")
			if keywords.Has(word) {
				tokens = append(tokens, Token{Type: "Keyword", Value: word})
			} else {
				tokens = append(tokens, Token{Type: "Identifier", Value: word})
			}
		} else if regexp.MustCompile(`\d`).MatchString(char.(string)) {
			var chars = []string{char.(string)}
			char = lexer.Read()
			for regexp.MustCompile(`\d`).MatchString(char.(string)) {
				chars = append(chars, char.(string))
				char = lexer.Read()
			}
			var word = strings.Join(chars, "")
			tokens = append(tokens, Token{Type: "Integer", Value: word})
		} else {
			char = lexer.Read()
		}
	}
	return tokens
}
