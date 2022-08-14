package parser

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// Parser for parsing text
type Parser interface {
	// ParseLine for parsing a single line of text to a collection of TokenValue
	ParseLine(text string) ([]TokenValue, error)

	// ParseAll contents from reader to a collection of TokenValue
	ParseAll(reader io.Reader) ([]TokenValue, error)

	// ParseAllText from text input to a collection of TokenValue
	ParseAllText(text string) ([]TokenValue, error)
}

type ParsingHelpers struct {
	target Parser
}

func (p ParsingHelpers) ParseLine(text string) ([]TokenValue, error) {
	return nil, errors.New("not implemented")
}

func (p ParsingHelpers) ParseAll(reader io.Reader) ([]TokenValue, error) {
	scanner := bufio.NewScanner(reader)
	result := make([]TokenValue, 0)
	iteration := 0

	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lineText := scanner.Text()
		lineSyntax, err := p.target.ParseLine(lineText)
		if err != nil {
			if err == io.EOF {
				break
			}
			return result, err
		}

		if iteration > 0 {
			result = append(result, NewLine)
		}

		result = append(result, lineSyntax...)
		iteration += 1
	}

	return result, nil
}

func (p ParsingHelpers) ParseAllText(text string) ([]TokenValue, error) {
	return p.ParseAll(strings.NewReader(text))
}

func NewParsingHelpers(p Parser) ParsingHelpers {
	return ParsingHelpers{target: p}
}

// Forces compilation error if interface contract changes (for any reflection use cases)
var (
	_ Parser = &ParsingHelpers{}
)
