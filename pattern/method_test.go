package pattern

import (
	"net/http"
	"reflect"
	"testing"

	"golang.org/x/net/context"
)

func TestWithMethods(t *testing.T) {
	t.Parallel()

	pat := WithMethods(boolPattern(true), "GET", "POST")
	req, _ := http.NewRequest("GET", "/", nil)
	if ctx := pat.Match(context.Background(), req); ctx == nil {
		t.Error("expected match on GET")
	}
	req, _ = http.NewRequest("POST", "/", nil)
	if ctx := pat.Match(context.Background(), req); ctx == nil {
		t.Error("expected match on POST")
	}
	req, _ = http.NewRequest("PUT", "/", nil)
	if ctx := pat.Match(context.Background(), req); ctx != nil {
		t.Error("expected no match on PUT")
	}

	hm := pat.(httpMethods).HTTPMethods()
	expected := map[string]struct{}{
		"GET":  struct{}{},
		"POST": struct{}{},
	}
	if !reflect.DeepEqual(hm, expected) {
		t.Errorf("expected %v, got %v", expected, hm)
	}
}
