package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func printHelpMessage() {
	log.SetFlags(0)
	log.Fatalln(("usage:\n go run . .\t\t\tCompiles all the .jack files in the current directory\n go run . <src.jack>\t\tCompiles the specified .jack file\n go run . src/\t\t\tCompiles all the .jack files in the specified directory"))
}

func main() {
	if len(os.Args) != 2 {
		printHelpMessage()
	}

	JackAnalyzer := NewAnalyzer()
	src := os.Args[1]
	stat, err := os.Stat(src)
	if err != nil {
		log.Fatal(err)
	}

	if stat.IsDir() {
		dir, err := os.ReadDir(src)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range dir {
			if strings.HasSuffix(file.Name(), ".jack") {
				fileName := filepath.Join(src, file.Name())
				JackAnalyzer.Run(fileName)
			}
		}
	} else {
		if !strings.HasSuffix(src, ".jack") {
			printHelpMessage()
		}
		JackAnalyzer.Run(src)
	}
}
