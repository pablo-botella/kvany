---
mkskill:
  pos: 34
  in: readme
---

## Casting a value

```go
func (kv Kv) Cast(like any, opt CastFlags) (any, error)
```

`V` is whatever `encoding/json` made of the wire (`float64` for every
number, `string` for dates…) and Go casts nothing by itself. `Cast`
returns `V` converted to the type of `like`, which is only a sample of the
wanted type; the result carries that exact type, so `.(int64)` on it is
safe. A `nil` like returns `V` untouched; a `nil` V is the zero value of
the wanted type (an empty DBF field). A named type over a basic kind
(`type Code string`) is fine on either side.

| Asked as | Accepted | Options (`CastFlags`) |
|---|---|---|
| `int`…`int64`, `uint`…`uint64` | any number, `json.Number`, bool (0/1), text | `Int.Flags`: fractions are truncated by default, `CastIntRound` rounds (half away from zero), `CastIntFractionError` fails; `CastIntErrOnString` refuses text; out of range fails unless `CastIntIsMask`, which keeps the low bits of the wanted width (`-1` as `uint8` is 255, `256` is 0, `255` as `int8` is -1) |
| `float32`, `float64` | the same | `Float.Flags`: `CastFloatErrOnString`; `Float.Decimals` rounds to that many |
| `string` | text, number, bool, `time.Time` | `String.Flags`, applied in this order: `CastStringLTrim` / `RTrim` / `AllTrim`, `CastStringUpper` / `Lower`, `CastStringPadL` / `PadR` to `String.Width` with `String.Fill` (a space when empty; PadL uses its first character, PadR its last); `String.Decimals` for a number written as text |
| `time.Time` | text, `time.Time` | `Date.In`: layouts tried in order, default `CastDateLayouts` (`20060102`, RFC 3339, MySQL DATETIME and DATE); `Date.Loc` for a text without zone (default `time.Local`); `""` is the zero time |
| `bool` | bool, number (0 = false), text (`true`/`1`, `false`/`0`/``) | none |

A `time.Time` asked as `string` uses `Date.Out` (default `20060102`); the
zero time gives `""`. Widths and pads count runes, not bytes.

```go
v, err := kv.Cast(int64(0), kvany.CastFlags{Int: kvany.CastIntOpt{Flags: kvany.CastIntIsMask}})
s, err := kv.Cast("", kvany.CastFlags{String: kvany.CastStringOpt{
	Flags: kvany.CastStringPadL, Width: 6, Fill: "0"}})           // "12" -> "000012"
t, err := kv.Cast(time.Time{}, kvany.CastFlags{})                  // "20260927" -> 2026-09-27 local
```
