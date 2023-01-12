package web

import (
	"net/http"
)

// Errorer ...
type Errorer interface {
	Reason() string
	HTTPCode() int
	Message() string
	Details() []string
}

// ResponseWriter ...
type ResponseWriter interface {
	JSON(code int, obj interface{})
	File(filepath string)
	Set(string, any)
}

// Success 通用成功返回
func Success(c ResponseWriter, bean interface{}) {
	c.JSON(http.StatusOK, bean)
}

const responseErr = "responseErr"

// Fail 通用错误返回
func Fail(c ResponseWriter, err Errorer) {
	r := map[string]interface{}{"reason": err.Reason(), "msg": err.Message()}
	d := err.Details()
	if len(d) > 0 {
		r["details"] = d
	}

	c.Set(responseErr, r)
	c.JSON(err.HTTPCode(), r)
}

// FailStr 有些错误没有详情，只有字符串，且展示给用户
func FailStr(c ResponseWriter, err string) {
	r := map[string]any{"msg": err}
	c.Set(responseErr, r)
	c.JSON(400, r)
}
