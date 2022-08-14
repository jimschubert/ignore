package strategy

import (
	"github.com/jimschubert/ignore/parser"
	"github.com/jimschubert/ignore/rules"
)

type Strategy interface {
	DefinitionPath() string
	Parser() parser.Parser
	RuleBuilder() RuleBuilder
}

type RuleBuilder interface {
	RuleFor(tokens []parser.TokenValue) (rules.Rule, error)
}

type BuildRulesFrom func(tokens []parser.TokenValue) (rules.Rule, error)

func (b BuildRulesFrom) RuleFor(tokens []parser.TokenValue) (rules.Rule, error) {
	return b(tokens)
}
