package golidator

import (
	"fmt"
	"reflect"
)

type Sample struct {
	ID      string `validate:"req"`
	Default int    `validate:"d=10"`
	Min     int    `validate:"min=10"`
	Max     int    `validate:"max=-10"`
	MinLen  string `validate:"minLen=3"`
	MaxLen  string `validate:"maxLen=3"`
	Email   string `validate:"email"`
	Enum    string `validate:"enum=A|B"`
}

func ExampleNewValidator() {
	v := NewValidator()

	{
		obj := &struct {
			Value string `validate:"req"`
		}{}
		err := v.Validate(obj)
		fmt.Println(err.Error())
	}
	{
		obj := &struct {
			Value int `validate:"d=10"`
		}{}
		v.Validate(obj)
		fmt.Println(obj.Value)
	}
	{
		obj := &struct {
			Value int `validate:"min=10"`
		}{
			Value: 0,
		}
		err := v.Validate(obj)
		fmt.Println(err.Error())
	}
	{
		obj := &struct {
			Value int `validate:"max=10"`
		}{
			Value: 100,
		}
		err := v.Validate(obj)
		fmt.Println(err.Error())
	}
	{
		obj := &struct {
			Value string `validate:"minLen=3"`
		}{
			Value: "ab",
		}
		err := v.Validate(obj)
		fmt.Println(err.Error())
	}
	{
		obj := &struct {
			Value string `validate:"maxLen=3"`
		}{
			Value: "abcd",
		}
		err := v.Validate(obj)
		fmt.Println(err.Error())
	}
	{
		obj := &struct {
			Value string `validate:"email"`
		}{
			Value: "foo@bar@buzz",
		}
		err := v.Validate(obj)
		fmt.Println(err.Error())
	}
	{
		obj := &struct {
			Value string `validate:"enum=A|B"`
		}{
			Value: "C",
		}
		err := v.Validate(obj)
		fmt.Println(err.Error())
	}

	// Output:
	// Value: required, actual ``
	// 10
	// Value: 0 less than 10
	// Value: 100 greater than 10
	// Value: less than 3: `ab`
	// Value: greater than 3: `abcd`
	// Value: unsupported email format foo@bar@buzz
	// Value: `C` is not member of [A, B]
}

func ExampleValidator_SetValidationFunc() {
	v := &Validator{}
	v.SetTag("validate")
	v.SetValidationFunc("req", func(t *Target, param string) error {
		val := t.FieldValue
		if str := val.String(); str == "" {
			return fmt.Errorf("unexpected value: '%s'", str)
		}

		return nil
	})

	obj := &struct {
		Value string `validate:"req"`
	}{}
	err := v.Validate(obj)
	fmt.Println(err.Error())

	// Output:
	// unexpected value: ''
}

func ExampleValidator_useCustomizedError() {
	v := &Validator{}
	v.SetTag("validate")
	v.SetValidationFunc("req", ReqFactory(&ReqErrorOption{
		ReqError: func(f reflect.StructField, actual interface{}) error {
			return fmt.Errorf("%s IS REQUIRED", f.Name)
		},
	}))

	obj := &struct {
		FooBar string `validate:"req"`
	}{}
	err := v.Validate(obj)
	fmt.Println(err.Error())

	// Output:
	// FooBar IS REQUIRED
}
