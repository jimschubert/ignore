package rules

import (
	"github.com/jimschubert/ignore/parser"
)

// invalidRule is a raw implementation of a rule, which can be used if an
// implementation wants to parse as much as possible from an ignore file.
type invalidRule struct {
	rule
}

func NewInvalidRule(raw string, syntax []parser.TokenValue) (Rule, error) {
	return &invalidRule{
		rule: rule{raw: raw, syntax: syntax},
	}, nil
}

var (
	_ Rule = &invalidRule{}
)
