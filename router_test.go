package goji

import (
	"net/http"
	"testing"

	"goji.io/internal"
	"golang.org/x/net/context"
)

func TestNoMatch(t *testing.T) {
	t.Parallel()

	m := NewMux()
	m.Handle(boolPattern(false), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("did not expect handler to be called")
	}))
	_, r := wr()
	ctx := m.router(context.Background(), r)

	if p := ctx.Value(internal.Pattern); p != nil {
		t.Errorf("unexpected pattern %v", p)
	}
	if h := ctx.Value(internal.Handler); h != nil {
		t.Errorf("unexpected handler %v", h)
	}
}
