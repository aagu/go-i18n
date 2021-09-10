package translation

import (
	"errors"
	lang "golang.org/x/text/language"
	"reflect"
	"testing"
	"text/template"
)

func Test_propagateTranslate(t *testing.T) {
	type args struct {
		l lang.Tag
		v interface{}
	}
	var iface interface{}
	iface = errors.New("test")
	tests := []struct {
		name     string
		args     args
		want     interface{}
		wantFail bool
	}{
		{
			name: "struct",
			args: args{
				l: lang.SimplifiedChinese,
				v: iface,
			},
			want:     iface,
			wantFail: true,
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
			if got := propagateTranslate(tt.args.l, tt.args.v); !reflect.DeepEqual(got, tt.want) != tt.wantFail {
				t.Errorf("propagateTranslate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_FormatTranslate(t *testing.T) {
	type fields struct {
		ID   string
		Text string
		tmpl *template.Template
	}
	type KeyPair struct {
		Key   string
		Value string
	}
	type args struct {
		l lang.Tag
		v interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "struct",
			fields: fields{
				ID:   "test",
				Text: "{{.}}",
			},
			args: args{
				l: lang.SimplifiedChinese,
				v: errors.New("error occur"),
			},
			want: "{error occur}",
		},
		{
			name: "struct filed",
			fields: fields{
				ID:   "test",
				Text: "Key:{{.Key}}, Value:{{.Value}}",
			},
			args: args{
				l: lang.SimplifiedChinese,
				v: KeyPair{Key: "test_k", Value: "test_val"},
			},
			want: "Key:test_k, Value:test_val",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Message{
				ID:   tt.fields.ID,
				Text: tt.fields.Text,
				tmpl: tt.fields.tmpl,
			}
			if got := m.FormatTranslate(tt.args.l, tt.args.v); got != tt.want {
				t.Errorf("FormatTranslate() = %v, want %v", got, tt.want)
			}
		})
	}
}
