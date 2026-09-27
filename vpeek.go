package kvany

import (
	"fmt"
	"strings"
)

// VPeek returns the value of one key, or defaultValue when the key is not
// in the list. It is PeekValues for a single value. With
// PeekValuesCastToDefault in opt.Peek the value found is converted to the
// type of defaultValue through Cast (the default is then also the sample
// of the wanted type, and opt carries the cast options); without the flag
// the value comes back as it came.
//
// The key is matched exactly, or case-insensitively with
// PeekValuesCaseInsensitive. A repeated key is an error unless a Dupe
// flag says what to do: DupeFirst and DupeLast take one; DupeArray
// returns every value as a []any, and is incompatible with
// CastToDefault. GivenCase, FirstFoundCase, LastFoundCase and
// SkipNotFound have no meaning here and are ignored.
func (l Lst) VPeek(key string, opt CastFlags, defaultValue any) (any, error) {
	ci := opt.Peek&PeekValuesCaseInsensitive != 0
	cast := opt.Peek&PeekValuesCastToDefault != 0
	dupe := opt.Peek & 0xF000
	if dupe == PeekValuesDupeArray && cast {
		return nil, fmt.Errorf("kvany: VPeek %q: PeekValuesDupeArray and PeekValuesCastToDefault are incompatible", key)
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
		vals := make([]any, len(found))
		for i, f := range found {
			vals[i] = f.V
		}
		return vals, nil
	case dupe == PeekValuesDupeFirst:
		kv = found[0]
	case dupe == PeekValuesDupeLast:
		kv = found[len(found)-1]
	case len(found) > 1:
		return nil, fmt.Errorf("kvany: key %q found %d times", key, len(found))
	default:
		kv = found[0]
	}
	if !cast {
		return kv.V, nil
	}
	return kv.Cast(defaultValue, opt)
}

func (l Lst) VPeekString(key string, opt CastFlags, defaultValue string) (string, error) {
	val, err := l.VPeek(key, opt, defaultValue)
	if err != nil {
		return "", err
	}
	if s, ok := val.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("kvany: VPeekString %q: value is not a string", key)
}
func (l Lst) VPeekInt(key string, opt CastFlags, defaultValue int) (int, error) {
	val, err := l.VPeek(key, opt, defaultValue)
	if err != nil {
		return 0, err
	}
	if i, ok := val.(int); ok {
		return i, nil
	}
	return 0, fmt.Errorf("kvany: VPeekInt %q: value is not an int", key)
}
func (l Lst) VPeekFloat64(key string, opt CastFlags, defaultValue float64) (float64, error) {
	val, err := l.VPeek(key, opt, defaultValue)
	if err != nil {
		return 0, err
	}
	if f, ok := val.(float64); ok {
		return f, nil
	}
	return 0, fmt.Errorf("kvany: VPeekFloat64 %q: value is not a float64", key)
}
func (l Lst) VPeekBool(key string, opt CastFlags, defaultValue bool) (bool, error) {
	val, err := l.VPeek(key, opt, defaultValue)
	if err != nil {
		return false, err
	}
	if b, ok := val.(bool); ok {
		return b, nil
	}
	return false, fmt.Errorf("kvany: VPeekBool %q: value is not a bool", key)
}
