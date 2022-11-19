package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"jack-compiler/lib"
)

func printHelpMessage() {
	log.SetFlags(0)
	log.Fatalln(
		errors.New("usage:\n go run main.go .\t\t\tCompiles all the .jack files in the current directory\n go run main.go <src.jack>\t\tCompiles the specified .jack file\n go run main.go src/\t\t\tCompiles all the .jack files in the specified directory"),
	)
}

func main() {
	if len(os.Args) != 2 {
		printHelpMessage()
	}

	src := os.Args[1]
	stat, err := os.Stat(src)
	if err != nil {
		log.Fatal(err)
	}

	Compiler := lib.NewCompiler()

	if stat.IsDir() {
		dir, err := os.ReadDir(src)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range dir {
			if strings.HasSuffix(file.Name(), ".jack") {
				fileName := filepath.Join(src, file.Name())
				Compiler.Compile(fileName)
			}
		}
	} else {
		if !strings.HasSuffix(src, ".jack") {
			printHelpMessage()
		}
		Compiler.Compile(src)
	}
}
