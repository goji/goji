package goji

import (
	"net/http"
	"testing"

	"github.com/weave-lab/goji/internal"

	"context"
)

func TestMuxExistingPath(t *testing.T) {
	m := NewMux()
	handler := func(w http.ResponseWriter, r *http.Request) {
		if path := r.Context().Value(internal.Path).(string); path != "/" {
			t.Errorf("expected path=/, got %q", path)
		}
	}
	m.HandleFunc(boolPattern(true), handler)
	w, r := wr()
	ctx := context.WithValue(r.Context(), internal.Path, "/hello")
	m.ServeHTTP(w, r.WithContext(ctx))
}

func TestSubMuxExistingPath(t *testing.T) {
	m := SubMux()
	handler := func(w http.ResponseWriter, r *http.Request) {
		if path := r.Context().Value(internal.Path).(string); path != "/hello" {
			t.Errorf("expected path=/hello, got %q", path)
		}
	}
	m.HandleFunc(boolPattern(true), handler)
	w, r := wr()
	ctx := context.WithValue(r.Context(), internal.Path, "/hello")
	m.ServeHTTP(w, r.WithContext(ctx))
}
