---
mkskill:
  pos: 40
  in: readme
---

## What it deliberately does not do

- **No key normalization.** Matching is exact unless you ask for
  case-insensitive; the list keeps whatever spelling came in.
- **No value conversion.** Numbers stay `float64`, dates stay whatever
  string the wire carried. Coercion is the consumer's business.
- **No JSON methods.** The list lives inside the `map[string]any` you
  already decode.
- **No locking.** A `Lst` is a slice; guard it as you would guard one.
