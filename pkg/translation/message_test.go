package translation

import (
	"errors"
	lang "golang.org/x/text/language"
	"reflect"
	"testing"
)

func Test_propagateTranslate(t *testing.T) {
	type args struct {
		l lang.Tag
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "struct",
			args: args{
				l: lang.SimplifiedChinese,
				v: errors.New("test"),
			},
			want: "test",
		},
		{
			name: "message",
			args: args{
				l: lang.SimplifiedChinese,
				v: Message{
					ID:   "test",
					Text: "test",
				},
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := propagateTranslate(tt.args.l, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("propagateTranslate() = %v, want %v", got, tt.want)
			}
		})
	}
}
