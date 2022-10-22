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
		log.Fatal(
			errors.New("usage: go run main.go [ src ]\n src: a .jack file or a directory with 1 or more .jack files"),
		)
	}

	var src = os.Args[1]
	var Lexer = lib.NewLexer()
	var tokens = Lexer.Tokenize(src)
	for i := 0; i < len(tokens); i++ {
		fmt.Printf("%+v\n", tokens[i])
	}
}
