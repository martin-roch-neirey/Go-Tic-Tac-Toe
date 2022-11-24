package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const translationsFile = "resources/translations.json"

var translationsMap = make(map[string]interface{})
var languageLoaded = "fr-FR"
var loaded = false

var fallbackCache = make(map[string]string)

func GetTranslation(id, lang string) string {
	if !loaded || languageLoaded != lang {
		loadTranslations(lang) // Todo get here language chosen by player
		loaded = true
		fmt.Println("accessed ")
	}

	var translation, ok = translationsMap[id]
	if !ok {
		fallback := translationsMap["fallback"]
		if fallback == "NO_FALLBACK" {
			return "Missing Translation (" + id + ", " + lang + ")"
		} else {
			var translationCache, ok = fallbackCache[id]
			if !ok {
				loadTranslations(fmt.Sprint(fallback))
				loaded = false
				fallbackCache[id] = GetTranslation(id, fmt.Sprint(fallback))
				return fallbackCache[id]
			} else {
				return translationCache
			}

		}
	}

	return fmt.Sprint(translation)
}

func loadTranslations(lang string) {
	jsonFile, err := os.Open(translationsFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(jsonFile)

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println(err)
		return
	}

	if result[lang] == nil {
		fmt.Println("Missing Language (" + lang + ")")
		return
	}

	var languageMap interface{}
	languageMap = result[lang]
	translationsMap = languageMap.(map[string]interface{})
	languageLoaded = lang

}
