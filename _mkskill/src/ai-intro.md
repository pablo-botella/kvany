---
mkskill:
  pos: 200
  in: ai*
---

# kvany — agent notes

Go package `github.com/pablo-botella/kvany`: `Kv{K string; V any}` pairs,
`Lst []Kv` and `PtrLst []*Kv`, kept in the order they came, with
`[["KEY", value], ...]` as their JSON form. Keys are not normalized and may
repeat; values are whatever `encoding/json` produced. Peeking values out is flag-driven (`PeekValues*`); `VPeek` gets
one value (cast to the type of its default on request); `Kv.Cast` is the
conversion on its own; the types marshal to and from the wire form. No
dependencies.
Files: `kvany.go` (types, PeekValues*), `flags.go` (every flag and option
structure), `cast.go`, `vpeek.go`, `marshall.go`.
