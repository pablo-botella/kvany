package kvany

import (
	"encoding/json"
	"fmt"
)

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
