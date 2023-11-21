package session

import (
	"WebFramework/web"
	"github.com/satori/go.uuid"
)

type Manager struct {
	Propagator
	Store
	CtxSessKey string
}

func (m *Manager) GetSession(ctx *web.Context) (Session, error) {
	if ctx.UserValue == nil {
		ctx.UserValue = make(map[string]any, 1)
	}
	value, ok := ctx.UserValue[m.CtxSessKey]
	if ok {
		return value.(Session), nil
	}
	//尝试缓存session
	sessionId, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	sess, err := m.Get(ctx.Req.Context(), sessionId)
	ctx.UserValue[m.CtxSessKey] = sess
	return sess, err
}
func (m *Manager) InitSession(ctx *web.Context) (Session, error) {
	id := uuid.NewV1().String()
	session, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	//注入到http响应里面
	err = m.Inject(id, ctx.Resp)
	if err != nil {
		return nil, err
	}
	return session, err
}
func (m *Manager) RemoveSession(ctx *web.Context) error {
	session, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	if err = m.Store.Remove(ctx.Req.Context(), session.ID()); err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Resp)
}
func (m *Manager) RefreshSession(ctx *web.Context) error {
	session, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	return m.Store.Refresh(ctx.Req.Context(), session.ID())
}
