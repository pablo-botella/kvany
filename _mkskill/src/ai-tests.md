---
mkskill:
  pos: 240
  in: ai*
---

## Tests

`test/kvany/`, package `kvany_test`, through the public API only. `go test ./...`.

- `kvany_test.go` - PeekValuesToLst: exact/insensitive matching, the three
  output-case flags, the four Dupe policies, absent keys with and without
  SkipNotFound, empty key_list, source untouched; PeekValuesToMap: the same
  through the map, last-wins collapse, nil map on error.
- `marshall_test.go` - pair and list to `[["KEY", value], ...]`, nil list as
  `[]`, round trip (numbers back as float64), `null` / `[]` in, the four
  malformed pairs rejected, a list nested in a document.
- `cast_test.go` - nil like and nil V; integers from every source, fraction
  (truncate, Round, FractionError), ErrOnString, range and IsMask; floats
  and Decimals; strings from number, bool, named type and time, trim, case,
  pad (fill char, runes); dates in (default layouts, In, Loc, "") and out
  (Out, zero time); bool; unsupported types.
- `vpeek_test.go` - cast to the default's type, nil default, absent key,
  nil value, repeated key (error, DupeFirst, DupeLast), DupeArray with nil,
  slice and non-slice defaults, ignored Peek flags.
