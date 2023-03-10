package tokenizer

import (
	"strings"
	"unicode"
)

// tokenize returns a slice of tokens for the given text.
func Tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		// Split on any character that is not a letter or a number.
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// analyze analyzes the text and returns a slice of tokens.
func Analyze(text string) []string {
	tokens := Tokenize(text)
	tokens = lowercaseFilter(tokens)
	tokens = stopwordFilter(tokens)
	tokens = stemmerFilter(tokens)
	// fmt.Println("====DEBUG tokens: ", tokens, "====")
	return tokens
}
