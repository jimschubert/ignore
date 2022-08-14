package rules

import (
	"reflect"
	"testing"

	"github.com/jimschubert/ignore/parser"
)

func TestNewDirectoryRule(t *testing.T) {
	type args struct {
		raw    string
		syntax []parser.TokenValue
	}

	fooText := "foo/"
	fooSyntax := parts(
		parser.TokenValue{Token: parser.Text, Value: "foo"},
		parser.TokenValue{Token: parser.DirectoryMarker},
	)

	tests := []struct {
		name    string
		args    args
		want    Rule
		wantErr bool
	}{
		{
			name: "new without error",
			args: args{
				raw:    fooText,
				syntax: fooSyntax,
			},
			want: &directoryRule{
				rule: rule{raw: fooText, syntax: fooSyntax},
			},
		},

		{
			name: "new with error",
			args: args{
				raw:    `/path/to/*g(-z]+ng`,
				syntax: fooSyntax,
			},
			want:    &rule{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDirectoryRule(tt.args.raw, tt.args.syntax)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewDirectoryRule() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				// we don't care about evaluating trash structures when error is returned.
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDirectoryRule() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_directoryRule_AppliesTo(t *testing.T) {
	type fields struct {
		rule rule
	}
	tests := []struct {
		name         string
		fields       fields
		relativePath string
		want         bool
	}{
		{
			name: "match path foo/ for foo/",
			fields: fields{rule: rule{raw: "foo/", syntax: parts(
				parser.TokenValue{Token: parser.Text, Value: "foo"},
				parser.TokenValue{Token: parser.DirectoryMarker},
			)}},
			relativePath: "foo/",
			want:         true,
		},
		{
			name: "match path foo/bar for foo/",
			fields: fields{rule: rule{raw: "foo/", syntax: parts(
				parser.TokenValue{Token: parser.Text, Value: "foo"},
				parser.TokenValue{Token: parser.DirectoryMarker},
			)}},
			relativePath: "foo/bar",
			want:         true,
		},
		{
			name: "match path foo/bar for foo/**",
			fields: fields{rule: rule{raw: "foo/**", syntax: parts(
				parser.TokenValue{Token: parser.Text, Value: "foo"},
				parser.TokenValue{Token: parser.DirectoryMarker},
				parser.TokenValue{Token: parser.MatchAll},
			)}},
			relativePath: "foo/bar",
			want:         true,
		},
		{
			name: "match path doc/frotz/ for doc/frotz/",
			fields: fields{rule: rule{raw: "doc/frotz/", syntax: parts(
				parser.TokenValue{Token: parser.Text, Value: "doc"},
				parser.TokenValue{Token: parser.DirectoryMarker},
				parser.TokenValue{Token: parser.Text, Value: "frotz"},
				parser.TokenValue{Token: parser.DirectoryMarker},
			)}},
			relativePath: "doc/frotz/",
			want:         true,
		},
		{
			name: "does not match path a/doc/frotz/ for doc/frotz/",
			fields: fields{rule: rule{raw: "doc/frotz/", syntax: parts(
				parser.TokenValue{Token: parser.Text, Value: "doc"},
				parser.TokenValue{Token: parser.DirectoryMarker},
				parser.TokenValue{Token: parser.Text, Value: "frotz"},
				parser.TokenValue{Token: parser.DirectoryMarker},
			)}},
			relativePath: "a/doc/frotz/",
			want:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := directoryRule{
				rule: tt.fields.rule,
			}
			if got := d.AppliesTo(tt.relativePath); got != tt.want {
				t.Errorf("AppliesTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_directoryRule_Evaluate(t *testing.T) {
	type fields struct {
		rule rule
	}
	type args struct {
		relativePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Operation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := directoryRule{
				rule: tt.fields.rule,
			}
			got, err := d.Evaluate(tt.args.relativePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Evaluate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
