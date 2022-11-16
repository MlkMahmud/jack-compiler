package main

import (
	"errors"
	"fmt"
	"jack-compiler/lib"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.SetFlags(0)
		log.Fatalln(
			errors.New("usage:\n go run main.go\t\t\t\tCompiles all the .jack files in the current directory\n go run main.go <src.jack>\t\tCompiles the specified .jack file\n go run main.go src/\t\t\tCompiles all the .jack files in the specified directory"),
		)
	}

	var src = os.Args[1]
	var Lexer = lib.NewLexer()
	var tokens = Lexer.Tokenize(src)
	var current = tokens.Front()
	for current != nil {
		var token = current.Value
		current = current.Next()
		fmt.Printf("%+v\n", token)
	}
}
