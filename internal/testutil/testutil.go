package testutil

import (
	"fmt"
	"slices"
	"testing"
)

func Eq[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func EqVal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func NonNilEmpty[T any](t *testing.T, got []T) {
	t.Helper()
	if got == nil || len(got) != 0 {
		t.Fatalf("got %#v, want non-nil empty slice", got)
	}
}

func MustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	fn()
}

// Ident is a fmt.Stringer int used to test String() formatting.
type Ident int

func (id Ident) String() string {
	return fmt.Sprintf("id=%d", id)
}
