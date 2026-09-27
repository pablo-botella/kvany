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

Four typed forms save the assertion when the type is known:

```go
func (l Lst) VPeekString(key string, opt CastFlags, defaultValue string) (string, error)
func (l Lst) VPeekInt(key string, opt CastFlags, defaultValue int) (int, error)
func (l Lst) VPeekFloat64(key string, opt CastFlags, defaultValue float64) (float64, error)
func (l Lst) VPeekBool(key string, opt CastFlags, defaultValue bool) (bool, error)
```

They are `VPeek` plus the assertion. `defaultValue` comes back only when
the key is absent. When the key is there, with `PeekValuesCastToDefault`
the value is converted first; without it the value must already be of
that type, and when it is not the call returns the zero value (`""`, `0`,
`false`) and an error ("value is not an int"), never the default. A
number decoded from JSON is a `float64`, so `VPeekInt` on it needs the
flag.

| `VPeekInt("K", opt, 7)` | without the flag | with `CastToDefault` |
|---|---|---|
| key absent | 7 | 7 |
| value `int` 4009 | 4009 | 4009 |
| value `float64` 4009 | 0, error | 4009 |
| value `"abc"` | 0, error | 0, error (not a number) |
| value `nil` | 0, error | 0 (the zero value) |
