package strategy

import (
	"github.com/jimschubert/ignore/parser"
	"github.com/jimschubert/ignore/rules"
)

// Strategy defines how to parse and build rules from an ignore file format.
type Strategy interface {
	DefinitionPath() string
	Parser() parser.Parser
	RuleBuilder() RuleBuilder
}

// RuleBuilder defines how to build rules from tokenized input.
type RuleBuilder interface {
	RuleFor(tokens []parser.TokenValue) (rules.Rule, error)
}

// BuildRulesFrom defines a function that builds rules from tokenized input.
type BuildRulesFrom func(tokens []parser.TokenValue) (rules.Rule, error)

// RuleFor builds a rule from tokenized input.
func (b BuildRulesFrom) RuleFor(tokens []parser.TokenValue) (rules.Rule, error) {
	return b(tokens)
}
