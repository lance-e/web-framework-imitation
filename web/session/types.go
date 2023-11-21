package session

import (
	"context"
	"net/http"
)

// Store 管理Session本身
type Store interface {
	// 生成Session
	Generate(ctx context.Context, id string) (Session, error)
	//刷新Session
	Refresh(ctx context.Context, id string) error
	//获取Session
	Get(ctx context.Context, id string) (Session, error)
	//删除Session
	Remove(ctx context.Context, id string) error
}
type Session interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value any) error
	ID() string
}
type Propagator interface {
	Inject(id string, writer http.ResponseWriter) error
	Extract(req *http.Request) (string, error)
	Remove(writer http.ResponseWriter) error
}
