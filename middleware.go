package goji

import (
	"net/http"

	"golang.org/x/net/context"
)

/*
Use appends a middleware to the Mux's middleware stack.

Middleware are composable pieces of functionality that augment Handlers. Common
examples of middleware include request loggers, authentication checkers, and
metrics gatherers.

Middleware are evaluated in the reverse order in which they were added, but the
resulting Handlers execute in "normal" order (i.e., the Handler returned by the
first Middleware to be added gets called first).

For instance, given middleware A, B, and C, added in that order, Goji's behavior
will look something like this:

	augmentedHandler := A(B(C(yourHandler)))
	augmentedHandler.ServeHTTPC(ctx, w, r)

Assuming each of A, B, and C look something like this:

	func A(inner goji.Handler) goji.Handler {
		log.Print("A: called")
		mw := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			log.Print("A: before")
			inner.ServeHTTPC(ctx, w, r)
			log.Print("A: after")
		}
		return goji.HandlerFunc(mw)
	}

we'd expect to see the following in the log:

	C: called
	B: called
	A: called
	---
	A: before
	B: before
	C: before
	yourHandler: called
	C: after
	B: after
	A: after

Note that augmentedHandler may be called many times. Put another way, you will
see many invocations of the portion of the log below the divider, and perhaps
only see the portion above the divider a single time. Also note that as an
implementation detail, net/http-style middleware will be called once per
request, even though the Goji-style middleware around them might only ever be
called a single time.

Middleware in Goji is called after routing has been performed. Therefore it is
possible to examine any routing information placed into the context by Patterns,
or to view or modify the Handler that will be routed to.

The http.Handler returned by the given middleware must be safe for concurrent
use by multiple goroutines. It is not safe to concurrently register middleware
from multiple goroutines.
*/
func (m *Mux) Use(middleware func(http.Handler) http.Handler) {
	if bridge, ok := m.h.(*handlerBridge); ok {
		bridge.hs = append(bridge.hs, middleware)
	} else {
		m.h = &handlerBridge{m.h, []func(http.Handler) http.Handler{middleware}}
	}
}

/*
UseC appends a context-aware middleware to the Mux's middleware stack. See the
documentation for Use for more information about the semantics of middleware.

The Handler returned by the given middleware must be safe for concurrent use by
multiple goroutines. It is not safe to concurrently register middleware from
multiple goroutines.
*/
func (m *Mux) UseC(middleware func(Handler) Handler) {
	m.h = middleware(m.h)
}

// handlerBridge allows us to up-convert http middleware to goji middleware.
type handlerBridge struct {
	h  Handler
	hs []func(http.Handler) http.Handler
}

func (b *handlerBridge) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b.h.ServeHTTPC(ctx, w, r)
	})
	for i := 0; i < len(b.hs); i++ {
		h = b.hs[len(b.hs)-i-1](h)
	}
	h.ServeHTTP(w, r)
}

var _ Handler = &handlerBridge{}
