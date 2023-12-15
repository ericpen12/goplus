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
		if IsCallRepeated(service, method, func() bool {
			return c.Query("time") == time.Now().Format("200601021504")
		}) {
			c.JSON(200, "请勿重复调用")
			return
		}
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
func IsCallRepeated(service interface{}, method string, isIdempotent func() bool) bool {
	fn := reflect.ValueOf(service).MethodByName("IdempotentFunc")
	if fn.Kind() != reflect.Func {
		return false
	}
	result := fn.Call([]reflect.Value{})
	if len(result) <= 0 {
		return false
	}
	returnData := result[0]
	if returnData.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < returnData.Len(); i++ {
		if returnData.Index(i).String() == method {
			return !isIdempotent()
		}
	}
	return false
}
