package errhandle

import "WebFramework/web"

type Middleware struct {
	//这种设计只能放回固定的值
	//不能做到动态渲染
	resp map[int][]byte
}

func NewMiddlewareBuilder() *Middleware {
	return &Middleware{
		resp: map[int][]byte{},
	}
}
func (m Middleware) AddCode(status int, data []byte) *Middleware {
	m.resp[status] = data
	return &m
}

func (m Middleware) Build() web.MiddleWare {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {

			next(ctx)
			resp, ok := m.resp[ctx.RespStatusCode]
			if ok {
				//篡改结果
				ctx.RespData = resp
			}
		}
	}
}
