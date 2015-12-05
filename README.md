Goji
====

Goji is a HTTP request multiplexer, similar to [`net/http.ServeMux`][servemux].
It compares incoming requests to a list of registered [Patterns][pattern], and
dispatches to the [Handler][handler] that corresponds to the first matching
Pattern. Goji also supports [Middleware][middleware] (composable shared
functionality applied to every request) and uses the de facto standard
[`x/net/context`][context] to store request-scoped values.

[servemux]: https://golang.org/pkg/net/http/#ServeMux
[pattern]: https://godoc.org/goji.io#Pattern
[handler]: https://godoc.org/goji.io#Handler
[middleware]: https://godoc.org/goji.io#Mux.Use
[context]: https://godoc.org/golang.org/x/net/context


Quick Start
-----------

```go
package main

import (
        "fmt"
        "net/http"

        "goji.io"
        "goji.io/pat"
        "golang.org/x/net/context"
)

func hello(ctx context.Context, w http.ResponseWriter, r *http.Request) {
        name := pat.Param(ctx, "name")
        fmt.Fprintf(w, "Hello, %s!", name)
}

func main() {
        mux := goji.NewMux()
        mux.HandleFuncC(pat.Get("/hello/:name"), hello)

        http.ListenAndServe("localhost:8000", mux)
}
```

Please refer to [Goji's GoDoc Documentation][godoc] for a full API reference.

[godoc]: https://godoc.org/goji.io


Stability
---------

As of this writing (late November 2015), this version of Goji is still very new,
and during this initial experimental stage it offers no API stability
guarantees. After the API has had a little time to bake, Goji expects to adhere
strictly to the Go project's [compatibility guidelines][compat], guaranteeing to
never break compatibility with existing code.

We expect to be able to make such a guarantee by early 2016. Although we reserve
the right to do so, there are no breaking API changes planned until that point,
and we are unlikely to accept any such breaking changes.

[compat]: https://golang.org/doc/go1compat


Community / Contributing
------------------------

Goji maintains a mailing list, [gojiberries][berries], where you should feel
welcome to ask questions about the project (no matter how simple!), to announce
projects or libraries built on top of Goji, or to talk about Goji more
generally. Goji's author (Carl Jackson) also loves to hear from users directly
at his personal email address, which is available on his GitHub profile page.

Contributions to Goji are welcome, however please be advised that due to Goji's
stability guarantees interface changes are unlikely to be accepted.

All interactions in the Goji community will be held to the high standard of the
broader Go community's [Code of Conduct][conduct].

[berries]: https://groups.google.com/forum/#!forum/gojiberries
[conduct]: https://golang.org/conduct
