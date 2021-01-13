package gini

import (
	"errors"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Context struct {
	*gin.Context
	resultFunc ResultFunc
	translator ut.Translator
}

type CtxOption func(c *Context)

type CtxOptions []CtxOption

func CtxResultFunc(resultFunc ResultFunc) CtxOption {
	return func(c *Context) {
		c.resultFunc = resultFunc
	}
}

func CtxTranslator(translator ut.Translator) CtxOption {
	return func(c *Context) {
		c.translator = translator
	}
}

func NewContext(origin *gin.Context, options ...CtxOption) *Context {
	c := &Context{
		Context:    origin,
		resultFunc: DefaultResult,
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

func (c *Context) JsonSuccess(data interface{}) {
	c.JSON(200, c.resultFunc(SUCCESS, "success", data))
}

func (c *Context) JsonError(code int, msg string) {
	c.JSON(200, c.resultFunc(code, msg, []string{}))
}

func (c *Context) JsonErrorMsg(msg string) {
	c.JsonError(FAIL, msg)
}

func (c *Context) ShouldBind(obj interface{}) error {
	err := c.Context.ShouldBind(obj)
	if err != nil && c.translator != nil {
		err = errors.New(err.(validator.ValidationErrors)[0].Translate(c.translator))
	}
	return err
}
