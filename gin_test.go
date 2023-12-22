package goplus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
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

func (q *quickStartTest) Get() ([]string, error) {
	return []string{"Hello"}, fmt.Errorf("ok")
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

func Test_getCallResponse(t *testing.T) {
	fn := reflect.ValueOf(&quickStartTest{}).MethodByName("Get")
	data, err := getCallResponse(fn.Call([]reflect.Value{}))
	t.Log(data, err)
}
