package golidator

import (
	"reflect"
	"testing"
	"unicode"
)

func TestReflect(t *testing.T) {
	value := struct {
		Test string `validate:"d=test" unknown:"test"`
		foo  string `validate:"d=test"`
	}{
		Test: "",
		foo:  "",
	}

	sv := reflect.ValueOf(&value)
	st := sv.Type()
	if sv.Kind() == reflect.Ptr && !sv.IsNil() {
		// recursive!
		sv = sv.Elem()
		st = sv.Type()
	} else if sv.Kind() != reflect.Struct {
		// error! unsupported
	}

	for i := 0; i < sv.NumField(); i++ {
		fv := sv.Field(i)
		ft := st.Field(i)
		for fv.Kind() == reflect.Ptr && !fv.IsNil() {
			fv = fv.Elem()
		}
		tag := ft.Tag.Get("validate")

		if !unicode.IsUpper([]rune(ft.Name)[0]) {
			t.Logf("private")
		}

		t.Logf(ft.Name)
		t.Logf(tag)

		if unicode.IsUpper([]rune(ft.Name)[0]) {
			fv.SetString("aaa")
		}
	}

	t.Logf("%#v", value)
}
