package framework

import (
	"net"
	"net/http"
)

// 确保HTTPServer 结构体 一定实现了Server接口
var _ Server = &HTTPServer{}

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	Start(add string) error

	//AddRoute 增加路由注册功能
	//method 是http方法、
	//path是路由
	//handleFunc是业务逻辑
	addRoute(method string, path string, handleFunc HandleFunc)
}
type HTTPServer struct {
	//*Router
	Router
}

// 用户可能不会使用NwHTTPServer,而是自己s := &HTTPServer{},会引起panic
func NewHTTPServer() HTTPServer {
	return HTTPServer{
		Router: NewRouter(),
	}

}

// 处理请求的入口：
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Resp: writer,
		Req:  request,
	}
	//接下来就是查找路由，并且执行命中的路由
	h.Serve(ctx)
}
func (h *HTTPServer) Serve(ctx *Context) {
	//接下来就是查找路由，并且执行命中的路由
	match, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	//如果路由未命中，或者没有业务逻辑（handler），则为404
	if !ok || match.n.handler == nil {
		ctx.Resp.WriteHeader(404)
		_, _ = ctx.Resp.Write([]byte("NOT FOUND"))
		return
	}
	ctx.PathParams = match.pathParams
	match.n.handler(ctx)
}
func (h *HTTPServer) Start(add string) error {
	listener, err := net.Listen("tcp", add)
	if err != nil {
		return err
	}
	//用户可以回调
	return http.Serve(listener, h)

	//另外一种写法
	//http.ListenAndServe()
}

func (h *HTTPServer) Get(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodGet, path, handleFunc)
}

func (h *HTTPServer) Post(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodPost, path, handleFunc)
}

func (h *HTTPServer) Put(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodPut, path, handleFunc)
}

func (h *HTTPServer) Delete(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodDelete, path, handleFunc)
}
