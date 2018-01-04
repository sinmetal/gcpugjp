package golidator

import (
	"testing"
)

func TestBuiltinReqValid(t *testing.T) {
	v := &Validator{}

	v.SetTag("validate")
	v.SetValidationFunc("req", ReqValidator)

	{ // string
		err := v.Validate(struct {
			Test string `validate:"req"`
		}{
			Test: "test",
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // int
		err := v.Validate(struct {
			TestA int   `validate:"req"`
			TestB int8  `validate:"req"`
			TestC int16 `validate:"req"`
			TestD int32 `validate:"req"`
			TestE int64 `validate:"req"`
		}{
			TestA: 1,
			TestB: 1,
			TestC: 1,
			TestD: 1,
			TestE: 1,
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // float
		err := v.Validate(struct {
			TestA float32 `validate:"req"`
			TestB float64 `validate:"req"`
		}{
			TestA: 1.0,
			TestB: 1.0,
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // bool
		err := v.Validate(struct {
			Test bool `validate:"req"`
		}{
			Test: false,
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // array, slice
		err := v.Validate(struct {
			Test []interface{} `validate:"req"`
		}{
			Test: []interface{}{""},
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // struct
		err := v.Validate(struct {
			Test *struct{} `validate:"req"`
		}{
			Test: &struct{}{},
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
}

func TestBuiltinReqInvalid(t *testing.T) {
	v := &Validator{}

	v.SetTag("validate")
	v.SetValidationFunc("req", ReqValidator)

	{ // string
		err := v.Validate(struct {
			Test string `validate:"req"`
		}{
			Test: "",
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
	{ // int
		err := v.Validate(struct {
			Test int `validate:"req"`
		}{
			Test: 0,
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
	{ // float
		err := v.Validate(struct {
			Test float32 `validate:"req"`
		}{
			Test: 0.0,
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
	{ // array, slice
		err := v.Validate(struct {
			Test []interface{} `validate:"req"`
		}{
			Test: []interface{}{},
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
	{ // struct
		err := v.Validate(struct {
			Test *struct{} `validate:"req"`
		}{
			Test: nil,
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
}

func TestBuiltinDefaultRequired(t *testing.T) {
	v := &Validator{}

	v.SetTag("validate")
	v.SetValidationFunc("d", DefaultValidator)

	{ // string
		value := &struct {
			Test string `validate:"d=test"`
		}{}
		err := v.Validate(value)
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}

		if value.Test != "test" {
			t.Fatalf("unexpected result: %#v", value)
		}
	}
	{ // int
		value := &struct {
			TestA int   `validate:"d=1"`
			TestB int8  `validate:"d=1"`
			TestC int16 `validate:"d=1"`
			TestD int32 `validate:"d=1"`
			TestE int64 `validate:"d=1"`
		}{}
		err := v.Validate(value)
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}

		if value.TestA != 1 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestB != 1 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestC != 1 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestD != 1 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestE != 1 {
			t.Fatalf("unexpected result: %#v", value)
		}
	}
	{ // uint
		value := &struct {
			TestA uint   `validate:"d=1"`
			TestB uint8  `validate:"d=1"`
			TestC uint16 `validate:"d=1"`
			TestD uint32 `validate:"d=1"`
			TestE uint64 `validate:"d=1"`
		}{}
		err := v.Validate(value)
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}

		if value.TestA != 1 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestB != 1 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestC != 1 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestD != 1 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestE != 1 {
			t.Fatalf("unexpected result: %#v", value)
		}
	}
	{ // float
		value := &struct {
			TestA float32 `validate:"d=1.5"`
			TestB float64 `validate:"d=1.5"`
		}{}
		err := v.Validate(value)
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}

		if value.TestA != 1.5 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestB != 1.5 {
			t.Fatalf("unexpected result: %#v", value)
		}
	}
}

func TestBuiltinDefaultNotRequired(t *testing.T) {
	v := &Validator{}

	v.SetTag("validate")
	v.SetValidationFunc("d", DefaultValidator)

	{ // string
		value := &struct {
			Test string `validate:"d=test"`
		}{
			Test: "",
		}
		err := v.Validate(value)
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}

		if value.Test != "test" {
			t.Fatalf("unexpected result: %#v", value)
		}
	}
	{ // int
		value := &struct {
			TestA int   `validate:"d=1"`
			TestB int8  `validate:"d=1"`
			TestC int16 `validate:"d=1"`
			TestD int32 `validate:"d=1"`
			TestE int64 `validate:"d=1"`
		}{
			TestA: 2,
			TestB: 2,
			TestC: 2,
			TestD: 2,
			TestE: 2,
		}
		err := v.Validate(value)
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}

		if value.TestA != 2 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestB != 2 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestC != 2 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestD != 2 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestE != 2 {
			t.Fatalf("unexpected result: %#v", value)
		}
	}
	{ // uint
		value := &struct {
			TestA uint   `validate:"d=1"`
			TestB uint8  `validate:"d=1"`
			TestC uint16 `validate:"d=1"`
			TestD uint32 `validate:"d=1"`
			TestE uint64 `validate:"d=1"`
		}{
			TestA: 2,
			TestB: 2,
			TestC: 2,
			TestD: 2,
			TestE: 2,
		}
		err := v.Validate(value)
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}

		if value.TestA != 2 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestB != 2 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestC != 2 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestD != 2 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestE != 2 {
			t.Fatalf("unexpected result: %#v", value)
		}
	}
	{ // float
		value := &struct {
			TestA float32 `validate:"d=1.5"`
			TestB float64 `validate:"d=1.5"`
		}{
			TestA: 2.5,
			TestB: 2.5,
		}
		err := v.Validate(value)
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}

		if value.TestA != 2.5 {
			t.Fatalf("unexpected result: %#v", value)
		}
		if value.TestB != 2.5 {
			t.Fatalf("unexpected result: %#v", value)
		}
	}
}

func TestBuiltinMinValid(t *testing.T) {
	v := &Validator{}

	v.SetTag("validate")
	v.SetValidationFunc("min", MinValidator)

	{ // int
		err := v.Validate(struct {
			TestA int   `validate:"min=10"`
			TestB int8  `validate:"min=10"`
			TestC int16 `validate:"min=10"`
			TestD int32 `validate:"min=10"`
			TestE int64 `validate:"min=10"`
		}{
			TestA: 10,
			TestB: 10,
			TestC: 10,
			TestD: 10,
			TestE: 10,
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // uint
		err := v.Validate(struct {
			TestA uint   `validate:"min=10"`
			TestB uint8  `validate:"min=10"`
			TestC uint16 `validate:"min=10"`
			TestD uint32 `validate:"min=10"`
			TestE uint64 `validate:"min=10"`
		}{
			TestA: 10,
			TestB: 10,
			TestC: 10,
			TestD: 10,
			TestE: 10,
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // float
		err := v.Validate(struct {
			TestA float32 `validate:"min=10.25"`
			TestB float64 `validate:"min=10.25"`
		}{
			TestA: 10.25,
			TestB: 10.25,
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
}

func TestBuiltinMinInvalid(t *testing.T) {
	v := &Validator{}

	v.SetTag("validate")
	v.SetValidationFunc("min", MinValidator)

	{ // int
		err := v.Validate(struct {
			TestA int `validate:"min=10"`
		}{
			TestA: 9,
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
	{ // uint
		err := v.Validate(struct {
			TestA uint `validate:"min=10"`
		}{
			TestA: 9,
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
	{ // float
		err := v.Validate(struct {
			TestA float32 `validate:"min=10.25"`
		}{
			TestA: 10.125,
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
}

func TestBuiltinMaxValid(t *testing.T) {
	v := &Validator{}

	v.SetTag("validate")
	v.SetValidationFunc("max", MaxValidator)

	{ // int
		err := v.Validate(struct {
			TestA int   `validate:"max=10"`
			TestB int8  `validate:"max=10"`
			TestC int16 `validate:"max=10"`
			TestD int32 `validate:"max=10"`
			TestE int64 `validate:"max=10"`
		}{
			TestA: 10,
			TestB: 10,
			TestC: 10,
			TestD: 10,
			TestE: 10,
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // uint
		err := v.Validate(struct {
			TestA uint   `validate:"max=10"`
			TestB uint8  `validate:"max=10"`
			TestC uint16 `validate:"max=10"`
			TestD uint32 `validate:"max=10"`
			TestE uint64 `validate:"max=10"`
		}{
			TestA: 10,
			TestB: 10,
			TestC: 10,
			TestD: 10,
			TestE: 10,
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // float
		err := v.Validate(struct {
			TestA float32 `validate:"max=10.25"`
			TestB float64 `validate:"max=10.25"`
		}{
			TestA: 10.25,
			TestB: 10.25,
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
}

func TestBuiltinMaxInvalid(t *testing.T) {
	v := &Validator{}

	v.SetTag("validate")
	v.SetValidationFunc("max", MaxValidator)

	{ // int
		err := v.Validate(struct {
			TestA int `validate:"max=10"`
		}{
			TestA: 11,
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
	{ // uint
		err := v.Validate(struct {
			TestA uint `validate:"max=10"`
		}{
			TestA: 11,
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
	{ // float
		err := v.Validate(struct {
			TestA float32 `validate:"max=10.25"`
		}{
			TestA: 10.375,
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
}

func TestBuiltinMinLenValid(t *testing.T) {
	v := &Validator{}

	v.SetTag("validate")
	v.SetValidationFunc("minLen", MinLenValidator)

	{ // string
		err := v.Validate(struct {
			Test string `validate:"minLen=3"`
		}{
			Test: "3文字",
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // array
		err := v.Validate(struct {
			Test [3]string `validate:"minLen=3"`
		}{
			Test: [3]string{},
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // slice
		err := v.Validate(struct {
			Test []string `validate:"minLen=3"`
		}{
			Test: []string{"1", "2", "3"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // map
		err := v.Validate(struct {
			Test map[string]string `validate:"minLen=3"`
		}{
			Test: map[string]string{"a": "A", "b": "B", "c": "C"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
}

func TestBuiltinMinLenInvalid(t *testing.T) {
	v := &Validator{}

	v.SetTag("validate")
	v.SetValidationFunc("minLen", MinLenValidator)

	{ // string
		err := v.Validate(struct {
			Test string `validate:"minLen=3"`
		}{
			Test: "2字",
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
	{ // array
		err := v.Validate(struct {
			Test [2]string `validate:"minLen=3"`
		}{
			Test: [2]string{},
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
	{ // slice
		err := v.Validate(struct {
			Test []string `validate:"minLen=3"`
		}{
			Test: []string{"1", "2"},
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
	{ // map
		err := v.Validate(struct {
			Test map[string]string `validate:"minLen=3"`
		}{
			Test: map[string]string{"a": "A", "b": "B"},
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
}

func TestBuiltinMaxLenValid(t *testing.T) {
	v := &Validator{}

	v.SetTag("validate")
	v.SetValidationFunc("maxLen", MaxLenValidator)

	{ // string
		err := v.Validate(struct {
			Test string `validate:"maxLen=3"`
		}{
			Test: "3文字",
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // array
		err := v.Validate(struct {
			Test [3]string `validate:"maxLen=3"`
		}{
			Test: [3]string{},
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // slice
		err := v.Validate(struct {
			Test []string `validate:"maxLen=3"`
		}{
			Test: []string{"1", "2", "3"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // map
		err := v.Validate(struct {
			Test map[string]string `validate:"maxLen=3"`
		}{
			Test: map[string]string{"a": "A", "b": "B", "c": "C"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
}

func TestBuiltinMaxLenInvalid(t *testing.T) {
	v := &Validator{}

	v.SetTag("validate")
	v.SetValidationFunc("maxLen", MaxLenValidator)

	{ // string
		err := v.Validate(struct {
			Test string `validate:"maxLen=3"`
		}{
			Test: "4文字！",
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
	{ // array
		err := v.Validate(struct {
			Test [4]string `validate:"maxLen=3"`
		}{
			Test: [4]string{},
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
	{ // slice
		err := v.Validate(struct {
			Test []string `validate:"maxLen=3"`
		}{
			Test: []string{"1", "2", "3", "4"},
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
	{ // map
		err := v.Validate(struct {
			Test map[string]string `validate:"maxLen=3"`
		}{
			Test: map[string]string{"a": "A", "b": "B", "c": "C", "d": "D"},
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
}

func TestBuiltinEmail(t *testing.T) {
	v := &Validator{}

	v.SetTag("validate")
	v.SetValidationFunc("email", EmailValidator)

	{
		list := []struct {
			Test string `validate:"email"`
		}{
			{
				"foo@bar",
			},
			{
				"foo@bar.com",
			},
			{
				"foo.bar@buzz.com",
			},
			{
				"foo@bar.co.jp",
			},
			{
				"foo+tag@bar.com",
			},
			{
				`"foo bar"@buzz.com`,
			},
		}
		for _, valid := range list {
			err := v.Validate(valid)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
		}
	}

	{
		list := []struct {
			Test string `validate:"email"`
		}{
			{
				"foo bar",
			},
			{
				"foo.bar",
			},
			{
				"foo..bar@buzz",
			},
			{
				"foo@bar@buzz",
			},
			{
				`foo+tag@bar.com\r\n\r\nTest Mail header Injection`,
			},
			{
				`"Test
				Test"@bar.com`,
			},
		}
		for _, invalid := range list {
			err := v.Validate(invalid)
			if err == nil {
				t.Fatalf("error expected")
			}
			if _, ok := err.(*ErrorReport); !ok {
				t.Fatal(err)
			}
		}
	}
}

func TestBuiltinEnumValid(t *testing.T) {
	v := &Validator{}

	v.SetTag("validate")
	v.SetValidationFunc("enum", EnumValidator)

	{ // string
		err := v.Validate(struct {
			Test string `validate:"enum=ok|ng"`
		}{
			Test: "ok",
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // int
		err := v.Validate(struct {
			TestA int   `validate:"enum=1|2"`
			TestB int8  `validate:"enum=1|2"`
			TestC int16 `validate:"enum=1|2"`
			TestD int32 `validate:"enum=1|2"`
			TestE int64 `validate:"enum=1|2"`
		}{
			TestA: 1,
			TestB: 1,
			TestC: 1,
			TestD: 1,
			TestE: 1,
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	{ // uint
		err := v.Validate(struct {
			TestA uint   `validate:"enum=1|2"`
			TestB uint8  `validate:"enum=1|2"`
			TestC uint16 `validate:"enum=1|2"`
			TestD uint32 `validate:"enum=1|2"`
			TestE uint64 `validate:"enum=1|2"`
		}{
			TestA: 1,
			TestB: 1,
			TestC: 1,
			TestD: 1,
			TestE: 1,
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
}

func TestBuiltinEnumInvalid(t *testing.T) {
	v := &Validator{}

	v.SetTag("validate")
	v.SetValidationFunc("enum", EnumValidator)

	{ // string
		err := v.Validate(struct {
			Test string `validate:"enum=ok|ng"`
		}{
			Test: "unknown",
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
	{ // int
		err := v.Validate(struct {
			Test int `validate:"enum=1|2"`
		}{
			Test: 3,
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
	{ // uint
		err := v.Validate(struct {
			Test uint `validate:"enum=1|2"`
		}{
			Test: 3,
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		if _, ok := err.(*ErrorReport); !ok {
			t.Fatal(err)
		}
	}
}
