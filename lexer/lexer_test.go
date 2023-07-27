package lexer_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"testing"

	. "github.com/MlkMahmud/jack-compiler/constants"
	. "github.com/MlkMahmud/jack-compiler/helpers"
	. "github.com/MlkMahmud/jack-compiler/lexer"
)

const testdataPath = "../testdata/"

func readFileContent(filename string) string {
	file, _ := os.Open(filename)
	bytes, err := io.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}

	fileContent := string(bytes)
	return fileContent
}

func writeTokensToXML(tokens []Token, dest string) {
	file, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString("<tokens>\n")
	for _, token := range tokens {
		file.WriteString(
			fmt.Sprintf(
				"  <%s> %s </%s>\n",
				token.TokenType.ToString(),
				WriteSymbol(token.Lexeme),
				token.TokenType.ToString(),
			),
		)
	}
	file.WriteString("</tokens>")
}

func TestLexer(t *testing.T) {
	files := []string{"Array", "Square", "SquareGame"}
	lexer := NewLexer()

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			filePath := path.Join(testdataPath, fmt.Sprintf("%s%s", file, ".jack"))
			cmpFilePath := path.Join(testdataPath, "expected", fmt.Sprintf("%s%s", file, "T.xml"))
			outFilePath := path.Join(testdataPath, "output", fmt.Sprintf("%s%s", file, "T.xml"))
			tokens := lexer.Tokenize(filePath)

			writeTokensToXML(tokens, outFilePath)

			if readFileContent(cmpFilePath) != readFileContent(outFilePath) {
				t.Errorf("Expected content of %s to match content of %s", outFilePath, cmpFilePath)
			}
		})
	}
}
