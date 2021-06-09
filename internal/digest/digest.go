package digest

import (
	"crypto/sha1"
	"encoding/base64"
	"gopkg.in/yaml.v3"
	"log"
	"regexp"
)

const DIGEST_RE = `<d:(.*?)-->`

func Find(source []byte) string {
	matched := regexp.MustCompile(DIGEST_RE).FindSubmatch(source)

	if matched == nil {
		return ""
	}

	return string(matched[1])
}

func Generate(content interface{}) string {
	d, err := yaml.Marshal(&content)

	if err != nil {
		log.Fatalf("Failed to hash, %v", err)
	}

	hasher := sha1.New()
	hasher.Write(d)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
