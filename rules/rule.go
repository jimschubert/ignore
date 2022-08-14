package rules

import (
	"github.com/jimschubert/ignore/parser"
)

// rule defines functionality
type rule struct {
	// the parsed syntax
	syntax []parser.TokenValue
	// the raw string which created this rule
	raw string
	// the Operation to take when a rule does not match (i.e. allows a path)
	include *Operation
	// the Operation to take when a rule matches (i.e. don't allow a path)
	exclude *Operation
}

func (b rule) Syntax() []parser.TokenValue {
	return b.syntax
}

func (b rule) Raw() string {
	return b.raw
}

func (b rule) Include() Operation {
	if b.include == nil {
		return Include
	}
	return *b.include
}

func (b rule) Exclude() Operation {
	if b.exclude == nil {
		return Exclude
	}
	return *b.exclude
}

// Negated determines if the syntax is negated, and the processor must invert Include() vs Exclude().
//
// Take, for example, an ignore file containing:
//  - /path/to/fileA
//  - !/path/to/fileB
//
// This would result in Exclude() for fileA (Negated() is false), but Include() for fileB (Negated() is true)
func (b rule) Negated() bool {
	syntax := b.Syntax()
	return len(syntax) > 0 && syntax[0].Token == parser.Negate
}

// func (b rule) Pattern() string {
// 	syntax := b.Syntax()
// 	if len(syntax) == 0 {
// 		return b.Raw()
// 	}
//
// 	buf := bytes.Buffer{}
//
// 	for _, part := range syntax {
// 		switch part.Token {
// 		case parser.Negate, parser.RootedMarker, parser.Comment:
// 			// cleanup based on certain tokens that may occur in a rule's definition
// 			continue
// 		default:
// 			buf.WriteString(part.Value)
// 		}
// 	}
//
// 	return buf.String()
// }

// Rule is an interface for implementing any rule
type Rule interface {
	Syntax() []parser.TokenValue
	Raw() string
	Include() Operation
	Exclude() Operation
	Negated() bool
}

// EvaluatingRule is a Rule which can be evaluated against a target path
type EvaluatingRule interface {
	Rule
	AppliesTo(relativePath string) bool
	Evaluate(relativePath string) (Operation, error)
}

// evaluateRule is a helper for checking if evaluatingRule is applicable, then determining Rule.Include vs Rule.Exclude (supporting Rule.Negated logic)
func evaluateRule(evaluatingRule EvaluatingRule, relativePath string) (Operation, error) {
	if evaluatingRule.AppliesTo(relativePath) {
		if evaluatingRule.Negated() {
			return evaluatingRule.Include(), nil
		}

		return evaluatingRule.Exclude(), nil
	}

	return Noop, nil
}

var (
	_ Rule = &rule{}
)
