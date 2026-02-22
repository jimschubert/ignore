package parser

import (
	"reflect"
	"testing"

	"github.com/jimschubert/ignore/internal/util"
)

func TestLineParser_parse(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    []TokenValue
		wantErr bool
	}{
		{
			name: "comment",
			args: args{"# This is a comment"},
			want: []TokenValue{
				{Token: Comment, Value: "This is a comment", Line: util.Ptr("# This is a comment")},
			},
		},
		{
			name: "directory marker",
			args: args{"foo/"},
			want: []TokenValue{
				{Token: Text, Value: "foo", Line: util.Ptr("foo/")},
				{Token: DirectoryMarker, Line: util.Ptr("foo/")},
			},
		},
		{
			name: "rooted",
			args: args{"/abcd"},
			want: []TokenValue{
				{Token: RootedMarker, Line: util.Ptr("/abcd")},
				{Token: Text, Value: "abcd", Line: util.Ptr("/abcd")},
			},
		},
		{
			name: "escaped comment",
			args: args{"\\#file.txt"},
			want: []TokenValue{
				{Token: Escape, Line: util.Ptr("\\#file.txt")},
				{Token: Text, Value: "#file.txt", Line: util.Ptr("\\#file.txt")},
			},
		},
		{
			name: "escaped negate",
			args: args{"\\!important!.txt"},
			want: []TokenValue{
				{Token: Escape, Line: util.Ptr("\\!important!.txt")},
				{Token: Text, Value: "!important!.txt", Line: util.Ptr("\\!important!.txt")},
			},
		},
		{
			name: "complex",
			args: args{"**/abcd/**/foo/bar/sample.txt"},
			want: []TokenValue{
				{Token: MatchAll, Line: util.Ptr("**/abcd/**/foo/bar/sample.txt")},
				{Token: PathDelim, Line: util.Ptr("**/abcd/**/foo/bar/sample.txt")},
				{Token: Text, Value: "abcd", Line: util.Ptr("**/abcd/**/foo/bar/sample.txt")},
				{Token: PathDelim, Line: util.Ptr("**/abcd/**/foo/bar/sample.txt")},
				{Token: MatchAll, Line: util.Ptr("**/abcd/**/foo/bar/sample.txt")},
				{Token: PathDelim, Line: util.Ptr("**/abcd/**/foo/bar/sample.txt")},
				{Token: Text, Value: "foo", Line: util.Ptr("**/abcd/**/foo/bar/sample.txt")},
				{Token: PathDelim, Line: util.Ptr("**/abcd/**/foo/bar/sample.txt")},
				{Token: Text, Value: "bar", Line: util.Ptr("**/abcd/**/foo/bar/sample.txt")},
				{Token: PathDelim, Line: util.Ptr("**/abcd/**/foo/bar/sample.txt")},
				{Token: Text, Value: "sample.txt", Line: util.Ptr("**/abcd/**/foo/bar/sample.txt")},
			},
		},
		{
			name:    "triple star",
			args:    args{"***"},
			wantErr: true,
		},
		{
			name: "negate pattern",
			args: args{"!important.txt"},
			want: []TokenValue{
				{Token: Negate, Line: util.Ptr("!important.txt")},
				{Token: Text, Value: "important.txt", Line: util.Ptr("!important.txt")},
			},
		},
		{
			name: "trailing spaces ignored",
			args: args{"file.txt   "},
			want: []TokenValue{
				{Token: Text, Value: "file.txt", Line: util.Ptr("file.txt   ")},
			},
		},
		{
			// Per https://git-scm.com/docs/gitignore: Trailing spaces are ignored unless quoted with backslash
			name: "escaped trailing space",
			args: args{"file.txt\\ "},
			want: []TokenValue{
				{Token: Text, Value: "file.txt", Line: util.Ptr("file.txt\\ ")},
				{Token: EscapedSpace, Line: util.Ptr("file.txt\\ ")},
			},
		},
		{
			name: "character range",
			args: args{"file[0-9].txt"},
			want: []TokenValue{
				{Token: Text, Value: "file[0-9].txt", Line: util.Ptr("file[0-9].txt")},
			},
		},
		{
			name: "question mark wildcard",
			args: args{"file?.txt"},
			want: []TokenValue{
				{Token: Text, Value: "file?.txt", Line: util.Ptr("file?.txt")},
			},
		},
		{
			name: "trailing double asterisk",
			args: args{"dir/**"},
			want: []TokenValue{
				{Token: Text, Value: "dir", Line: util.Ptr("dir/**")},
				{Token: PathDelim, Line: util.Ptr("dir/**")},
				{Token: MatchAll, Line: util.Ptr("dir/**")},
			},
		},
		{
			name: "middle double asterisk",
			args: args{"a/**/b"},
			want: []TokenValue{
				{Token: Text, Value: "a", Line: util.Ptr("a/**/b")},
				{Token: PathDelim, Line: util.Ptr("a/**/b")},
				{Token: MatchAll, Line: util.Ptr("a/**/b")},
				{Token: PathDelim, Line: util.Ptr("a/**/b")},
				{Token: Text, Value: "b", Line: util.Ptr("a/**/b")},
			},
		},
		{
			name: "blank line",
			args: args{""},
			want: []TokenValue{},
		},
		{
			// Per https://git-scm.com/docs/gitignore: A blank line matches no files
			name: "only spaces",
			args: args{"   "},
			want: []TokenValue{},
		},
		{
			// Per https://git-scm.com/docs/gitignore: A backslash can escape special characters
			name: "escaped asterisk",
			args: args{"\\*file.txt"},
			want: []TokenValue{
				{Token: Escape, Line: util.Ptr("\\*file.txt")},
				{Token: Text, Value: "*file.txt", Line: util.Ptr("\\*file.txt")},
			},
		},
		{
			// Per https://git-scm.com/docs/gitignore: A backslash can escape special characters
			name: "escaped question mark",
			args: args{"file\\?.txt"},
			want: []TokenValue{
				{Token: Text, Value: "file", Line: util.Ptr("file\\?.txt")},
				{Token: Escape, Line: util.Ptr("file\\?.txt")},
				{Token: Text, Value: "?.txt", Line: util.Ptr("file\\?.txt")},
			},
		},
		{
			// Per https://git-scm.com/docs/gitignore: A backslash can escape special characters
			name: "escaped backslash",
			args: args{"file\\\\name.txt"},
			want: []TokenValue{
				{Token: Text, Value: "file", Line: util.Ptr("file\\\\name.txt")},
				{Token: Escape, Line: util.Ptr("file\\\\name.txt")},
				{Token: Text, Value: "\\name.txt", Line: util.Ptr("file\\\\name.txt")},
			},
		},
		{
			name: "rooted directory",
			args: args{"/root/"},
			want: []TokenValue{
				{Token: RootedMarker, Line: util.Ptr("/root/")},
				{Token: Text, Value: "root", Line: util.Ptr("/root/")},
				{Token: DirectoryMarker, Line: util.Ptr("/root/")},
			},
		},
		{
			name: "negated directory",
			args: args{"!keep/"},
			want: []TokenValue{
				{Token: Negate, Line: util.Ptr("!keep/")},
				{Token: Text, Value: "keep", Line: util.Ptr("!keep/")},
				{Token: DirectoryMarker, Line: util.Ptr("!keep/")},
			},
		},
		{
			name: "path with multiple slashes",
			args: args{"foo/bar/baz"},
			want: []TokenValue{
				{Token: Text, Value: "foo", Line: util.Ptr("foo/bar/baz")},
				{Token: PathDelim, Line: util.Ptr("foo/bar/baz")},
				{Token: Text, Value: "bar", Line: util.Ptr("foo/bar/baz")},
				{Token: PathDelim, Line: util.Ptr("foo/bar/baz")},
				{Token: Text, Value: "baz", Line: util.Ptr("foo/bar/baz")},
			},
		},
		{
			name: "wildcard in path",
			args: args{"foo/*/bar"},
			want: []TokenValue{
				{Token: Text, Value: "foo", Line: util.Ptr("foo/*/bar")},
				{Token: PathDelim, Line: util.Ptr("foo/*/bar")},
				{Token: MatchAny, Line: util.Ptr("foo/*/bar")},
				{Token: PathDelim, Line: util.Ptr("foo/*/bar")},
				{Token: Text, Value: "bar", Line: util.Ptr("foo/*/bar")},
			},
		},
		{
			name: "leading double asterisk",
			args: args{"**/foo"},
			want: []TokenValue{
				{Token: MatchAll, Line: util.Ptr("**/foo")},
				{Token: PathDelim, Line: util.Ptr("**/foo")},
				{Token: Text, Value: "foo", Line: util.Ptr("**/foo")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := gitignoreParser{}
			got, err := l.ParseLine(tt.args.text)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ParseLine() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseLine() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLineParser_parseSingleTokens(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    []TokenValue
		wantErr bool
	}{
		{
			name: "match all",
			args: args{"**"},
			want: []TokenValue{{Token: MatchAll, Line: util.Ptr("**")}},
		},
		{
			name: "match any",
			args: args{"*"},
			want: []TokenValue{{Token: MatchAny, Line: util.Ptr("*")}},
		},
		{
			name: "escaped space",
			args: args{`\ `},
			want: []TokenValue{{Token: EscapedSpace, Line: util.Ptr(`\ `)}},
		},
		{
			name:    "negate",
			args:    args{"!"},
			want:    []TokenValue{{Token: Negate, Line: util.Ptr("!")}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := gitignoreParser{}
			got, err := l.ParseLine(tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if len(got) != 1 {
					t.Errorf("ParseLine() expected a single returned token, got = %v", got)
					return
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ParseLine() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
