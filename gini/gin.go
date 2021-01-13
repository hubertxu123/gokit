package gini

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerFunc func(ctx *Context)

func (fn HandlerFunc) Unwrap(options CtxOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fn(NewContext(ctx, options...))
	}
}

type HandlersChain []HandlerFunc

func (fns HandlersChain) Unwrap(options CtxOptions) (handlers gin.HandlersChain) {
	for _, fn := range fns {
		handlers = append(handlers, fn.Unwrap(options))
	}
	return
}

type RouteInfo struct {
	Method      string
	Path        string
	Handler     string
	HandlerFunc HandlerFunc
}

type RoutesInfo []RouteInfo

type Engine struct {
	*gin.Engine
	*RouterGroup
}

func New(options ...GroupOption) *Engine {
	engine := gin.New()
	return &Engine{
		Engine:      engine,
		RouterGroup: NewRouterGroup(&engine.RouterGroup, options...),
	}
}

func (e *Engine) Delims(left, right string) *Engine {
	e.Engine.Delims(left, right)
	return e
}

func (e *Engine) SecureJsonPrefix(prefix string) *Engine {
	e.Engine.SecureJsonPrefix(prefix)
	return e
}

func (e *Engine) NoRoute(handlers ...HandlerFunc) {
	e.Engine.NoRoute(UnwrapHandlersChain(handlers, e.ctxOptions)...)
}

func (e *Engine) NoMethod(handlers ...HandlerFunc) {
	e.Engine.NoMethod(UnwrapHandlersChain(handlers, e.ctxOptions)...)
}

func (e *Engine) Use(middleware ...HandlerFunc) IRoutes {
	e.Engine.Use(UnwrapHandlersChain(middleware, e.ctxOptions)...)
	return e
}

func (e *Engine) Routes() (routes RoutesInfo) {
	return e.routeInfosWrap(e.Engine.Routes())
}

func (e *Engine) routeInfoWrap(ri gin.RouteInfo) RouteInfo {
	return RouteInfo{
		Method:      ri.Method,
		Path:        ri.Path,
		Handler:     ri.Handler,
		HandlerFunc: WrapHandlerFunc(ri.HandlerFunc),
	}
}

func (e *Engine) routeInfosWrap(ris gin.RoutesInfo) (ret RoutesInfo) {
	for _, v := range ris {
		ret = append(ret, e.routeInfoWrap(v))
	}
	return
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.Engine.ServeHTTP(w, req)
}

func (e *Engine) HandleContext(c *Context) {
	e.Engine.HandleContext(c.Context)
}
