package gee

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

var METHOD = []string{GET, POST, PUT, DELETE}

type HandlerFunc func(c *Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware, like filter in java web
	engine      *Engine       // all groups share an Engine instance
}

type Engine struct {
	*RouterGroup
	router        *router
	groups        []*RouterGroup
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

func New() (engine *Engine) {
	engine = &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = append(engine.groups, engine.RouterGroup)
	return
}

func Default() *Engine {
	engine := New()
	engine.Use(TimeRecord(), Recovery())
	return engine
}

func (group *RouterGroup) defaultStaticHandler(pathPrefix string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, pathPrefix)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (group *RouterGroup) Static(pathPrefix, fileRoot string) {
	handler := group.defaultStaticHandler(pathPrefix, http.Dir(fileRoot))
	pattern := path.Join(group.prefix, pathPrefix, "/*filepath")
	group.GET(pattern, handler)
}

func (group *RouterGroup) Group(prefix string) (newGroup *RouterGroup) {
	engine := group.engine
	newGroup = &RouterGroup{prefix: group.prefix + prefix, engine: engine}
	engine.groups = append(engine.groups, newGroup)
	return
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute(GET, pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute(POST, pattern, handler)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHtmlGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// find all middlewares for the group
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	context := newContext(w, req)
	context.handlers = middlewares
	context.engine = engine
	engine.router.handle(context)
}
