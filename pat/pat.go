/*
Package pat is a URL-matching domain-specific language for Goji.


Quick Reference

The following table gives an overview of the language this package accepts. See
the subsequent sections for a more detailed explanation of what each pattern
does.

	Pattern			Matches			Does Not Match

	/			/			/hello

	/hello			/hello			/hi
							/hello/

	/user/:name		/user/carl		/user/carl/photos
				/user/alice		/user/carl/
							/user/

	/:file.:ext		/data.json		/.json
				/info.txt		/data.
				/data.tar.gz		/data.json/download

	/user/*			/user/			/user
				/user/carl
				/user/carl/photos


Static Paths

Most URL paths may be specified directly: the pattern "/hello" matches URLs with
precisely that path ("/hello/", for instance, is treated as distinct).

Note that this package operates on raw (i.e., escaped) paths (see the
documentation for net/url.URL.EscapedPath). In order to match a character that
can appear escaped in a URL path, use its percent-encoded form.


Named Matches

Named matches allow URL paths to contain any value in a particular path segment.
Such matches are denoted by a leading ":", for example ":name" in the rule
"/user/:name", and permit any non-empty value in that position. For instance, in
the previous "/user/:name" example, the path "/user/carl" is matched, while
"/user/" or "/user/carl/" (note the trailing slash) are not matched. Pat rules
can contain any number of named matches.

Named matches set URL variables by comparing pattern names to the segments they
matched. In our "/user/:name" example, a request for "/user/carl" would bind the
"name" variable to the value "carl". Use the Param function to extract these
variables from the request context.

Matches are ordinarily delimited by slashes ("/"), but several other characters
are accepted as delimiters (with slightly different semantics): the period
("."), semicolon (";"), and comma (",") characters. For instance, given the
pattern "/:file.:ext", the request "/data.json" would match, binding "file" to
"data" and "ext" to "json". Note that these special characters are treated
slightly differently than slashes: the above pattern also matches the path
"/data.tar.gz", with "ext" getting set to "tar.gz"; and the pattern "/:file"
matches names with dots in them (like "data.json").


Prefix Matches

Pat can also match prefixes of routes using wildcards. Prefix wildcard routes
end with "/*", and match just the path segments preceding the asterisk. For
instance, the pattern "/user/*" will match "/user/" and "/user/carl/photos" but
not "/user" (note the lack of a trailing slash).

The unmatched suffix, including the leading slash ("/"), are placed into the
request context, which allows subsequent routing (e.g., a subrouter) to continue
from where this pattern left off. For instance, in the "/user/*" pattern from
above, a request for "/user/carl/photos" will consume the "/user" prefix,
leaving the path "/carl/photos" for subsequent patterns to handle. A subrouter
pattern for "/:name/photos" would match this remaining path segment, for
instance.
*/
package pat

import (
	"net/http"
	"regexp"
	"strings"
	"sync"

	"goji.io/pattern"
	"golang.org/x/net/context"
)

/*
Pattern implements goji.Pattern using a path-matching domain specific language.
See the package documentation for more information about the semantics of this
object.
*/
type Pattern struct {
	raw string
	// These are parallel arrays of each pattern string (sans ":"), the
	// breaks each expect afterwords (used to support e.g., "." dividers),
	// and the string literals in between every pattern. There is always one
	// more literal than pattern, and they are interleaved like this:
	// <literal> <pattern> <break> <literal> <pattern> <break> <literal> etc...
	pats     []string
	breaks   []byte
	literals []string
	wildcard bool
	// This is used to store scratch space for the Match function.
	pool sync.Pool
}

// "Break characters" are characters that can end patterns. They are not allowed
// to appear in pattern names. "/" was chosen because it is the standard path
// separator, and "." was chosen because it often delimits file extensions. ";"
// and "," were chosen because Section 3.3 of RFC 3986 suggests their use.
const bc = "/.;,"

var patternRe = regexp.MustCompile(`[` + bc + `]:([^` + bc + `]+)`)

