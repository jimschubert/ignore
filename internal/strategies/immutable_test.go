package strategies

import (
	"reflect"
	"testing"

	"github.com/jimschubert/ignore/strategy"
)

func TestAsImmutable(t *testing.T) {
	simpleEmpty := simpleStrategy{}
	immutable := immutableStrategy{strategy: simpleEmpty}

	type args struct {
		strategy strategy.Strategy
	}
	tests := []struct {
		name    string
		args    args
		want    Immutable
		wantErr bool
	}{
		// note: &immutable gives us het interface to compare
		{name: "converts", args: args{strategy: simpleEmpty}, want: &immutable, wantErr: false},
		// note: lack of reference (&immutable) causes deep equal of the referenced struct itself
		{name: "reuses reference", args: args{strategy: immutable}, want: immutable, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AsImmutable(tt.args.strategy)
			if (err != nil) != tt.wantErr {
				t.Errorf("AsImmutable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsImmutable() got = %v, want %v", got, tt.want)
			}
		})
	}
}
