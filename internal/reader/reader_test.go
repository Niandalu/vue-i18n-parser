package reader

import (
	"fmt"
	"path/filepath"
	"testing"
)

func validateTranslation(t *testing.T, path string) {
	realpath, err := filepath.Abs(filepath.Join("testdata", path))

	if err != nil {
		t.Fail()
	}

	translation := Format(realpath)

	enTranslation := translation.Content["en"]["pets"].(map[interface{}]interface{})
	sanitized := make(map[string]string)

	for k, v := range enTranslation {
		sanitized[fmt.Sprintf("%v", k)] = fmt.Sprintf("%v", v)
	}

	hasContent := sanitized["dog"] == "Dog"

	nextDigestGenerated := translation.NextDigest != translation.PrevDigest
	pathReserved := translation.Path == realpath

	if !hasContent || !nextDigestGenerated || !pathReserved {
		t.Fail()
	}
}

func TestReadYaml(t *testing.T) {
	validateTranslation(t, "main.yml")
}

func TestReadVue(t *testing.T) {
	validateTranslation(t, "main.vue")
}
