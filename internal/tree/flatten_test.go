package tree

import (
	"testing"
)

func TestFlatten(t *testing.T) {
	tree := map[string]interface{}{
		"parent": map[string]string{
			"son":      "Son",
			"daughter": "Daughter",
		},

		"uncle": "Uncle",
	}

	result := Flatten(tree)

	if result["parent.son"] != "Son" || result["parent.daughter"] != "Daughter" {
		t.Fail()
	}
}
