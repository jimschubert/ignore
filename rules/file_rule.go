package rules

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jimschubert/ignore/parser"
)

// fileRule is a rule which applies to files, given a defined extension and target filename pattern
type fileRule struct {
	rule
	definedExt      string
	filenamePattern *regexp.Regexp
}

func (f fileRule) Evaluate(relativePath string) (Operation, error) {
	return evaluateRule(f, relativePath)
}

func (f fileRule) AppliesTo(relativePath string) bool {
	// Directories aren't files and should be evaluated only via directoryRule
	if fileInfo, err := os.Stat(relativePath); (err == nil || os.IsExist(err)) && fileInfo.IsDir() {
		return false
	}
	// todo: consider filepath.Match
	evaluatedExt := strings.TrimPrefix(filepath.Ext(relativePath), ".")
	if extensionPattern, err := filePattern(strings.TrimPrefix(f.definedExt, ".")); err == nil {
		if !extensionPattern.MatchString(evaluatedExt) {
			return false
		}
	}

	return f.filenamePattern.MatchString(relativePath)
}

func (f fileRule) GoString() string {
	b := bytes.Buffer{}
	b.WriteString("fileRule {")
	b.WriteString(fmt.Sprintf("\trule:\t%#v", f.rule))
	b.WriteString(fmt.Sprintf("\tdefinedExt:\t%#v", f.definedExt))
	b.WriteString(fmt.Sprintf("\tfilenamePattern:\t%s", f.filenamePattern.String()))
	b.WriteString("}")
	return b.String()
}

func NewFileRule(raw string, syntax []parser.TokenValue) (Rule, error) {
	definedExt := filepath.Ext(raw)
	pattern, err := filePattern(raw)
	if err != nil {
		return rule{}, err
	}

	return &fileRule{
		rule:            rule{raw: raw, syntax: syntax},
		definedExt:      definedExt,
		filenamePattern: pattern,
	}, nil
}

// Forces compilation error if interface contract changes (for any reflection use cases)
var (
	_ Rule           = &fileRule{}
	_ EvaluatingRule = &fileRule{}
)
