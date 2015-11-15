package goji

import (
	"net/http"
	"net/http/httptest"

	"golang.org/x/net/context"
)

type boolPattern bool

func (b boolPattern) Match(ctx context.Context, _ *http.Request) context.Context {
	if b {
		return ctx
	}
	return nil
}

func wr() (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		panic(err)
	}
	return w, r
}
