package golidator

import (
	"reflect"
	"testing"
)

type ErrorReportSample struct {
	String         string         `validate:"ng"`
	Int            int            `validate:"ng"`
	Int8           int8           `validate:"ng"`
	Int16          int16          `validate:"ng"`
	Int32          int32          `validate:"ng"`
	Int64          int64          `validate:"ng"`
	Float32        float32        `validate:"ng"`
	Float64        float64        `validate:"ng"`
	Bool           bool           `validate:"ng"`
	SliceInterface []interface{}  `validate:"ng"`
	ArrayInterface [3]interface{} `validate:"ng"`
	PtrStruct      *struct{}      `validate:"ng"`
}

func TestErrorReport_Error(t *testing.T) {
	v := &Validator{}

	v.SetTag("validate")
	v.SetValidationFunc("ng", func(param string, value reflect.Value) (ValidationResult, error) {
		return ValidationNG, nil
	})

	err := v.Validate(&ErrorReportSample{})
	if err == nil {
		t.Fatal("unexpected")
	}
	t.Log(err.Error())
}
