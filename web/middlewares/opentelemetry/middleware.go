package opentelemetry

import (
	"WebFramework/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

const instrumentotionName = "github.com/lance547/frameworkLeaning/web/middlewares/opentelemetry"

func (m MiddlewareBuilder) Build() web.MiddleWare {
	if m.Tracer == nil {
		m.Tracer = otel.GetTracerProvider().Tracer(instrumentotionName)
	}
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			reqCtx := ctx.Req.Context()
			//尝试与客户端的trace结合在一起
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier{})

			reqCtx, span := m.Tracer.Start(reqCtx, "unknown")
			defer span.End()
			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.scheme", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("http.host", ctx.Req.Host))
			//这里还可以继续加
			ctx.Req = ctx.Req.WithContext(reqCtx)
			//直接调用下一步
			next(ctx)
			//只有执行完next才可能有值
			span.SetName(ctx.MatchedRoute)
			//把响应码加上去
			span.SetAttributes(attribute.Int("http.status", ctx.RespStatusCode))
		}
	}
}
