package constants

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestOutJsonPath(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "included {}",
			args: args{
				query: "{.something}",
			},
			want: "-ojsonpath={.something}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, OutJsonPath(tt.args.query))
		})
	}
}
