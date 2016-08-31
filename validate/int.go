package validate

import (
	"errors"
	"fmt"
	"github.com/graniticio/granitic/ioc"
	rt "github.com/graniticio/granitic/reflecttools"
	"github.com/graniticio/granitic/types"
	"reflect"
	"strconv"
	"strings"
)

type ExternalInt64Validator interface {
	ValidInt64(int64) bool
}

const IntRuleCode = "INT"

const (
	IntOpIsRequiredCode = commonOpRequired
	IntOpIsStopAllCode  = commonOpStopAll
	IntOpInCode         = commonOpIn
	IntOpBreakCode      = commonOpBreak
	IntOpExtCode        = commonOpExt
)

type intValidationOperation uint

const (
	IntOpUnsupported = iota
	IntOpRequired
	IntOpStopAll
	IntOpIn
	IntOpBreak
	IntOpExt
)

func NewIntValidator(field, defaultErrorCode string) *IntValidator {
	iv := new(IntValidator)
	iv.defaultErrorCode = defaultErrorCode
	iv.field = field
	iv.codesInUse = types.NewOrderedStringSet([]string{})
	iv.dependsFields = determinePathFields(field)
	iv.operations = make([]*intOperation, 0)

	iv.codesInUse.Add(iv.defaultErrorCode)

	return iv
}

type IntValidator struct {
	stopAll             bool
	codesInUse          types.StringSet
	dependsFields       types.StringSet
	defaultErrorCode    string
	field               string
	missingRequiredCode string
	required            bool
	operations          []*intOperation
}

type intOperation struct {
	OpType   intValidationOperation
	ErrCode  string
	InSet    types.StringSet
	External ExternalInt64Validator
}

func (iv *IntValidator) Validate(vc *validationContext) (result *ValidationResult, unexpected error) {

	f := iv.field

	if vc.OverrideField != "" {
		f = vc.OverrideField
	}

	sub := vc.Subject

	fv, err := rt.FindNestedField(rt.ExtractDotPath(f), sub)

	if err != nil {
		return nil, err
	}

	r := new(ValidationResult)

	value, err := iv.extractValue(fv, f)

	if err != nil {
		return nil, err
	}

	if value == nil || !value.IsSet() {

		r.Unset = true

		if iv.required {
			r.ErrorCodes = []string{iv.missingRequiredCode}
		} else {
			r.ErrorCodes = []string{}
		}

		return r, nil
	}

	return iv.runOperations(value.Int64())
}

func (iv *IntValidator) runOperations(i int64) (*ValidationResult, error) {

	ec := new(types.OrderedStringSet)

OpLoop:
	for _, op := range iv.operations {

		switch op.OpType {
		case IntOpIn:
			if !iv.checkIn(i, op) {
				ec.Add(op.ErrCode)
			}
		case IntOpBreak:

			if ec.Size() > 0 {
				break OpLoop
			}
		case IntOpExt:

			if !op.External.ValidInt64(i) {
				ec.Add(op.ErrCode)
			}
		}

	}

	r := new(ValidationResult)
	r.ErrorCodes = ec.Contents()

	return r, nil

}

func (iv *IntValidator) checkIn(i int64, o *intOperation) bool {
	s := strconv.FormatInt(i, 10)

	return o.InSet.Contains(s)
}

func (iv *IntValidator) Break() *IntValidator {

	o := new(intOperation)
	o.OpType = IntOpBreak

	iv.addOperation(o)

	return iv

}

func (iv *IntValidator) addOperation(o *intOperation) {
	iv.operations = append(iv.operations, o)
	iv.codesInUse.Add(o.ErrCode)
}

func (iv *IntValidator) extractValue(v reflect.Value, f string) (*types.NilableInt64, error) {

	if rt.NilPointer(v) {
		return nil, nil
	}

	var ex int64

	switch i := v.Interface().(type) {
	case *types.NilableInt64:
		return i, nil
	case int:
		ex = int64(i)
	case int8:
		ex = int64(i)
	case int16:
		ex = int64(i)
	case int32:
		ex = int64(i)
	case int64:
		ex = i
	default:
		m := fmt.Sprintf("%s is type %T, not an int, int8, int16, int32, int64 or *NilableInt.", f, i)
		return nil, errors.New(m)

	}

	return types.NewNilableInt64(ex), nil

}

func (iv *IntValidator) StopAllOnFail() bool {
	return iv.stopAll
}

func (iv *IntValidator) CodesInUse() types.StringSet {
	return iv.codesInUse
}

func (iv *IntValidator) DependsOnFields() types.StringSet {

	return iv.dependsFields
}

