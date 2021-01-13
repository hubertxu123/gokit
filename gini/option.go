package gini

import "gokit/vali"

type GroupOption func(g *RouterGroup)

func GroupCtxOptions(ctxOptions ...CtxOption) GroupOption {
	return func(g *RouterGroup) {
		g.ctxOptions = ctxOptions
	}
}

func GroupValidate(validate *vali.Validate) GroupOption {
	return func(g *RouterGroup) {
		g.validate = validate
	}
}
