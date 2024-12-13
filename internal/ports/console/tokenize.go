package console

import "strings"

// tokenize command line with quoted arguments support
// example: `command [subcommand] parameter "parameter with spaces" --flag value`
func tokenize(input string) []string {
	var tokens []string
	var currentToken strings.Builder
	inQuotes := false

	for _, char := range input {
		switch char {
		case '"':
			inQuotes = !inQuotes
		case ' ':
			if inQuotes {
				currentToken.WriteRune(char) // add space to token if we're quoted
			} else {
				if currentToken.Len() > 0 {
					tokens = append(tokens, currentToken.String()) // close current token
					currentToken.Reset()
				}
			}
		default:
			currentToken.WriteRune(char) // add rune to current token
		}
	}

	// add last token
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}
