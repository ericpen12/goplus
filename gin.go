package goplus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"time"
)

type register struct {
	g       *gin.Engine
	gCtx    *gin.Context
	service interface{}
	method  string
}

var r *register

func RegisterHandler(e *gin.Engine, serviceName string, service interface{}) {
	r = &register{
		g:       e,
		service: service,
	}
	r.g.POST(fmt.Sprintf("/%s/quick-start/do", serviceName), r.handler)
}

func (r *register) handler(c *gin.Context) {
	r.gCtx = c
	r.method = c.Query("method")
	if r.method == "" {
		response(nil, fmt.Errorf("method 不能为空"))
		return
	}
	defer func() {
		if err := recover(); err != nil {
			response(nil, err.(error))
		}
	}()
	err := r.checkRepeated(r.checkRepeatedByTime)
	if err != nil {
		response(nil, err)
		return
	}
	response(r.callByFuncName())
}

func (r *register) callByFuncName() (interface{}, error) {
	fn := reflect.ValueOf(r.service).MethodByName(r.method)
	if fn.Kind() != reflect.Func {
		return nil, fmt.Errorf("method is not exist: %s", r.method)
	}
	params, err := r.parseParams(fn.Type())
	if err != nil {
		return nil, fmt.Errorf("parse params error: %s", err)
	}
	return getCallResponse(fn.Call(params))
}

func getCallResponse(list []reflect.Value) (interface{}, error) {
	var (
		result []interface{}
		err    error
	)
	for _, item := range list {
		if v, ok := item.Interface().(error); ok {
			err = v
		} else {
			result = append(result, item.Interface())
		}
	}
	if len(result) == 1 {
		return result[0], err
	}
	return result, err
}

func (r *register) parseParams(fnType reflect.Type) ([]reflect.Value, error) {
	result := make([]reflect.Value, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		in := fnType.In(i)

		if in.Elem().AssignableTo(reflect.TypeOf(gin.Context{})) {
			result[i] = reflect.ValueOf(r.gCtx)
		} else {
			tv := reflect.New(in.Elem())
			err := r.gCtx.ShouldBindJSON(tv.Interface())
			if err != nil {
				return nil, err
			}
			result[i] = reflect.ValueOf(tv.Interface())
		}

	}
	return result, nil
}

type responseModel struct {
	Code int
	Data interface{}
	Msg  string
}

func response(data interface{}, err error) {
	resp := responseModel{
		Data: data,
	}
	if err != nil {
		resp.Msg = err.Error()
		resp.Code = -1
	}
	r.gCtx.JSON(200, resp)
}

type CallRepeat interface {
	CheckList() []string
}

func (r *register) checkRepeated(isRepeated func() bool) error {
	c, ok := r.service.(CallRepeat)
	if !ok {
		return nil
	}
	for _, name := range c.CheckList() {
		if name == r.method && isRepeated() {
			return fmt.Errorf("请勿重复调用")
		}
	}
	return nil
}

func (r *register) checkRepeatedByTime() bool {
	return r.gCtx.Query("time") != time.Now().Format("200601021504")
}
