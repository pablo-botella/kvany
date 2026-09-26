---
mkskill:
  pos: 210
  in: ai*
---

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
