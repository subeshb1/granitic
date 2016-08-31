package validate

import (
	"github.com/graniticio/granitic/test"
	"github.com/graniticio/granitic/types"
	"testing"
)

func TestMissingRequiredStringField(t *testing.T) {

	sb := newStringValidatorBuilder("DEF")

	sv, err := sb.parseStringRule("S", []string{"REQ:MISSING", "LEN:5-10:SHORT"})

	test.ExpectNil(t, err)

	sub := new(StringTest)
	vc := new(validationContext)
	vc.Subject = sub

	r, err := sv.Validate(vc)
	c := r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectBool(t, r.Unset, true)
	test.ExpectInt(t, len(c), 1)
	test.ExpectString(t, c[0], "MISSING")

	nsSub := new(NillableStringTest)
	vc.Subject = nsSub

	r, err = sv.Validate(vc)
	c = r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectBool(t, r.Unset, true)
	test.ExpectInt(t, len(c), 1)
	test.ExpectString(t, c[0], "MISSING")

	nsSub.S = new(types.NilableString)

	r, err = sv.Validate(vc)
	c = r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectBool(t, r.Unset, true)
	test.ExpectInt(t, len(c), 1)
	test.ExpectString(t, c[0], "MISSING")

}

func TestUnsetButOptional(t *testing.T) {
	sb := newStringValidatorBuilder("DEF")
	sv, err := sb.parseStringRule("S", []string{"LEN:5-10:SHORT"})

	test.ExpectNil(t, err)

	sub := new(StringTest)
	sub.S = ""
	vc := new(validationContext)
	vc.Subject = sub

	r, err := sv.Validate(vc)
	c := r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectBool(t, r.Unset, true)
	test.ExpectInt(t, len(c), 0)

	sub.S = "A"

	r, err = sv.Validate(vc)
	c = r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectBool(t, r.Unset, false)
	test.ExpectInt(t, len(c), 1)
}

func TestHardTrim(t *testing.T) {
	sb := newStringValidatorBuilder("DEF")

	sv, err := sb.parseStringRule("S", []string{"REQ:MISSING", "HARDTRIM"})

	test.ExpectNil(t, err)

	sub := new(StringTest)
	sub.S = " A "

	vc := new(validationContext)
	vc.Subject = sub

	r, err := sv.Validate(vc)
	c := r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 0)
	test.ExpectString(t, sub.S, "A")

	subNs := new(NillableStringTest)
	subNs.S = types.NewNilableString("  B  ")
	vc.Subject = subNs

	r, err = sv.Validate(vc)
	c = r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 0)
	test.ExpectString(t, subNs.S.String(), "B")

}

func TestSoftTrim(t *testing.T) {
	sb := newStringValidatorBuilder("DEF")

	sv, err := sb.parseStringRule("S", []string{"REQ:MISSING", "TRIM", "LEN:2-"})

	test.ExpectNil(t, err)

	sub := new(StringTest)
	sub.S = "  A  "

	vc := new(validationContext)
	vc.Subject = sub

	r, err := sv.Validate(vc)
	c := r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 1)
	test.ExpectString(t, sub.S, "  A  ")

	subNs := new(NillableStringTest)
	subNs.S = types.NewNilableString("  B  ")
	vc.Subject = subNs

	r, err = sv.Validate(vc)
	c = r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 1)
	test.ExpectString(t, subNs.S.String(), "  B  ")

}

func TestInSet(t *testing.T) {
	sb := newStringValidatorBuilder("DEF")

	sv, err := sb.parseStringRule("S", []string{"REQ:MISSING", "IN:AA,BB:NOTIN"})

	test.ExpectNil(t, err)

	sub := new(StringTest)
	sub.S = "A"

	vc := new(validationContext)
	vc.Subject = sub

	r, err := sv.Validate(vc)
	c := r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 1)
	test.ExpectString(t, c[0], "NOTIN")

	sub.S = "AA"

	r, err = sv.Validate(vc)
	c = r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 0)
}

