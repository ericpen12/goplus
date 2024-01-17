package main

import (
	"errors"
	"github.com/gin-gonic/gin"
)
import . "github.com/ericpen12/goplus"

func main() {
	r := gin.Default()
	// 注册服务
	RegisterHandler(r, "srv/quickstart", &service{})
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
