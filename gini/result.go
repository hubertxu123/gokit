package gini

const (
	SUCCESS = 0
	FAIL    = -1
)

type ResultFunc func(code int, msg string, data interface{}) interface{}

type StdResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func DefaultResult(code int, msg string, data interface{}) interface{} {
	return StdResult{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
