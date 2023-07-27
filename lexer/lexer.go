package lexer

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/MlkMahmud/jack-compiler/constants"
)



type Lexer struct {
	colNum  int
	lineNum int
	source  *os.File
}


func NewLexer() *Lexer {
	return new(Lexer)
}

func (lexer *Lexer) appendToken(tokens *[]constants.Token, entry constants.Token) {
	if len(entry.Lexeme) > 1 {
		// Set the current token's colNum to the position of its first character.
		entry.ColNum = lexer.colNum - len(entry.Lexeme)
	} else {
		entry.ColNum = lexer.colNum
	}
	entry.LineNum = lexer.lineNum
	*tokens = append(*tokens, entry)
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

func (lexer *Lexer) Tokenize(src string) []constants.Token {
	defer func() {
		err := lexer.source.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	tokens := make([]constants.Token, 0)
	file, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}

	lexer.colNum = 0
	lexer.lineNum = 1
	lexer.source = file
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
				lexer.appendToken(&tokens, constants.Token{
					TokenType: constants.SYMBOL,
					Lexeme:    "/",
				})
				char = nextChar
			}
		} else if constants.SYMBOLS[char.(string)] {
			lexer.appendToken(&tokens, constants.Token{
				TokenType: constants.SYMBOL,
				Lexeme:    char.(string),
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

			lexer.appendToken(&tokens, constants.Token{
				TokenType: constants.STRING_CONSTANT,
				Lexeme:    word,
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
			token := constants.Token{Lexeme: word}
			if constants.KEYWORDS[word] {
				token.TokenType = constants.KEYWORD
			} else {
				token.TokenType = constants.IDENTIFIER
			}
			lexer.appendToken(&tokens, token)
		} else if regexp.MustCompile(`\d`).MatchString(char.(string)) {
			chars := []string{char.(string)}
			char = lexer.read()
			for regexp.MustCompile(`\d`).MatchString(char.(string)) {
				chars = append(chars, char.(string))
				char = lexer.read()
			}
			word := strings.Join(chars, "")
			lexer.appendToken(&tokens, constants.Token{
				TokenType: constants.INTEGER_CONSTANT,
				Lexeme:    word,
			})
		} else {
			char = lexer.read()
		}
	}
	return tokens
}
