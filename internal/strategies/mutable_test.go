package strategies

import (
	"reflect"
	"testing"

	"github.com/jimschubert/ignore/parser"
	"github.com/jimschubert/ignore/strategy"
)

func TestAsMutable(t *testing.T) {
	simpleEmpty := simpleStrategy{}
	simpleFull := simpleStrategy{
		fullPath: ".ignore",
		parser:   parser.NewGitignoreParser(),
	}

	type args struct {
		strategy strategy.Strategy
	}
	tests := []struct {
		name    string
		args    args
		want    Mutable
		wantErr bool
	}{
		{name: "can make empty struct mutable", args: args{strategy: simpleEmpty}, want: &mutableStrategy{
			fullPath: &simpleEmpty.fullPath, parser: &simpleEmpty.parser, ruleBuilder: &simpleEmpty.ruleBuilder,
		}, wantErr: false},
		{name: "can make gitignore struct mutable", args: args{strategy: simpleFull}, want: &mutableStrategy{
			fullPath: &simpleFull.fullPath, parser: &simpleFull.parser, ruleBuilder: &simpleFull.ruleBuilder,
		}, wantErr: false},
		{name: "fail if immutable", args: args{strategy: immutableStrategy{simpleFull}}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AsMutable(tt.args.strategy)
			if (err != nil) != tt.wantErr {
				t.Errorf("AsMutable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsMutable() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
