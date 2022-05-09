package gee

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute(GET, "/", nil)
	r.addRoute(GET, "/hello/:name", nil)
	r.addRoute(GET, "/hello/b/c", nil)
	r.addRoute(GET, "/hi/:name", nil)
	r.addRoute(GET, "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern fail")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute(GET, "/hello/lzhseu")

	if n == nil {
		t.Fatal("nil should['n be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "lzhseu" {
		t.Fatal("name should be equal 'lzhseu'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])
}
