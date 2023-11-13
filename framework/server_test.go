package framework

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

	h.Start(":8080")
}
