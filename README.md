# kvany

Key/value pairs whose values may be anything, kept **in the order they
came**, with `[["KEY", value], ...]` as their JSON form. Nothing is
normalized and nothing is deduplicated: a key may appear more than once. Picking values out of the
list is a flag-driven question — which case, which repeat, what when absent —
and turning a value into the type you want (a `float64` into an `int64`, a
`"20260927"` into a `time.Time`) is a separate, explicit step.

```go
rec := kvany.Lst{{K: "CODE", V: "AB12"}, {K: "QTY", V: 3.0}, {K: "code", V: "dup"}}

m, err := rec.PeekValuesToMap([]string{"code", "qty", "note"},
	kvany.PeekValuesCaseInsensitive|kvany.PeekValuesDupeFirst)
// m: {"code": "AB12", "qty": 3.0, "note": nil}
```

## The types

```go
type Kv struct {
	K string   // the key, as it came
	V any      // the value: string, float64, bool, nil, []any, map[string]any…
}

type Lst    []Kv    // the pairs, in the order they came; keys may repeat
type PtrLst []*Kv   // the same, by pointer: for pairs shared and mutated while the list lives
```

A `Lst` is a plain slice: range over it, index it, slice it.

## Peeking values

```go
func (l Lst) PeekValuesToLst(key_list []string, flags PeekValuesFlag) (Lst, error)
func (l Lst) PeekValuesToMap(key_list []string, flags PeekValuesFlag) (map[string]any, error)
```

Both take the keys you want and return a **new** value; the source list is
never touched. Keys are visited in `key_list` order — that is the order of
the result — and for each one the list is scanned in its own order. The map
version is the list version folded into a map: same flags, same errors;
pairs sharing an output key collapse into one entry and the last wins.

| Flag | Meaning |
|---|---|
| `PeekValuesNone` | exact key match; every repeat comes out; an absent key comes out as `Kv{key, nil}` |
| `PeekValuesCaseInsensitive` | match ignoring case (`strings.EqualFold`) |
| `PeekValuesGivenCase` | with case-insensitive matching, the output key is the one you gave (the default: same value as `CaseInsensitive`) |
| `PeekValuesFirstFoundCase` / `PeekValuesLastFoundCase` | …or the spelling of the first / last pair found (both imply `CaseInsensitive`) |
| `PeekValuesDupeFirst` / `PeekValuesDupeLast` | a repeated key keeps its first / last value only |
| `PeekValuesDupeArray` | a repeated key comes out once, its `V` a `[]any` with every value — a single match too, as a one-element array |
| `PeekValuesDupeError` | a repeated key is an error naming the key and how many times it appeared |
| `PeekValuesSkipNotFound` | an absent key adds nothing instead of `Kv{key, nil}` |

## One value: VPeek

```go
func (l Lst) VPeek(key string, opt CastFlags, defaultValue any) (any, error)
```

The single-value question: the value of `key`, or `defaultValue` when the
key is not there. With `PeekValuesCastToDefault` in `opt.Peek` the value
found is converted to the type of `defaultValue` through `Cast`, so the
default is also the sample of the wanted type and the assertion on the
result is safe; without the flag the value comes back as it came.

```go
raw, err := rec.VPeek("INUMB", kvany.CastFlags{}, nil)                 // 4009.0, untouched
n, err := rec.VPeek("INUMB", kvany.CastFlags{Peek: kvany.PeekValuesCastToDefault}, int64(0))   // n.(int64)
s, err := rec.VPeek("BCODE", kvany.CastFlags{Peek: kvany.PeekValuesCastToDefault,
	String: kvany.CastStringOpt{Flags: kvany.CastStringAllTrim | kvany.CastStringPadL, Width: 15}}, "")
```

`opt.Peek` drives the search: `PeekValuesCaseInsensitive` matches ignoring
case; a repeated key is an error unless `PeekValuesDupeFirst` or
`PeekValuesDupeLast` picks one, or `PeekValuesDupeArray` asks for all of
them as a `[]any` (a single match too) - which cannot be cast, so it is
incompatible with `CastToDefault`. The other Peek flags (output case,
SkipNotFound) have no meaning here and are ignored. The rest of `opt` is
the cast.

## Casting a value

```go
func (kv Kv) Cast(like any, opt CastFlags) (any, error)
```

`V` is whatever `encoding/json` made of the wire (`float64` for every
number, `string` for dates…) and Go casts nothing by itself. `Cast`
returns `V` converted to the type of `like`, which is only a sample of the
wanted type; the result carries that exact type, so `.(int64)` on it is
safe. A `nil` like returns `V` untouched; a `nil` V is the zero value of
the wanted type (an empty DBF field). A named type over a basic kind
(`type Code string`) is fine on either side.

| Asked as | Accepted | Options (`CastFlags`) |
|---|---|---|
| `int`…`int64`, `uint`…`uint64` | any number, `json.Number`, bool (0/1), text | `Int.Flags`: fractions are truncated by default, `CastIntRound` rounds (half away from zero), `CastIntFractionError` fails; `CastIntErrOnString` refuses text; out of range fails unless `CastIntIsMask`, which keeps the low bits of the wanted width (`-1` as `uint8` is 255, `256` is 0, `255` as `int8` is -1) |
| `float32`, `float64` | the same | `Float.Flags`: `CastFloatErrOnString`; `Float.Decimals` rounds to that many |
| `string` | text, number, bool, `time.Time` | `String.Flags`, applied in this order: `CastStringLTrim` / `RTrim` / `AllTrim`, `CastStringUpper` / `Lower`, `CastStringPadL` / `PadR` to `String.Width` with `String.Fill` (a space when empty; PadL uses its first character, PadR its last); `String.Decimals` for a number written as text |
| `time.Time` | text, `time.Time` | `Date.In`: layouts tried in order, default `CastDateLayouts` (`20060102`, RFC 3339, MySQL DATETIME and DATE); `Date.Loc` for a text without zone (default `time.Local`); `""` is the zero time |
| `bool` | bool, number (0 = false), text (`true`/`1`, `false`/`0`/``) | none |

A `time.Time` asked as `string` uses `Date.Out` (default `20060102`); the
zero time gives `""`. Widths and pads count runes, not bytes.

```go
v, err := kv.Cast(int64(0), kvany.CastFlags{Int: kvany.CastIntOpt{Flags: kvany.CastIntIsMask}})
s, err := kv.Cast("", kvany.CastFlags{String: kvany.CastStringOpt{
	Flags: kvany.CastStringPadL, Width: 6, Fill: "0"}})           // "12" -> "000012"
t, err := kv.Cast(time.Time{}, kvany.CastFlags{})                  // "20260927" -> 2026-09-27 local
```

## JSON

`Kv`, `Lst` and `PtrLst` marshal to the wire form and back:

```go
b, _ := json.Marshal(kvany.Lst{{"A", 1}, {"b", nil}})   // [["A",1],["b",null]]
var l kvany.Lst
err := json.Unmarshal([]byte(`[["K",1],["K",true]]`), &l)
```

A pair is a two-element array `["KEY", value]`; anything else (an object,
a longer array, a key that is not a string) fails with a `kvany:` error. A
`nil` list marshals as `[]`, never `null`, because the other side always
expects an array; `null` unmarshals to a `nil` list and `[]` to an empty
one. Values are whatever `encoding/json` makes of them.

## Install

```
go get github.com/pablo-botella/kvany
```

## License

MIT — see [LICENSE](LICENSE).
