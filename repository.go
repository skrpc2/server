package rysrv

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/fxamacker/cbor/v2"
)

// HandlerFunc is a Hertz request handler.
type HandlerFunc = app.HandlerFunc

// Repository is a JSON-RPC 2.0 methods repository.
type Repository struct {
	obj      map[string]interface{}
	handlers map[string]HandlerFunc
}

// NewRepository returns empty repository.
//
// It's safe to use Repository default value.
func NewRepository() *Repository {
	repo := &Repository{}
	repo.obj = map[string]interface{}{}
	repo.handlers = map[string]HandlerFunc{}
	return repo
}

func (r *Repository) Use(name string, obj interface{}) {
	r.obj[name] = obj
}

func (r *Repository) preHandler() HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if !c.IsPost() {
			c.SetStatusCode(consts.StatusMethodNotAllowed)
			c.Abort()
			return
		}

		c.Set("lang", c.GetHeader("lang"))
		for k, v := range r.obj {
			c.Set(k, v)
		}
		c.Next(ctx)
	}
}

// Handler returns a Hertz handler that dispatches requests to registered method handlers.
func (r *Repository) Handler() HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if !c.IsPost() {
			c.SetStatusCode(consts.StatusMethodNotAllowed)
			return
		}

		c.Set("lang", c.GetHeader("lang"))
		for k, v := range r.obj {
			c.Set(k, v)
		}

		handler, ok := r.handlers[string(c.Path())]
		if !ok {
			c.SetStatusCode(consts.StatusNotFound)
			return
		}
		handler(ctx, c)
	}
}

// RegisterRoutes registers all method handlers on the given Hertz server.
func (r *Repository) RegisterRoutes(h *server.Hertz) {
	//group := h.Group("/", r.preHandler())
	for method, handler := range r.handlers {
		fmt.Println("method = ", method)
		h.POST(method, handler)
	}
}

// Register registers new method handler.
func (r *Repository) Register(method string, handler HandlerFunc) {
	r.handlers[method] = handler
}

func Unmarshal(c *app.RequestContext, v interface{}) error {
	return cbor.Unmarshal(c.Request.Body(), v)
}
