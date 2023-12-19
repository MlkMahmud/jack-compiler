package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/MlkMahmud/jack-compiler/lexer"
	"github.com/MlkMahmud/jack-compiler/parser"
)

func printHelpMessage() {
	log.SetFlags(0)
	log.Fatalln(("usage:\n go run main.go --src .\t\t\tCompiles all the .jack files in the current directory\n go run main.go --src <fileName.jack>\tCompiles the specified .jack file\n go run main.go --src <dirName>\t\tCompiles all the .jack files in the specified directory"))
}

func main() {
	var source string
	flag.StringVar(&source, "src", "", "Path to a '.jack' file or a directory containing one or more '.jack' files.")
	flag.Parse()

	info, err := os.Stat(source)

	if err != nil {
		log.Fatal(err)
	}

	lexer := lexer.NewLexer()
	parser := parser.NewParser()

	filePaths := []string{}

	if info.IsDir() {
		entries, err := os.ReadDir(source)

		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range entries {
			if fileName := entry.Name(); strings.HasSuffix(fileName, ".jack") {
				filePaths = append(filePaths, filepath.Join(source, fileName))
			}
		}
	} else {
		if !strings.HasSuffix(source, ".jack") {
			printHelpMessage()
		}
		filePaths = append(filePaths, source)
	}

	for _, src := range filePaths {
		tokens := lexer.Tokenize(src)
		class := parser.Parse(tokens)
		fmt.Printf("ClassName: %s\nVar Count: %d\nSubroutine Count: %d\n", class.Name, len(class.Vars), len(class.Subroutines))
	}
}
