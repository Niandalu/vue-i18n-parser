package digest

import (
	"fmt"
	"testing"
)

func TestFind(t *testing.T) {
	all := `
<d:abcd-->
Hello World
`
	s := Find([]byte(all))

	if s != "abcd" {
		t.Fail()
	}
}

func TestGenerate(t *testing.T) {
	dummy := make(map[string]string)
	dummy["windows"] = "sucks"

	s := Generate(dummy)

	fmt.Printf("%v", s)
}
