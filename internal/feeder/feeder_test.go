package feeder

import (
	"testing"
)

func TestGroupByFile(t *testing.T) {
	candidates := groupByFile("testdata/example.csv")

	switch r := candidates["testdata/main.vue"].Content["en"]["pets"].(type) {
	case map[string]interface{}:
		if r["cat"] != "Cat" {
			t.Fail()
		}
	default:
		t.Logf("malformed %v", r)
		t.Fail()
	}
}
