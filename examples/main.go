package main

import (
	"fmt"
	//ut "github.com/go-playground/universal-translator"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"gokit/gini"
	"gokit/vali"
	"gokit/zapi"
)

type user struct {
	Username string `json:"username" binding:"required" zh:"用户名"`
	//Age      int    `json:"age" binding:"gte=1,lte=10" zh:"年龄"`
	Age int `json:"age" binding:"test=aa,gte=1,lte=10" zh:"年龄"`
}

func main() {
	vt := vali.New(
		vali.Origin(binding.Validator.Engine().(*validator.Validate)),
		vali.Locale("zh"),
	)
	vt.RegisterValidationTrans("test", func(fl validator.FieldLevel) bool {
		//140 Age aa test 年龄
		fmt.Println(fl.Field(), fl.StructFieldName(), fl.Param(), fl.GetTag(), fl.FieldName())
		return false
	}, func(ut ut.Translator) error {
		if ut.Locale() == "en" {
			return ut.Add("test", "test rul {0} fail", false)
		} else {
			return ut.Add("test", "测试规则{0}未通过", false)
		}
	})

	r := gini.New(
		gini.GroupValidate(vt),
		gini.GroupCtxOptions(
			gini.CtxTranslator(vt.Translator()),
		),
	)

	logger := zapi.DefaultLogger
	r.Use(
		//trace log
		func(c *gini.Context) {
			if c.Keys == nil {
				c.Keys = map[string]interface{}{}
			}
			c.Keys[logger.TraceIdKey()] = zapi.UUIDV4()
		},
		gini.Logger(logger),
		gini.Recovery(logger),
	)
	// curl -XPOST --header 'Content-Type: application/json' -d '{"username":"asong","age":140}' http://localhost:8361/aa/cc
	g := r.Group("/aa")
	g.POST("cc", func(c *gini.Context) {
		u := new(user)
		err := c.ShouldBind(u)
		c.Writer.WriteString(err.Error())
	})

	r.Run(":8361")
}
