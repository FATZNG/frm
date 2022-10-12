# gee framework
## http/base 
**net/http**包

 原来的样子
```
 import "net/http"
 
 func main(){
    http.HandleFunc("/",indexHandler)
    http.HandleFunc("/hello",helloHandler)
    log.Fatal(http.ListenAndServe(":9999",nil))
 }
 
 func indexHandler(w http.ResponseWriter, req *http.Request){
    //逻辑代码
 }
 
  func HandleFunc(w http.ResponseWriter, req *http.Request){
    //逻辑代码
 }
 
```

封装后 @gee/gee.go
```
package gee
import(
    "net/http"
    "fmt"
) 

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct{
    router map[string]HandlerFunc
}

func New() *Engine{
    return &Engine{router:make(map[string]HandlerFunc)} 
}

func (engine *Engine) addRoute(method string,pattern string,handler HandlerFunc){
    key := method + "-" + pattern
    engine.router[key] = handler
}

func (engine *Engine) GET(pattern string,handler HandlerFunc){
    engine.addRoute("GET",pattern,handler)
}

func (engine *Engine) POST(pattern string,handler HandlerFunc){
    engine.addRoute("POST",pattern,handler)
}

func (engine *Engine) Run(addr string)(err errror){
    return http.ListenAndServe(addr,engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.response){
    key := req.Method + "-" + req.URL.Path
    if handle,ok := engine.router[key];ok{
        handle(w,req)
    }else{
        fmt.Fprint("not found")
    }
}
```
作者的解释：
```
* 首先定义了类型HandlerFunc，这是提供给框架用户的，用来定义路由映射的处理方法。我们在Engine中，添加了一张路由映射表router，key 由请求方法和静态路由地址构成，例如GET-/、GET-/hello、POST-/hello，这样针对相同的路由，如果请求方法不同,可以映射不同的处理方法(Handler)，value 是用户映射的处理方法。

* 当用户调用(*Engine).GET()方法时，会将路由和处理方法注册到映射表 router 中，(*Engine).Run()方法，是 ListenAndServe 的包装。

* Engine实现的 ServeHTTP 方法的作用就是，解析请求的路径，查找路由映射表，如果查到，就执行注册的处理方法。如果查不到，就返回 404 NOT FOUND 。
```

问题或者还不理解的地方
``` 
ServeHTTP 是重写了 net/http 的ServeHTTP
相当于是对net/http进行了封装。
strcut
```

## context
[@content.go](../gee/context.go)

```go
package gee

import (
	"encoding/json"
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

func (c *Context) String(code int, value string) {
	c.SetHeader("context-type", "txt/plain")
	c.StatusCode = code
	c.Writer.Write([]byte(value))
}


```

调整handler的参数
[@router.go](../gee/router.go)
```go
package gee

import (
	"log"
	"net/http"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
```

[main.go](../gee/main.go)
```go
package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	//r.GET("/", func(w http.ResponseWriter, req *http.Request) { fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path) })
	//r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
	//	for k, v := range req.Header {
	//		fmt.Fprintf(w, "Header[%q] = %q \n", k, v)
	//	}
	//})

	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}

```