package golidator

import (
	"bytes"
	"fmt"
	mailp "net/mail"
	"reflect"
	re "regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// atext in RFC5322
// http://www.hde.co.jp/rfc/rfc5322.php?page=12
var atextChars = []byte(
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789" +
		"!#$%&'*+-/=?^_`{|}~")

// isAtext is return atext contains `c`
func isAtext(c rune) bool {
	return bytes.IndexRune(atextChars, c) >= 0
}

// UnsupportedTypeErrorOption is a factory of error caused by unsupported type.
type UnsupportedTypeErrorOption struct {
	UnsupportedTypeError func(in string, f reflect.StructField) error
}

// EmptyParamErrorOption is a factory of error caused by empty parameter.
type EmptyParamErrorOption struct {
	EmptyParamError func(in string, f reflect.StructField) error
}

// ParamParseErrorOption is a factory of error caused by parsing parameter.
type ParamParseErrorOption struct {
	ParamParseError func(in string, f reflect.StructField, expected string) error
}

// ReqErrorOption is a factory of error caused by "req" validator.
type ReqErrorOption struct {
	*UnsupportedTypeErrorOption
	ReqError func(f reflect.StructField, actual interface{}) error
}

// ReqFactory returns function validate parameter. it must not be empty.
func ReqFactory(e *ReqErrorOption) ValidationFunc {
	if e == nil {
		e = &ReqErrorOption{}
	}
	if e.UnsupportedTypeErrorOption == nil {
		e.UnsupportedTypeErrorOption = &UnsupportedTypeErrorOption{UnsupportedTypeError}
	}
	if e.ReqError == nil {
		e.ReqError = ReqError
	}
	return func(t *Target, param string) error {
		v := t.FieldValue

		switch v.Kind() {
		case reflect.String:
			if v.String() == "" {
				return e.ReqError(t.FieldInfo, v.String())
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v.Int() == 0 {
				return e.ReqError(t.FieldInfo, fmt.Sprintf("%d", v.Int()))
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if v.Uint() == 0 {
				return e.ReqError(t.FieldInfo, fmt.Sprintf("%d", v.Uint()))
			}
		case reflect.Float32, reflect.Float64:
			if v.Float() == 0 {
				return e.ReqError(t.FieldInfo, fmt.Sprintf("%f", v.Float()))
			}
		case reflect.Bool:
		// :)

		case reflect.Array, reflect.Slice:
			if v.Len() == 0 {
				return e.ReqError(t.FieldInfo, v.Len())
			}
		case reflect.Struct:
			// t.FieldValue is de-referenced
			ofv := t.StructValue.Field(t.FieldIndex)
			if ofv.Kind() == reflect.Ptr && ofv.IsNil() {
				return e.ReqError(t.FieldInfo, "nil")
			}
		default:
			return e.UnsupportedTypeError("req", t.FieldInfo)
		}
		return nil
	}
}

// DefaultErrorOption is a factory of error caused by "d" validator.
type DefaultErrorOption struct {
	*UnsupportedTypeErrorOption
	*ParamParseErrorOption
	DefaultError func(f reflect.StructField) error
}

// DefaultFactory returns function set default value.
func DefaultFactory(e *DefaultErrorOption) ValidationFunc {
	if e == nil {
		e = &DefaultErrorOption{}
	}
	if e.UnsupportedTypeErrorOption == nil {
		e.UnsupportedTypeErrorOption = &UnsupportedTypeErrorOption{UnsupportedTypeError}
	}
	if e.ParamParseErrorOption == nil {
		e.ParamParseErrorOption = &ParamParseErrorOption{ParamParseError}
	}
	if e.DefaultError == nil {
		e.DefaultError = DefaultError
	}
	return func(t *Target, param string) error {
		v := t.FieldValue
		switch v.Kind() {
		case reflect.String:
			if !v.CanAddr() {
				return e.DefaultError(t.FieldInfo)
			}
			if v.String() == "" {
				v.SetString(param)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if !v.CanAddr() {
				return e.DefaultError(t.FieldInfo)
			}
			pInt, err := strconv.ParseInt(param, 0, 64)
			if err != nil {
				return e.ParamParseError("d", t.FieldInfo, "int")
			}
			if v.Int() == 0 {
				v.SetInt(pInt)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if !v.CanAddr() {
				return e.DefaultError(t.FieldInfo)
			}
			pUint, err := strconv.ParseUint(param, 0, 64)
			if err != nil {
				return e.ParamParseError("d", t.FieldInfo, "uint")
			}
			if v.Uint() == 0 {
				v.SetUint(pUint)
			}
		case reflect.Float32, reflect.Float64:
			if !v.CanAddr() {
				return e.DefaultError(t.FieldInfo)
			}
			pFloat, err := strconv.ParseFloat(param, 64)
			if err != nil {
				return e.ParamParseError("d", t.FieldInfo, "float")
			}
			if v.Float() == 0 {
				v.SetFloat(pFloat)
			}
		default:
			return e.UnsupportedTypeError("default", t.FieldInfo)
		}

		return nil
	}
}

// MinErrorOption is a factory of error caused by "min" validator.
type MinErrorOption struct {
	*UnsupportedTypeErrorOption
	*ParamParseErrorOption
	MinError func(f reflect.StructField, actual, min interface{}) error
}

// MinFactory returns function check minimum value.
func MinFactory(e *MinErrorOption) ValidationFunc {
	if e == nil {
		e = &MinErrorOption{}
	}
	if e.UnsupportedTypeErrorOption == nil {
		e.UnsupportedTypeErrorOption = &UnsupportedTypeErrorOption{UnsupportedTypeError}
	}
	if e.ParamParseErrorOption == nil {
		e.ParamParseErrorOption = &ParamParseErrorOption{ParamParseError}
	}
	if e.MinError == nil {
		e.MinError = MinError
	}
	return func(t *Target, param string) error {
		v := t.FieldValue
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			pInt, err := strconv.ParseInt(param, 0, 64)
			if err != nil {
				return e.ParamParseError("min", t.FieldInfo, "number")
			}
			if v.Int() < pInt {
				return e.MinError(t.FieldInfo, v.Int(), pInt)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			pUint, err := strconv.ParseUint(param, 0, 64)
			if err != nil {
				return e.ParamParseError("min", t.FieldInfo, "number")
			}
			if v.Uint() < pUint {
				return e.MinError(t.FieldInfo, v.Uint(), pUint)
			}
		case reflect.Float32, reflect.Float64:
			pFloat, err := strconv.ParseFloat(param, 64)
			if err != nil {
				return e.ParamParseError("min", t.FieldInfo, "number")
			}
			if v.Float() < pFloat {
				return e.MinError(t.FieldInfo, v.Float(), pFloat)
			}
		default:
			return e.UnsupportedTypeError("min", t.FieldInfo)
		}

		return nil
	}
}

// MaxErrorOption is a factory of error caused by "max" validator.
type MaxErrorOption struct {
	*UnsupportedTypeErrorOption
	*ParamParseErrorOption
	MaxError func(f reflect.StructField, actual, max interface{}) error
}

// MaxFactory returns function check maximum value.
func MaxFactory(e *MaxErrorOption) ValidationFunc {
	if e == nil {
		e = &MaxErrorOption{}
	}
	if e.UnsupportedTypeErrorOption == nil {
		e.UnsupportedTypeErrorOption = &UnsupportedTypeErrorOption{UnsupportedTypeError}
	}
	if e.ParamParseErrorOption == nil {
		e.ParamParseErrorOption = &ParamParseErrorOption{ParamParseError}
	}
	if e.MaxError == nil {
		e.MaxError = MaxError
	}
	return func(t *Target, param string) error {
		v := t.FieldValue
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			pInt, err := strconv.ParseInt(param, 0, 64)
			if err != nil {
				return e.ParamParseError("max", t.FieldInfo, "number")
			}
			if v.Int() > pInt {
				return e.MaxError(t.FieldInfo, v.Int(), pInt)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			pUint, err := strconv.ParseUint(param, 0, 64)
			if err != nil {
				return e.ParamParseError("max", t.FieldInfo, "number")
			}
			if v.Uint() > pUint {
				return e.MaxError(t.FieldInfo, v.Uint(), pUint)
			}
		case reflect.Float32, reflect.Float64:
			pFloat, err := strconv.ParseFloat(param, 64)
			if err != nil {
				return e.ParamParseError("max", t.FieldInfo, "number")
			}
			if v.Float() > pFloat {
				return e.MaxError(t.FieldInfo, v.Float(), pFloat)
			}
		default:
			return e.UnsupportedTypeError("max", t.FieldInfo)
		}

		return nil
	}
}

// MinLenErrorOption is a factory of error caused by "minLen" validator.
type MinLenErrorOption struct {
	*UnsupportedTypeErrorOption
	*ParamParseErrorOption
	MinLenError func(f reflect.StructField, actual, min interface{}) error
}

// MinLenFactory returns function check minimum length.
func MinLenFactory(e *MinLenErrorOption) ValidationFunc {
	if e == nil {
		e = &MinLenErrorOption{}
	}
	if e.UnsupportedTypeErrorOption == nil {
		e.UnsupportedTypeErrorOption = &UnsupportedTypeErrorOption{UnsupportedTypeError}
	}
	if e.ParamParseErrorOption == nil {
		e.ParamParseErrorOption = &ParamParseErrorOption{ParamParseError}
	}
	if e.MinLenError == nil {
		e.MinLenError = MinLenError
	}
	return func(t *Target, param string) error {
		v := t.FieldValue

		switch v.Kind() {
		case reflect.String:
			if v.String() == "" {
				return nil // emptyの場合は無視 これはreqの役目だ
			}
			p, err := strconv.ParseInt(param, 0, 64)
			if err != nil {
				return e.ParamParseError("minLen", t.FieldInfo, "number")
			}
			if int64(utf8.RuneCountInString(v.String())) < p {
				return e.MinLenError(t.FieldInfo, v.String(), p)
			}
		case reflect.Array, reflect.Map, reflect.Slice:
			if v.Len() == 0 {
				return nil
			}
			p, err := strconv.ParseInt(param, 0, 64)
			if err != nil {
				return e.ParamParseError("minLen", t.FieldInfo, "number")
			}
			if int64(v.Len()) < p {
				return e.MinLenError(t.FieldInfo, v.Len(), p)
			}
		default:
			return e.UnsupportedTypeError("minLen", t.FieldInfo)
		}
		return nil
	}
}

// MaxLenErrorOption is a factory of error caused by "maxLen" validator.
type MaxLenErrorOption struct {
	*UnsupportedTypeErrorOption
	*ParamParseErrorOption
	MaxLenError func(f reflect.StructField, actual, max interface{}) error
}

// MaxLenFactory returns function check maximum length.
func MaxLenFactory(e *MaxLenErrorOption) ValidationFunc {
	if e == nil {
		e = &MaxLenErrorOption{}
	}
	if e.UnsupportedTypeErrorOption == nil {
		e.UnsupportedTypeErrorOption = &UnsupportedTypeErrorOption{UnsupportedTypeError}
	}
	if e.ParamParseErrorOption == nil {
		e.ParamParseErrorOption = &ParamParseErrorOption{ParamParseError}
	}
	if e.MaxLenError == nil {
		e.MaxLenError = MaxLenError
	}
	return func(t *Target, param string) error {
		v := t.FieldValue

		switch v.Kind() {
		case reflect.String:
			if v.String() == "" {
				return nil // emptyの場合は無視 これはreqの役目だ
			}
			p, err := strconv.ParseInt(param, 0, 64)
			if err != nil {
				return e.ParamParseError("maxLen", t.FieldInfo, "number")
			}
			if int64(utf8.RuneCountInString(v.String())) > p {
				return e.MaxLenError(t.FieldInfo, v.String(), p)
			}
		case reflect.Array, reflect.Map, reflect.Slice:
			if v.Len() == 0 {
				return nil
			}
			p, err := strconv.ParseInt(param, 0, 64)
			if err != nil {
				return e.ParamParseError("maxLen", t.FieldInfo, "number")
			}
			if int64(v.Len()) > p {
				return e.MaxLenError(t.FieldInfo, v.Len(), p)
			}
		default:
			return e.UnsupportedTypeError("maxLen", t.FieldInfo)
		}
		return nil
	}
}

// EmailErrorOption is a factory of error caused by "email" validator.
type EmailErrorOption struct {
	*UnsupportedTypeErrorOption
	EmailError func(f reflect.StructField, actual string) error
}

// EmailFactory returns function check value is `Email` format.
func EmailFactory(e *EmailErrorOption) ValidationFunc {
	if e == nil {
		e = &EmailErrorOption{}
	}
	if e.UnsupportedTypeErrorOption == nil {
		e.UnsupportedTypeErrorOption = &UnsupportedTypeErrorOption{UnsupportedTypeError}
	}
	if e.EmailError == nil {
		e.EmailError = EmailError
	}
	return func(t *Target, param string) error {
		v := t.FieldValue

		switch v.Kind() {
		case reflect.String:
			if v.String() == "" {
				return nil
			}
			addr := v.String()
			// do validation by RFC5322
			// http://www.hde.co.jp/rfc/rfc5322.php?page=17

			// screening
			if _, err := mailp.ParseAddress(addr); err != nil {
				return e.EmailError(t.FieldInfo, addr)
			}

			addrSpec := strings.Split(addr, "@")
			if len(addrSpec) != 2 {
				return e.EmailError(t.FieldInfo, addr)
			}
			// check local part
			localPart := addrSpec[0]
			// divided by quoted-string style or dom-atom style
			if match, err := re.MatchString(`"[^\t\n\f\r\\]*"`, localPart); err == nil && match { // "\"以外の表示可能文字を認める
				// OK
			} else if match, err := re.MatchString(`^([^.\s]+\.)*([^.\s]+)$`, localPart); err != nil || !match { // (hoge.)*hoge
				return e.EmailError(t.FieldInfo, addr)
			} else {
				// atext check for local part
				for _, c := range localPart {
					if string(c) == "." {
						// "." is already checked by regexp
						continue
					}
					if !isAtext(c) {
						e.EmailError(t.FieldInfo, addr)
					}
				}
			}
			// check domain part
			domain := addrSpec[1]
			if match, err := re.MatchString(`^([^.\s]+\.)*[^.\s]+$`, domain); err != nil || !match { // (hoge.)*hoge
				return e.EmailError(t.FieldInfo, addr)
			}
			// atext check for domain part
			for _, c := range domain {
				if string(c) == "." {
					// "." is already checked by regexp
					continue
				}
				if !isAtext(c) {
					return e.EmailError(t.FieldInfo, addr)
				}
			}
			return nil
		default:
			return e.UnsupportedTypeError("email", t.FieldInfo)
		}
	}
}

// EnumErrorOption is a factory of error caused by "enum" validator.
type EnumErrorOption struct {
	*UnsupportedTypeErrorOption
	*EmptyParamErrorOption
	EnumError func(f reflect.StructField, actual interface{}, enum []string) error
}

// EnumFactory returns function check value is established value.
func EnumFactory(e *EnumErrorOption) ValidationFunc {
	if e == nil {
		e = &EnumErrorOption{}
	}
	if e.UnsupportedTypeErrorOption == nil {
		e.UnsupportedTypeErrorOption = &UnsupportedTypeErrorOption{UnsupportedTypeError}
	}
	if e.EmptyParamErrorOption == nil {
		e.EmptyParamErrorOption = &EmptyParamErrorOption{EmptyParamError}
	}
	if e.EnumError == nil {
		e.EnumError = EnumError
	}
	return func(t *Target, param string) error {
		params := strings.Split(param, "|")
		if param == "" {
			return e.EmptyParamError("enum", t.FieldInfo)
		}

		v := t.FieldValue

		var enum func(v reflect.Value) error
		enum = func(v reflect.Value) error {
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			switch v.Kind() {
			case reflect.String:
				val := v.String()
				if val == "" {
					// need empty checking? use req :)
					return nil
				}
				for _, value := range params {
					if val == value {
						return nil
					}
				}
				return e.EnumError(t.FieldInfo, val, params)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				val := v.Int()
				for _, value := range params {
					value2, err := strconv.ParseInt(value, 10, 0)
					if err != nil {
						return nil
					}
					if val == value2 {
						return nil
					}
				}
				return e.EnumError(t.FieldInfo, val, params)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				val := v.Uint()
				for _, value := range params {
					value2, err := strconv.ParseUint(value, 10, 0)
					if err != nil {
						return nil
					}
					if val == value2 {
						return nil
					}
				}
				return e.EnumError(t.FieldInfo, val, params)
			case reflect.Array, reflect.Slice:
				for i := 0; i < v.Len(); i++ {
					e := v.Index(i)
					err := enum(e)
					if err != nil {
						return err
					}
				}
			default:
				return e.UnsupportedTypeError("enum", t.FieldInfo)
			}

			return nil
		}

		return enum(v)
	}
}
