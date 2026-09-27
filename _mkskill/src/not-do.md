---
mkskill:
  pos: 40
  in: readme
---

## What it deliberately does not do

- **No key normalization.** Matching is exact unless you ask for
  case-insensitive; the list keeps whatever spelling came in.
- **No value conversion on its own.** Numbers stay `float64`, dates stay
  whatever string the wire carried; `PeekValues*` return them untouched.
  Coercion happens only when you ask for it, through `Cast` or `VPeek`.
- **No locking.** A `Lst` is a slice; guard it as you would guard one.
