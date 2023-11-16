package web

import (
	"fmt"
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
type HTTPServerOption func(server *HTTPServer)

type HTTPServer struct {
	//*router
	router
	middlewares []MiddleWare
	log         func(message string, arg ...any)
}

// 用户可能不会使用NwHTTPServer,而是自己s := &HTTPServer{},会引起panic
func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	res := &HTTPServer{
		router: NewRouter(),
		log: func(message string, arg ...any) {
			fmt.Printf(message, arg...)
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res

}
func ServerWithMiddleware(meddlewares ...MiddleWare) HTTPServerOption {
	return func(server *HTTPServer) {
		server.middlewares = meddlewares
	}
}

// 处理请求的入口：
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Resp: writer,
		Req:  request,
	}
	//这是最后一个
	root := h.Serve
	//然后这里就是利用最后一个不断向前回溯组装链条
	//然后向前，把后一个作为前一个的next 构造好链条
	for i := len(h.middlewares) - 1; i >= 0; i-- {
		root = h.middlewares[i](root)
	}
	//这里最后一个步骤，就是把 RespData 和 RespStatusCode 刷新到响应里面
	var m MiddleWare = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			//就设置好了 RespData 和RespStatusCode
			next(ctx)
			h.flashResp(ctx)
		}
	}
	root = m(root)
	root(ctx)
}

func (h *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode != 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	n, err := ctx.Resp.Write(ctx.RespData)
	if err != nil || n != len(ctx.RespData) {
		h.log("写入响应失败 %v", err)
	}
}

func (h *HTTPServer) Serve(ctx *Context) {
	//接下来就是查找路由，并且执行命中的路由
	match, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	//如果路由未命中，或者没有业务逻辑（handler），则为404
	if !ok || match.n.handler == nil {
		ctx.RespStatusCode = 404
		ctx.RespData = []byte("NOT FOUND")
		return
	}
	ctx.PathParams = match.pathParams
	ctx.MatchedRoute = match.n.route
	match.n.handler(ctx)
}

// Start 启动服务器，用户指定端口
// 这种就是编程接口
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
