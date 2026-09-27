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

## API

```go
type Kv struct{ K string; V any }
type Lst []Kv           // plain slice: range, index, slice it
type PtrLst []*Kv

lst, err := l.PeekValuesToLst(keys []string, flags PeekValuesFlag)   // new Lst, source untouched
m, err   := l.PeekValuesToMap(keys []string, flags PeekValuesFlag)   // the same folded into a map
v, err   := l.VPeek(key string, opt CastFlags, defaultValue any)     // one value; the default when absent
s, err   := l.VPeekString(key, opt, "")                              // VPeek + assertion; also VPeekInt,
                                                                     // VPeekFloat64, VPeekBool
v, err   := kv.Cast(like any, opt CastFlags)                        // V as the type of like

// PeekValuesFlag (bitmask, OR them):
PeekValuesNone              // exact match, every repeat out, absent → Kv{key, nil}
PeekValuesCaseInsensitive   // 0x01: EqualFold matching
PeekValuesGivenCase         // 0x01: output key = the given one (default)   [ignored by VPeek]
PeekValuesFirstFoundCase    // 0x11: output key = first found spelling      [ignored by VPeek]
PeekValuesLastFoundCase     // 0x21: output key = last found spelling       [ignored by VPeek]
PeekValuesDupeFirst         // 0x1000: repeated key → first value
PeekValuesDupeLast          // 0x2000: repeated key → last value
PeekValuesDupeArray         // 0x4000: repeated key → one pair, V = []any (VPeek: incompatible with CastToDefault)
PeekValuesDupeError         // 0x8000: repeated key → error (VPeek: the default when no Dupe flag)
PeekValuesSkipNotFound      // 0x10000: absent key adds nothing               [ignored by VPeek]
PeekValuesCastToDefault     // 0x100000: VPeek only: Cast the value to the type of defaultValue

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
- Absent key: `Kv{key, nil}` unless SkipNotFound.
- Empty or nil key_list → empty, non-nil result. Source never modified.
- PeekValues* never convert: float64 numbers, wire strings for dates, nil
  for null. Conversion is Cast / VPeek only.
- VPeek: absent key -> defaultValue as given; found -> the value as it
  came, or Cast to the type of defaultValue with PeekValuesCastToDefault
  (nil default = untouched); repeated key -> error unless DupeFirst /
  DupeLast / DupeArray. DupeArray returns a []any and is incompatible with
  CastToDefault. VPeekString / Int / Float64 / Bool add the assertion:
  the default comes back only for an absent key; a present value of the
  wrong type gives the zero value and an error, never the default (a JSON
  number is float64: VPeekInt on it needs the flag).
- Cast: nil like → V untouched; nil V → zero of the wanted type; the result
  has the exact type of like (named types included). Integers: truncate by
  default, Round or FractionError by flag; text parsed unless ErrOnString;
  out of range fails unless IsMask (low bits kept, no overflow: -1 as uint8
  = 255, 256 as uint8 = 0, 255 as int8 = -1). Strings: trim, then case,
  then pad (Width, Fill: first char for PadL, last for PadR, space when
  empty); runes, not bytes. Dates: text tried against In or CastDateLayouts
  in Loc (default time.Local); "" <-> zero time; Out for time -> text.
  Bool: bool, number (0 = false), "true"/"1"/"false"/"0"/"".
