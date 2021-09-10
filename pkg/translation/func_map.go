package translation

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

var templateFuncMap = make(template.FuncMap)

func RegisterTemplateFunc(name string, fn interface{}) {
	templateFuncMap[name] = fn
}

func Stringer(v interface{}) string {
	value := reflect.ValueOf(v)
	switch value.Type().Kind() {
	case reflect.Struct:
		sprintf := fmt.Sprintf("%s", v)
		sprintf = strings.TrimLeft(sprintf, "{")
		sprintf = strings.TrimRight(sprintf, "}")
		return sprintf
	default:
		return fmt.Sprintf("%s", v)
	}
}
