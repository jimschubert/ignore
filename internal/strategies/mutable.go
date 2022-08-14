package strategies

import (
	"errors"

	"github.com/jimschubert/ignore/parser"
	"github.com/jimschubert/ignore/strategy"
)

// Mutable definition of a strategy which can be modified programmaticall
type Mutable interface {
	strategy.Strategy
	SetDefinitionPath(path string) error
	SetParser(parser parser.Parser) error
	SetRuleBuilder(builder strategy.RuleBuilder) error
}

// mutableStrategy is an implementation of Strategy and Mutable
type mutableStrategy struct {
	fullPath    *string
	parser      *parser.Parser
	ruleBuilder *strategy.RuleBuilder
}

// RuleBuilder …
func (m *mutableStrategy) RuleBuilder() strategy.RuleBuilder {
	return *m.ruleBuilder
}

// SetDefinitionPath …
func (m *mutableStrategy) SetDefinitionPath(path string) error {
	m.fullPath = &path
	return nil
}

// SetParser …
func (m *mutableStrategy) SetParser(parser parser.Parser) error {
	m.parser = &parser
	return nil
}

// SetRuleBuilder …
func (m *mutableStrategy) SetRuleBuilder(builder strategy.RuleBuilder) error {
	m.ruleBuilder = &builder
	return nil
}

// DefinitionPath …
func (m *mutableStrategy) DefinitionPath() string {
	return *m.fullPath
}

// Parser …
func (m *mutableStrategy) Parser() parser.Parser {
	return *m.parser
}

// AsMutable evaluates a strategy, if it can mutate it returns a mutable object. If Immutable, this raises an error.
func AsMutable(strategy strategy.Strategy) (Mutable, error) {
	switch value := strategy.(type) {
	case Immutable:
		return nil, errors.New("cannot modify immutable strategy")
	case Mutable:
		return value, nil
	default:
		d := strategy.DefinitionPath()
		p := strategy.Parser()
		b := strategy.RuleBuilder()
		return &mutableStrategy{
			fullPath:    &d,
			parser:      &p,
			ruleBuilder: &b,
		}, nil
	}
}
