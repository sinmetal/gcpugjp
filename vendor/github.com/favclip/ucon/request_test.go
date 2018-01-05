package ucon

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func noopMiddleware(b *Bubble) error {
	return b.Next()
}

func TestBubbleInit(t *testing.T) {
	b, _ := MakeMiddlewareTestBed(t, noopMiddleware, func(ctx context.Context) error {
		return nil
	}, nil)

	if v := len(b.ArgumentTypes); v != 1 {
		t.Errorf("unexpected: %v", v)
	}
	if v := b.ArgumentTypes[0]; v != contextType {
		t.Errorf("unexpected: %v", v)
	}
	if v := len(b.Arguments); v != 1 {
		t.Errorf("unexpected: %v", v)
	}
	if v := len(b.Returns); v != 0 {
		t.Errorf("unexpected: %v", v)
	}
}

func TestBubbleDo(t *testing.T) {
	b, _ := MakeMiddlewareTestBed(t, noopMiddleware, func(ctx context.Context) error {
		return nil
	}, nil)

	b.Arguments = nil
	if v := b.Next(); v != ErrInvalidArgumentLength {
		t.Errorf("unexpected: %v", v)
	}

	b.Arguments = make([]reflect.Value, 1)
	b.Arguments[0] = reflect.Value{}
	if v := b.Next(); v != ErrInvalidArgumentValue {
		t.Errorf("unexpected: %v", v)
	}

	b.Arguments[0] = reflect.ValueOf(time.Time{})
	if v := b.Next(); v != ErrInvalidArgumentValue {
		t.Errorf("unexpected: %v", v)
	}

	if v := b.Handled; v {
		t.Errorf("unexpected: %v", v)
	}

	b.Arguments[0] = reflect.ValueOf(context.Background())
	if v := b.Next(); v != nil {
		t.Errorf("unexpected: %v", v)
	}

	if v := len(b.Returns); v != 1 {
		t.Errorf("unexpected: %v", v)
	}
	if v := b.Returns[0]; !v.IsNil() {
		t.Errorf("unexpected: %v", v)
	}
	if v := b.Handled; !v {
		t.Errorf("unexpected: %v", v)
	}
}
