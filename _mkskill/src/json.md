---
mkskill:
  pos: 36
  in: readme
---

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
