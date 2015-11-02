package pattern

import (
	"net/http"

	"goji.io"
	"golang.org/x/net/context"
)

/*
WithMethods returns a Pattern that matches the same requests as the given
pattern, but additionally only accepts requests with the given HTTP methods.

The returned Pattern supports both the HTTPMethods and PathPrefix Pattern
optimizations.
*/
func WithMethods(p goji.Pattern, methods ...string) goji.Pattern {
	ms := make(withMethods)
	for _, m := range methods {
		ms[m] = struct{}{}
	}
	return Compose(ms, p)
}

type withMethods map[string]struct{}

func (w withMethods) Match(ctx context.Context, r *http.Request) context.Context {
	if _, ok := w[r.Method]; ok {
		return ctx
	}
	return nil
}

func (w withMethods) HTTPMethods() map[string]struct{} {
	return w
}
