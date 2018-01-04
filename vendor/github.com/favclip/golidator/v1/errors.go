package golidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	// ErrUnsupportedValue is error of unsupported value type
	ErrUnsupportedValue = errors.New("unsupported value")
)

// UnsupportedTypeError returns an error from unsupported type.
func UnsupportedTypeError(in string, f reflect.StructField) error {
	return fmt.Errorf("%s: [%s] unsupported type %s", f.Name, in, f.Type.Name())
}

// EmptyParamError returns an error from empty param.
func EmptyParamError(in string, f reflect.StructField) error {
	return fmt.Errorf("%s: %s value is required", f.Name, in)
}

// ParamParseError returns an error from parsing param.
func ParamParseError(in string, f reflect.StructField, expected string) error {
	return fmt.Errorf("%s: %s value must be %s", f.Name, in, expected)
}

// ReqError returns an error from "req" validator.
func ReqError(f reflect.StructField, actual interface{}) error {
	return fmt.Errorf("%s: required, actual `%v`", f.Name, actual)
}

// DefaultError returns an error from "d" validator.
func DefaultError(f reflect.StructField) error {
	return fmt.Errorf("%s is not pointer. can't set default value", f.Name)
}

// MinError returns an error from "min" validator.
func MinError(f reflect.StructField, actual, min interface{}) error {
	return fmt.Errorf("%s: %v less than %v", f.Name, actual, min)
}

// MaxError returns an error from "max" validator.
func MaxError(f reflect.StructField, actual, max interface{}) error {
	return fmt.Errorf("%s: %v greater than %v", f.Name, actual, max)
}

// MinLenError returns an error from "minLen" validator.
func MinLenError(f reflect.StructField, actual, min interface{}) error {
	return fmt.Errorf("%s: less than %v: `%v`", f.Name, min, actual)
}

// MaxLenError returns an error from "maxLen" validator.
func MaxLenError(f reflect.StructField, actual, max interface{}) error {
	return fmt.Errorf("%s: greater than %v: `%v`", f.Name, max, actual)
}

// EmailError returns an error from "email" validator.
func EmailError(f reflect.StructField, actual string) error {
	return fmt.Errorf("%s: unsupported email format %s", f.Name, actual)
}

// EnumError returns an error from "enum" validator.
func EnumError(f reflect.StructField, actual interface{}, enum []string) error {
	return fmt.Errorf("%s: `%v` is not member of [%v]", f.Name, actual, strings.Join(enum, ", "))
}
