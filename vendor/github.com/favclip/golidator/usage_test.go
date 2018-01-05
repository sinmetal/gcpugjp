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
		fmt.Println("value:", obj.Value)
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
	{
		obj := &struct {
			Min int `validate:"min=1"`
			Max int `validate:"max=10"`
		}{
			Min: 0,
			Max: 11,
		}
		err := v.Validate(obj)
		fmt.Println(err.Error())
	}

	// Output:
	// invalid. #1 Value: req actual: ''
	// value: 10
	// invalid. #1 Value: min=10 actual: 0
	// invalid. #1 Value: max=10 actual: 100
	// invalid. #1 Value: minLen=3 actual: 'ab'
	// invalid. #1 Value: maxLen=3 actual: 'abcd'
	// invalid. #1 Value: email actual: 'foo@bar@buzz'
	// invalid. #1 Value: enum=A|B actual: 'C'
	// invalid. #1 Min: min=1 actual: 0, #2 Max: max=10 actual: 11
}

func ExampleValidator_SetValidationFunc() {
	v := &Validator{}
	v.SetTag("validate")
	v.SetValidationFunc("req", func(param string, val reflect.Value) (ValidationResult, error) {
		if str := val.String(); str == "" {
			return ValidationNG, nil
		}

		return ValidationOK, nil
	})

	obj := &struct {
		Value string `validate:"req"`
	}{}
	err := v.Validate(obj)
	fmt.Println(err.Error())

	// Output:
	// invalid. #1 Value: req actual: ''
}

type Example struct {
	String string `validate:"req"`
}

func ExampleValidator_Validate() {
	v := NewValidator()

	err := v.Validate(&Example{})
	fmt.Println(err.Error())

	// Output:
	// invalid. Example #1 String: req actual: ''
}
