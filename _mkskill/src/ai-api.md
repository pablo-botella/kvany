---
mkskill:
  pos: 210
  in: ai*
---

## API

```go
type Kv struct{ K string; V any }
type Lst []Kv           // plain slice: range, index, slice it
type PtrLst []*Kv

lst, err := l.PeekValuesToLst(keys []string, flags PeekValuesFlag)   // new Lst, source untouched
m, err   := l.PeekValuesToMap(keys []string, flags PeekValuesFlag)   // the same folded into a map
v, err   := l.VPeek(key string, opt CastFlags, defaultValue any)     // one value; the default when absent
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
