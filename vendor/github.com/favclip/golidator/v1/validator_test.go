package golidator

import (
	"errors"
	"reflect"
	"testing"
)

type Request struct {
	Restrict
	Title string `validate:"d=unknown"`
}

type Restrict struct {
	Limit   int `validate:"d=10"`
	Curstor string
}

func TestUsage(t *testing.T) {
	v := NewValidator()

	{
		err := v.Validate(struct {
			Test string `validate:"req,enum=ok|ng"`
		}{
			Test: "unknown",
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		t.Logf(err.Error())
	}
	{
		err := v.Validate(struct {
			Test string `validate:"req,enum=ok|ng"`
		}{
			Test: "",
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		t.Logf(err.Error())
	}
	{
		err := v.Validate(struct {
			Test string `validate:"req,maxLen=3"`
		}{
			Test: "1234",
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		t.Logf(err.Error())
	}
	{
		err := v.Validate(struct {
			Test string `validate:"req,minLen=2,maxLen=3"`
		}{
			Test: "1",
		})
		if err == nil {
			t.Fatalf("error expected")
		}
		t.Logf(err.Error())
	}
}

func TestEmbedStruct(t *testing.T) {
	v := NewValidator()
	value := &Request{
		Restrict: Restrict{},
	}
	err := v.Validate(value)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

func TestUseCustomError(t *testing.T) {
	type Test struct {
		A string `validate:"req" customError:"Test.A"`
		B string `validate:"req"`
	}
	v := &Validator{}
	v.SetTag("validate")
	v.SetValidationFunc("req", ReqFactory(&ReqErrorOption{
		ReqError: func(f reflect.StructField, actual interface{}) error {
			if f.Tag.Get("customError") == "Test.A" {
				return errors.New("Test.A")
			}
			return errors.New(f.Name)
		},
	}))

	err := v.Validate(&Test{
		A: "",
		B: "1111",
	})
	if err == nil {
		t.Fatalf("unexpected")
	}
	if err.Error() != "Test.A" {
		t.Fatalf("unexpected: %s", err.Error())
	}
}
