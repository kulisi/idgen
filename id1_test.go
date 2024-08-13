package idgen

import (
	"fmt"
	"testing"
)

func TestId4(t *testing.T) {
	opts := DefaultOptions(63)
	opts.BaseTime = 1582136402000
	idGen, err := NewIdGenerator(opts)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for i := 0; i < 100; i++ {
		fmt.Println(idGen.NewID())
	}
}
