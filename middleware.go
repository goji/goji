package goji

import "net/http"

/*
Use appends a middleware to the Mux's middleware stack.
Middleware are composable pieces of functionality that augment Handlers. Common
examples of middleware include request loggers, authentication checkers, and
metrics gatherers.
Middleware are evaluated in the reverse order in which they were added, but the
resulting Handlers execute in "normal" order (i.e., the Handler returned by the
first Middleware to be added gets called first).
For instance, given middleware A, B, and C, added in that order, Goji will
behave similarly to this snippet:
	augmentedHandler := A(B(C(yourHandler)))
	augmentedHandler.ServeHTTPC(ctx, w, r)
Assuming each of A, B, and C look something like this:
	func A(inner http.Handler) http.Handler {
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
or to view or modify the Handler that will be routed to. Middleware authors
should read the documentation for the "middleware" subpackage for more
information about how this is done.
The http.Handler returned by the given middleware must be safe for concurrent
use by multiple goroutines. It is not safe to concurrently register middleware
from multiple goroutines, or to register middleware concurrently with requests.
*/
func (m *Mux) Use(middleware func(http.Handler) http.Handler) {
	m.middleware = append(m.middleware, func(h http.Handler) http.Handler {
		return outerBridge{middleware, h}
	})
	m.buildChain()
}

// Pre-compile a Handler for us to use during dispatch. Yes, this means that
// adding middleware is quadratic, but it (a) happens during configuration time,
// not at "runtime", and (b) n should ~always be small.
func (m *Mux) buildChain() {
	m.handler = dispatch{}
	for i := len(m.middleware) - 1; i >= 0; i-- {
		m.handler = m.middleware[i](m.handler)
	}
}

type innerBridge struct {
	inner http.Handler
}

func (b innerBridge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.inner.ServeHTTP(w, r)
}

type outerBridge struct {
	mware func(http.Handler) http.Handler
	inner http.Handler
}

func (b outerBridge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.mware(innerBridge{b.inner}).ServeHTTP(w, r)
}
