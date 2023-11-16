//go:build e2e

package opentelemetry

import (
	"WebFramework/web"
	"go.opentelemetry.io/otel"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer(instrumentotionName)
	builder := MiddlewareBuilder{
		Tracer: tracer,
	}
	server := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))
	server.Get("/user", func(ctx *web.Context) {
		c, span := tracer.Start(ctx.Req.Context(), "first_layer")
		defer span.End()
		secondC, second := tracer.Start(c, "second_layer")
		time.Sleep(time.Second)
		_, third1 := tracer.Start(secondC, "third_layer_1")
		time.Sleep(100 * time.Millisecond)
		third1.End()
		_, third2 := tracer.Start(secondC, "third_layer_2")
		time.Sleep(300 * time.Millisecond)
		third2.End()
		second.End()

		_, span = tracer.Start(ctx.Req.Context(), "first_layer_2")
		defer span.End()
		time.Sleep(100 * time.Millisecond)
		ctx.RespJSON(202, User{
			Name: "Tom",
		})
	})

	//initZipkin(t)

	server.Start(":8081")
}

type User struct {
	Name string
}

//func initZipkin(t *testing.T){
//
//}
