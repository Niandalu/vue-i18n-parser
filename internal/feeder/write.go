package feeder

import (
	"fmt"
	"github.com/niandalu/vue-i18n-parser/internal/reader"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

func Write(path string, code []byte, tf reader.TranslationFile) {
	isVue := strings.HasSuffix(path, ".vue")
	ioutil.WriteFile(path, content(isVue, code, tf), 0644)
}

func content(isVue bool, code []byte, tf reader.TranslationFile) []byte {
	newContent := marshal(tf)

	if isVue {
		matched := regexp.MustCompile(reader.VUE_BLOCK_RE).FindAllSubmatch(code, -1)
		prevContent := string(matched[len(matched)-1][2])
		newContent = []byte(strings.Replace(string(code), prevContent, string(newContent), 1))
	}

	return newContent
}

func marshal(tf reader.TranslationFile) []byte {
	s, err := yaml.Marshal(&tf.Content)

	if err != nil {
		log.Fatalf("Failed to marshal content of %s", tf.Path)
	}

	return []byte(
		fmt.Sprintf("\n# AUTO GENERATED\n# <d:%s-->\n---\n%s", tf.NextDigest, s),
	)
}
