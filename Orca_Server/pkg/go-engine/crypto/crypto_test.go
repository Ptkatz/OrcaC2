package crypto

import (
	"fmt"
	"testing"
)

func Test0001(t *testing.T) {
	fmt.Println(Algo())
	if !TestAllSum() {
		t.Error("TestSum fail")
	}
}
