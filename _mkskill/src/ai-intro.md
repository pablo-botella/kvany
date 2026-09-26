---
mkskill:
  pos: 200
  in: ai*
---

# kvany — agent notes

Go package `github.com/pablo-botella/kvany`: `Kv{K string; V any}` pairs,
`Lst []Kv` and `PtrLst []*Kv`, kept in the order they came — the Go side of
the `[["KEY", value], ...]` JSON arrays exchanged with Xbase++ (ot4xb). Keys
are not normalized and may repeat; values are whatever `encoding/json`
produced. Peeking values out is flag-driven. No dependencies.
