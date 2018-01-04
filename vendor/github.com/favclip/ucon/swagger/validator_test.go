package swagger

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/favclip/ucon"
)

type TargetRequestValidate struct {
	Text string `swagger:"enum=ok|ng"`
}

func TestRequestValidator_valid(t *testing.T) {
	b, _ := ucon.MakeMiddlewareTestBed(t, RequestValidator(), func(req *TargetRequestValidate) {
	}, nil)
	b.Arguments[0] = reflect.ValueOf(&TargetRequestValidate{Text: "ok"})

	err := b.Next()
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestRequestValidator_invalid(t *testing.T) {
	b, _ := ucon.MakeMiddlewareTestBed(t, RequestValidator(), func(req *TargetRequestValidate) {
	}, nil)
	b.Arguments[0] = reflect.ValueOf(&TargetRequestValidate{Text: "invalid"})

	err := b.Next()
	if err == nil {
		t.Fatalf("unexpected: %v", err)
	}
	if v := err.(ucon.HTTPErrorResponse).StatusCode(); v != http.StatusBadRequest {
		t.Errorf("unexpected: %v", v)
	}
}
