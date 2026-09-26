---
mkskill:
  pos: 30
  in: readme
---

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
