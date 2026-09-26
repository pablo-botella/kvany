# kvany

Key/value pairs whose values may be anything, kept **in the order they
came** — the Go side of the `[["KEY", value], ...]` arrays that travel in
JSON between Go and Xbase++ (ot4xb). Nothing is normalized and nothing is
deduplicated: a key may appear more than once. Picking values out of the
list is a flag-driven question — which case, which repeat, what when absent.

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

## What it deliberately does not do

- **No key normalization.** Matching is exact unless you ask for
  case-insensitive; the list keeps whatever spelling came in.
- **No value conversion.** Numbers stay `float64`, dates stay whatever
  string the wire carried. Coercion is the consumer's business.
- **No JSON methods.** The list lives inside the `map[string]any` you
  already decode.
- **No locking.** A `Lst` is a slice; guard it as you would guard one.

## Install

```
go get github.com/pablo-botella/kvany
```

## License

MIT — see [LICENSE](LICENSE).
