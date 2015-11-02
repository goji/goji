package pattern

import (
	"reflect"
	"testing"

	"golang.org/x/net/context"
)

func TestStorage(t *testing.T) {
	t.Parallel()

	var s Storage
	s.Set("hello", "world")
	s.Set("number", 4)
	s.SetPath("/test")

	ctx := s.Bind(context.Background())
	if path := Path(ctx); path != "/test" {
		t.Errorf("expected path of /test, got %q", path)
	}

	if hello := ctx.Value(Variable("hello")).(string); hello != "world" {
		t.Errorf("expected hello=world, got %q", hello)
	}
	if number := ctx.Value(Variable("number")).(int); number != 4 {
		t.Errorf("expected number=4, got %d", number)
	}

	expected := map[Variable]interface{}{
		"hello":  "world",
		"number": 4,
	}
	if all := ctx.Value(AllVariables); !reflect.DeepEqual(expected, all) {
		t.Errorf("expected all=%v, got %v", expected, all)
	}

	s.Set("hello", "this should have no effect")
	if hello := ctx.Value(Variable("hello")).(string); hello != "world" {
		t.Errorf("expected hello=world, got %q", hello)
	}
}

func TestStorageOverflow(t *testing.T) {
	t.Parallel()

	var s Storage
	s.Set("one", 1)
	s.Set("two", 2)
	s.Set("three", 3)
	s.Set("four", 4)
	s.Set("five", 5)
	s.Set("six", 6)
	s.Set("seven", 7)

	ctx := s.Bind(context.Background())
	if path := Path(ctx); path != "" {
		t.Errorf("expected path of /test, got %q", path)
	}

	expected := map[Variable]interface{}{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
		"six":   6,
		"seven": 7,
	}
	if all := ctx.Value(AllVariables); !reflect.DeepEqual(expected, all) {
		t.Errorf("expected all=%v, got %v", expected, all)
	}
}

func TestStorageNesting(t *testing.T) {
	t.Parallel()

	var s1, s2 Storage
	s1.Set("a", "a")
	s1.Set("b", "b")
	s1.Set("c", "c")
	s1.Set("d", "d")
	s1.Set("e", "e")
	s1.Set("f", "f")
	s2.Set("g", "g")
	s2.Set("h", "h")

	ctx := s2.Bind(context.Background())
	ctx = s1.Bind(ctx)

	if v := ctx.Value(Variable("g")); v != "g" {
		t.Errorf("expected g, got %q", v)
	}
	if v := ctx.Value(Variable("d")); v != "d" {
		t.Errorf("expected d, got %q", v)
	}

	expected := map[Variable]interface{}{
		"a": "a",
		"b": "b",
		"c": "c",
		"d": "d",
		"e": "e",
		"f": "f",
		"g": "g",
		"h": "h",
	}
	if all := ctx.Value(AllVariables); !reflect.DeepEqual(expected, all) {
		t.Errorf("expected all=%v, got %v", expected, all)
	}
}
