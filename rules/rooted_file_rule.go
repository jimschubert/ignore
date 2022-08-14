package rules

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jimschubert/ignore/parser"
)

// rootedFileRule is a special case which applies to only files in the root (same path as ignore file)
type rootedFileRule struct {
	fileRule
}

func (r rootedFileRule) Evaluate(relativePath string) (Operation, error) {
	return evaluateRule(r, relativePath)
}

func (r rootedFileRule) AppliesTo(relativePath string) bool {
	subPath := strings.Contains(strings.TrimPrefix(relativePath, "/"), "/")
	if subPath {
		return false
	}
	return r.fileRule.AppliesTo(relativePath)
}

func (r rootedFileRule) GoString() string {
	b := bytes.Buffer{}
	b.WriteString("fileRule {")
	b.WriteString(fmt.Sprintf("\tfileRule:\t%#v", r.fileRule))
	b.WriteString("}")
	return b.String()
}

func NewRootedFileRule(raw string, syntax []parser.TokenValue) (Rule, error) {
	f, err := NewFileRule(raw, syntax)
	if err != nil {
		return rule{}, err
	}

	return &rootedFileRule{
		fileRule: *(f.(*fileRule)),
	}, nil
}

// Forces compilation error if interface contract changes (for any reflection use cases)
var (
	_ Rule           = &rootedFileRule{}
	_ EvaluatingRule = &rootedFileRule{}
)
