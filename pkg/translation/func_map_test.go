package translation

import "testing"

func TestStringer(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "struct",
			args: args{v: struct {
				Field string
			}{Field: "abcd"}},
			want: "abcd",
		},
		{
			name: "int",
			args: args{v: 33},
			want: "33",
		},
		{
			name: "float",
			args: args{v: 1.23},
			want: "1.23",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Stringer(tt.args.v); got != tt.want {
				t.Errorf("Stringer() = %v, want %v", got, tt.want)
			}
		})
	}
}
