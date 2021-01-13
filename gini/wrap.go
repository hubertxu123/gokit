package gini

import (
	"github.com/gin-gonic/gin"
)

// gin improved

func WrapHandlerFunc(fn gin.HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		fn(ctx.Context)
	}
}

func WrapHandlersChain(handlers gin.HandlersChain) (ret HandlersChain) {
	for _, h := range handlers {
		ret = append(ret, WrapHandlerFunc(h))
	}
	return
}

func UnwrapHandlersChain(handlers HandlersChain, options CtxOptions) gin.HandlersChain {
	return handlers.Unwrap(options)
}
