package pattern

import (
	"net/http"
	"reflect"
	"testing"

	"golang.org/x/net/context"
)

func TestComposeBasic(t *testing.T) {
	t.Parallel()

	pat := Compose()
	req, _ := http.NewRequest("GET", "/", nil)
	if ctx := pat.Match(context.Background(), req); ctx == nil {
		t.Error("expected composition of nothing to succeed")
	}

	pat = Compose(boolPattern(true))
	if ctx := pat.Match(context.Background(), req); ctx == nil {
		t.Error("expected composition of true to succeed")
	}

	pat = Compose(boolPattern(false))
	if ctx := pat.Match(context.Background(), req); ctx != nil {
		t.Error("expected composition of false to fail")
	}

	pat = Compose(boolPattern(true), boolPattern(true))
	if ctx := pat.Match(context.Background(), req); ctx == nil {
		t.Error("expected composition of true and true to succeed")
	}

	pat = Compose(boolPattern(true), boolPattern(false))
	if ctx := pat.Match(context.Background(), req); ctx != nil {
		t.Error("expected composition of true and false to fail")
	}

	pat = Compose(boolPattern(false), boolPattern(true))
	if ctx := pat.Match(context.Background(), req); ctx != nil {
		t.Error("expected composition of false and true to fail")
	}
}

func TestComposeHTTPMethods(t *testing.T) {
	t.Parallel()

	get := withMethods{"GET": struct{}{}}
	getpost := withMethods{"GET": struct{}{}, "POST": struct{}{}}
	getput := withMethods{"GET": struct{}{}, "PUT": struct{}{}}
	pat := Compose(get, getpost, getput, boolPattern(true), withMethods(nil))

	hm := pat.(httpMethods).HTTPMethods()
	expected := map[string]struct{}{"GET": struct{}{}}
	if !reflect.DeepEqual(hm, expected) {
		t.Errorf("expected %v, got %v", expected, hm)
	}

	pat = Compose(boolPattern(true), withMethods(nil))
	hm = pat.(httpMethods).HTTPMethods()
	if hm != nil {
		t.Errorf("expected nil, got %v", hm)
	}
}

func TestComposePathPrefix(t *testing.T) {
	t.Parallel()

	pat := Compose(boolPattern(true), prefixPattern("hi"), prefixPattern("bye"))
	if pp := pat.(pathPrefix).PathPrefix(); pp != "hi" {
		t.Errorf("expected hi, got %q", pp)
	}

	pat = Compose(boolPattern(true))
	if pp := pat.(pathPrefix).PathPrefix(); pp != "" {
		t.Errorf("expected empty prefix, got %q", pp)
	}
}
