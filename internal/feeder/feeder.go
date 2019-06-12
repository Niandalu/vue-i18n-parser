package feeder

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
	"github.com/niandalu/vue-i18n-parser/internal/digest"
	"github.com/niandalu/vue-i18n-parser/internal/reader"
)

const RESERVED = 4

func Run(projectRoot string, csvFile string) {
	for path, tf := range groupByFile(csvFile) {
		realpath := filepath.Join(projectRoot, path)

		code, err := ioutil.ReadFile(realpath)

		if err != nil {
			log.Fatalf("Cannot open %s, %v", realpath, err)
		}

		currentDigest := digest.Find(code)
		if currentDigest != "" && currentDigest != tf.PrevDigest {
			log.Fatalf("%s translation has been mutated since last change", realpath)
		}

		Write(realpath, code, tf)
	}
}

func groupByFile(source string) map[string]reader.TranslationFile {
	langs, records := read(source)
	candiates := map[string]reader.TranslationFile{}

	for _, r := range records {
		prevDigest := r[1]
		path := r[2]
		key := r[3]

		nextContent := reader.Translation{}
		if _, ok := candiates[path]; ok {
			nextContent = candiates[path].Content
		}

		tmp := reader.Translation{}
		for i, lang := range langs {
			tmp[lang] = convertToNested(key, r[i+RESERVED])
		}

		if err := mergo.Merge(&nextContent, tmp); err != nil {
			log.Fatalf("Failed to merge %v", err)
		}

		candiates[path] = reader.TranslationFile{
			Content:    nextContent,
			PrevDigest: prevDigest,
			NextDigest: digest.Generate(nextContent),
			Path:       path,
		}
	}

	return candiates
}

func read(source string) ([]string, [][]string) {
	file, err := os.Open(source)

	if err != nil {
		log.Fatalf("Cannot open file %s, %v", source, err)
	}

	r := csv.NewReader(file)
	rows, err := r.ReadAll()

	if err != nil {
		log.Fatalf("Please make sure %s is a valid csv file, %v", source, err)
	}

	if len(rows) <= 1 {
		log.Fatalf("%s seems empty", source)
	}

	langs := rows[0][RESERVED:len(rows[0])]
	records := rows[1:len(rows)]

	return langs, records
}

func convertToNested(key string, value string) map[string]interface{} {
	nested := make(map[string]interface{})
	parts := strings.Split(key, ".")
	cached := nested

	for i, part := range parts {
		if i == len(parts)-1 {
			cached[part] = value
		} else {
			cached[part] = make(map[string]interface{})
			cached = cached[part].(map[string]interface{})
		}
	}

	return nested
}
