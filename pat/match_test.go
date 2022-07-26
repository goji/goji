package pat

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"goji.io/pattern"
)

type testContextKey struct{}

func TestExistingContext(t *testing.T) {
	t.Parallel()

	pat := New("/hi/:c/:a/:r/:l")
	req, err := http.NewRequest("GET", "/hi/foo/bar/baz/quux", nil)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	ctx = pattern.SetPath(ctx, req.URL.EscapedPath())
	ctx = context.WithValue(ctx, pattern.AllVariables, map[pattern.Variable]interface{}{
		"hello": "world",
		"c":     "nope",
	})
	ctx = context.WithValue(ctx, pattern.Variable("user"), "carl")

	req = req.WithContext(ctx)
	req = pat.Match(req)
	if req == nil {
		t.Fatalf("expected pattern to match")
	}
	ctx = req.Context()

	expected := map[pattern.Variable]interface{}{
		"c": "foo",
		"a": "bar",
		"r": "baz",
		"l": "quux",
	}
	for k, v := range expected {
		if p := Param(req, string(k)); p != v {
			t.Errorf("expected %s=%q, got %q", k, v, p)
		}
	}

	expected["hello"] = "world"
	all := ctx.Value(pattern.AllVariables).(map[pattern.Variable]interface{})
	if !reflect.DeepEqual(all, expected) {
		t.Errorf("expected %v, got %v", expected, all)
	}

	if path := pattern.Path(ctx); path != "" {
		t.Errorf("expected path=%q, got %q", "", path)
	}

	if user := ctx.Value(pattern.Variable("user")); user != "carl" {
		t.Errorf("expected user=%q, got %q", "carl", user)
	}
}

func TestMatchValueDoesntAllocate(t *testing.T) {
	t.Parallel()

	pat := New("/*")
	req, err := http.NewRequest("GET", "/hithere", nil)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	ctx = pattern.SetPath(ctx, req.URL.EscapedPath())
	req = req.WithContext(ctx)
	req = pat.Match(req)
	if req == nil {
		t.Fatalf("expected pattern to match")
	}
	ctx = req.Context()
	// add an extra context layer to ensure all our Value requests
	// bounce through the generic context.Context.Value implementation
	ctx = context.WithValue(ctx, testContextKey{}, "huzzah!")

	if all := ctx.Value(pattern.AllVariables); all != nil {
		t.Errorf("expected all variable to be nil, got %v", all)
	}

	allocs := testing.AllocsPerRun(1, func() {
		if path := pattern.Path(ctx); path != "/hithere" {
			t.Errorf("expected path=%q, got %q", "", path)
		}
	})

	if allocs != 0 {
		t.Errorf("expected 0 allocs, not %f", allocs)
	}
}
