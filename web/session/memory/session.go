package memory

import (
	"WebFramework/web/session"
	"context"
	"errors"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var (
	errorKeyNotFound     = errors.New("session :can't find the key")
	errorSessionNotFound = errors.New("session:can't find the session")
)

type Store struct {
	mutex      sync.RWMutex
	sessions   cache.Cache
	expiration time.Duration
}

func NewStore(expiration time.Duration) *Store {
	return &Store{
		sessions:   *cache.New(expiration, time.Second),
		expiration: expiration,
	}
}
func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sess := &Session{
		id: id,
	}
	s.sessions.Set(id, sess, s.expiration)
	return sess, nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	value, ok := s.sessions.Get(id)
	if !ok {
		return errors.New("session :id对应session不存在")
	}
	s.sessions.Set(id, value, s.expiration)
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	sess, ok := s.sessions.Get(id)
	if !ok {
		return nil, errorSessionNotFound
	}
	return sess.(*Session), nil
}

func (s Store) Remove(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessions.Delete(id)
	return nil
}

type Session struct {
	id     string
	values sync.Map
}

func (s *Session) Get(ctx context.Context, key string) (any, error) {
	value, ok := s.values.Load(key)
	if !ok {
		return nil, errorKeyNotFound
	}
	return value, nil
}

func (s *Session) Set(ctx context.Context, key string, value any) error {
	s.values.Store(key, value)
	return nil
}

func (s *Session) ID() string {
	return s.id
}
