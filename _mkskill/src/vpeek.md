---
mkskill:
  pos: 32
  in: readme
---

## One value: VPeek

```go
func (l Lst) VPeek(key string, opt CastFlags, defaultValue any) (any, error)
```

The single-value question, with the cast built in: the value of `key`
converted to the type of `defaultValue`, or `defaultValue` itself when the
key is not there. The default is the sample of the wanted type, so the
assertion on the result is safe; a `nil` default asks for no conversion
and returns the value as it came.

```go
n, err := rec.VPeek("INUMB", kvany.CastFlags{}, int64(0))   // n.(int64), from 4009.0 or "4009"
s, err := rec.VPeek("BCODE", kvany.CastFlags{String: kvany.CastStringOpt{
	Flags: kvany.CastStringAllTrim | kvany.CastStringPadL, Width: 15}}, "")
raw, err := rec.VPeek("ANY", kvany.CastFlags{}, nil)            // untouched
```

`opt.Peek` drives the search: `PeekValuesCaseInsensitive` matches ignoring
case; a repeated key is an error unless `PeekValuesDupeFirst` or
`PeekValuesDupeLast` picks one, or `PeekValuesDupeArray` asks for all of
them — then `defaultValue` must be a slice (or `nil`, for a `[]any`) and
every value is cast to its element type (`[]int64{}` from `1.0`, `"2"`
and `3` gives `[]int64{1, 2, 3}`). The other Peek flags (output case,
SkipNotFound) have no meaning here and are ignored. The rest of `opt` is
the cast.
