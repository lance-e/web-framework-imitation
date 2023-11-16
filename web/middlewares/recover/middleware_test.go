package recover

import (
	"WebFramework/web"
	"fmt"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := MiddlewareBuilder{
		StatusCode: 500,
		Data:       []byte("你 panic 了"),
		LogFunc: func(ctx *web.Context) {
			fmt.Printf("panic 路径: %s", ctx.Req.URL.String())
		},
	}
	server := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))
	server.Get("/user", func(ctx *web.Context) {
		panic("发生 panic 了")
	})
	server.Start(":8080")
}
