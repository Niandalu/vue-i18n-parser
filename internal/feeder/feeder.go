package feeder

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
	"github.com/niandalu/vue-i18n-parser/internal/digest"
	"github.com/niandalu/vue-i18n-parser/internal/reader"
	"github.com/niandalu/vue-i18n-parser/internal/utils"
)

const RESERVED = 4

type Options struct {
	PatchMode bool
	Indent    int
}

func Run(projectRoot string, file string, options Options) {
	for path, tf := range groupByFile(file, options) {
		realpath := filepath.Join(projectRoot, path)

		code, err := ioutil.ReadFile(realpath)

		if err != nil {
			log.Fatalf("Cannot open %s, %v", realpath, err)
		}

		currentDigest := digest.Find(code)
		if !options.PatchMode && currentDigest != "" && currentDigest != tf.PrevDigest {
			log.Fatalf("%s translation has been mutated since last change", realpath)
		}

		Write(realpath, code, tf, options)
	}
}

func groupByFile(sourceFile string, options Options) map[string]reader.TranslationFile {
	langs, records := read(sourceFile)
	candiates := map[string]reader.TranslationFile{}

	for _, r := range records {
		prevDigest := r[1]
		path := r[2]
		key := r[3]

		nextContent := reader.Translation{}
		_, hasSomething := candiates[path]
		if hasSomething {
			nextContent = candiates[path].Content
		} else if options.PatchMode {
			nextContent = reader.Format(path).Content
		}

		tmp := reader.Translation{}
		for i, lang := range langs {
			tmp[lang] = convertToNested(key, r[i+RESERVED])
		}

		if err := mergo.Merge(&nextContent, tmp, mergo.WithOverride); err != nil {
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

func read(sourceFile string) ([]string, [][]string) {
	file, err := os.Open(sourceFile)

	if err != nil {
		log.Fatalf("Cannot open file %s, %v", sourceFile, err)
	}
	defer file.Close()

	if filepath.Ext(sourceFile) == ".csv" {
		return readCSV(file)
	}

	return readJSON(file)
}

func readJSON(file *os.File) ([]string, [][]string) {
	byteValue, e := ioutil.ReadAll(file)
	if e != nil {
		log.Fatalf("Invalid JSON file %s, %v", file.Name(), e)
	}

	var result map[string]reader.Translation
	json.Unmarshal([]byte(byteValue), &result)

	// collect available languages
	langs := []string{}
	for _, chunk := range result {
		for lang := range chunk {
			if utils.Contains(langs, lang) {
				continue
			}
			langs = append(langs, lang)
		}
	}

	records := [][]string{}
	for path, chunk := range result {
		createBase := func() []string {
			base := []string{
				"x",  // changed
				"",   // prevDigest
				path, // path
				"",   // key
			}
			for range langs {
				base = append(base, "")
			}
			return base
		}

		memo := map[string][]string{}
		for lang, translations := range chunk {

			langIdx := 0
			for i := 0; i < len(langs); i++ {
				if langs[i] == lang {
					langIdx = i
					break
				}
			}

			for k, v := range translations {
				one := []string{}
				if lastOne, ok := memo[k]; ok {
					one = lastOne
				} else {
					one = createBase()
				}
				one[3] = k
				one[langIdx+4] = v.(string)
				memo[k] = one
			}
		}

		for _, one := range memo {
			records = append(records, one)
		}
	}

	return langs, records
}

func readCSV(file *os.File) ([]string, [][]string) {
	r := csv.NewReader(file)
	rows, err := r.ReadAll()

	if err != nil {
		log.Fatalf("Please make sure %s is a valid csv file, %v", file.Name(), err)
	}

	if len(rows) <= 1 {
		log.Fatalf("%s seems empty", file.Name())
	}

	langs := rows[0][RESERVED:len(rows[0])]
	records := rows[1:]

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
