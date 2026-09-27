package kvany

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// num is a value read as a number. isInt / isUint say it is an exact
// integer (then i or u is the one to use); otherwise f is all there is.
type num struct {
	f      float64
	i      int64
	u      uint64
	isInt  bool
	isUint bool
}

var timeType = reflect.TypeOf(time.Time{})

// ---- Cast

// Cast returns V converted to the type of like, which is only a sample of
// the wanted type; its value is ignored. A nil like asks for no
// conversion: V comes back as it is. A nil V is the zero value of the
// wanted type (an empty DBF field). The result carries the exact type of
// like, so the assertion on it is safe:
//
//	n := kv.Cast(int64(0), kvany.CastFlags{})   // n.(int64)
//
// What converts to what: any number, json.Number, bool (0 / 1) or text
// to any integer or float; anything to string; text or time.Time to
// time.Time; bool, number (0 = false) or text ("true", "1", "false",
// "0", "") to bool. A type with a name over one of those kinds (type
// Code string) is fine on either side. Anything else must already be of
// the wanted type.
func (kv Kv) Cast(like any, opt CastFlags) (any, error) {
	if like == nil {
		return kv.V, nil
	}
	want := reflect.TypeOf(like)
	v := kv.V
	if v == nil {
		return reflect.Zero(want).Interface(), nil
	}
	k := want.Kind()
	if reflect.TypeOf(v) == want && k != reflect.String && k != reflect.Float32 && k != reflect.Float64 {
		return v, nil // already the type; strings and floats still go through their options (trim, pad, decimals)
	}
	var out any
	var err error
	switch {
	case want == timeType:
		out, err = kv.castDate(v, opt.Date)
	case k == reflect.String:
		out, err = kv.castString(v, opt)
	case k == reflect.Bool:
		out, err = kv.castBool(v)
	case isIntKind(k):
		out, err = kv.castInt(v, want, opt.Int)
	case k == reflect.Float32 || k == reflect.Float64:
		out, err = kv.castFloat(v, want, opt.Float)
	default:
		err = fmt.Errorf("kvany: %q is %T, cannot cast to %s", kv.K, v, want)
	}
	if err != nil {
		return nil, err
	}
	if reflect.TypeOf(out) != want { // a named type over a basic kind
		out = reflect.ValueOf(out).Convert(want).Interface()
	}
	return out, nil
}

func isIntKind(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}
	return false
}

// numOf reads a Go number, a json.Number or a bool. Text is not read here.
func numOf(v any) (num, bool) {
	switch x := v.(type) {
	case json.Number:
		return parseNum(string(x))
	case bool:
		if x {
			return num{f: 1, i: 1, isInt: true}, true
		}
		return num{isInt: true}, true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := rv.Int()
		return num{f: float64(i), i: i, isInt: true}, true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := rv.Uint()
		return num{f: float64(u), u: u, isUint: true}, true
	case reflect.Float32, reflect.Float64:
		return num{f: rv.Float()}, true
	case reflect.Bool:
		if rv.Bool() {
			return num{f: 1, i: 1, isInt: true}, true
		}
		return num{isInt: true}, true
	}
	return num{}, false
}

// parseNum reads a number written as text, spaces around tolerated.
func parseNum(s string) (num, bool) {
	s = strings.TrimSpace(s)
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return num{f: float64(i), i: i, isInt: true}, true
	}
	if u, err := strconv.ParseUint(s, 10, 64); err == nil {
		return num{f: float64(u), u: u, isUint: true}, true
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return num{}, false
	}
	return num{f: f}, true
}

// stringOf reads text: a string, or a named type over string.
func stringOf(v any) (string, bool) {
	if s, ok := v.(string); ok {
		return s, true
	}
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.String {
		return rv.String(), true
	}
	return "", false
}

