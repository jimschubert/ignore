package rules

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/jimschubert/ignore/parser"
)

func TestNewFileRule(t *testing.T) {
	type args struct {
		raw    string
		syntax []parser.TokenValue
	}

	fooText := "/foo"
	fooSyntax := parts(
		parser.TokenValue{Token: parser.RootedMarker},
		parser.TokenValue{Token: parser.Text, Value: "foo"},
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
			want: &fileRule{
				rule:            rule{raw: fooText, syntax: fooSyntax},
				filenamePattern: regexp.MustCompile("^foo$"),
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
			got, err := NewFileRule(tt.args.raw, tt.args.syntax)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewFileRule() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				// we don't care about evaluating trash structures when error is returned.
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileRule()\ngot=\n%#v\nwant=\n%#v", got, tt.want)
			}
		})
	}
}

func Test_fileRule_AppliesTo(t *testing.T) {
	type fields struct {
		rule            rule
		definedExt      string
		filenamePattern *regexp.Regexp
	}
	type args struct {
		relativePath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{

		// region no extension
		{
			name: "match: no extension to valid target",
			fields: fields{
				rule: rule{
					raw: "/foo",
					syntax: parts(
						parser.TokenValue{Token: parser.RootedMarker},
						parser.TokenValue{Token: parser.Text, Value: "foo"},
					),
				},
				filenamePattern: regexp.MustCompile(`foo`),
			},
			args: args{relativePath: "foo"},
			want: true,
		},
		{
			name: "no match: no extension to invalid target",
			fields: fields{
				rule: rule{
					raw: "/foo",
					syntax: parts(
						parser.TokenValue{Token: parser.RootedMarker},
						parser.TokenValue{Token: parser.Text, Value: "foo"},
					),
				},
				filenamePattern: regexp.MustCompile(`foo`),
			},
			args: args{relativePath: "bar"},
			want: false,
		},
		// endregion no extension

		// region with extension
		{
			name: "match: with extension to valid target",
			fields: fields{
				rule: rule{
					raw: "/foo.txt",
					syntax: parts(
						parser.TokenValue{Token: parser.RootedMarker},
						parser.TokenValue{Token: parser.Text, Value: "foo.txt"},
					),
				},
				definedExt:      ".txt",
				filenamePattern: regexp.MustCompile(`foo\.txt`),
			},
			args: args{relativePath: "foo.txt"},
			want: true,
		},
		{
			name: "no match: with extension to invalid target",
			fields: fields{
				rule: rule{
					raw: "/foo.txt",
					syntax: parts(
						parser.TokenValue{Token: parser.RootedMarker},
						parser.TokenValue{Token: parser.Text, Value: "foo.txt"},
					),
				},
				definedExt:      ".txt",
				filenamePattern: regexp.MustCompile(`foo\.txt`),
			},
			args: args{relativePath: "bar.txt"},
			want: false,
		},
		// endregion with extension

		// region nested files
		{
			name: "no match: nested file paths target",
			fields: fields{
				rule: rule{
					raw: "/foo.txt",
					syntax: parts(
						parser.TokenValue{Token: parser.RootedMarker},
						parser.TokenValue{Token: parser.Text, Value: "foo.txt"},
					),
				},
				definedExt:      ".txt",
				filenamePattern: regexp.MustCompile(`foo\.txt`),
			},
			args: args{relativePath: "foo/bar.txt"},
			want: false,
		},
		// endregion nested files

		// region pattern matching files
		{
			name: "match: pattern matching filename to valid target",
			fields: fields{
				rule: rule{
					raw: "/fo*.txt",
					syntax: parts(
						parser.TokenValue{Token: parser.RootedMarker},
						parser.TokenValue{Token: parser.Text, Value: "fo*.txt"},
					),
				},
				definedExt:      ".txt",
				filenamePattern: regexp.MustCompile(`fo.*?\.txt`),
			},
			args: args{relativePath: "fop.txt"},
			want: true,
		},
		{
			name: "match: pattern matching extension to valid target",
			fields: fields{
				rule: rule{
					raw: "/foo.*",
					syntax: parts(
						parser.TokenValue{Token: parser.RootedMarker},
						parser.TokenValue{Token: parser.Text, Value: "foo.*"},
					),
				},
				definedExt:      ".*",
				filenamePattern: regexp.MustCompile(`foo\..*?`),
			},
			args: args{relativePath: "foo.bar"},
			want: true,
		},
		{
			name: "match: pattern matching extension to valid target with complex expression",
			fields: fields{
				rule: rule{
					raw: "/foo.t*t",
					syntax: parts(
						parser.TokenValue{Token: parser.RootedMarker},
						parser.TokenValue{Token: parser.Text, Value: "foo.t*t"},
					),
				},
				definedExt:      ".*",
				filenamePattern: regexp.MustCompile(`foo\.t.*?t`),
			},
			args: args{relativePath: "foo.tot"},
			want: true,
		},
		{
			name: "no match: pattern matching filename to invalid target",
			fields: fields{
				rule: rule{
					raw: "/fo*.txt",
					syntax: parts(
						parser.TokenValue{Token: parser.RootedMarker},
						parser.TokenValue{Token: parser.Text, Value: "fo*.txt"},
					),
				},
				definedExt:      ".txt",
				filenamePattern: regexp.MustCompile(`fo.*?\.txt`),
			},
			args: args{relativePath: "bar.txt"},
			want: false,
		},
		{
			name: "no match: pattern matching extension to invalid target",
			fields: fields{
				rule: rule{
					raw: "/foo.t*t",
					syntax: parts(
						parser.TokenValue{Token: parser.RootedMarker},
						parser.TokenValue{Token: parser.Text, Value: "foo.t*t"},
					),
				},
				definedExt:      ".*",
				filenamePattern: regexp.MustCompile(`foo\.t.*?t`),
			},
			args: args{relativePath: "foo.bar"},
			want: false,
		},
		// endregion pattern matching files
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := fileRule{
				rule:            tt.fields.rule,
				definedExt:      tt.fields.definedExt,
				filenamePattern: tt.fields.filenamePattern,
			}
			if got := r.AppliesTo(tt.args.relativePath); got != tt.want {
				t.Errorf("AppliesTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileRule_Evaluate(t *testing.T) {
	type fields struct {
		rule            rule
		definedExt      string
		filenamePattern *regexp.Regexp
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
			f := fileRule{
				rule:            tt.fields.rule,
				definedExt:      tt.fields.definedExt,
				filenamePattern: tt.fields.filenamePattern,
			}
			got, err := f.Evaluate(tt.args.relativePath)
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
