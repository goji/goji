package goji

import (
	"net/http"

	"goji.io/internal"
	"golang.org/x/net/context"
)

func (m *Mux) router(ctx context.Context, r *http.Request) context.Context {
	if ctx.Value(internal.Path) == nil {
		ctx = context.WithValue(ctx, internal.Path, r.URL.EscapedPath())
	}
	for _, route := range m.routes {
		if ctx := route.Match(ctx, r); ctx != nil {
			ctx = context.WithValue(ctx, internal.Pattern, route.Pattern)
			ctx = context.WithValue(ctx, internal.Handler, route.Handler)
			return ctx
		}
	}
	return ctx
}
