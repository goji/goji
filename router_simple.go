// +build goji_router_simple

package goji

import (
	"net/http"

	"context"
)

/*
This is the simplest correct router implementation for Goji.
*/

type router []route

type route struct {
	Pattern
	Handler
}

func (rt *router) add(p Pattern, h http.Handler) {
	*rt = append(*rt, route{p, h})
}

func (rt *router) route(ctx context.Context, r *http.Request) context.Context {
	for _, route := range *rt {
		if ctx := route.Match(ctx, r); ctx != nil {
			return &match{ctx, route.Pattern, route.Handler}
		}
	}
	return &match{Context: ctx}
}
