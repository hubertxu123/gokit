package gini

import (
	"github.com/gin-gonic/gin"
	"gokit/vali"
	"net/http"
)

type IRouter interface {
	IRoutes
	Group(string, ...HandlerFunc) *RouterGroup
}

type IRoutes interface {
	Use(...HandlerFunc) IRoutes

	Handle(string, string, ...HandlerFunc) IRoutes
	Any(string, ...HandlerFunc) IRoutes
	GET(string, ...HandlerFunc) IRoutes
	POST(string, ...HandlerFunc) IRoutes
	DELETE(string, ...HandlerFunc) IRoutes
	PATCH(string, ...HandlerFunc) IRoutes
	PUT(string, ...HandlerFunc) IRoutes
	OPTIONS(string, ...HandlerFunc) IRoutes
	HEAD(string, ...HandlerFunc) IRoutes

	StaticFile(string, string) IRoutes
	Static(string, string) IRoutes
	StaticFS(string, http.FileSystem) IRoutes
}

type RouterGroup struct {
	origin     *gin.RouterGroup
	Handlers   HandlersChain
	ctxOptions CtxOptions
	validate   *vali.Validate
}

func NewRouterGroup(origin *gin.RouterGroup, options ...GroupOption) *RouterGroup {
	g := &RouterGroup{
		origin: origin,
	}
	for _, opt := range options {
		opt(g)
	}
	return g
}

func (g *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
	g.origin.Use(UnwrapHandlersChain(middleware, g.ctxOptions)...)
	g.Handlers = WrapHandlersChain(g.origin.Handlers)
	return g
}

func (g *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		origin:     g.origin.Group(relativePath, UnwrapHandlersChain(handlers, g.ctxOptions)...),
		Handlers:   WrapHandlersChain(g.origin.Handlers),
		ctxOptions: g.ctxOptions,
		validate:   g.validate,
	}
}

func (g *RouterGroup) BasePath() string {
	return g.origin.BasePath()
}

func (g *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) IRoutes {
	g.origin.Handle(httpMethod, relativePath, UnwrapHandlersChain(handlers, g.ctxOptions)...)
	return g
}

func (g *RouterGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
	return g.handle(httpMethod, relativePath, handlers)
}

func (g *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) IRoutes {
	return g.handle(http.MethodPost, relativePath, handlers)
}

func (g *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	return g.handle(http.MethodGet, relativePath, handlers)
}

func (g *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes {
	return g.handle(http.MethodDelete, relativePath, handlers)
}

func (g *RouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) IRoutes {
	return g.handle(http.MethodPatch, relativePath, handlers)
}

func (g *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) IRoutes {
	return g.handle(http.MethodPut, relativePath, handlers)
}

func (g *RouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) IRoutes {
	return g.handle(http.MethodOptions, relativePath, handlers)
}

func (g *RouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) IRoutes {
	return g.handle(http.MethodHead, relativePath, handlers)
}

func (g *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) IRoutes {
	g.origin.Any(relativePath, UnwrapHandlersChain(handlers, g.ctxOptions)...)
	return g
}

func (g *RouterGroup) StaticFile(relativePath, filepath string) IRoutes {
	g.origin.StaticFile(relativePath, filepath)
	return g
}

func (g *RouterGroup) Static(relativePath, root string) IRoutes {
	g.origin.Static(relativePath, root)
	return g
}

func (g *RouterGroup) StaticFS(relativePath string, fs http.FileSystem) IRoutes {
	g.origin.StaticFS(relativePath, fs)
	return g
}
