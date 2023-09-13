package call

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
)

var globalService interface{}

func RegisterHandler(e *gin.Engine, uri string, service interface{}) {
	globalService = service
	e.POST(uri, handler)
}

func handler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			response(c, nil, err.(error))
		}
	}()
	ret, err := callByFuncName(c.Query("method"), bindParams(c))
	response(c, ret, err)
}

func callByFuncName(method string, bind bindOption) (interface{}, error) {
	if method == "" {
		return nil, fmt.Errorf("method 不能为空")
	}
	fn := reflect.ValueOf(globalService).MethodByName(method)
	if fn.Kind() != reflect.Func {
		return nil, fmt.Errorf("method is not exist: %s", method)
	}
	params, err := reflectParams(fn.Type(), bind)
	if err != nil {
		return nil, fmt.Errorf("parse params error: %s", err)
	}
	return formatResult(fn.Call(params))
}

func formatResult(list []reflect.Value) (interface{}, error) {
	var (
		result interface{}
		err    error
	)
	for _, item := range list {
		if item.Type().String() == "error" {
			err, _ = item.Interface().(error)
		} else {
			result = item.Interface()
		}
	}
	return result, err
}

type ginResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func response(ginCtx *gin.Context, data interface{}, err error) {
	resp := ginResponse{
		Data: data,
	}
	if err != nil {
		resp.Msg = err.Error()
		resp.Code = -1
	}
	ginCtx.JSON(200, resp)
}

type bindOption func(param any) error

func reflectParams(fnType reflect.Type, bind bindOption) ([]reflect.Value, error) {
	var result []reflect.Value
	if fnType.NumIn() == 0 {
		return result, nil
	}
	tv := reflect.New(fnType.In(0).Elem())
	if err := bind(tv.Interface()); err != nil {
		return result, err
	}
	result = append(result, reflect.ValueOf(tv.Interface()))
	return result, nil
}

func bindParams(ginCtx *gin.Context) bindOption {
	return func(param any) error {
		return ginCtx.ShouldBindJSON(param)
	}
}
