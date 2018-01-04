package ucon

import (
	"reflect"
	"testing"
	"time"
)

type valueStringTime time.Time

func (t valueStringTime) ParseString(value string) (interface{}, error) {
	v, err := time.Parse("2006-01-02", value)
	if err != nil {
		return valueStringTime(time.Time{}), err
	}
	return valueStringTime(v), nil
}

type ValueStringMapperSample struct {
	AString string
	BString string `json:"bStr"`

	DInt8    int8
	EInt64   int64
	FUint8   uint8
	GUint64  uint64
	HFloat32 float32
	IFloat64 float64
	JBool    bool
	KTime    valueStringTime

	ValueStringMapperSampleInner
}

type ValueStringMapperSampleInner struct {
	CString string
}

type ValueStringSliceMapperSample struct {
	AStrings []string
	BStrings []string `json:"bStrs"`

	DInt8s    []int8
	EInt64s   []int64
	FUint8s   []uint8
	GUint64s  []uint64
	HFloat32s []float32
	IFloat64s []float64
	JBools    []bool
	KTimes    []valueStringTime

	ValueStringSliceMapperSampleInner

	YString string
	ZString string
}

type ValueStringSliceMapperSampleInner struct {
	CStrings []string
}

func TestValueStringMapper(t *testing.T) {
	obj := &ValueStringMapperSample{}
	target := reflect.ValueOf(obj)
	valueStringMapper(target, "AString", "This is A")
	valueStringMapper(target, "bStr", "This is B")
	valueStringMapper(target, "CString", "This is C")
	valueStringMapper(target, "DInt8", "1")
	valueStringMapper(target, "EInt64", "2")
	valueStringMapper(target, "FUint8", "3")
	valueStringMapper(target, "GUint64", "4")
	valueStringMapper(target, "HFloat32", "1.25")
	valueStringMapper(target, "IFloat64", "2.75")
	valueStringMapper(target, "JBool", "true")
	valueStringMapper(target, "KTime", "2016-01-07")

	if obj.AString != "This is A" {
		t.Errorf("unexpected A: %v", obj.AString)
	}
	if obj.BString != "This is B" {
		t.Errorf("unexpected B: %v", obj.BString)
	}
	if obj.CString != "This is C" {
		t.Errorf("unexpected C: %v", obj.CString)
	}
	if obj.DInt8 != 1 {
		t.Errorf("unexpected D: %v", obj.DInt8)
	}
	if obj.EInt64 != 2 {
		t.Errorf("unexpected E: %v", obj.EInt64)
	}
	if obj.FUint8 != 3 {
		t.Errorf("unexpected F: %v", obj.FUint8)
	}
	if obj.GUint64 != 4 {
		t.Errorf("unexpected G: %v", obj.GUint64)
	}
	if obj.HFloat32 != 1.25 {
		t.Errorf("unexpected H: %v", obj.HFloat32)
	}
	if obj.IFloat64 != 2.75 {
		t.Errorf("unexpected I: %v", obj.IFloat64)
	}
	if obj.JBool != true {
		t.Errorf("unexpected J: %v", obj.JBool)
	}
	if y, m, d := time.Time(obj.KTime).Date(); y != 2016 {
		t.Errorf("unexpected KTime.y: %v", y)
	} else if m != 1 {
		t.Errorf("unexpected KTime.m: %v", m)
	} else if d != 7 {
		t.Errorf("unexpected KTime.d: %v", d)
	}
}

