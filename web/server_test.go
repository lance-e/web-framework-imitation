package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	h := NewHTTPServer()
	//h.addRoute(http.MethodGet, "/user/login", func(ctx *Context) {
	//	fmt.Println("first")
	//})

	h.addRoute(http.MethodGet, "/user/login", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello world"))
	})
	h.Get("/user/abc", func(ctx *Context) {
		ctx.Resp.Write([]byte(fmt.Sprintf("Hello %s", ctx.Req.URL.Path)))
	})
	h.Get("/user/:id", func(ctx *Context) {
		id, err := ctx.PathValue("id").AsInt64()
		if err != nil {
			ctx.Resp.WriteHeader(400)
			ctx.Resp.Write([]byte("id 输入不对"))
			return
		}
		ctx.Resp.Write([]byte(fmt.Sprintf("hello %d", id)))
	})
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	h.Get("/user/longxu", func(ctx *Context) {
		ctx.RespJSON(200, Person{Name: "龙旭"})
	})

	h.Start(":8080")
}

func TestHTTPServer_ServeHTTP(t *testing.T) {
	server := NewHTTPServer()
	server.middlewares = []MiddleWare{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("这是第一个before")
				next(ctx)
				fmt.Println("这是第一个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("这是第二个before")
				next(ctx)
				fmt.Println("这是第二个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第三个中断")

			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第四个你看不到这句话")

			}
		},
	}
	server.ServeHTTP(nil, &http.Request{})
}
