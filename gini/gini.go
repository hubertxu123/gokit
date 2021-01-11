package gini

import (
	"github.com/gin-gonic/gin"
)
// gin improved
const (
	SUCCESS = 0
	FAIL = -1
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResult(code int, msg string, data interface{}) *Result {
	return &Result{
		Code: code,
		Msg: msg,
		Data: data,
	}
}

func JsonSuccess(c *gin.Context,data interface{})  {
	c.JSON(200,NewResult(SUCCESS, "success", data))
}

func JsonError(c *gin.Context,code int,msg string)  {
	c.JSON(200,NewResult(code, msg, []string{}))
}

func JsonErrorMsg(c *gin.Context,msg string)  {
	JsonError(c,FAIL,msg)
}
