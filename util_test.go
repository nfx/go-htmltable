package htmltable

import (
	"reflect"
	"testing"
)

func assertError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("expected no error, got %s", err.Error())
	}
}

func assertEqualError(t *testing.T, err error, msg string) {
	assertError(t, err)
	got := err.Error()
	if got != msg {
		t.Errorf("%#v (expected) != %#v (got)", msg, err.Error())
	}
}

func assertEqual(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%#v (expected) != %#v (got)", a, b)
	}
}

type comparable interface {
    int | string
}

func assertGreaterOrEqual[T comparable](t *testing.T, a, b T) {
	if !(a >= b) {
		t.Errorf("%#v (expected) >= %#v (got)", a, b)
	}
}
