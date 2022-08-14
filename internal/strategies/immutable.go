package strategies

import (
	"github.com/jimschubert/ignore/parser"
	"github.com/jimschubert/ignore/strategy"
)

// Immutable decorates a Strategy such that it communicates itself as immutable to other internal implementation.
// Be aware that this interface simply shows intent and can't actually enforce immutability in Go.
type Immutable interface {
	strategy.Strategy
	Immutable() bool
}

// immutableStrategy is an implementation of Strategy and Immutable
type immutableStrategy struct {
	strategy strategy.Strategy
}

// RuleBuilder …
func (i immutableStrategy) RuleBuilder() strategy.RuleBuilder {
	return i.strategy.RuleBuilder()
}

// DefinitionPath …
func (i immutableStrategy) DefinitionPath() string {
	return i.strategy.DefinitionPath()
}

// Parser …
func (i immutableStrategy) Parser() parser.Parser {
	return i.strategy.Parser()
}

// Immutable …
func (i immutableStrategy) Immutable() bool {
	return true
}

// AsImmutable decorates an existing strategy as Immutable
func AsImmutable(strategy strategy.Strategy) (Immutable, error) {
	if is, ok := strategy.(Immutable); ok {
		return is, nil
	}

	return &immutableStrategy{strategy: strategy}, nil
}
