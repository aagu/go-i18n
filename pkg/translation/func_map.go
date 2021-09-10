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

var stringerFuncNames = []string{"Error", "String"}

func Stringer(v interface{}) string {
	value := reflect.ValueOf(v)
	switch value.Type().Kind() {
	case reflect.Struct:
		for ix := range stringerFuncNames {
			if stringerFunc, ok := value.Type().MethodByName(stringerFuncNames[ix]); ok {
				if canProduceString(value.Type(), stringerFuncNames[ix]) {
					return stringerFunc.Func.Call([]reflect.Value{})[0].Interface().(string)
				}
			}
		}
		sprintf := fmt.Sprintf("%v", v)
		sprintf = strings.TrimLeft(sprintf, "{")
		sprintf = strings.TrimRight(sprintf, "}")
		return sprintf
	default:
		return fmt.Sprintf("%s", v)
	}
}

func canProduceString(val reflect.Type, funcName string) bool {
	kind := val.Kind() == reflect.Func
	name := val.Name() == funcName
	zeroIn := val.NumIn() == 0
	oneOut := val.NumOut() == 1
	outString := val.Out(0).Kind() == reflect.String
	return kind && name && zeroIn && oneOut && outString
}
