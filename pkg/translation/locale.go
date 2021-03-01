package translation

import (
	lang "golang.org/x/text/language"
	"sync"
)

type LocaleManager struct {
}

var defaultLocaleManager *LocaleManager
var localeManagerInitiator sync.Once
var defaultLocale = lang.English

func DefaultLocaleManager() *LocaleManager {
	localeManagerInitiator.Do(func() {
		defaultLocaleManager = &LocaleManager{}
	})
	return defaultLocaleManager
}

func SetDefaultLocale(l lang.Tag) {
	defaultLocale = l
}

// get current Locale
func (l LocaleManager) locale() lang.Tag {
	return l.defaultLocale()
}

func (l LocaleManager) defaultLocale() lang.Tag {
	return defaultLocale
}

func parseLocale(localeString string) lang.Tag {
	parsed, err := lang.Parse(localeString)
	if err != nil {
		panic(err)
	}
	return parsed
}
