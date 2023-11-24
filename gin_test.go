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

func (q *quickStartTest) IdempotentFunc(method string) bool {
	return map[string]bool{
		"Hello": false,
	}[method]
}

func TestCheckIdempotent(t *testing.T) {
	q := &quickStartTest{}
	ret := IsCallRepeated(q, "Hello", func() bool {
		return false
	})
	t.Log(ret)
}
