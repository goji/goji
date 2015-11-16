package goji

import (
	"net/http"

	"goji.io/internal"

	"golang.org/x/net/context"
)

/*
Mux is a HTTP multiplexer / router similar to net/http.ServeMux.

Muxes multiplex traffic between many Handlers by selecting the first applicable
Pattern. They then call a common middleware stack, finally passing control to
the selected Handler. See the documentation on the Handle function for more
information about how routing is performed, the documentation on the Pattern
type for more information about request matching, and the documentation for the
Use method for more about middleware.

Muxes cannot be configured concurrently from multiple goroutines, nor can they
be configured concurrently with requests.
*/
type Mux struct {
	handler    Handler
	middleware []func(Handler) Handler
	router     router
}

/*
NewMux returns a new Mux with no configured middleware or routes.

A common pattern is to organize your application similarly to how you structure
your URLs. For instance, a photo-sharing site might have URLs that start with
"/users/" and URLs that start with "/albums/"; such a site might have three
Muxes: one to manage the URL hierarchy for users, one to manage albums, and a
third top-level Mux to select between the other two.
*/
func NewMux() *Mux {
	m := &Mux{}
	m.buildChain()
	return m
}

/*
ServeHTTP implements net/http.Handler. It uses context.TODO as the root context
in order to ease the conversion of non-context-aware Handlers to context-aware
ones using static analysis.

Users who know that their mux sits at the top of the request hierarchy should
consider creating a small helper http.Handler that calls this Mux's ServeHTTPC
function with context.Background.
*/
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.ServeHTTPC(context.TODO(), w, r)
}

/*
ServeHTTPC implements Handler.
*/
func (m *Mux) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if ctx.Value(internal.Path) == nil {
		ctx = context.WithValue(ctx, internal.Path, r.URL.EscapedPath())
	}
	ctx = m.router.route(ctx, r)
	m.handler.ServeHTTPC(ctx, w, r)
}

var _ http.Handler = &Mux{}
var _ Handler = &Mux{}
