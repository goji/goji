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

For instance, given middleware A, B, and C, added in that order, Goji's behavior
will look something like this:

	augmentedHandler := A(B(C(yourHandler)))
	augmentedHandler.ServeHTTP(w, r)

Assuming each of A, B, and C look something like this:

	func A(inner http.Handler) http.Handler {
		log.Print("A: called")
		mw := func(w http.ResponseWriter, r *http.Request) {
			log.Print("A: before")
			inner.ServeHTTP(w, r)
			log.Print("A: after")
		}
		return http.HandlerFunc(mw)
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
only see the portion above the divider a single time.

Middleware in Goji is called after routing has been performed. Therefore it is
possible to examine any routing information placed into the context by Patterns,
or to view or modify the Handler that will be routed to.

The http.Handler returned by the given middleware must be safe for concurrent
use by multiple goroutines. It is not safe to concurrently register middleware
from multiple goroutines.
*/
func (m *Mux) Use(middleware func(http.Handler) http.Handler) {
}

/*
UseC appends a context-aware middleware to the Mux's middleware stack. See the
documentation for Use for more information about the semantics of middleware.

The Handler returned by the given middleware must be safe for concurrent use by
multiple goroutines. It is not safe to concurrently register middleware from
multiple goroutines.
*/
func (m *Mux) UseC(middleware func(Handler) Handler) {
}
