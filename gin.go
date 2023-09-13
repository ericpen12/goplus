package goplus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"time"
)

func RegisterHandler(e *gin.Engine, serviceName string, service interface{}) {
	e.POST(fmt.Sprintf("/%s/quick-start/do", serviceName), func(c *gin.Context) {
		method := c.Query("method")
		if method == "" {
			c.JSON(200, "method is empty")
			return
		}
		defer func() {
			if err := recover(); err != nil {
				c.JSON(200, err)
			}
		}()
		callByFuncName(c, service, method)
	})
}

func callByFuncName(c *gin.Context, service interface{}, method string) {
	fn := reflect.ValueOf(service).MethodByName(method)
	if fn.Kind() != reflect.Func {
		c.JSON(200, fmt.Sprintf("method is not exist: %s", method))
		return
	}
	params, err := parseParams(c, fn.Type())
	if err != nil {
		c.JSON(200, fmt.Sprintf("parse params error: %s", err))
		return
	}
	fn.Call(params)
}

func parseParams(c *gin.Context, fnType reflect.Type) ([]reflect.Value, error) {
	result := make([]reflect.Value, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		in := fnType.In(i)

		if in.Elem().AssignableTo(reflect.TypeOf(gin.Context{})) {
			result[i] = reflect.ValueOf(c)
		} else {
			tv := reflect.New(in.Elem())
			err := c.ShouldBindJSON(tv.Interface())
			if err != nil {
				return nil, err
			}
			result[i] = reflect.ValueOf(tv.Interface())
		}

	}
	return result, nil
}

func CheckSafe(c *gin.Context) {
	if !isSafe(c) {
		panic("current environment is unsafe")
	}
}

func isSafe(c *gin.Context) bool {
	timeStr := time.Now().Format("200601021504")
	inputTime := c.Query("time")
	if inputTime != timeStr {
		return false
	}
	return true
}

type quickStart struct{}

type PingReq struct {
	Msg string
}

func (m *quickStart) Ping(c *gin.Context, input *PingReq) {
	c.JSON(200, input.Msg)
}
