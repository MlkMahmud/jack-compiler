package parser_test

import (
	"encoding/json"
	"fmt"
	"path"
	"testing"

	. "github.com/MlkMahmud/jack-compiler/lexer"
	. "github.com/MlkMahmud/jack-compiler/parser"
)

const TEST_DATA_PATH = "../testdata"

func TestParser(t *testing.T) {
	lexer := NewLexer()
	filePath := path.Join(TEST_DATA_PATH, "SquareGame.jack")
	tokens := lexer.Tokenize(filePath)

	parser := NewParser()
	class := parser.Parse(tokens)

	jsonData, err := json.Marshal(class)
	if err != nil {
		fmt.Println(err)
		t.Fatal()
	}

	fmt.Print(string(jsonData))
}
