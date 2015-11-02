package pattern

import (
	"net/http"

	"goji.io"
	"golang.org/x/net/context"
)

/*
Compose returns a new Pattern which is the composition of the given Patterns.
The patterns are run in the order in which they are provided, and the returned
Pattern only matches if every input pattern matches.

The returned Pattern supports both the HTTPMethods and PathPrefix Pattern
optimizations, taking the intersection of all supported methods and the first
PathPrefix.
*/
func Compose(patterns ...goji.Pattern) goji.Pattern {
	return compose(patterns)
}

type compose []goji.Pattern

func (c compose) Match(ctx context.Context, r *http.Request) context.Context {
	for _, p := range c {
		ctx = p.Match(ctx, r)
		if ctx == nil {
			return nil
		}
	}
	return ctx
}

type httpMethods interface {
	HTTPMethods() map[string]struct{}
}

type pathPrefix interface {
	PathPrefix() string
}

func (c compose) HTTPMethods() map[string]struct{} {
	out := make(map[string]struct{})
	hms := make([]map[string]struct{}, 0, len(c))
	for _, p := range c {
		if hm, ok := p.(httpMethods); ok {
			ms := hm.HTTPMethods()
			if ms == nil {
				continue
			}
			hms = append(hms, ms)
			for m := range ms {
				out[m] = struct{}{}
			}
		}
	}
	if len(hms) == 0 {
		return nil
	}

	for _, hm := range hms {
		for m := range out {
			if _, ok := hm[m]; !ok {
				delete(out, m)
			}
		}
	}
	return out
}

func (c compose) PathPrefix() string {
	for _, p := range c {
		if pp, ok := p.(pathPrefix); ok {
			return pp.PathPrefix()
		}
	}
	return ""
}

var _ httpMethods = compose{}
var _ pathPrefix = compose{}
