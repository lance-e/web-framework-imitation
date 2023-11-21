package test

import (
	"WebFramework/web"
	"WebFramework/web/session"
	"WebFramework/web/session/cookie"
	"WebFramework/web/session/memory"
	"net/http"
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	//做一个简单的登录校验
	var m *session.Manager = &session.Manager{
		Propagator: cookie.NewPropagator(),
		Store:      memory.NewStore(15 * time.Minute),
	}
	server := web.NewHTTPServer(web.ServerWithMiddleware(func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			if ctx.Req.URL.Path == "/login" {
				//放过去，用户准备登录
				next(ctx)
				return
			}

			_, err := m.GetSession(ctx)
			if err != nil {
				ctx.RespStatusCode = http.StatusUnauthorized
				ctx.RespData = []byte("请重新登录")
			}
			//刷新session
			_ = m.RefreshSession(ctx)

			next(ctx)
		}
	}))
	server.Post("/login", func(ctx *web.Context) {
		//要在这之前校验用户名和登录密码
		sess, err := m.InitSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("登录失败")
			return
		}
		err = sess.Set(ctx.Req.Context(), "nickname", "lance47")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("登录失败")
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("de")
	})
	server.Get("/user", func(ctx *web.Context) {
		sess, _ := m.GetSession(ctx)
		value, _ := sess.Get(ctx.Req.Context(), "nickname")
		ctx.RespData = []byte(value.(string))
	})
	server.Post("/logout", func(ctx *web.Context) {
		//清理数据
		err := m.RemoveSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("退出失败")
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("退出成功")
		return
	})

	server.Start(":8081")
}