func (iv *IntValidator) StopAll() *IntValidator {

	iv.stopAll = true

	return iv
}

func (iv *IntValidator) Required(code ...string) *IntValidator {

	iv.required = true
	iv.missingRequiredCode = iv.chooseErrorCode(code)

	return iv
}

func (iv *IntValidator) In(set []string, code ...string) *IntValidator {

	ss := types.NewUnorderedStringSet(set)

	ec := iv.chooseErrorCode(code)

	o := new(intOperation)
	o.OpType = IntOpIn
	o.ErrCode = ec
	o.InSet = ss

	iv.addOperation(o)

	return iv

}

func (iv *IntValidator) ExternalValidation(v ExternalInt64Validator, code ...string) *IntValidator {
	ec := iv.chooseErrorCode(code)

	o := new(intOperation)
	o.OpType = IntOpExt
	o.ErrCode = ec
	o.External = v

	iv.addOperation(o)

	return iv
}

func (iv *IntValidator) chooseErrorCode(v []string) string {

	if len(v) > 0 {
		iv.codesInUse.Add(v[0])
		return v[0]
	} else {
		return iv.defaultErrorCode
	}

}

func (iv *IntValidator) Operation(c string) (boolValidationOperation, error) {
	switch c {
	case IntOpIsRequiredCode:
		return IntOpRequired, nil
	case IntOpIsStopAllCode:
		return IntOpStopAll, nil
	case IntOpInCode:
		return IntOpIn, nil
	case IntOpBreakCode:
		return IntOpBreak, nil
	case IntOpExtCode:
		return IntOpExt, nil
	}

	m := fmt.Sprintf("Unsupported int validation operation %s", c)
	return IntOpUnsupported, errors.New(m)

}

func NewIntValidatorBuilder(ec string, cf ioc.ComponentByNameFinder) *intValidatorBuilder {
	iv := new(intValidatorBuilder)
	iv.componentFinder = cf
	iv.defaultErrorCode = ec

	return iv
}

type intValidatorBuilder struct {
	defaultErrorCode string
	componentFinder  ioc.ComponentByNameFinder
}

func (vb *intValidatorBuilder) parseRule(field string, rule []string) (Validator, error) {

	defaultErrorcode := DetermineDefaultErrorCode(IntRuleCode, rule, vb.defaultErrorCode)
	bv := NewIntValidator(field, defaultErrorcode)

	for _, v := range rule {

		ops := DecomposeOperation(v)
		opCode := ops[0]

		if IsTypeIndicator(IntRuleCode, opCode) {
			continue
		}

		op, err := bv.Operation(opCode)

		if err != nil {
			return nil, err
		}

		switch op {
		case IntOpRequired:
			err = vb.markRequired(field, ops, bv)
		case IntOpIn:
			err = vb.addIntInOperation(field, ops, bv)
		case IntOpStopAll:
			bv.StopAll()
		case IntOpBreak:
			bv.Break()
		case IntOpExt:
			err = vb.addIntExternalOperation(field, ops, bv)
		}

		if err != nil {

			return nil, err
		}

	}

	return bv, nil

}

func (vb *intValidatorBuilder) markRequired(field string, ops []string, iv *IntValidator) error {

	pCount, err := paramCount(ops, "Required", field, 1, 2)

	if err != nil {
		return err
	}

	if pCount == 1 {
		iv.Required()
	} else {
		iv.Required(ops[1])
	}

	return nil
}

func (vb *intValidatorBuilder) addIntExternalOperation(field string, ops []string, iv *IntValidator) error {

	pCount, i, err := validateExternalOperation(vb.componentFinder, field, ops)

	if err != nil {
		return err
	}

	ev, found := i.Instance.(ExternalInt64Validator)

	if !found {
		m := fmt.Sprintf("Component %s to validate field %s does not implement ExternalInt64Validator", i.Name, field)
		return errors.New(m)
	}

	if pCount == 2 {
		iv.ExternalValidation(ev)
	} else {
		iv.ExternalValidation(ev, ops[2])
	}

	return nil

}

func (vb *intValidatorBuilder) addIntInOperation(field string, ops []string, sv *IntValidator) error {

	pCount, err := paramCount(ops, "In Set", field, 2, 3)

	if err != nil {
		return err
	}

	members := strings.SplitN(ops[1], setMemberSep, -1)

	for _, m := range members {

		_, err := strconv.ParseInt(m, 10, 64)

		if err != nil {
			m := fmt.Sprintf("%s defined as a valid value when validating field %s cannot be parsed as an int64")
			return errors.New(m)
		}

	}

	if pCount == 2 {
		sv.In(members)
	} else {
		sv.In(members, ops[2])
	}

	return nil

}
