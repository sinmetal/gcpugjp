package golidator

import (
	"encoding/json"
	"testing"
)

type Target struct {
	Name string `validate:"req,minLen=2"`
	Age  int    `json:"age" validate:"req,min=0"`
	Like string `json:"like" validate:"maxLen=10"`

	TargetSub `validate:"req"`
	Sub       *TargetSub `validate:"req"`
}

type TargetSub struct {
	Address string `validate:"req"`
}

func TestUsage(t *testing.T) {
	v := &Validator{}
	v.tag = "validate"
	v.funcs = make(map[string]ValidationFunc)
	v.funcs["req"] = ReqValidator
	v.funcs["minLen"] = MinLenValidator
	v.funcs["maxLen"] = MaxLenValidator
	v.funcs["min"] = MinValidator

	err := v.Validate(&Target{
		Name: "foobar",
		Sub:  &TargetSub{},
	})
	report, ok := err.(*ErrorReport)
	if err != nil && !ok {
		t.Fatal(err)
	}

	b, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}
