package translation

import (
	lang "golang.org/x/text/language"
	"sync"
	"text/template"
)

type language struct {
	l          lang.Tag
	translated map[string]string
	templated  map[string]*template.Template
}

type languageManager struct {
	mapping map[lang.Tag]language
}

var managerInstance *languageManager = nil
var languageManagerInitiator sync.Once

func defaultLanguageManager() *languageManager {
	languageManagerInitiator.Do(func() {
		managerInstance = &languageManager{mapping: make(map[lang.Tag]language)}
	})
	return managerInstance
}

func registerLang(l lang.Tag, trans []Message) {
	lang2use := language{l: l, translated: map[string]string{}, templated: map[string]*template.Template{}}
	for idx := range trans {
		lang2use.translated[trans[idx].ID] = trans[idx].Text
	}
	defaultLanguageManager().mapping[l] = lang2use
}

func (lm languageManager) localizedString(l lang.Tag, t Message) (str string, ok bool) {
	if _, ok = lm.mapping[l]; !ok {
		return "", false
	}
	str, ok = lm.mapping[l].translated[t.ID]
	return
}

func (lm *languageManager) localizedTemplate(l lang.Tag, t Message) (tmpl *template.Template, ok bool) {
	if _, ok = lm.mapping[l]; !ok {
		return nil, false
	}
	tmpl, ok = lm.mapping[l].templated[t.ID]
	if !ok {
		str, find := lm.localizedString(l, t)
		if !find {
			return nil, false
		}
		parse, err := template.New(t.ID).Parse(str)
		if err != nil {
			return nil, false
		}
		lm.mapping[l].templated[t.ID] = parse
		tmpl = parse
		ok = true
	}
	return
}