func (kv Kv) castInt(v any, want reflect.Type, opt CastIntOpt) (any, error) {
	var n num
	if s, ok := stringOf(v); ok {
		if opt.Flags&CastIntErrOnString != 0 {
			return nil, fmt.Errorf("kvany: %q is text, cannot cast to %s", kv.K, want)
		}
		if n, ok = parseNum(s); !ok {
			return nil, fmt.Errorf("kvany: %q = %q is not a number", kv.K, s)
		}
	} else if n, ok = numOf(v); !ok {
		return nil, fmt.Errorf("kvany: %q is %T, cannot cast to %s", kv.K, v, want)
	}

	// an exact integer, as int64 or uint64; a float is truncated or rounded first
	if !n.isInt && !n.isUint {
		f := n.f
		if f != math.Trunc(f) {
			switch {
			case opt.Flags&CastIntFractionError != 0:
				return nil, fmt.Errorf("kvany: %q = %v has a fraction, cannot cast to %s", kv.K, v, want)
			case opt.Flags&CastIntRound != 0:
				f = math.Round(f)
			default:
				f = math.Trunc(f)
			}
		}
		switch {
		case math.IsNaN(f) || math.IsInf(f, 0):
			return nil, fmt.Errorf("kvany: %q = %v cannot cast to %s", kv.K, v, want)
		case f >= math.MinInt64 && f < math.MaxInt64:
			n.i, n.isInt = int64(f), true
		case f >= 0 && f < math.MaxUint64:
			n.u, n.isUint = uint64(f), true
		default:
			return nil, fmt.Errorf("kvany: %q = %v does not fit %s", kv.K, v, want)
		}
	}

	out := reflect.New(want).Elem()
	if opt.Flags&CastIntIsMask != 0 { // keep the low bits of the wanted width, whatever the sign
		bits := want.Bits()
		u := n.u
		if n.isInt {
			u = uint64(n.i)
		}
		if bits < 64 {
			u &= 1<<uint(bits) - 1
		}
		if isUnsignedKind(want.Kind()) {
			out.SetUint(u)
		} else {
			i := int64(u)
			if bits < 64 && u&(1<<uint(bits-1)) != 0 { // sign extend
				i -= 1 << uint(bits)
			}
			out.SetInt(i)
		}
		return out.Interface(), nil
	}
	if isUnsignedKind(want.Kind()) {
		u := n.u
		if n.isInt {
			if n.i < 0 {
				return nil, fmt.Errorf("kvany: %q = %v is negative, cannot cast to %s", kv.K, v, want)
			}
			u = uint64(n.i)
		}
		if out.OverflowUint(u) {
			return nil, fmt.Errorf("kvany: %q = %v does not fit %s", kv.K, v, want)
		}
		out.SetUint(u)
		return out.Interface(), nil
	}
	i := n.i
	if n.isUint {
		if n.u > math.MaxInt64 {
			return nil, fmt.Errorf("kvany: %q = %v does not fit %s", kv.K, v, want)
		}
		i = int64(n.u)
	}
	if out.OverflowInt(i) {
		return nil, fmt.Errorf("kvany: %q = %v does not fit %s", kv.K, v, want)
	}
	out.SetInt(i)
	return out.Interface(), nil
}

func isUnsignedKind(k reflect.Kind) bool {
	switch k {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}
	return false
}

func (kv Kv) castFloat(v any, want reflect.Type, opt CastFloatOpt) (any, error) {
	var n num
	if s, ok := stringOf(v); ok {
		if opt.Flags&CastFloatErrOnString != 0 {
			return nil, fmt.Errorf("kvany: %q is text, cannot cast to %s", kv.K, want)
		}
		if n, ok = parseNum(s); !ok {
			return nil, fmt.Errorf("kvany: %q = %q is not a number", kv.K, s)
		}
	} else if n, ok = numOf(v); !ok {
		return nil, fmt.Errorf("kvany: %q is %T, cannot cast to %s", kv.K, v, want)
	}
	f := n.f
	if opt.Decimals > 0 {
		p := math.Pow(10, float64(opt.Decimals))
		f = math.Round(f*p) / p
	}
	out := reflect.New(want).Elem()
	if out.OverflowFloat(f) {
		return nil, fmt.Errorf("kvany: %q = %v does not fit %s", kv.K, v, want)
	}
	out.SetFloat(f)
	return out.Interface(), nil
}

