---
mkskill:
  pos: 240
  in: ai*
---

## Tests

`test/kvany/` — PeekValuesToLst: exact/insensitive matching, the three
output-case flags, the four Dupe policies, absent keys with and without
SkipNotFound, empty key_list, source untouched; PeekValuesToMap: the same
through the map, last-wins collapse, nil map on error. `go test ./...`.
