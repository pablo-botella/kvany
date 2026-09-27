---
mkskill:
  pos: 32
  in: readme
---

## One value: VPeek

```go
func (l Lst) VPeek(key string, opt CastFlags, defaultValue any) (any, error)
```

The single-value question: the value of `key`, or `defaultValue` when the
key is not there. With `PeekValuesCastToDefault` in `opt.Peek` the value
found is converted to the type of `defaultValue` through `Cast`, so the
default is also the sample of the wanted type and the assertion on the
result is safe; without the flag the value comes back as it came.

```go
raw, err := rec.VPeek("INUMB", kvany.CastFlags{}, nil)                 // 4009.0, untouched
n, err := rec.VPeek("INUMB", kvany.CastFlags{Peek: kvany.PeekValuesCastToDefault}, int64(0))   // n.(int64)
s, err := rec.VPeek("BCODE", kvany.CastFlags{Peek: kvany.PeekValuesCastToDefault,
	String: kvany.CastStringOpt{Flags: kvany.CastStringAllTrim | kvany.CastStringPadL, Width: 15}}, "")
```

`opt.Peek` drives the search: `PeekValuesCaseInsensitive` matches ignoring
case; a repeated key is an error unless `PeekValuesDupeFirst` or
`PeekValuesDupeLast` picks one, or `PeekValuesDupeArray` asks for all of
them as a `[]any` (a single match too) - which cannot be cast, so it is
incompatible with `CastToDefault`. The other Peek flags (output case,
SkipNotFound) have no meaning here and are ignored. The rest of `opt` is
the cast.
