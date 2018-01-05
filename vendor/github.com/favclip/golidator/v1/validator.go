package golidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// Validator is holder of validation information.
type Validator struct {
	tag   string
	funcs map[string]ValidationFunc
}

// Target is information of validation target.
type Target struct {
	StructType  reflect.Type
	StructValue reflect.Value
	FieldIndex  int
	FieldInfo   reflect.StructField
	FieldValue  reflect.Value
}

// ValidationFunc is validation function itself.
type ValidationFunc func(t *Target, param string) error

// NewValidator create and setup new Validator.
func NewValidator() *Validator {
	v := &Validator{}
	v.SetTag("validate")
	v.SetValidationFunc("req", ReqFactory(nil))
	v.SetValidationFunc("d", DefaultFactory(nil))
	v.SetValidationFunc("min", MinFactory(nil))
	v.SetValidationFunc("max", MaxFactory(nil))
	v.SetValidationFunc("minLen", MinLenFactory(nil))
	v.SetValidationFunc("maxLen", MaxLenFactory(nil))
	v.SetValidationFunc("email", EmailFactory(nil))
	v.SetValidationFunc("enum", EnumFactory(nil))
	return v
}

// SetTag is setup tag name in struct field tags.
func (vl *Validator) SetTag(tag string) {
	vl.tag = tag
}

// SetValidationFunc is setup tag name with ValidationFunc.
func (vl *Validator) SetValidationFunc(name string, vf ValidationFunc) {
	if vl.funcs == nil {
		vl.funcs = make(map[string]ValidationFunc, 0)
	}
	vl.funcs[name] = vf
}

// Validate do validate.
func (vl *Validator) Validate(v interface{}) error {
	return vl.validateStruct(reflect.ValueOf(v))
}

func (vl *Validator) validateStruct(sv reflect.Value) error {
	st := sv.Type()

	for sv.Kind() == reflect.Ptr && !sv.IsNil() {
		sv = sv.Elem()
		st = sv.Type()
	}

	if sv.Kind() != reflect.Struct {
		return ErrUnsupportedValue
	}

	for i := 0; i < sv.NumField(); i++ {
		fv := sv.Field(i)
		ft := st.Field(i)
		for fv.Kind() == reflect.Ptr && !fv.IsNil() {
			fv = fv.Elem()
		}
		if !unicode.IsUpper([]rune(ft.Name)[0]) {
			// private field!
			continue
		}

		tag := ft.Tag.Get(vl.tag)

		if tag == "-" {
			continue
		}
		if tag == "" {
			if fv.Kind() == reflect.Struct {
				err := vl.validateStruct(fv)
				if err != nil {
					return err
				}
			}
			continue
		}

		target := &Target{
			StructType:  st,
			StructValue: sv,
			FieldIndex:  i,
			FieldInfo:   ft,
			FieldValue:  fv,
		}
		params, err := parseTag(tag)
		if err != nil {
			return err
		}

		err = vl.validateField(target, params)
		if err != nil {
			// TODO custom error handler
			return err
		}

		if fv.Kind() == reflect.Struct {
			err := vl.validateStruct(fv)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (vl *Validator) validateField(t *Target, params map[string]string) error {
	for k, v := range params {
		f, ok := vl.funcs[k]
		if !ok {
			return fmt.Errorf("%s: unknown rule %s in %s", t.StructType.Name(), k, t.FieldInfo.Name)
		}
		err := f(t, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseTag(tagBody string) (map[string]string, error) {
	result := make(map[string]string, 0)
	if tagBody == "" {
		return result, nil
	}
	ss := strings.Split(tagBody, ",")
	for _, s := range ss {
		if s == "" {
			continue
		}
		p := strings.SplitN(s, "=", 2)
		name := p[0]
		if name == "" {
			return nil, errors.New("validator - empty name")
		}
		if len(p) == 1 {
			result[name] = ""
			continue
		}

		result[name] = p[1]
	}

	return result, nil
}
