package web

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestJSONError(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	type Engine struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// 错误的 json 语法
	s := `{"name":"}`
	var e Engine
	err := WarpJSONErr(json.Unmarshal([]byte(s), &e))
	fmt.Println(err.(*Error))

	// 空串错误
	s = ``
	err = WarpJSONErr(json.NewDecoder(strings.NewReader(s)).Decode(&e))
	fmt.Println(err.(*Error))

	// 类型错误
	s = `{"age":"10"}`
	err = WarpJSONErr(json.Unmarshal([]byte(s), &e))
	fmt.Println(err.(*Error))

	// panic
	s = `{"name":"test"}`
	err = WarpJSONErr(json.Unmarshal([]byte(s), e))
	fmt.Println(err.(*Error))
}
