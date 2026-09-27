---
mkskill:
  pos: 230
  in: ai*
---

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
