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
