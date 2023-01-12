// 自定义错误

package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// 业务常用错误
var (
	ErrUnknown           = NewError("UnKnow", "未知错误")
	BadRequest           = NewError("BadRequest", "参数有误")
	ErrDB                = NewError("StoreErr", "数据发生错误")
	ErrUnauthorizedToken = NewError("UnauthorizedToken", "Token 已过期或错误")
	ErrJSON              = NewError("UnmarshalErr", "JSON 编解码出错")
	ErrNameOrPasswd      = NewError("NameOrPasswd", "用户名或密码错误")
	ErrNotFound          = NewError("NotFound", "资源未找到或服务未开启")

	ErrJSONSyntaxError = NewError("JSONSyntaxError", "正文包含错误的 JSON 语法")
	ErrJSONEmpty       = NewError("JSONEmpty", "正文内容为空，请检查 JSON")
	ErrorJSONType      = NewError("JSONType", "正文包含错误的 JSON 类型字段")
)

// WarpJSONErr 封装 JSON 错误
func WarpJSONErr(err error) Errorer {
	if err == nil {
		return nil
	}

	var (
		syntaxError           *json.SyntaxError
		unmarshalTypeError    *json.UnmarshalTypeError
		invalidUnmarshalError *json.InvalidUnmarshalError
	)

	switch {
	// 检查错误，并获取具体类型
	// 如果是，则返回错误消息和问题位置
	case errors.As(err, &syntaxError):
		return ErrJSONSyntaxError.With(fmt.Sprintf("包含错误的 JSON 格式 (at character %d)", syntaxError.Offset))
	// https://github.com/golang/go/issues/25956.
	case errors.Is(err, io.ErrUnexpectedEOF):
		return ErrJSONSyntaxError.With("包含错误的 JSON 格式")
	// json 字段类型 与 go 结构体参数类型不匹配
	case errors.As(err, &unmarshalTypeError):
		if unmarshalTypeError.Field != "" {
			return ErrorJSONType.With(fmt.Sprintf("请检查 %q", unmarshalTypeError.Field))
		}
		return ErrorJSONType.With(fmt.Sprintf("at character %d", unmarshalTypeError.Offset))
	// 空 json 正文
	case errors.Is(err, io.EOF):
		return ErrJSONEmpty.With("数据不能为空")
	// 绑定的目标非指针
	case errors.As(err, &invalidUnmarshalError):
		panic(err)
	}
	return BadRequest.With(err.Error())
}

// Error ...
type Error struct {
	reason  string   // 错误原因
	msg     string   // 错误信息，用户可读
	details []string // 错误扩展，开发可读
}

// E 可反序列化的 err
type E struct {
	Reason  string   `json:"reason"`
	Msg     string   `json:"msg"`
	Details []string `json:"details"`
}

func (e E) String() string {
	var msg strings.Builder
	msg.WriteString(e.Msg)
	msg.WriteByte('\n')

	for _, v := range e.Details {
		msg.WriteString("  " + v)
	}
	return msg.String()
}

// Unmarshal ...
func Unmarshal(b []byte) (e E) {
	_ = json.Unmarshal(b, &e)
	return
}

var codes = make(map[string]string, 8)

// NewError 创建自定义错误
func NewError(reason, msg string) *Error {
	if _, ok := codes[reason]; ok {
		panic(fmt.Sprintf("错误码 %s 已经存在，请更换一个", reason))
	}
	codes[reason] = msg
	return &Error{reason: reason, msg: msg}
}

// Reason ..
func (e *Error) Reason() string {
	return e.reason
}

// Message ..
func (e *Error) Message() string {
	return e.msg
}

// Details 错误
func (e *Error) Details() []string {
	return e.details
}

// Map ..
func (e *Error) Map() map[string]any {
	return map[string]any{
		"msg":    e.Message(),
		"reason": e.Reason(),
	}
}

// With 错误详情
func (e *Error) With(args ...string) *Error {
	newErr := *e
	newErr.details = make([]string, 0, len(args))
	newErr.details = append(newErr.details, args...)
	return &newErr
}

// HTTPCode http status code
// 权限相关错误 401
// 程序错误 500
// 其它错误 400
func (e *Error) HTTPCode() int {
	switch e.reason {
	case "":
		return http.StatusOK
	case ErrUnauthorizedToken.reason:
		return http.StatusUnauthorized
	}
	return http.StatusBadRequest
}
