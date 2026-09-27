package kvany

import "time"

// ---- Peek

type PeekValuesFlag int

const (
	PeekValuesNone            PeekValuesFlag = 0x00000000
	PeekValuesCaseInsensitive PeekValuesFlag = 0x00000001
	PeekValuesGivenCase       PeekValuesFlag = 0x00000001 // Ignored by VPeek
	PeekValuesFirstFoundCase  PeekValuesFlag = 0x00000011 // Ignored by VPeek
	PeekValuesLastFoundCase   PeekValuesFlag = 0x00000021 // Ignored by VPeek
	PeekValuesDupeFirst       PeekValuesFlag = 0x00001000
	PeekValuesDupeLast        PeekValuesFlag = 0x00002000
	PeekValuesDupeArray       PeekValuesFlag = 0x00004000
	PeekValuesDupeError       PeekValuesFlag = 0x00008000
	PeekValuesSkipNotFound    PeekValuesFlag = 0x00010000 // Ignored by VPeek

)

// ---- the option structures

// CastFlags is the options of Cast and of VPeek: one structure for both,
// each one ignoring the members that are not its business, grouped by
// the type of the value asked for. The zero value is "no flags, nothing
// else": the defaults described with each group of flags below.
type CastFlags struct {
	Peek   PeekValuesFlag // VPeek only: how the key is searched
	Int    CastIntOpt     // asked as int, int8..int64, uint, uint8..uint64
	Float  CastFloatOpt   // asked as float32 or float64
	String CastStringOpt  // asked as string
	Date   CastDateOpt    // asked as time.Time, or a time.Time asked as string
}

type CastIntOpt struct {
	Flags CastIntFlag
}

type CastFloatOpt struct {
	Flags    CastFloatFlag
	Decimals int // round to this many decimals; 0 = untouched
}

type CastStringOpt struct {
	Flags    CastStringFlag
	Width    int    // PadL / PadR: pad to this many characters; 0 = no padding
	Fill     string // the pad character; " " when empty
	Decimals int    // a number written as text carries this many decimals; 0 = the shortest exact form
}

// An empty text asked as a date is the zero time.Time, not an error; a
// zero time.Time asked as text is "".
type CastDateOpt struct {
	In  []string       // layouts tried in order; nil = CastDateLayouts
	Out string         // layout of a time.Time asked as text; "" = "20060102"
	Loc *time.Location // zone of a date parsed from text without one; nil = time.Local
}

// CastDateLayouts is what a text asked as a date is tried against, in
// this order, when In is not given: the DBF date, what JSON makes of a
// time.Time, MySQL DATETIME and MySQL DATE.
var CastDateLayouts = []string{"20060102", time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"}

// ---- the flags

type CastIntFlag int

const (
	CastIntNone CastIntFlag = 0x0000

	// A value with a fraction asked as an integer is truncated by default
	// (12.7 -> 12, -12.7 -> -12). Round rounds it instead (12.7 -> 13,
	// 12.5 -> 13, -12.5 -> -13: half away from zero). FractionError fails
	// instead: only an exact integer converts.
	CastIntRound         CastIntFlag = 0x0001
	CastIntFractionError CastIntFlag = 0x0002

	// Text asked as an integer is parsed by default ("4009" -> 4009, with
	// spaces around tolerated). ErrOnString fails instead: only a real
	// number converts.
	CastIntErrOnString CastIntFlag = 0x0004

	// A value that does not fit the wanted type fails by default (300 asked
	// as int8, -1 asked as uint). IsMask says the value is a bit mask: it
	// is AND-ed with the mask of the wanted width (0xFF for 8 bits, 0xFFFF
	// for 16...), signed or unsigned alike: asked as uint8 -1 is 255 and
	// 256 is 0; asked as int8 255 is -1. What falls outside the mask is
	// dropped: there is no overflow and no underflow.
	CastIntIsMask CastIntFlag = 0x0008
)

type CastFloatFlag int

const (
	CastFloatNone CastFloatFlag = 0x0000

	// Text asked as a float is parsed by default (" 12.5 " -> 12.5).
	// ErrOnString fails instead: only a real number converts.
	CastFloatErrOnString CastFloatFlag = 0x0001
)

type CastStringFlag int

const (
	CastStringNone CastStringFlag = 0x0000

	// What is done to the text on the way out, in this order: trim, then
	// case, then pad. AllTrim is both LTrim and RTrim. PadL and PadR pad to
	// Width characters with Fill (a space when Fill is empty); a text
	// longer than Width is left as it is.
	CastStringLTrim   CastStringFlag = 0x0001
	CastStringRTrim   CastStringFlag = 0x0002
	CastStringAllTrim CastStringFlag = 0x0003
	CastStringPadL    CastStringFlag = 0x0004
	CastStringPadR    CastStringFlag = 0x0008
	CastStringUpper   CastStringFlag = 0x0010
	CastStringLower   CastStringFlag = 0x0020
)
