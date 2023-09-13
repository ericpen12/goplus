package call

import (
	"errors"
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

func (q *quickStartTest) Print(s string) {
	fmt.Println(s)
}

func (q *quickStartTest) CheckList() []string {
	return []string{"Hello"}
}

func (q *quickStartTest) Get() ([]string, error) {
	return []string{"Hello"}, fmt.Errorf("ok")
}

func (q *quickStartTest) Get2() ([]string, error) {
	return []string{"Hello"}, nil
}

func Test_getCallResponse(t *testing.T) {
	fn := reflect.ValueOf(&quickStartTest{}).MethodByName("Get")
	data, err := formatResult(fn.Call([]reflect.Value{}))
	t.Log(data, err)

	fn = reflect.ValueOf(&quickStartTest{}).MethodByName("Get2")
	data, err = formatResult(fn.Call([]reflect.Value{}))
	t.Log(data, err)

	fn = reflect.ValueOf(&quickStartTest{}).MethodByName("Hello")
	data, err = formatResult(fn.Call([]reflect.Value{}))
	t.Log(data, err)
}

func TestServer(t *testing.T) {
	r := gin.Default()
	// 注册服务
	RegisterHandler(r, "quickstart", &service{})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

type service struct {
}

func (q *service) Hello() string {
	return "hello"
}

func (q *service) Hello1() (string, error) {
	return "hello", errors.New("error")
}

func (q *service) Hellow() (string, error) {
	return "hello", nil
}

type InputReq struct {
	Data interface{} `json:"data"`
}

func (q *service) Inputss(input *InputReq) string {
	fmt.Println(input.Data)
	return fmt.Sprintf("%v", input.Data)
}
