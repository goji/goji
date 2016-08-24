package goji

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/weave-lab/goji/internal"
)

func TestDispatch(t *testing.T) {
	t.Parallel()

	var d dispatch

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	d.ServeHTTP(w, r)
	if w.Code != 404 {
		t.Errorf("status: expected %d, got %d", 404, w.Code)
	}

	w = httptest.NewRecorder()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(123)
	})

	ctx := context.WithValue(context.Background(), internal.Handler, h)
	d.ServeHTTP(w, r.WithContext(ctx))
	if w.Code != 123 {
		t.Errorf("status: expected %d, got %d", 123, w.Code)
	}
}
