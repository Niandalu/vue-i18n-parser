package feeder

import (
	"github.com/niandalu/vue-i18n-parser/internal/reader"
	"strings"
	"testing"
)

func TestWrite(t *testing.T) {
	sourceCode := `
<script>
</script>

<i18n>
haha
</i18n>
`
	tf := reader.TranslationFile{
		Content: reader.Translation{
			"en": reader.KV{
				"a": "b",
			},
		},
		PrevDigest: "",
		NextDigest: "abc",
		Path:       "kk",
	}

	replaced := string(content(true, []byte(sourceCode), tf))

	for _, k := range []string{"<script>", "en", "abc"} {
		if !strings.Contains(replaced, k) {
			t.Fail()
		}
	}
}
