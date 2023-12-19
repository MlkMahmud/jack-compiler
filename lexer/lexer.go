package lexer

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/MlkMahmud/jack-compiler/types"
)

type LexerError struct {
	message string
}

func (e LexerError) Error() string {
	return e.message
}

type Lexer struct {
	colNum  int
	lineNum int
	source  *os.File
}

func NewLexer() *Lexer {
	return new(Lexer)
}

func (lexer *Lexer) emitError(col, line int, message string) {
	panic(&LexerError{message: fmt.Sprintf(
		"<%s:%d:%d>\tError: %s",
		lexer.source.Name(),
		line,
		col,
		message,
	)})
}

func (lexer *Lexer) appendToken(tokens *[]types.Token, entry types.Token) {
	if len(entry.Lexeme) > 1 {
		// Set the current token's colNum to the position of its first character.
		entry.ColNum = lexer.colNum - len(entry.Lexeme)
	} else {
		entry.ColNum = lexer.colNum
	}
	entry.Filename = lexer.source.Name()
	entry.LineNum = lexer.lineNum
	*tokens = append(*tokens, entry)
}

func (lexer *Lexer) read() string {
	buffer := make([]byte, 1)
	bytes, err := lexer.source.Read(buffer)
	if err != nil {
		if err == io.EOF {
			return "\000"
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

func (lexer *Lexer) Tokenize(src string) []types.Token {
	defer func() {
		err := lexer.source.Close()
		if err != nil {
			log.Fatal(err)
		}

		if r := recover(); r != nil {
			var lexerError *LexerError
			if errors.As(r.(error), &lexerError) {
				fmt.Println(r)
			} else {
				debug.PrintStack()
			}
			os.Exit(1)
		}
	}()

	tokens := make([]types.Token, 0)
	file, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}

	lexer.colNum = 0
	lexer.lineNum = 1
	lexer.source = file
	char := lexer.read()

	for char != "\000" {
		if char == "/" {
			nextChar := lexer.read()
			if nextChar == "/" {
				// This is a single line comment
				// Advance until we hit the next newline char or EOF
				newlineChar := lexer.read()

				for {
					if newlineChar == "\n" || newlineChar == "\000" {
						break
					}
					newlineChar = lexer.read()
				}

				char = newlineChar
			} else if nextChar == "*" {
				// This is a multiline comment
				// Advance until we hit the "*/" terminator
				// Deduct one from the "colNum" to account for the read operation on "nextChar".
				startCol := lexer.colNum - 1
				startLine := lexer.lineNum
				asteriskChar := lexer.read()
				forwardSlashChar := lexer.read()

				for fmt.Sprintf("%s%s", asteriskChar, forwardSlashChar) != "*/" {
					if forwardSlashChar == "\000" {
						lexer.emitError(startCol, startLine, "Unterminated multiline comment.")	
					}
					asteriskChar = forwardSlashChar
					forwardSlashChar = lexer.read()
				}
				char = lexer.read()
			} else {
				// This is a division symbol
				lexer.appendToken(&tokens, types.Token{
					TokenType: types.SYMBOL,
					Lexeme:    "/",
				})
				char = nextChar
			}
		} else if types.SYMBOLS[char] {
			lexer.appendToken(&tokens, types.Token{
				TokenType: types.SYMBOL,
				Lexeme:    char,
			})
			char = lexer.read()
		} else if char == `"` {
			char = lexer.read()
			chars := []string{}

			startCol := lexer.colNum
			startLine := lexer.lineNum

			for {
				if char == "\n" || char == "\000" {
					lexer.emitError(startCol, startLine, "Unterminated string literal.")
				}

				if char == `"` {
					break
				}

				chars = append(chars, char)
				char = lexer.read()
			}

			word := strings.Join(chars, "")

			lexer.appendToken(&tokens, types.Token{
				TokenType: types.STRING_CONSTANT,
				Lexeme:    word,
			})

			char = lexer.read()
		} else if regexp.MustCompile(`(?i)[_a-z]`).MatchString(char) {
			chars := []string{char}
			char = lexer.read()

			for regexp.MustCompile(`\w`).MatchString(char) {
				chars = append(chars, char)
				char = lexer.read()
			}

			word := strings.Join(chars, "")
			token := types.Token{Lexeme: word}
			if types.KEYWORDS[word] {
				token.TokenType = types.KEYWORD
			} else {
				token.TokenType = types.IDENTIFIER
			}
			lexer.appendToken(&tokens, token)
		} else if regexp.MustCompile(`\d`).MatchString(char) {
			chars := []string{char}
			char = lexer.read()
			for regexp.MustCompile(`\d`).MatchString(char) {
				chars = append(chars, char)
				char = lexer.read()
			}
			word := strings.Join(chars, "")
			lexer.appendToken(&tokens, types.Token{
				TokenType: types.INTEGER_CONSTANT,
				Lexeme:    word,
			})
		} else {
			char = lexer.read()
		}
	}
	return tokens
}
