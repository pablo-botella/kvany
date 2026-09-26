# kvany — agent notes

Go package `github.com/pablo-botella/kvany`: `Kv{K string; V any}` pairs,
`Lst []Kv` and `PtrLst []*Kv`, kept in the order they came — the Go side of
the `[["KEY", value], ...]` JSON arrays exchanged with Xbase++ (ot4xb). Keys
are not normalized and may repeat; values are whatever `encoding/json`
produced. Peeking values out is flag-driven. No dependencies.

## API

```go
type Kv struct{ K string; V any }
type Lst []Kv           // plain slice: range, index, slice it
type PtrLst []*Kv

lst, err := l.PeekValuesToLst(keys []string, flags PeekValuesFlag)   // new Lst, source untouched
m, err   := l.PeekValuesToMap(keys []string, flags PeekValuesFlag)   // the same folded into a map

// PeekValuesFlag (bitmask, OR them):
PeekValuesNone              // exact match, every repeat out, absent → Kv{key, nil}
PeekValuesCaseInsensitive   // 0x01: EqualFold matching
PeekValuesGivenCase         // 0x01: output key = the given one (default)
PeekValuesFirstFoundCase    // 0x11: output key = first found spelling
PeekValuesLastFoundCase     // 0x21: output key = last found spelling
PeekValuesDupeFirst         // 0x1000: repeated key → first value
PeekValuesDupeLast          // 0x2000: repeated key → last value
PeekValuesDupeArray         // 0x4000: repeated key → one pair, V = []any (single match too)
PeekValuesDupeError         // 0x8000: repeated key → error
PeekValuesSkipNotFound      // 0x10000: absent key adds nothing
```

## Semantics that matter

- Result order is `key_list` order; within one key, source order.
- Without a Dupe flag EVERY repeat comes out of PeekValuesToLst; in
  PeekValuesToMap they collapse and the LAST wins. DupeFirst/Last/Array
  decide it explicitly; DupeError fails naming the key and the count.
- The case flags only matter with CaseInsensitive; every pair of one key
  carries the same output key (given, first found or last found).
- Absent key: `Kv{key, nil}` (mpsetget style) unless SkipNotFound.
- Empty or nil key_list → empty, non-nil result. Source never modified.
- Values untouched: float64 numbers, wire strings for dates, nil for null.

## Tests

`test/kvany/` — PeekValuesToLst: exact/insensitive matching, the three
output-case flags, the four Dupe policies, absent keys with and without
SkipNotFound, empty key_list, source untouched; PeekValuesToMap: the same
through the map, last-wins collapse, nil map on error. `go test ./...`.