/*
New returns a new Pattern from the given Pat route. See the package
documentation for more information about what syntax is accepted by this
function.
*/
func New(pat string) *Pattern {
	p := &Pattern{raw: pat}

	if strings.HasSuffix(pat, "/*") {
		pat = pat[:len(pat)-1]
		p.wildcard = true
	}

	matches := patternRe.FindAllStringSubmatchIndex(pat, -1)
	numMatches := len(matches)
	p.pats = make([]string, numMatches)
	p.breaks = make([]byte, numMatches)
	p.literals = make([]string, numMatches+1)

	n := 0
	for i, match := range matches {
		a, b := match[2], match[3]
		p.literals[i] = pat[n : a-1] // Need to leave off the colon
		p.pats[i] = pat[a:b]
		if b == len(pat) {
			p.breaks[i] = '/'
		} else {
			p.breaks[i] = pat[b]
		}
		n = b
	}
	p.literals[numMatches] = pat[n:]

	p.pool.New = func() interface{} {
		return make([]string, 0, numMatches)
	}

	return p
}

/*
Match runs the Pat pattern on the given request, returning a non-nil context if
the request matches the request.

This function satisfies goji.Pattern.
*/
func (p *Pattern) Match(ctx context.Context, r *http.Request) context.Context {
	// In order to avoid doing a second pass over the path, there is a
	// period of time where we need scratch space to store matches we have
	// accumulated so far, but where we don't know that the path matches
	// completely (and therefore before we are forced to commit resources to
	// store the match). This of course isn't necessary, but gives us a
	// substantial bump in benchmarks, so we may as well.
	v := p.pool.Get()
	scratch := v.([]string)
	// Ordinarily we'd defer this, but doing so appreciably adds to the
	// runtime of this function in benchmarks.
	cleanup := func() {
		for i := range scratch {
			scratch[i] = ""
		}
		p.pool.Put(v)
	}
	path := pattern.Path(ctx)

	for i := range p.pats {
		sli := p.literals[i]
		if !strings.HasPrefix(path, sli) {
			cleanup()
			return nil
		}
		path = path[len(sli):]

		m := 0
		bc := p.breaks[i]
		for ; m < len(path); m++ {
			if path[m] == bc || path[m] == '/' {
				break
			}
		}
		if m == 0 {
			// Empty strings are not matches, otherwise routes like
			// "/:foo" would match the path "/"
			cleanup()
			return nil
		}
		scratch = append(scratch, path[:m])
		path = path[m:]
	}

	// There's exactly one more literal than pat.
	tail := p.literals[len(p.pats)]
	if p.wildcard {
		if !strings.HasPrefix(path, tail) {
			cleanup()
			return nil
		}
		scratch = append(scratch, path[len(tail)-1:])
	} else if path != tail {
		cleanup()
		return nil
	}

	var storage pattern.Storage
	for i, pat := range p.pats {
		unescaped, err := unescape(scratch[i])
		if err != nil {
			// If we encounter an encoding error here, there's
			// really not much we can do about it with our current
			// API, and I'm not really interested in supporting
			// clients that misencode URLs anyways.
			cleanup()
			return nil
		}
		storage.Set(pat, unescaped)
	}
	if p.wildcard {
		storage.SetPath(scratch[len(p.pats)])
	}

	cleanup()
	return storage.Bind(ctx)
}

/*
PathPrefix returns a string prefix that the Paths of all requests that this
Pattern accepts must contain.

This function satisfies goji's PathPrefix Pattern optimization.
*/
func (p *Pattern) PathPrefix() string {
	return p.literals[0]
}

/*
String returns the pattern string that was used to create this Pattern.
*/
func (p *Pattern) String() string {
	return p.raw
}

/*
Param returns the bound parameter with the given name. For instance, given the
route:

	/user/:name

and the URL Path:

	/user/carl

a call to Param(ctx, "name") would return the string "carl". It is the caller's
responsibility to ensure that the variable has been bound. Attempts to access
variables that have not been set (or which have been invalidly set) are
considered programmer errors and will trigger a panic.
*/
func Param(ctx context.Context, name string) string {
	return ctx.Value(pattern.Variable(name)).(string)
}
