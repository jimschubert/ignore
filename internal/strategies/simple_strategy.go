package strategies

import (
	"github.com/jimschubert/ignore/parser"
	"github.com/jimschubert/ignore/strategy"
)

type simpleStrategy struct {
	fullPath    string
	parser      parser.Parser
	ruleBuilder strategy.RuleBuilder
}

func (s simpleStrategy) RuleBuilder() strategy.RuleBuilder {
	return s.ruleBuilder
}

func (s simpleStrategy) DefinitionPath() string {
	return s.fullPath
}

func (s simpleStrategy) Parser() parser.Parser {
	return s.parser
}
