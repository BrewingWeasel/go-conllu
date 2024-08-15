package conllu

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Parse parses conllu data via the io.Reader interface and returns all of the tokens found
// Parse doesn't close the reader when finished, that must be done manually
func Parse(r io.Reader) ([]Sentence, []error) {
	errors := []error{}
	sentences := []Sentence{}
	currentSentence := Sentence{
		Tokens: []Token{},
	}
	reader := bufio.NewReader(r)
	done := false
	lineNumber := 0
	for !done {
		lineNumber++
		line, err := reader.ReadString('\n')
		if err != nil {
			done = true
		}
		token, isComment, commentValue, isSep, err := parseLine(line)
		if err != nil {
			errors = append(errors, fmt.Errorf("error on line %v, err: %v", lineNumber, err))
			continue
		}
		if isComment {
			if text, exists := strings.CutPrefix(strings.Trim(commentValue, " "), "text = "); exists {
				currentSentence.Text = text
			}
			continue
		}
		if isSep {
			if len(currentSentence.Tokens) == 0 {
				continue
			}
			sentences = append(sentences, currentSentence)
			currentSentence = Sentence{
				Tokens: []Token{},
			}
			continue
		}
		currentSentence.Tokens = append(currentSentence.Tokens, token...)
	}
	return sentences, errors
}

// ParseFile opens, reads, and parses a file in conllu format and returns all of the tokens found.
// ParseFile is a convencience wrapper for the Parse() function when working with files on disk
func ParseFile(filepath string) ([]Sentence, []error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, []error{err}
	}
	defer file.Close()
	return Parse(file)
}
