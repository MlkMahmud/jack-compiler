package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"testing"
)

func stripLineBreakAndTabCharacters(filename string) string {
	file, _ := os.Open(filename)
	bytes, err := io.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}

	fileContent := string(bytes)
	re := regexp.MustCompile("[\r\t\n]")
	return re.ReplaceAllLiteralString(fileContent, "")
}

func TestJackAnalyzer(t *testing.T) {
	files := []string{"testdata/Array.jack", "testdata/Square.jack", "testdata/SquareGame.jack"}
	jackAnalyzer := NewAnalyzer()

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			outFile := replaceFileExt(file, ".xml")
			jackAnalyzer.Run(file)

			info, err := os.Stat(outFile)
			if err != nil {
				t.Error(err)
			}

			cmpFile := fmt.Sprintf("testdata/xml/%s", info.Name())

			if stripLineBreakAndTabCharacters(cmpFile) != stripLineBreakAndTabCharacters(outFile) {
				t.Errorf("Expected content of %s to match content of %s", outFile, cmpFile)
			}
		})
	}
}
