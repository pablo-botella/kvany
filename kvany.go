// Package kvany holds key/value pairs whose values may be anything, kept
// in the order they came: the Go side of the [["KEY", value], ...] arrays
// that travel in JSON between Go and Xbase++ (ot4xb). Nothing is
// normalized and nothing is deduplicated: a key may appear more than once.
package kvany

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PeekValuesFlag int

const (
	PeekValuesNone            PeekValuesFlag = 0x00000000
	PeekValuesCaseInsensitive PeekValuesFlag = 0x00000001
	PeekValuesGivenCase       PeekValuesFlag = 0x00000001
	PeekValuesFirstFoundCase  PeekValuesFlag = 0x00000011
	PeekValuesLastFoundCase   PeekValuesFlag = 0x00000021
	PeekValuesDupeFirst       PeekValuesFlag = 0x00001000
	PeekValuesDupeLast        PeekValuesFlag = 0x00002000
	PeekValuesDupeArray       PeekValuesFlag = 0x00004000
	PeekValuesDupeError       PeekValuesFlag = 0x00008000
	PeekValuesSkipNotFound    PeekValuesFlag = 0x00010000
)

// The docs are composed from _mkskill/ - edit the sources there, then:
//go:generate go run github.com/pablo-botella/mkskill/cmd/mkskill@latest -q build
//go:generate go run github.com/pablo-botella/mkskill/cmd/mkskill@latest -q -vbuild

// Kv is one pair: a string key and any value (nil included).
type Kv struct {
	K string
	V any
}

// Lst is a list of pairs in the order they came; keys may repeat.
type Lst []Kv

// PtrLst is a list of pointers to pairs, for when the pairs are shared
// and mutated while the list lives on.
type PtrLst []*Kv

// PeekValues returns a new Lst with the pairs matching the given keys,
// honoring the flags. Keys are visited in key_list order; for each one the
// list is scanned in its own order. A key with no match comes out as
// Kv{key, nil}, or is skipped with PeekValuesSkipNotFound.
//
// Matching: exact, or case-insensitive with PeekValuesCaseInsensitive; in
// that case the output key is the one given (GivenCase), the first
// found (FirstFoundCase) or the last found (LastFoundCase).
//
// Repeats: with no Dupe flag every match comes out, in order. DupeFirst
// and DupeLast keep one; DupeArray keeps one pair whose V is a []any with
// every value; DupeError fails when a key matches more than once.
func (l Lst) PeekValuesToLst(key_list []string, flags PeekValuesFlag) (Lst, error) {
	out := Lst{}
	ci := flags&PeekValuesCaseInsensitive != 0
	caseMode := flags & 0x30
	dupeMode := flags & 0xF000
	for _, key := range key_list {
		var found Lst
		for _, kv := range l {
			if kv.K == key || (ci && strings.EqualFold(kv.K, key)) {
				found = append(found, kv)
			}
		}
		if len(found) == 0 {
			if flags&PeekValuesSkipNotFound == 0 {
				out = append(out, Kv{K: key, V: nil})
			}
			continue
		}
		outKey := key
		if ci {
			switch caseMode {
			case PeekValuesFirstFoundCase & 0x30:
				outKey = found[0].K
			case PeekValuesLastFoundCase & 0x30:
				outKey = found[len(found)-1].K
			}
		}
		switch dupeMode {
		case PeekValuesDupeFirst:
			out = append(out, Kv{K: outKey, V: found[0].V})
		case PeekValuesDupeLast:
			out = append(out, Kv{K: outKey, V: found[len(found)-1].V})
		case PeekValuesDupeArray:
			vals := make([]any, len(found))
			for i, kv := range found {
				vals[i] = kv.V
			}
			out = append(out, Kv{K: outKey, V: vals})
		case PeekValuesDupeError:
			if len(found) > 1 {
				return nil, fmt.Errorf("kvany: key %q found %d times", key, len(found))
			}
			out = append(out, Kv{K: outKey, V: found[0].V})
		default:
			for _, kv := range found {
				out = append(out, Kv{K: outKey, V: kv.V})
			}
		}
	}
	return out, nil
}

// PeekValuesToMap is PeekValuesToLst folded into a map: same keys, same flags,
// same errors. Pairs that share an output key collapse into one entry and
// the last one wins, so without a Dupe flag a repeated key keeps its last
// value; DupeFirst, DupeLast and DupeArray decide it explicitly.
func (l Lst) PeekValuesToMap(key_list []string, flags PeekValuesFlag) (map[string]any, error) {
	lst, err := l.PeekValuesToLst(key_list, flags)
	if err != nil {
		return nil, err
	}
	out := make(map[string]any, len(lst))
	for _, kv := range lst {
		out[kv.K] = kv.V
	}
	return out, nil
}

// MarshalJSON writes the pair as ["KEY", value]: the two-element array
// that travels to Xbase++.
func (kv Kv) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]any{kv.K, kv.V})
}

// UnmarshalJSON reads ["KEY", value]. Anything else fails: an object, an
// array that is not two elements long, or a key that is not a string.
// The value is whatever encoding/json makes of it (numbers as float64).
func (kv *Kv) UnmarshalJSON(data []byte) error {
	var pair []json.RawMessage
	if err := json.Unmarshal(data, &pair); err != nil {
		return fmt.Errorf("kvany: pair is not an array: %w", err)
	}
	if len(pair) != 2 {
		return fmt.Errorf("kvany: pair has %d elements, want 2", len(pair))
	}
	var k string
	if err := json.Unmarshal(pair[0], &k); err != nil {
		return fmt.Errorf("kvany: pair key is not a string: %w", err)
	}
	var v any
	if err := json.Unmarshal(pair[1], &v); err != nil {
		return fmt.Errorf("kvany: pair %q: %w", k, err)
	}
	kv.K, kv.V = k, v
	return nil
}

// MarshalJSON writes the list as [["KEY", value], ...]. A nil list is [],
// never null: the other side always expects an array.
func (l Lst) MarshalJSON() ([]byte, error) {
	if l == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]Kv(l))
}

// UnmarshalJSON reads [["KEY", value], ...] in order. null leaves the list
// nil; [] leaves it empty but not nil.
func (l *Lst) UnmarshalJSON(data []byte) error {
	var pairs []Kv
	if err := json.Unmarshal(data, &pairs); err != nil {
		return err
	}
	*l = Lst(pairs)
	return nil
}

// MarshalJSON writes the list as [["KEY", value], ...]; a nil list is [].
func (l PtrLst) MarshalJSON() ([]byte, error) {
	if l == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]*Kv(l))
}
