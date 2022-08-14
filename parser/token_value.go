package parser

// NewLine is a common token representing a line feed character (\n)
var NewLine = TokenValue{
	Token: LineFeed,
	Value: string(LineFeed),
}

// TokenValue represents the Token and the raw string value declared by this token.
// This includes Line as a pointer to the original _full_ line for easier post-processing.
type TokenValue struct {
	Token Token
	Value string
	Line  *string
}
