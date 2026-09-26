---
mkskill:
  pos: 230
  in: ai*
---

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
