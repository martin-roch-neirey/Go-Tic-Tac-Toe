// Copyright (c) 2022 Haute école d'ingénierie et d'architecture de Fribourg
// SPDX-License-Identifier: Apache-2.0
// Author:  William Margueron & Martin Roch-Neirey

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const translationsFile = "resources/translations.json"

var translationsMap = make(map[string]interface{})
var languageLoaded = "fr-FR" // default language
var loaded = false

var fallbackCache = make(map[string]string) // cache of fallback translations to query each of them only one time

// GetTranslation returns translation content of given translation id, in given lang
func GetTranslation(id, lang string) string {
	if !loaded || languageLoaded != lang {
		loadTranslations(lang)
		loaded = true
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
				fallbackCache[id] = GetTranslation(id, fmt.Sprint(fallback)) // recursive appeal
				return fallbackCache[id]
			} else {
				return translationCache
			}

		}
	}

	return fmt.Sprint(translation)
}

// loadTranslations loads translations of given lang in translationsMap object
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
