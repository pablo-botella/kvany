package kvany

import (
	"fmt"
	"reflect"
	"strings"
)

// VPeek returns the value of one key converted to the type of
// defaultValue, or defaultValue itself when the key is not in the list.
// It is PeekValues for a single value with the cast built in: the
// default is also the sample of the wanted type, and opt carries the
// cast options. A nil defaultValue asks for no conversion: the value
// comes back as it is.
//
// The key is matched exactly, or case-insensitively with
// PeekValuesCaseInsensitive. A repeated key is an error unless a Dupe
// flag says what to do: DupeFirst and DupeLast take one; DupeArray
// returns every value in a slice, so defaultValue must be a slice (or
// nil, for a []any): each value is cast to its element type, and a
// missing key returns defaultValue as with any other type.
// GivenCase, FirstFoundCase, LastFoundCase and SkipNotFound have no
// meaning here and are ignored.
func (l Lst) VPeek(key string, opt CastFlags, defaultValue any) (any, error) {
	ci := opt.Peek&PeekValuesCaseInsensitive != 0
	dupe := opt.Peek & 0xF000
	if dupe == PeekValuesDupeArray && defaultValue != nil && reflect.TypeOf(defaultValue).Kind() != reflect.Slice {
		return nil, fmt.Errorf("kvany: VPeek %q: PeekValuesDupeArray needs a slice as default, not %T", key, defaultValue)
	}
	var found []*Kv
	for i := range l {
		kv := &l[i]
		if kv.K == key || (ci && strings.EqualFold(kv.K, key)) {
			found = append(found, kv)
		}
	}
	if len(found) == 0 {
		return defaultValue, nil
	}
	var kv *Kv
	switch {
	case dupe == PeekValuesDupeArray:
		if defaultValue == nil {
			vals := make([]any, len(found))
			for i, f := range found {
				vals[i] = f.V
			}
			return vals, nil
		}
		st := reflect.TypeOf(defaultValue)
		like := reflect.Zero(st.Elem()).Interface() // the sample of each element
		out := reflect.MakeSlice(st, 0, len(found))
		for _, f := range found {
			v, err := f.Cast(like, opt)
			if err != nil {
				return nil, err
			}
			out = reflect.Append(out, reflect.ValueOf(v))
		}
		return out.Interface(), nil
	case dupe == PeekValuesDupeFirst:
		kv = found[0]
	case dupe == PeekValuesDupeLast:
		kv = found[len(found)-1]
	case len(found) > 1:
		return nil, fmt.Errorf("kvany: key %q found %d times", key, len(found))
	default:
		kv = found[0]
	}
	return kv.Cast(defaultValue, opt)
}
