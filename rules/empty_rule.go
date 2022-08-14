package rules

import "github.com/jimschubert/ignore/parser"

// NewEmptyRule creates an empty rule for a given syntax. Useful for when one wants to retain a line
// without resulting in a parsing error
func NewEmptyRule(raw string, syntax []parser.TokenValue) (Rule, error) {
	return rule{raw: raw, syntax: syntax}, nil
}