func TestBreak(t *testing.T) {
	sb := newStringValidatorBuilder("DEF")

	sv, err := sb.parseStringRule("S", []string{"REQ:MISSING", "LEN:2-2:LENGTH", "BREAK", "IN:AA,BB:NOTIN"})

	test.ExpectNil(t, err)

	sub := new(StringTest)
	sub.S = "A"

	vc := new(validationContext)
	vc.Subject = sub

	r, err := sv.Validate(vc)
	c := r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 1)
	test.ExpectString(t, c[0], "LENGTH")

}

func TestStopAll(t *testing.T) {
	sb := newStringValidatorBuilder("DEF")

	sv, _ := sb.parseStringRule("S", []string{"REQ:MISSING", "LEN:2-:LENGTH"})

	test.ExpectBool(t, sv.StopAllOnFail(), false)

	sv, _ = sb.parseStringRule("S", []string{"REQ:MISSING", "LEN:2-:LENGTH", "STOPALL"})

	test.ExpectBool(t, sv.StopAllOnFail(), true)
}

func TestRegex(t *testing.T) {
	sb := newStringValidatorBuilder("DEF")

	sv, err := sb.parseStringRule("S", []string{"REQ:MISSING", "REG:^::A$:REGFAIL"})

	test.ExpectNil(t, err)

	sub := new(StringTest)
	sub.S = "B"

	vc := new(validationContext)
	vc.Subject = sub

	r, err := sv.Validate(vc)
	c := r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 1)
	test.ExpectString(t, c[0], "REGFAIL")

	sub.S = ":A"

	r, err = sv.Validate(vc)
	c = r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 0)

}

func TestLength(t *testing.T) {
	sb := newStringValidatorBuilder("DEF")

	sv, err := sb.parseStringRule("S", []string{"REQ:MISSING", "LEN:2-:LENGTH"})

	test.ExpectNil(t, err)

	sub := new(StringTest)
	sub.S = "A"

	vc := new(validationContext)
	vc.Subject = sub

	r, err := sv.Validate(vc)
	c := r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 1)
	test.ExpectString(t, c[0], "LENGTH")

	sub.S = "AA"

	r, err = sv.Validate(vc)
	c = r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 0)

	sv, err = sb.parseStringRule("S", []string{"REQ:MISSING", "LEN:2-3:LENGTH"})

	sub.S = "AAA"

	r, err = sv.Validate(vc)
	c = r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 0)

	sub.S = "AAAA"

	r, err = sv.Validate(vc)
	c = r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 1)
	test.ExpectString(t, c[0], "LENGTH")

	sv, err = sb.parseStringRule("S", []string{"REQ:MISSING", "LEN:-3:LENGTH"})

	sub.S = ""

	r, err = sv.Validate(vc)
	c = r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 0)

	sub.S = "AAA"

	r, err = sv.Validate(vc)
	c = r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 0)

	sub.S = "AAAA"

	r, err = sv.Validate(vc)
	c = r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 1)
	test.ExpectString(t, c[0], "LENGTH")

}

func TestExternal(t *testing.T) {
	sb := newStringValidatorBuilder("DEF")

	_, err := sb.parseStringRule("S", []string{"EXT:extComp"})

	test.ExpectNotNil(t, err)

	sb.componentFinder = new(CompFinder)

	_, err = sb.parseStringRule("S", []string{"EXT:unknown"})

	test.ExpectNotNil(t, err)

	sv, err := sb.parseStringRule("S", []string{"EXT:extChecker:EXTFAIL"})

	test.ExpectNil(t, err)

	sub := new(StringTest)
	sub.S = "A"

	vc := new(validationContext)
	vc.Subject = sub

	r, err := sv.Validate(vc)
	c := r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 1)
	test.ExpectString(t, c[0], "EXTFAIL")

	sub.S = "valid"
	r, err = sv.Validate(vc)
	c = r.ErrorCodes

	test.ExpectNil(t, err)
	test.ExpectInt(t, len(c), 0)
}

type StringTest struct {
	S string
}

type NillableStringTest struct {
	S *types.NilableString
}
