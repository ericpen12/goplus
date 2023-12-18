package goplus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
)

type quickStartTest struct {
}

func (q *quickStartTest) Hello() {
	fmt.Println("hello")
}

func (q *quickStartTest) CheckList() []string {
	return []string{"Hello"}
}

func Test_checkRepeated(t *testing.T) {
	r = &register{
		gCtx:    &gin.Context{},
		service: &quickStartTest{},
		method:  "Hello",
	}
	ret := r.checkRepeated(r.checkRepeatedByTime)
	t.Log(ret)
}
