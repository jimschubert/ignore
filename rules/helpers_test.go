package rules

import "github.com/jimschubert/ignore/parser"

func mustRootedFileRule(raw string, syntax []parser.TokenValue) rootedFileRule {
	result, err := NewRootedFileRule(raw, syntax)
	if err != nil {
		panic(`test: mustRootedFileRule(raw="` + raw + `"): ` + err.Error())
	}

	return *(result.(*rootedFileRule))
}
func parts(values ...parser.TokenValue) []parser.TokenValue {
	return values
}
