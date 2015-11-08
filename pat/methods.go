package pat

import (
	"goji.io"
	"goji.io/pattern"
)

/*
Delete returns a Pat route that only matches the DELETE HTTP method.
*/
func Delete(pat string) goji.Pattern {
	return pattern.WithMethods(New(pat), "DELETE")
}

/*
Get returns a Pat route that only matches the GET and HEAD HTTP method. HEAD
requests are handled transparently by net/http.
*/
func Get(pat string) goji.Pattern {
	return pattern.WithMethods(New(pat), "GET", "HEAD")
}

/*
Head returns a Pat route that only matches the HEAD HTTP method.
*/
func Head(pat string) goji.Pattern {
	return pattern.WithMethods(New(pat), "HEAD")
}

/*
Options returns a Pat route that only matches the OPTIONS HTTP method.
*/
func Options(pat string) goji.Pattern {
	return pattern.WithMethods(New(pat), "OPTIONS")
}

/*
Patch returns a Pat route that only matches the PATCH HTTP method.
*/
func Patch(pat string) goji.Pattern {
	return pattern.WithMethods(New(pat), "PATCH")
}

/*
Post returns a Pat route that only matches the POST HTTP method.
*/
func Post(pat string) goji.Pattern {
	return pattern.WithMethods(New(pat), "POST")
}

/*
Put returns a Pat route that only matches the PUT HTTP method.
*/
func Put(pat string) goji.Pattern {
	return pattern.WithMethods(New(pat), "PUT")
}
