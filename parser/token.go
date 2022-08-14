package parser

// Token decorates string as a ignore file token
type Token string

// MatchRune determines if a token matches a single rune
func (t Token) MatchRune(r rune) bool {
	return string(t) == string(r)
}

// MatchRunes determines if a token matches multiple runes
func (t Token) MatchRunes(r ...rune) bool {
	return string(t) == string(r)
}

// anyTokenMatch returns true if any token matches the target run
func anyTokenMatch(value rune, t ...Token) bool {
	for _, token := range t {
		if token.MatchRune(value) {
			return true
		}
	}
	return false
}

const (
	MatchAll        Token = "**"
	MatchAny        Token = "*"
	EscapedSpace    Token = "\\ "
	Escape          Token = "\\"
	PathDelim       Token = "/"
	Negate          Token = "!"
	Text            Token = ""
	DirectoryMarker Token = "/"
	RootedMarker    Token = "/"
	Comment         Token = "#"
	LineFeed        Token = "\n"
)
