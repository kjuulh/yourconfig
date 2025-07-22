package yourconfig

import (
	"context"
	"os"
	"sync/atomic"
)

var defaultLogger atomic.Pointer[Provider]

func init() {
	defaultLogger.Store(newProvider())
}

func SetDefault(provider *Provider) {
	defaultLogger.Store(provider)
}

func Default() *Provider {
	return defaultLogger.Load()
}

type Handler interface {
	Get(ctx context.Context, key string) (string, error)
}

type Provider struct {
	handler Handler
}

func (p *Provider) Get(ctx context.Context, key string) (string, error) {
	return p.handler.Get(ctx, key)
}

func newProvider() *Provider {
	return &Provider{
		handler: defaultHandler(),
	}
}

func New(handler Handler) *Provider {
	return &Provider{
		handler: handler,
	}
}

type envHandler struct{}

func (e *envHandler) Get(ctx context.Context, key string) (string, error) {
	val := os.Getenv(key)
	return val, nil
}

func defaultHandler() Handler {
	return &envHandler{}
}
