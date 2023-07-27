package analyzer_test

import (
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	. "github.com/MlkMahmud/jack-compiler/analyzer"
)

const TEST_DATA_PATH = "../testdata"

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
	files := []string{"Array", "Square", "SquareGame"}
	jackAnalyzer := NewAnalyzer()

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			srcFilePath := path.Join(TEST_DATA_PATH, strings.Join([]string{file, "jack"}, "."))
			outFilePath := path.Join(TEST_DATA_PATH, strings.Join([]string{file, "xml"}, "."))
			cmpFilePath := path.Join(TEST_DATA_PATH, "expected", strings.Join([]string{file, "xml"}, "."))

			jackAnalyzer.Run(srcFilePath)

			if stripLineBreakAndTabCharacters(cmpFilePath) != stripLineBreakAndTabCharacters(outFilePath) {
				t.Errorf("Expected content of %s to match content of %s", outFilePath, cmpFilePath)
			}
		})
	}
}
