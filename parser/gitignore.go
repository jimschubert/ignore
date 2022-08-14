package parser

import (
	"bytes"
	"io"
	"strings"
)

type gitignoreParser struct {
}

// ParseAll contents from reader to a collection of TokenValue
func (l gitignoreParser) ParseAll(reader io.Reader) ([]TokenValue, error) {
	helpers := NewParsingHelpers(l)
	return helpers.ParseAll(reader)
}

// ParseAllText from text input to a collection of TokenValue
func (l gitignoreParser) ParseAllText(text string) ([]TokenValue, error) {
	helpers := NewParsingHelpers(l)
	return helpers.ParseAllText(text)
}

// ParseLine parses a line of text
func (l gitignoreParser) ParseLine(text string) ([]TokenValue, error) {
	parts := make([]TokenValue, 0)

	switch {
	case text == ".", text == "!.", strings.HasPrefix(text, ".."):
		return parts, &InvalidPatternError{pattern: text}
	}

	runes := []rune(text)

	if len(runes) == 0 {
		return parts, nil
	}

	buf := bytes.Buffer{}
	last := len(runes) - 1
	for i := 0; i < len(runes); i++ {
		current := runes[i]
		var next rune
		if i < last {
			next = runes[i+1]
		}

		if i == 0 {
			if Comment.MatchRune(current) {
				commentText := strings.TrimSpace(strings.TrimPrefix(text, string(Comment)))
				parts = append(parts, TokenValue{Token: Comment, Value: commentText, Line: &text})
				break
			}

			if Negate.MatchRune(current) {
				if i == last {
					return parts, newParsingError("negation with no negated pattern")
				}
				parts = append(parts, TokenValue{Token: Negate, Line: &text})
				continue
			}

			if RootedMarker.MatchRune(current) {
				parts = append(parts, TokenValue{Token: RootedMarker, Line: &text})
				continue
			}

			if Escape.MatchRune(current) && anyTokenMatch(next, Comment, Negate) {
				// : Put a backslash ("`\`") in front of the first hash for patterns
				// : that begin with a hash.
				// NOTE: Just push forward and "drop" the escape character. Falls through to TEXT token.
				// we still track the escape character so the parser can eventually recreate documents
				parts = append(parts, TokenValue{Token: Escape, Line: &text})
				current = next
				next = 0
				i++
			}
		}

		if MatchAny.MatchRune(current) {
			if MatchAny.MatchRune(next) {
				if (i+2) < len(runes) && MatchAny.MatchRune(runes[i+2]) {
					return parts, &InvalidPatternError{pattern: "***"}
				}

				parts = append(parts, TokenValue{Token: MatchAll, Line: &text})
				i++
				continue
			}

			if buf.Len() > 0 {
				parts = append(parts, TokenValue{Token: Text, Value: buf.String(), Line: &text})
				buf.Reset()
			}

			parts = append(parts, TokenValue{Token: MatchAny, Line: &text})
			continue
		}

		if EscapedSpace.MatchRunes(current, next) {
			parts = append(parts, TokenValue{Token: EscapedSpace, Line: &text})
			i++
			continue
		}

		if PathDelim.MatchRune(current) {
			if i == last {
				parts = append(parts, TokenValue{Token: Text, Value: buf.String(), Line: &text})
				buf.Reset()
				parts = append(parts, TokenValue{Token: DirectoryMarker, Line: &text})
				continue
			} else {
				if buf.Len() > 0 {
					parts = append(parts, TokenValue{Token: Text, Value: buf.String(), Line: &text})
					buf.Reset()
				}

				parts = append(parts, TokenValue{Token: PathDelim, Line: &text})
				if PathDelim.MatchRune(next) {
					// ignore doubled path delims. NOTE: doesn't do full lookahead, so /// will result in //
					i++
				}
				continue
			}
		}

		buf.WriteRune(current)
	}

	if buf.Len() > 0 {
		// NOTE: All spaces escaped spaces are a special token, ESCAPED_SPACE
		// : Trailing spaces are ignored unless they are quoted with backslash ("`\`")
		parts = append(parts, TokenValue{Token: Text, Value: strings.TrimSpace(buf.String()), Line: &text})
		buf.Reset()
	}

	return parts, nil
}

// NewGitignoreParser is a strategy which parses input text line-by-line
func NewGitignoreParser() Parser {
	return gitignoreParser{}
}

// Forces compilation error if interface contract changes (for any reflection use cases)
var (
	_ Parser = &gitignoreParser{}
)
