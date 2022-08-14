package rules

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/jimschubert/ignore/parser"
)

func TestNewRootedFileRule(t *testing.T) {
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
			want: &rootedFileRule{
				fileRule: fileRule{
					rule:            rule{raw: fooText, syntax: fooSyntax},
					filenamePattern: regexp.MustCompile(`^foo$`),
				},
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
			got, err := NewRootedFileRule(tt.args.raw, tt.args.syntax)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewRootedFileRule() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				// we don't care about evaluating trash structures when error is returned.
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRootedFileRule()\ngot=\n%#v\nwant=\n%#v", got, tt.want)
			}
		})
	}
}

func Test_rootedFileRule_AppliesTo(t *testing.T) {
	type fields struct {
		raw    string
		syntax []parser.TokenValue
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
				raw: "/foo",
				syntax: parts(
					parser.TokenValue{Token: parser.RootedMarker},
					parser.TokenValue{Token: parser.Text, Value: "foo"},
				),
			},
			args: args{relativePath: "foo"},
			want: true,
		},
		{
			name: "no match: no extension to invalid target",
			fields: fields{
				raw: "/foo",
				syntax: parts(
					parser.TokenValue{Token: parser.RootedMarker},
					parser.TokenValue{Token: parser.Text, Value: "foo"},
				),
			},
			args: args{relativePath: "bar"},
			want: false,
		},
		// endregion no extension

		// region with extension
		{
			name: "match: with extension to valid target",
			fields: fields{
				raw: "/foo.txt",
				syntax: parts(
					parser.TokenValue{Token: parser.RootedMarker},
					parser.TokenValue{Token: parser.Text, Value: "foo.txt"},
				),
			},
			args: args{relativePath: "foo.txt"},
			want: true,
		},
		{
			name: "no match: with extension to invalid target",
			fields: fields{
				raw: "/foo.txt",
				syntax: parts(
					parser.TokenValue{Token: parser.RootedMarker},
					parser.TokenValue{Token: parser.Text, Value: "foo.txt"},
				),
			},
			args: args{relativePath: "bar.txt"},
			want: false,
		},
		// endregion with extension

		// region nested files
		{
			name: "no match: nested file paths target",
			fields: fields{
				raw: "/foo.txt",
				syntax: parts(
					parser.TokenValue{Token: parser.RootedMarker},
					parser.TokenValue{Token: parser.Text, Value: "foo.txt"},
				),
			},
			args: args{relativePath: "foo/bar.txt"},
			want: false,
		},
		// endregion nested files

		// region pattern matching files
		{
			name: "match: pattern matching filename to valid target",
			fields: fields{
				raw: "/fo*.txt",
				syntax: parts(
					parser.TokenValue{Token: parser.RootedMarker},
					parser.TokenValue{Token: parser.Text, Value: "fo*.txt"},
				),
			},
			args: args{relativePath: "fop.txt"},
			want: true,
		},
		{
			name: "match: pattern matching extension to valid target",
			fields: fields{
				raw: "/foo.*",
				syntax: parts(
					parser.TokenValue{Token: parser.RootedMarker},
					parser.TokenValue{Token: parser.Text, Value: "foo.*"},
				),
			},
			args: args{relativePath: "foo.bar"},
			want: true,
		},
		{
			name: "match: pattern matching extension to valid target with complex expression",
			fields: fields{
				raw: "/foo.t*t",
				syntax: parts(
					parser.TokenValue{Token: parser.RootedMarker},
					parser.TokenValue{Token: parser.Text, Value: "foo.t*t"},
				),
			},
			args: args{relativePath: "foo.tot"},
			want: true,
		},
		{
			name: "no match: pattern matching filename to invalid target",
			fields: fields{
				raw: "/fo*.txt",
				syntax: parts(
					parser.TokenValue{Token: parser.RootedMarker},
					parser.TokenValue{Token: parser.Text, Value: "fo*.txt"},
				),
			},
			args: args{relativePath: "bar.txt"},
			want: false,
		},
		{
			name: "no match: pattern matching extension to invalid target",
			fields: fields{
				raw: "/foo.t*t",
				syntax: parts(
					parser.TokenValue{Token: parser.RootedMarker},
					parser.TokenValue{Token: parser.Text, Value: "foo.t*t"},
				),
			},
			args: args{relativePath: "foo.bar"},
			want: false,
		},
		// endregion pattern matching files
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mustRootedFileRule(tt.fields.raw, tt.fields.syntax)
			if got := r.AppliesTo(tt.args.relativePath); got != tt.want {
				t.Errorf("AppliesTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

// todo implement
//nolint:unused
func Test_rootedFileRule_Evaluate(t *testing.T) {
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
			r := rootedFileRule{
				// rule:            tt.fields.rule,
				// definedExt:      tt.fields.definedExt,
				// filenamePattern: tt.fields.filenamePattern,
			}
			got, err := r.Evaluate(tt.args.relativePath)
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
