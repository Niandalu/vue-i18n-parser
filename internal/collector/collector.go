package collector

import (
	"log"

	"github.com/niandalu/vue-i18n-parser/internal/reader"
	"github.com/niandalu/vue-i18n-parser/internal/tree"
)

func Run(files []reader.TranslationFile, languages []string, diffOnly bool) [][]string {
	if len(files) == 0 {
		return [][]string{}
	}

	var content [][]string
	for _, f := range files {
		if diffOnly && f.PrevDigest == f.NextDigest {
			continue
		}

		log.Printf("Parsing %s...", f.Path)
		content = append(content, prepareBody(f, languages)...)
	}

	return append([][]string{prepareHeader(languages)}, content...)
}

func prepareHeader(languages []string) []string {
	return append([]string{"changed", "digest", "file", "key"}, languages...)
}

func prepareBody(f reader.TranslationFile, languages []string) [][]string {
	translated := make(map[string]map[string]string)

	for lang, mapping := range f.Content {
		contracted := tree.Flatten(mapping)

		for k, v := range contracted {
			if _, ok := translated[k]; !ok {
				translated[k] = map[string]string{}
			}

			translated[k][lang] = v
		}
	}

	var result [][]string

	for k, v := range translated {
		var recordLanguages []string
		changed := "x"

		for _, lang := range languages {
			isEmpty := len(v[lang]) == 0

			if isEmpty || (len(recordLanguages) > 0 && changed == "x" && v[lang] == recordLanguages[len(recordLanguages)-1]) {
				changed = "o"
			}

			recordLanguages = append(recordLanguages, v[lang])
		}

		base := []string{changed, f.PrevDigest, f.Path, k}
		recordInString := append(base, recordLanguages...)

		result = append(result, recordInString)
	}

	return result
}
