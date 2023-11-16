package errhandle

import (
	"WebFramework/web"
	"net/http"
	"testing"
)

func TestMiddleware_Build(t *testing.T) {
	builder := NewMiddlewareBuilder()
	builder.AddCode(http.StatusNotFound, []byte(`
<html>
	<body>
		<h1>哈哈哈哈，笨蛋！走失了！</h1>
	<body>
</html>
`)).
		AddCode(http.StatusBadRequest, []byte(`
<html>
	<body>
		<h1>请求不对头</h1>
	<body>
</html>
`))
	server := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))
	server.Start(":8081")
}
