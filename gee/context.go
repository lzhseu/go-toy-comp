package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// store middlewares and service handler during one request
	handlers   []HandlerFunc
	handlerIdx int
	// engine
	engine *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:     w,
		Req:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
		handlerIdx: -1,
	}
}

func (c *Context) Next() {
	c.handlerIdx++
	count := len(c.handlers)
	for ; c.handlerIdx < count; c.handlerIdx++ {
		c.handlers[c.handlerIdx](c)
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) SetContentType(value string) {
	c.SetHeader("Content-Type", value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetContentType("text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetContentType("application/json")
	c.Status(code)

	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		panic(err) // net/http 包本身有捕获 panic，同时自己实现的中间件也会捕获并恢复
	}
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetContentType("text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(http.StatusInternalServerError, err.Error())
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) Fail(code int, err string) {
	c.handlerIdx = len(c.handlers)
	c.JSON(code, H{"message": err})
}
