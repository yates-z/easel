package uuid

import (
	"fmt"
	"testing"
)

func TestUUID(t *testing.T) {
	for i := 0; i < 20; i++ {
		fmt.Println(NewV7String())
	}
}
