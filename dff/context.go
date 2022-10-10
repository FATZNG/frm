package dff

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 先封装了一个接收json格式的interface
type H map[string]interface{}

// 再封装一个接收请求的
type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Method     string
	Path       string
	StatusCode int
}

// 构造context?
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// 在请求里面get为key的值
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 获去请求里面的formData
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 设置状态以及设置请求头的状态
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置请求头,即请求投里面的 content-type
func (c *Context) SetHeader(key string, value string) {
	c.SetHeader("Content-Type", "text/html")
	c.Writer.Header().Set(key, value)
}

func (c *Context) HTML(code int, html string) {
	c.StatusCode = code
	c.Writer.Write([]byte(html))
}

func (c *Context) Data(code int, data []byte) {
	c.StatusCode = code
	c.Writer.Write(data)
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("context-type", "application/json")
	c.StatusCode = code
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Context-type", "txt/plain")
	c.StatusCode = code
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}
