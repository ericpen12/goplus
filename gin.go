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
	r.callByFuncName()
}

func (r *register) callByFuncName() {
	fn := reflect.ValueOf(r.service).MethodByName(r.method)
	if fn.Kind() != reflect.Func {
		response(nil, fmt.Errorf("method is not exist: %s", r.method))
		return
	}
	params, err := r.parseParams(fn.Type())
	if err != nil {
		response(nil, fmt.Errorf("parse params error: %s", err))
		return
	}
	fn.Call(params)
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