func (kv Kv) castString(v any, opt CastFlags) (any, error) {
	var s string
	switch x := v.(type) {
	case json.Number:
		s = x.String()
		if opt.String.Decimals > 0 {
			if f, err := x.Float64(); err == nil {
				s = strconv.FormatFloat(f, 'f', opt.String.Decimals, 64)
			}
		}
	case time.Time:
		var err error
		if s, err = kv.dateOut(x, opt.Date); err != nil {
			return nil, err
		}
	default:
		if t, ok := stringOf(v); ok {
			s = t
			break
		}
		n, ok := numOf(v)
		if !ok {
			return nil, fmt.Errorf("kvany: %q is %T, cannot cast to string", kv.K, v)
		}
		switch {
		case reflect.ValueOf(v).Kind() == reflect.Bool:
			s = strconv.FormatBool(n.i != 0)
		case n.isInt:
			s = strconv.FormatInt(n.i, 10)
		case n.isUint:
			s = strconv.FormatUint(n.u, 10)
		case opt.String.Decimals > 0:
			s = strconv.FormatFloat(n.f, 'f', opt.String.Decimals, 64)
		default:
			s = strconv.FormatFloat(n.f, 'f', -1, 64)
		}
	}
	return applyStringFlags(s, opt.String), nil
}

// applyStringFlags does what CastStringOpt asks, in this order: trim,
// case, pad. Widths and pads count runes, not bytes.
func applyStringFlags(s string, opt CastStringOpt) string {
	f := opt.Flags
	if f&CastStringLTrim != 0 {
		s = strings.TrimLeftFunc(s, unicode.IsSpace)
	}
	if f&CastStringRTrim != 0 {
		s = strings.TrimRightFunc(s, unicode.IsSpace)
	}
	if f&CastStringUpper != 0 {
		s = strings.ToUpper(s)
	}
	if f&CastStringLower != 0 {
		s = strings.ToLower(s)
	}
	if f&(CastStringPadL|CastStringPadR) != 0 && opt.Width > 0 {
		fill := []rune(opt.Fill)
		if len(fill) == 0 {
			fill = []rune{' '}
		}
		if missing := opt.Width - utf8.RuneCountInString(s); missing > 0 {
			if f&CastStringPadL != 0 {
				s = strings.Repeat(string(fill[0]), missing) + s
			} else {
				s += strings.Repeat(string(fill[len(fill)-1]), missing)
			}
		}
	}
	return s
}

func (kv Kv) castBool(v any) (any, error) {
	if s, ok := stringOf(v); ok {
		switch strings.ToLower(strings.TrimSpace(s)) {
		case "true", "1":
			return true, nil
		case "false", "0", "":
			return false, nil
		}
		return nil, fmt.Errorf("kvany: %q = %q is not a bool", kv.K, s)
	}
	if n, ok := numOf(v); ok {
		return n.f != 0, nil
	}
	return nil, fmt.Errorf("kvany: %q is %T, cannot cast to bool", kv.K, v)
}

func (kv Kv) castDate(v any, opt CastDateOpt) (any, error) {
	if t, ok := v.(time.Time); ok {
		return t, nil
	}
	s, ok := stringOf(v)
	if !ok {
		return nil, fmt.Errorf("kvany: %q is %T, cannot cast to time.Time", kv.K, v)
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, nil
	}
	loc := opt.Loc
	if loc == nil {
		loc = time.Local
	}
	layouts := opt.In
	if layouts == nil {
		layouts = CastDateLayouts
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, loc); err == nil {
			return t, nil
		}
	}
	return nil, fmt.Errorf("kvany: %q = %q is not a date", kv.K, s)
}

func (kv Kv) dateOut(t time.Time, opt CastDateOpt) (string, error) {
	if t.IsZero() {
		return "", nil
	}
	layout := opt.Out
	if layout == "" {
		layout = "20060102"
	}
	return t.Format(layout), nil
}
