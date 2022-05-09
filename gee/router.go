package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       // key: method, value: trie root
	handlers map[string]HandlerFunc // key: 'method'- 'pattern', value: HandlerFunc
}

func newRouter() (r *router) {
	r = &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}

	for _, method := range METHOD {
		r.roots[method] = &node{}
	}
	return
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	root, ok := r.roots[method]
	if !ok {
		panic("unsupported method")
	}
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	root.insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method, path string) (node *node, params map[string]string) {
	root, ok := r.roots[method]
	if !ok {
		return
	}
	params = make(map[string]string)
	searchParts := parsePattern(path)
	node = root.search(searchParts, 0)
	if node != nil {
		parts := parsePattern(node.pattern)
		for i, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[i]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
	}
	return
}

func (r *router) handle(c *Context) {
	route, params := r.getRoute(c.Method, c.Path)
	if route != nil { // found
		c.Params = params
		key := c.Method + "-" + route.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 Not Found: %s\n", c.Path)
		})
	}
	c.Next()
}

func parsePattern(pattern string) (parts []string) {
	tmp := strings.Split(pattern, "/")
	for _, part := range tmp {
		if part != "" {
			parts = append(parts, part)
			if part[0] == '*' {
				break // if pattern is '/a/*b/c', return parts[a, *b]
			}
		}
	}
	return
}
