package parser

import "fmt"

// InvalidPatternError should be used for patterns which can't be parsed.
type InvalidPatternError struct {
	pattern string
}

// Error string representation of InvalidPatternError
func (i *InvalidPatternError) Error() string {
	return fmt.Sprintf("Pattern '%s' is invalid.", i.pattern)
}

// ParsingError provides details of any error raised in the parser package
type ParsingError struct {
	data string
}

// Error string representation of ParsingError
func (p *ParsingError) Error() string {
	return fmt.Sprintf("parsing error: %s", p.data)
}

func newParsingError(message string) error {
	return &ParsingError{data: message}
}
