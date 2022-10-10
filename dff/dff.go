package dff

import "net/http"

// 定义一个handlerFunc 作为方法 接收http.ResponseWriter,*http.Request变成了接收 *Context，即Context.go里面的
type HandlerFunc func(*Context)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

// 添加路由重构成了，跳转到router.go里面的addRouter了。
func (engine *Engine) addRouter(method string, pattern string, handlerFunc HandlerFunc) {
	engine.router.addRouter(method, pattern, handlerFunc)
}

func (engine *Engine) GET(pattern string, handlerFunc HandlerFunc) {
	engine.addRouter("GET", pattern, handlerFunc)
}

func (engine *Engine) POST(pattern string, handlerFunc HandlerFunc) {
	engine.addRouter("POST", pattern, handlerFunc)
}

func (engine *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c := newContext(writer, request)
	engine.router.handle(c)
}

//为什么要有 ServeHTTP ，因为将自定义的handler放在http.ListenAndServe时，需要自己实现ServeHTTP.
//对比http-base下base1和base2的代码。不用net/http的handleFunc，就需要自己实现ServeHTTP

func (engine *Engine) run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
