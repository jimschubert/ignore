package rules

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jimschubert/ignore/parser"
)

// filePatternFromTokens builds a regular expression from tokenized gitignore syntax.
// handles escape sequences, wildcards, special patterns, etc
// see https://git-scm.com/docs/gitignore
func filePatternFromTokens(syntax []parser.TokenValue) (*regexp.Regexp, error) {
	if len(syntax) == 0 {
		return regexp.Compile("^$")
	}

	var pattern strings.Builder
	skipNext := false
	prevWasEscape := false

	for i, token := range syntax {
		if skipNext {
			skipNext = false
			continue
		}

		// note: RootedMarker, DirectoryMarker, and PathDelim all have value "/"
		// we know which is which by index
		if token.Token == parser.RootedMarker && i == 0 {
			// leading slash handled at rule level (see rule.go), skip
			prevWasEscape = false
			continue
		}

		if token.Token == parser.DirectoryMarker && i == len(syntax)-1 {
			prevWasEscape = false
			// foo/ means anythin gin the foo directory
			if os.PathSeparator == '\\' {
				pattern.WriteString(regexp.QuoteMeta("\\"))
			} else {
				pattern.WriteString(`\/`)
			}

			// now match 0..n chars following the directory
			pattern.WriteString(`.*?`)
			continue
		}

		switch token.Token {
		case parser.Negate:
			// Negation handled at rule level (rule.go), skip here
			prevWasEscape = false
			continue

		case parser.PathDelim:
			// PathDelim and DirectoryMarker are hte same char "/"
			// This handles parser.DirectoryMarker in the middle of patterns
			if os.PathSeparator == '\\' {
				pattern.WriteString(regexp.QuoteMeta("\\"))
			} else {
				pattern.WriteString(`\/`)
			}
			prevWasEscape = false

		case parser.MatchAll:
			// ** matches zero or more directories
			// see:  https://git-scm.com/docs/gitignore: "**" followed by slash matches zero or more directories
			if i+1 < len(syntax) && syntax[i+1].Token == parser.PathDelim {
				// **/ at start or in middle
				if os.PathSeparator == '\\' {
					pattern.WriteString(`(?:.*?\\)?`)
				} else {
					pattern.WriteString(`(?:.*?\/)?`)
				}
				skipNext = true
			} else if i > 0 && syntax[i-1].Token == parser.PathDelim {
				// /** at end (already handled path delim, match everything after)
				pattern.WriteString(`.*?`)
			} else {
				// standalone ** - treat as * for compatibility
				pattern.WriteString(`.*?`)
			}
			prevWasEscape = false

		case parser.MatchAny:
			// * matches anything except path separator
			if os.PathSeparator == '\\' {
				pattern.WriteString(`[^\\]*?`)
			} else {
				pattern.WriteString(`[^\/]*?`)
			}
			prevWasEscape = false

		case parser.EscapedSpace:
			// Escaped space matches literal space
			// see:  https://git-scm.com/docs/gitignore: trailing spaces ignored unless quoted with backslash
			pattern.WriteString(` `)
			prevWasEscape = false

		case parser.Escape:
			// now we're in an escape pattern
			prevWasEscape = true
			continue

		case parser.Text:
			// Text values can contain:
			// 1. Literal characters (especially after Escape token)
			// 2. Gitignore wildcards like ? and [...] (if not after Escape)
			text := token.Value

			if prevWasEscape {
				// Previous token was Escape, so this text is all literal characters
				pattern.WriteString(regexp.QuoteMeta(text))

				prevWasEscape = false
				continue
			}

			// Not escaped - convert gitignore patterns to regex
			for j := 0; j < len(text); j++ {
				ch := text[j]
				switch ch {
				case '?':
					// ? matches any single character except /
					if os.PathSeparator == '\\' {
						pattern.WriteString(`[^\\]`)
					} else {
						pattern.WriteString(`[^\/]`)
					}
				case '[':
					// Character range - find next ]
					endIdx := strings.IndexByte(text[j+1:], ']')
					if endIdx >= 0 {
						// Valid character range - use as-is in regex
						rangeExpr := text[j : j+2+endIdx]
						pattern.WriteString(rangeExpr)
						j += 1 + endIdx
					} else {
						// No closing ], treat as literal
						pattern.WriteString(regexp.QuoteMeta(string(ch)))
					}
				default:
					// regular characters
					pattern.WriteString(regexp.QuoteMeta(string(ch)))
				}
			}

			prevWasEscape = false

		case parser.Comment:
			// shouldn't be here, but might if constructed manually
			prevWasEscape = false
			continue

		default:
			if token.Value != "" {
				pattern.WriteString(regexp.QuoteMeta(token.Value))
			}
			prevWasEscape = false
		}
	}

	return regexp.Compile(fmt.Sprintf("^%s$", pattern.String()))
}