func TestValueStringSliceMapper(t *testing.T) {
	obj := &ValueStringSliceMapperSample{}
	target := reflect.ValueOf(obj)
	valueStringSliceMapper(target, "AStrings", []string{"This is A1", "This is A2"})
	valueStringSliceMapper(target, "bStrs", []string{"This is B1", "This is B2"})
	valueStringSliceMapper(target, "CStrings", []string{"This is C1", "This is C2"})
	valueStringSliceMapper(target, "DInt8s", []string{"1", "11"})
	valueStringSliceMapper(target, "EInt64s", []string{"2", "22"})
	valueStringSliceMapper(target, "FUint8s", []string{"3", "33"})
	valueStringSliceMapper(target, "GUint64s", []string{"4", "44"})
	valueStringSliceMapper(target, "HFloat32s", []string{"1.25", "11.25"})
	valueStringSliceMapper(target, "IFloat64s", []string{"2.75", "22.75"})
	valueStringSliceMapper(target, "JBools", []string{"true", "false"})
	valueStringSliceMapper(target, "KTimes", []string{"2016-01-07", "2016-04-05"})
	valueStringSliceMapper(target, "YString", []string{"This is Y"})
	valueStringSliceMapper(target, "ZString", []string{})

	if len(obj.AStrings) != 2 {
		t.Errorf("unexpected A len: %v", len(obj.AStrings))
	}
	if obj.AStrings[0] != "This is A1" {
		t.Errorf("unexpected A[0]: %v", obj.AStrings[0])
	}
	if obj.AStrings[1] != "This is A2" {
		t.Errorf("unexpected A[1]: %v", obj.AStrings[1])
	}

	if len(obj.BStrings) != 2 {
		t.Errorf("unexpected B len: %v", len(obj.BStrings))
	}
	if obj.BStrings[0] != "This is B1" {
		t.Errorf("unexpected B[0]: %v", obj.BStrings[0])
	}
	if obj.BStrings[1] != "This is B2" {
		t.Errorf("unexpected B[1]: %v", obj.BStrings[1])
	}

	if len(obj.CStrings) != 2 {
		t.Errorf("unexpected C len: %v", len(obj.CStrings))
	}
	if obj.CStrings[0] != "This is C1" {
		t.Errorf("unexpected C[0]: %v", obj.CStrings[0])
	}
	if obj.CStrings[1] != "This is C2" {
		t.Errorf("unexpected C[1]: %v", obj.CStrings[1])
	}

	if len(obj.DInt8s) != 2 {
		t.Errorf("unexpected D len: %v", len(obj.DInt8s))
	}
	if obj.DInt8s[0] != 1 {
		t.Errorf("unexpected D[0]: %v", obj.DInt8s[0])
	}
	if obj.DInt8s[1] != 11 {
		t.Errorf("unexpected D[1]: %v", obj.DInt8s[1])
	}

	if len(obj.EInt64s) != 2 {
		t.Errorf("unexpected E len: %v", len(obj.EInt64s))
	}
	if obj.EInt64s[0] != 2 {
		t.Errorf("unexpected E[0]: %v", obj.EInt64s[0])
	}
	if obj.EInt64s[1] != 22 {
		t.Errorf("unexpected E[1]: %v", obj.EInt64s[1])
	}

	if len(obj.FUint8s) != 2 {
		t.Errorf("unexpected F len: %v", len(obj.FUint8s))
	}
	if obj.FUint8s[0] != 3 {
		t.Errorf("unexpected F[0]: %v", obj.FUint8s[0])
	}
	if obj.FUint8s[1] != 33 {
		t.Errorf("unexpected F[1]: %v", obj.FUint8s[1])
	}

	if len(obj.GUint64s) != 2 {
		t.Errorf("unexpected G len: %v", len(obj.GUint64s))
	}
	if obj.GUint64s[0] != 4 {
		t.Errorf("unexpected G[0]: %v", obj.GUint64s[0])
	}
	if obj.GUint64s[1] != 44 {
		t.Errorf("unexpected G[1]: %v", obj.GUint64s[1])
	}

	if len(obj.HFloat32s) != 2 {
		t.Errorf("unexpected H len: %v", len(obj.HFloat32s))
	}
	if obj.HFloat32s[0] != 1.25 {
		t.Errorf("unexpected H[0]: %v", obj.HFloat32s[0])
	}
	if obj.HFloat32s[1] != 11.25 {
		t.Errorf("unexpected H[1]: %v", obj.HFloat32s[1])
	}

	if len(obj.IFloat64s) != 2 {
		t.Errorf("unexpected I len: %v", len(obj.IFloat64s))
	}
	if obj.IFloat64s[0] != 2.75 {
		t.Errorf("unexpected I[0]: %v", obj.IFloat64s[0])
	}
	if obj.IFloat64s[1] != 22.75 {
		t.Errorf("unexpected I[1]: %v", obj.IFloat64s[1])
	}

	if len(obj.JBools) != 2 {
		t.Errorf("unexpected J len: %v", len(obj.JBools))
	}
	if obj.JBools[0] != true {
		t.Errorf("unexpected J[0]: %v", obj.JBools[0])
	}
	if obj.JBools[1] != false {
		t.Errorf("unexpected J[1]: %v", obj.JBools[1])
	}

	if len(obj.KTimes) != 2 {
		t.Errorf("unexpected K len: %v", len(obj.KTimes))
	}
	if y, m, d := time.Time(obj.KTimes[0]).Date(); y != 2016 {
		t.Errorf("unexpected KTime[0].y: %v", y)
	} else if m != 1 {
		t.Errorf("unexpected KTime[0].m: %v", m)
	} else if d != 7 {
		t.Errorf("unexpected KTime[0].d: %v", d)
	}
	if y, m, d := time.Time(obj.KTimes[1]).Date(); y != 2016 {
		t.Errorf("unexpected KTime[1].y: %v", y)
	} else if m != 4 {
		t.Errorf("unexpected KTime[1].m: %v", m)
	} else if d != 5 {
		t.Errorf("unexpected KTime[1].d: %v", d)
	}

	if obj.YString != "This is Y" {
		t.Errorf("unexpected Y: %v", obj.YString)
	}
	if obj.ZString != "" {
		t.Errorf("unexpected Z: %v", obj.ZString)
	}
}
