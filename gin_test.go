package goplus

import (
	"fmt"
	"testing"
)

type quickStartTest struct {
}

func (q *quickStartTest) Hello() {
	fmt.Println("hello")
}

func (q *quickStartTest) IdempotentFunc() []string {
	return []string{"Hello"}
}

func TestCheckIdempotent(t *testing.T) {
	q := &quickStartTest{}
	ret := IsCallRepeated(q, "Hello", func() bool {
		return true
	})
	t.Log(ret)
}
