package dff

import (
	"log"
	"net/http"
)

// 先定义一个router的结构体，里面放handleFunc，相当于重写net/http的handleFunc
type router struct {
	handlers map[string]HandlerFunc
}

// 实例化？？
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// 添加路由，即把拼接请求的方法，路径作为router.handle的key，请求作为value，
func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handle, ok := r.handlers[key]; ok {
		handle(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
