package sqlhelpers

import (
	"fmt"
	"testing"
)

func TestUpdateBuilder(tst *testing.T) {
	in := InQueryPlaceholder(3, 3)
	fmt.Println(in)
	return
}
