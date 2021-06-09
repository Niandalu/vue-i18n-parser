package reader

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"

	"github.com/niandalu/vue-i18n-parser/internal/digest"
	"gopkg.in/yaml.v3"
)

type KV map[string]interface{}

// map[Lang]KV
type Translation map[string]KV

type TranslationFile struct {
	Content    Translation
	PrevDigest string
	NextDigest string
	Path       string
}

const VUE_BLOCK_RE = `<i18n(.*?)>((.|\s)*?)</i18n>`

func Run(dir string, ignore string) []TranslationFile {
	var result []TranslationFile

	for _, path := range Walk(dir) {
		matched, _ := regexp.MatchString(ignore, path)

		if matched {
			continue
		}

		t := Format(path)
		if len(t.Content) > 0 {
			result = append(result, t)
		}
	}

	return result
}

func Format(path string) TranslationFile {
	source := extractSource(path)
	content := unmarshal(source)

	return TranslationFile{
		content,
		digest.Find(source),
		digest.Generate(content),
		path,
	}
}

func unmarshal(source []byte) Translation {
	var t Translation
	err := yaml.Unmarshal(source, &t)

	if err != nil {
		log.Fatalf("Syntax error: %v", err)
	}

	return t
}

func extractSource(path string) []byte {
	if filepath.Ext(path) == ".vue" {
		return extractVueSource(path)
	}

	return extractYmlSource(path)
}

func extractVueSource(path string) []byte {
	source := extractYmlSource(path)
	matched := regexp.MustCompile(VUE_BLOCK_RE).FindAllSubmatch(source, -1)

	if len(matched) == 0 {
		return []byte{}
	}
	return matched[len(matched)-1][2]
}

func extractYmlSource(path string) []byte {
	source, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatalf("Failed to read %s", path)
	}

	return source
}
