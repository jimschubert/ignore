package rules

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/jimschubert/ignore/parser"
)

// directoryRule is a rule which applies to directories
type directoryRule struct {
	rule
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

	noTrail := strings.TrimSuffix(d.rule.Raw(), "/")
	if strings.Count(noTrail, `/`) == 0 {
		if singleDirectory, err := filePattern(`(*/)?` + noTrail + `/*`); err == nil {
			return singleDirectory.MatchString(relativePath)
		}
	} else {
		// This logic taken from .gitignore logic:
		// For example, a pattern doc/frotz/ matches doc/frotz directory, but not a/doc/frotz directory; however
		// frotz/ matches frotz and a/frotz that is a directory (all paths are relative from the .gitignore file).
		if multiDirectory, err := filePattern(`^` + noTrail + `/?*`); err == nil {
			return multiDirectory.MatchString(relativePath)
		}
	}

	return false
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
	// check if the raw definition can be treated as a regex…
	if _, err := filePattern(raw); err != nil {
		return rule{}, err
	}
	return &directoryRule{
		rule: rule{raw: raw, syntax: syntax},
	}, nil
}

// Forces compilation error if interface contract changes (for any reflection use cases)
var (
	_ Rule           = &directoryRule{}
	_ EvaluatingRule = &directoryRule{}
)
