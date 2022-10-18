package main

import (
	"fmt"
	"jack-compiler/lib"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("Add useful message: ")
	}

	var src = os.Args[1]
	var Lexer = lib.NewLexer()
	var tokens = Lexer.Tokenize(src)
	for i := 0; i < len(tokens); i++ {
		fmt.Printf("%+v\n", tokens[i])
	}
}
