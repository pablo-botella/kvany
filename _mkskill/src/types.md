---
mkskill:
  pos: 20
  in: readme
---

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
