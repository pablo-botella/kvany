package kvany_test

import (
	"reflect"
	"testing"

	"github.com/pablo-botella/kvany"
)

// src has a repeated key (CODE, three times, mixed case) and a lone one.
var src = kvany.Lst{
	{K: "CODE", V: "a"},
	{K: "QTY", V: 3.0},
	{K: "code", V: "b"},
	{K: "Code", V: "c"},
}

func peek(t *testing.T, keys []string, flags kvany.PeekValuesFlag) kvany.Lst {
	t.Helper()
	out, err := src.PeekValuesToLst(keys, flags)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func want(t *testing.T, got, want kvany.Lst) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got  %#v\nwant %#v", got, want)
	}
}

func TestExactMatchKeepsEveryRepeat(t *testing.T) {
	want(t, peek(t, []string{"QTY", "CODE"}, kvany.PeekValuesNone),
		kvany.Lst{{K: "QTY", V: 3.0}, {K: "CODE", V: "a"}})
	want(t, peek(t, []string{"code"}, kvany.PeekValuesNone),
		kvany.Lst{{K: "code", V: "b"}})
}

func TestNotFound(t *testing.T) {
	want(t, peek(t, []string{"NOPE", "QTY"}, kvany.PeekValuesNone),
		kvany.Lst{{K: "NOPE", V: nil}, {K: "QTY", V: 3.0}})
	want(t, peek(t, []string{"NOPE", "QTY"}, kvany.PeekValuesSkipNotFound),
		kvany.Lst{{K: "QTY", V: 3.0}})
}

func TestEmptyKeys(t *testing.T) {
	out := peek(t, nil, kvany.PeekValuesNone)
	if out == nil || len(out) != 0 {
		t.Fatalf("got %#v, want empty non-nil", out)
	}
}

func TestCaseInsensitiveOutputKey(t *testing.T) {
	keys := []string{"cOdE"}
	want(t, peek(t, keys, kvany.PeekValuesGivenCase),
		kvany.Lst{{K: "cOdE", V: "a"}, {K: "cOdE", V: "b"}, {K: "cOdE", V: "c"}})
	want(t, peek(t, keys, kvany.PeekValuesFirstFoundCase|kvany.PeekValuesDupeFirst),
		kvany.Lst{{K: "CODE", V: "a"}})
	want(t, peek(t, keys, kvany.PeekValuesLastFoundCase|kvany.PeekValuesDupeLast),
		kvany.Lst{{K: "Code", V: "c"}})
}

func TestDupePolicies(t *testing.T) {
	keys := []string{"code"}
	ci := kvany.PeekValuesCaseInsensitive
	want(t, peek(t, keys, ci|kvany.PeekValuesDupeFirst), kvany.Lst{{K: "code", V: "a"}})
	want(t, peek(t, keys, ci|kvany.PeekValuesDupeLast), kvany.Lst{{K: "code", V: "c"}})
	want(t, peek(t, keys, ci|kvany.PeekValuesDupeArray),
		kvany.Lst{{K: "code", V: []any{"a", "b", "c"}}})
	// a single match under DupeArray still comes out as a one-element array
	want(t, peek(t, []string{"QTY"}, kvany.PeekValuesDupeArray),
		kvany.Lst{{K: "QTY", V: []any{3.0}}})
}

func TestDupeError(t *testing.T) {
	if _, err := src.PeekValuesToLst([]string{"code"}, kvany.PeekValuesCaseInsensitive|kvany.PeekValuesDupeError); err == nil {
		t.Fatal("expected an error on a repeated key")
	}
	// exact match: "code" appears once, no error
	want(t, peek(t, []string{"code"}, kvany.PeekValuesDupeError), kvany.Lst{{K: "code", V: "b"}})
}

func TestSourceUntouched(t *testing.T) {
	before := append(kvany.Lst{}, src...)
	peek(t, []string{"CODE", "QTY"}, kvany.PeekValuesCaseInsensitive|kvany.PeekValuesDupeArray)
	want(t, src, before)
}

func peekMap(t *testing.T, keys []string, flags kvany.PeekValuesFlag) map[string]any {
	t.Helper()
	out, err := src.PeekValuesToMap(keys, flags)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func wantMap(t *testing.T, got, want map[string]any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got  %#v\nwant %#v", got, want)
	}
}

func TestMapFollowsLst(t *testing.T) {
	wantMap(t, peekMap(t, []string{"QTY", "CODE"}, kvany.PeekValuesNone),
		map[string]any{"QTY": 3.0, "CODE": "a"})
	wantMap(t, peekMap(t, []string{"NOPE", "QTY"}, kvany.PeekValuesNone),
		map[string]any{"NOPE": nil, "QTY": 3.0})
	wantMap(t, peekMap(t, []string{"NOPE", "QTY"}, kvany.PeekValuesSkipNotFound),
		map[string]any{"QTY": 3.0})
}

func TestMapRepeatsCollapseLastWins(t *testing.T) {
	// no Dupe flag: the Lst carries a, b, c under one key; the map keeps c
	wantMap(t, peekMap(t, []string{"code"}, kvany.PeekValuesCaseInsensitive),
		map[string]any{"code": "c"})
	wantMap(t, peekMap(t, []string{"code"}, kvany.PeekValuesCaseInsensitive|kvany.PeekValuesDupeFirst),
		map[string]any{"code": "a"})
	wantMap(t, peekMap(t, []string{"code"}, kvany.PeekValuesCaseInsensitive|kvany.PeekValuesDupeArray),
		map[string]any{"code": []any{"a", "b", "c"}})
	// FirstFoundCase: the map key is the found spelling, not the given one
	wantMap(t, peekMap(t, []string{"code"}, kvany.PeekValuesFirstFoundCase|kvany.PeekValuesDupeFirst),
		map[string]any{"CODE": "a"})
}

func TestMapEmptyAndError(t *testing.T) {
	out := peekMap(t, nil, kvany.PeekValuesNone)
	if out == nil || len(out) != 0 {
		t.Fatalf("got %#v, want empty non-nil", out)
	}
	if m, err := src.PeekValuesToMap([]string{"code"}, kvany.PeekValuesCaseInsensitive|kvany.PeekValuesDupeError); err == nil || m != nil {
		t.Fatalf("got %#v, %v; want nil map and an error", m, err)
	}
}
