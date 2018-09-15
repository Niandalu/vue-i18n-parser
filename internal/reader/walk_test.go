package reader

import (
	"testing"
)

func TestWalk(t *testing.T) {
	files := Walk("./testdata")

	if len(files) != 2 {
		t.Fail()
	}
}
