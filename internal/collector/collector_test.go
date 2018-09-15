package collector

import (
	"github.com/niandalu/vue-i18n-parser/internal/reader"
	"testing"
)

func TestPrepareBody(t *testing.T) {
	f := reader.TranslationFile{
		Content: reader.Translation{
			"cn": reader.KV{
				"shop": "宠物商店",
				"pet": map[string]string{
					"dog": "狗",
					"cat": "猫",
				},
			},
			"en": reader.KV{
				"shop": "Pet Shop",
				"pet": map[string]string{
					"dog": "Dooog",
					"cat": "Caaat",
				},
			},
		},

		PrevDigest: "prevDigest",
		NextDigest: "nextDigest",
		Path:       "/tmp/test.yaml",
	}

	languages := []string{"cn", "en"}

	records := prepareBody(f, languages)

	if len(records) != 3 {
		t.Fail()
	}
}
