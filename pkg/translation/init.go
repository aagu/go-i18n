package translation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func LoadTranslations(i18nPath string) {
	stat, err := os.Stat(i18nPath)
	if err != nil {
		panic(fmt.Sprintf("failed to open translations dir, %v", err))
	}
	if !stat.IsDir() {
		panic(fmt.Sprintf("failed to open translations dir, path %s is not directory", i18nPath))
	}
	dir, _ := ioutil.ReadDir(i18nPath)
	for _, translationFile := range dir {
		if filepath.Ext(translationFile.Name()) != ".json" {
			continue
		}
		readFile, err := ioutil.ReadFile(filepath.Join(i18nPath, translationFile.Name()))
		if err != nil {
			panic(fmt.Sprintf("failed to open translation file %s, %v", translationFile.Name(), err))
		}
		var trans []Message
		err = json.Unmarshal(readFile, &trans)
		if err != nil {
			panic(fmt.Sprintf("failed to open translation file %s, %v", translationFile.Name(), err))
		}
		registerLang(parseLocale(strings.TrimSuffix(translationFile.Name(), ".json")), trans)
	}
}
