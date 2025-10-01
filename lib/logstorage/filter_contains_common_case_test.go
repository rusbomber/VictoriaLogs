package logstorage

import (
	"reflect"
	"testing"
)

func TestGetCommonCasePhrases_Success(t *testing.T) {
	f := func(phrases, resultExpected []string) {
		t.Helper()

		result, err := getCommonCasePhrases(phrases)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if !reflect.DeepEqual(result, resultExpected) {
			t.Fatalf("unexpected result\ngot\n%q\nwant\n%q", result, resultExpected)
		}
	}

	f(nil, nil)
	f([]string{""}, []string{""})
	f([]string{"foo"}, []string{"FOO", "foo"})
	f([]string{"Foo"}, []string{"FOO", "Foo", "foo"})
	f([]string{"foo", "Foo"}, []string{"FOO", "Foo", "foo"})
	f([]string{"FOO"}, []string{"F O O", "F O o", "F OO", "F Oo", "F o O", "F o o", "F oO", "F oo", "FO O", "FO o", "FOO", "FOo", "Fo O", "Fo o", "FoO", "Foo", "f O O", "f O o", "f OO", "f Oo", "f o O", "f o o", "f oO", "f oo", "fO O", "fO o", "fOO", "fOo", "fo O", "fo o", "foO", "foo"})

	f([]string{"FooBar"}, []string{"FOO BAR", "FOOBAR", "Foo Bar", "Foo bar", "FooBar", "Foobar", "foo Bar", "foo bar", "fooBar", "foobar"})
	f([]string{"fooBar"}, []string{"FOO BAR", "FOOBAR", "foo Bar", "foo bar", "fooBar", "foobar"})
}

func TestGetCommonCasePhrases_Failure(t *testing.T) {
	f := func(phrases []string) {
		t.Helper()

		result, err := getCommonCasePhrases(phrases)
		if err == nil {
			t.Fatalf("expecting non-nil error")
		}
		if result != nil {
			t.Fatalf("expecting nil result; got %q", result)
		}
	}

	// More than 6 uppercase chars
	f([]string{"FOOBARB"})
	f([]string{"FoOOBbARrBzsdf"})
}
