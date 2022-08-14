package rules

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// filePattern builds up a regular expression from an ignore-pattern style glob
func filePattern(input string) (*regexp.Regexp, error) {
	// TODO: Consider platform specific library for converting input glob to regex.
	globCleanup := input
	if strings.HasPrefix(input, "!") {
		// drop leading negation character
		globCleanup = strings.TrimPrefix(globCleanup, "!")
	}
	if strings.HasPrefix(input, "/") {
		// drop leading path separator character
		globCleanup = strings.TrimPrefix(globCleanup, "/")
	}
	// dots in file pattern are explicit characters
	globCleanup = strings.ReplaceAll(globCleanup, ".", `\Q.\E`)
	// $ symbols in file pattern are explicit characters
	globCleanup = strings.ReplaceAll(globCleanup, "$", `\Q$\E`)
	// double-asterisk is redundant
	globCleanup = strings.ReplaceAll(globCleanup, "**", "*")
	// non-greedy match * as 0..n characters
	globCleanup = strings.ReplaceAll(globCleanup, "*", ".*?")
	globCleanup = strings.TrimPrefix(globCleanup, "/")

	if os.PathSeparator == '\\' {
		globCleanup = strings.ReplaceAll(globCleanup, "/", regexp.QuoteMeta("\\"))
	} else {
		// *nix path separator needs to be escaped
		globCleanup = strings.ReplaceAll(globCleanup, `/`, `\/`)
	}
	return regexp.Compile(fmt.Sprintf("^%s$", globCleanup))
}
