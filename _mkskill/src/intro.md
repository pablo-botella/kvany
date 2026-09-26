---
mkskill:
  pos: 10
  in: readme
---

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
