# kvany — agent notes

Go package `github.com/pablo-botella/kvany`: `Kv{K string; V any}` pairs,
`Lst []Kv` and `PtrLst []*Kv`, kept in the order they came — the Go side of
the `[["KEY", value], ...]` JSON arrays exchanged with Xbase++ (ot4xb). Keys
are not normalized and may repeat; values are whatever `encoding/json`
produced. Peeking values out is flag-driven (`PeekValues*`); `VPeek` gets
one value cast to the type of its default; `Kv.Cast` is the conversion on
its own; the types marshal to and from the wire form. No dependencies.
Files: `kvany.go` (types, PeekValues*), `flags.go` (every flag and option
structure), `cast.go`, `vpeek.go`, `marshall.go`.

## API

```go
type Kv struct{ K string; V any }
type Lst []Kv           // plain slice: range, index, slice it
type PtrLst []*Kv

lst, err := l.PeekValuesToLst(keys []string, flags PeekValuesFlag)   // new Lst, source untouched
m, err   := l.PeekValuesToMap(keys []string, flags PeekValuesFlag)   // the same folded into a map
v, err   := l.VPeek(key string, opt CastFlags, defaultValue any)     // one value, cast to the default's type
v, err   := kv.Cast(like any, opt CastFlags)                        // V as the type of like

// PeekValuesFlag (bitmask, OR them):
PeekValuesNone              // exact match, every repeat out, absent → Kv{key, nil}
PeekValuesCaseInsensitive   // 0x01: EqualFold matching
PeekValuesGivenCase         // 0x01: output key = the given one (default)   [ignored by VPeek]
PeekValuesFirstFoundCase    // 0x11: output key = first found spelling      [ignored by VPeek]
PeekValuesLastFoundCase     // 0x21: output key = last found spelling       [ignored by VPeek]
PeekValuesDupeFirst         // 0x1000: repeated key → first value
PeekValuesDupeLast          // 0x2000: repeated key → last value
PeekValuesDupeArray         // 0x4000: repeated key → one pair, V = []any (VPeek: a slice of the default's type)
PeekValuesDupeError         // 0x8000: repeated key → error (VPeek: the default when no Dupe flag)
PeekValuesSkipNotFound      // 0x10000: absent key adds nothing               [ignored by VPeek]

// CastFlags: one structure for Cast and VPeek, grouped by the type asked for
type CastFlags struct {
	Peek   PeekValuesFlag   // VPeek only
	Int    CastIntOpt       // Flags: CastIntRound | CastIntFractionError | CastIntErrOnString | CastIntIsMask
	Float  CastFloatOpt     // Flags: CastFloatErrOnString; Decimals int
	String CastStringOpt    // Flags: CastStringLTrim | RTrim | AllTrim | PadL | PadR | Upper | Lower; Width, Fill, Decimals
	Date   CastDateOpt      // In []string (layouts; default CastDateLayouts), Out string, Loc *time.Location
}

// JSON (marshall.go): Kv <-> ["KEY", value]; Lst / PtrLst <-> [["KEY", value], ...]; nil list -> []
```

## Semantics that matter

- Result order is `key_list` order; within one key, source order.
- Without a Dupe flag EVERY repeat comes out of PeekValuesToLst; in
  PeekValuesToMap they collapse and the LAST wins. DupeFirst/Last/Array
  decide it explicitly; DupeError fails naming the key and the count.
- The case flags only matter with CaseInsensitive; every pair of one key
  carries the same output key (given, first found or last found).
- Absent key: `Kv{key, nil}` (mpsetget style) unless SkipNotFound.
- Empty or nil key_list → empty, non-nil result. Source never modified.
- PeekValues* never convert: float64 numbers, wire strings for dates, nil
  for null. Conversion is Cast / VPeek only.
- VPeek: absent key → defaultValue as given; found → Cast to the type of
  defaultValue (nil default = untouched); repeated key → error unless
  DupeFirst / DupeLast / DupeArray. DupeArray needs a slice default (or
  nil → []any) and casts each value to its element type.
- Cast: nil like → V untouched; nil V → zero of the wanted type; the result
  has the exact type of like (named types included). Integers: truncate by
  default, Round or FractionError by flag; text parsed unless ErrOnString;
  out of range fails unless IsMask (low bits kept, no overflow: -1 as uint8
  = 255, 256 as uint8 = 0, 255 as int8 = -1). Strings: trim, then case,
  then pad (Width, Fill: first char for PadL, last for PadR, space when
  empty); runes, not bytes. Dates: text tried against In or CastDateLayouts
  in Loc (default time.Local); "" <-> zero time; Out for time -> text.
  Bool: bool, number (0 = false), "true"/"1"/"false"/"0"/"".

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
