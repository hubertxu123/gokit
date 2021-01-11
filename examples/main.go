package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gokit/gini"
	"gokit/zapi"
)

func main() {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		//c.Keys["traceId"] = fmt.Sprintf("%s", uuid.NewV4())
		fmt.Println(fmt.Sprintf("%s", uuid.NewV4()))
	},gini.Logger(zapi.DefaultLogger), gini.Recovery(zapi.DefaultLogger), )
	r.GET("/test", func(c *gin.Context) {
		panic(1111)
	})
	zapi.Error("222222222")
	r.Run(":8361")
}

