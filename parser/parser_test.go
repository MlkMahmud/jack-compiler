package parser_test

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"testing"

	. "github.com/MlkMahmud/jack-compiler/lexer"
	. "github.com/MlkMahmud/jack-compiler/parser"
	. "github.com/nsf/jsondiff"
)

const TEST_DATA_PATH = "../testdata"

func readFileContent(filename string) []byte {
	file, _ := os.Open(filename)
	bytes, err := io.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func TestParser(t *testing.T) {
	files := []string{"Array", "Square", "SquareGame"}
	lexer := NewLexer()
	parser := NewParser()

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			filePath := path.Join(TEST_DATA_PATH, strings.Join([]string{file, "jack"}, "."))
			cmpFilePath := path.Join(TEST_DATA_PATH, "expected", strings.Join([]string{file, "json"}, "."))
			
			tokens := lexer.Tokenize(filePath)
			class := parser.Parse(tokens)
			
			expected := readFileContent(cmpFilePath)

			actual, err := json.Marshal(class)

			if err != nil {
				t.Fatal(err)
			}
			
			difference, desc := Compare(expected, actual, &Options{ SkipMatches: true })

			if difference != FullMatch {
				t.Fatal(desc)
			}
		})
	}
}
