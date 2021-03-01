package main

import (
	"fmt"
	"github.com/aagu/go-i18n/pkg/translation"
	lang "golang.org/x/text/language"
)

var (
	jan      = translation.Message{ID: "January", Text: "January"}
	greeting = translation.Message{ID: "Hello", Text: "Hello"}
	format   = translation.Message{ID: "DayOfMonth", Text: "The {{.Day}}(th) day of {{.Month}}"}
	ways     = translation.Message{ID: "TwoWay", Text: "One way is to {{.One}}, the other is to {{.Other}}"}
	roma     = translation.Message{ID: "roma", Text: "roma"}
	paris    = translation.Message{ID: "paris", Text: "paris"}
	text     = translation.Message{ID: "text", Text: "{{.}}"}
)

func main() {
	translation.SetDefaultLocale(lang.English)
	translation.LoadTranslations(`C:\Users\aagui\IdeaProjects\go-i18n\i18n`)
	fmt.Println(greeting.Translate(lang.SimplifiedChinese))
	fmt.Println(format.FormatTranslate(lang.SimplifiedChinese, translation.TemplateData{"Day": 2, "Month": jan}))
	fmt.Println(ways.Format(translation.TemplateData{"One": roma, "Other": paris}))
	fmt.Println(text.FormatTranslate(lang.SimplifiedChinese, &greeting))
}
