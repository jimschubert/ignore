package rules

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/jimschubert/ignore/parser"
)

// directoryRule is a rule which applies to directories
type directoryRule struct {
	rule
	// reuse pattern for performance
	directoryPattern *regexp.Regexp
}

func (d directoryRule) Evaluate(relativePath string) (Operation, error) {
	return evaluateRule(d, relativePath)
}

func (d directoryRule) AppliesTo(relativePath string) bool {
	// if path exists, but is _not_ a directory, we won't apply a directory rule.
	// That is, if our rule is /path/to/cupcakes/ and there's a file at /path/to/cupcakes, the rule won't evaluate.
	if fileInfo, err := os.Stat(relativePath); (err == nil || os.IsExist(err)) && !fileInfo.IsDir() {
		return false
	}

	return d.directoryPattern.MatchString(relativePath)
}

func (d directoryRule) GoString() string {
	b := bytes.Buffer{}
	b.WriteString("directoryRule {")
	b.WriteString(fmt.Sprintf("\trule:\t%#v", d.rule))
	b.WriteString("}")
	return b.String()
}

// NewDirectoryRule constructs a new directory rule from raw syntax, exposing an error if the raw pattern is invalid.
func NewDirectoryRule(raw string, syntax []parser.TokenValue) (Rule, error) {
	// Directory patterns are weird (see https://git-scm.com/docs/gitignore)
	pattern, err := filePatternFromTokens(syntax)
	if err != nil {
		return rule{}, err
	}

	return &directoryRule{
		rule:             rule{raw: raw, syntax: syntax},
		directoryPattern: pattern,
	}, nil
}

// Forces compilation error if interface contract changes (for any reflection use cases)
var (
	_ Rule           = &directoryRule{}
	_ EvaluatingRule = &directoryRule{}
)
