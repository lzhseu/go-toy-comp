package main

import (
	"flag"
	"fmt"
	"gee"
	"geecache"
	"log"
	"net/http"
	"strings"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *geecache.Group {
	return geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))
}

func startCacheServer(addr string, addrs []string, group *geecache.Group) {
	peers := geecache.NewHTTPPool(addr)
	peers.Set(addrs...)
	group.RegisterPeerPicker(peers)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startApiServer(apiAddr string, group *geecache.Group) {
	http.Handle("/api", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		key := req.URL.Query().Get("key")
		view, err := group.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(view.ByteSlice())
	}))

	log.Println("font-end server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {

	// test geecache
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "GeeCache server port")
	flag.BoolVar(&api, "api", false, "Start an api server")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	geeCache := createGroup()
	if api {
		go startApiServer(apiAddr, geeCache)
	}
	startCacheServer(addrMap[port], addrs, geeCache)

	//geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(func(key string) ([]byte, error) {
	//	log.Println("[SlowDB] search key", key)
	//	if v, ok := db[key]; ok {
	//		return []byte(v), nil
	//	}
	//	return nil, fmt.Errorf("%s not exist", key)
	//}))
	//addr := "localhost:9999"
	//peer := geecache.NewHTTPPool(addr)
	//log.Println("geecache is running at", addr)
	//log.Fatal(http.ListenAndServe(addr, peer))

	//r := gee.Default()

	// Recovery test
	//r.GET("/panic", func(c *gee.Context) {
	//	arr := []string{"1","2","3"}
	//	c.String(http.StatusOK, arr[10])
	//})

	// static template test
	//r := gee.New()
	//r.Use(gee.TimeRecord())
	//r.SetFuncMap(template.FuncMap{
	//	"FormatAsDate": FormatAsDate,
	//})
	//r.LoadHtmlGlob("templates/*")
	//r.Static("/assets", "./static")
	//
	//stu1 := &student{Name: "John", Age: 18}
	//stu2 := &student{Name: "Mary", Age: 20}
	//
	//r.GET("/", func(c *gee.Context) {
	//	c.HTML(http.StatusOK, "css.tmpl", nil)
	//})
	//
	//r.GET("/stu", func(c *gee.Context) {
	//	c.HTML(http.StatusOK, "arr.tmpl", gee.H{
	//		"title": "gee",
	//		"stuArr": [2]*student{stu1, stu2},
	//	})
	//})
	//
	//r.GET("/date", func(c *gee.Context) {
	//	c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
	//		"title": "gee",
	//		"now": time.Date(2022, 5, 6, 17, 5, 0, 0, time.UTC),
	//	})
	//})

	//v1 := r.Group("/v1")
	//v1.Use(onlyForV1())
	//{
	//	v1.GET("/", func(c *gee.Context) {
	//		c.String(http.StatusOK, "<h1>Hello Gee<h1>")
	//	})
	//
	//	v1.GET("/hello", func(c *gee.Context) {
	//		c.String(http.StatusOK, "hello %s, you are at %s\n", c.Query("name"), c.Path)
	//	})
	//}
	//
	//v2 := v1.Group("/v2")
	//v2.Use(onlyForV2())
	//{
	//	v2.GET("/hello/:name", func(c *gee.Context) {
	//		c.String(http.StatusOK, "hello %s, you are at %s\n", c.Param("name"), c.Path)
	//	})
	//
	//	v2.POST("/login", func(c *gee.Context) {
	//		c.JSON(http.StatusOK, gee.H{
	//			"username": c.PostForm("username"),
	//			"password": c.PostForm("password"),
	//		})
	//	})
	//}

	// test dynamic router
	//r.GET("/", func(c *gee.Context) {
	//	c.HTML(http.StatusOK, "<h1>hello Gee</h1>")
	//})
	//r.GET("/hello", func(c *gee.Context) {
	//	c.String(http.StatusOK, "hello %s, you are at %s\n", c.Query("name"), c.Path)
	//})
	//r.GET("/hello/:name", func(c *gee.Context) {
	//	c.String(http.StatusOK, "hello %s, you are at %s\n", c.Param("name"), c.Path)
	//})
	//r.GET("assets/*filepath", func(c *gee.Context) {
	//	c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	//})
	//r.GET("/a/:a/b/:b/*c", func(c *gee.Context) {
	//	c.JSON(http.StatusOK, gee.H{"a": c.Param("a"), "b": c.Param("b"), "c": c.Param("c")})
	//})
	//r.POST("/login", func(c *gee.Context) {
	//	c.JSON(http.StatusOK, gee.H{
	//		"username": c.PostForm("username"),
	//		"password": c.PostForm("password"),
	//	})
	//})
	//
	//r.Run(":9999")

	//s1 := make([]int, 1)
	//s2 := s1
	//fmt.Println("s1: ", s1)
	//fmt.Println("s2: ", s2)
	//
	//s2[0] = 100
	//fmt.Println("s1: ", s1)
	//fmt.Println("s2: ", s2)
	//
	//s1 = append(s1, 200)
	//s2 = s1
	//fmt.Println("s1: ", s1)
	//fmt.Println("s2: ", s2)

	//var s1 []int
	//fmt.Println(s1)
	//if s1 == nil {
	//	fmt.Println("s1 = ", s1)
	//
	//}
	//s2 := append(s1, 0)
	//fmt.Println(s2)
	//test(",a,b,c")
	//fmt.Println("===")
	//res := test2()
	//if res == nil {
	//	fmt.Printf("map: %#v\n", res)
	//}
	//var m = make(map[string]string)
	//m["a"] = "aaa"
	//fmt.Println(m["a"])
	//fmt.Println(m["b"])
	//for i, i2 := range m {
	//	fmt.Printf("i: %v, i2: %v\n", i, i2)
	//}
	//
	//a1, a2 := m["b"]
	//fmt.Println(a1)
	//fmt.Println(a2)
	//if a1 == "" {
	//	fmt.Println("---")
	//}
}

func test(s string) (res []string) {
	fmt.Println(res)
	tmp := strings.Split(s, ",")
	for _, part := range tmp {
		res = append(res, part)
	}
	fmt.Println(res)
	return
}

func test2() (res map[string]string) {
	return
}

func onlyForV1() gee.HandlerFunc {
	return func(c *gee.Context) {
		log.Printf("begin handler onlyForV1")
		c.Next()
		log.Printf("[%d] %s for group v1", c.StatusCode, c.Req.RequestURI)
	}
}

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		log.Printf("begin handler onlyForV2z")
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s for group v2", c.StatusCode, c.Req.RequestURI)
	}
}
