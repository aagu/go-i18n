package translation

import (
	"bytes"
	lang "golang.org/x/text/language"
	"reflect"
	"text/template"
)

type TemplateData map[string]interface{}

type Message struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	// fallback template
	tmpl *template.Template
}

func (m Message) String() string {
	return m.Translate(DefaultLocaleManager().locale())
}

func (m Message) Translate(l lang.Tag) string {
	return m.localeStringWithFallback(l)
}

func (m Message) Format(v interface{}) string {
	return m.FormatTranslate(DefaultLocaleManager().locale(), v)
}

func (m Message) FormatTranslate(l lang.Tag, v interface{}) string {
	tmpl := m.localeTemplateWithFallback(l)
	var buf bytes.Buffer
	v = propagateTranslate(l, v)
	err := tmpl.Execute(&buf, v)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// propagateTranslate translate all fields on which type is Message inside v
func propagateTranslate(l lang.Tag, v interface{}) interface{} {
	typ := reflect.TypeOf(v)
	switch typ.Kind() {
	case reflect.Struct:
		if typ.PkgPath() == "github.com/aagu/go-i18n/pkg/translation" && typ.Name() == "Message" {
			message := v.(Message)
			return message.Translate(l)
		}
	case reflect.Map:
		value := reflect.ValueOf(v)
		for _, k := range value.MapKeys() {
			value.SetMapIndex(k, reflect.ValueOf(propagateTranslate(l, value.MapIndex(k).Interface())))
		}
	case reflect.Ptr:
		return propagateTranslate(l, reflect.ValueOf(v).Elem().Interface())
	}
	return v
}

func (m Message) localeStringWithFallback(l lang.Tag) string {
	if s, ok := defaultLanguageManager().localizedString(l, m); ok {
		return s
	} else {
		// fallback
		return m.Text
	}
}

func (m *Message) localeTemplateWithFallback(l lang.Tag) *template.Template {
	if tmpl, ok := defaultLanguageManager().localizedTemplate(l, *m); ok {
		return tmpl
	} else {
		if m.tmpl == nil {
			parse, err := template.New(m.ID).Funcs(templateFuncMap).Parse(m.Text)
			if err != nil {
				panic(err)
			}
			m.tmpl = parse
		}
		return m.tmpl
	}
}
