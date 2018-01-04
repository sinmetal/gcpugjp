package ucon

import "testing"

func TestMiddleware(t *testing.T) {
	DefaultMux = NewServeMux()

	if v := len(DefaultMux.middlewares); v != 0 {
		t.Fatalf("unexpected: %v", v)
	}

	Middleware(func(b *Bubble) error {
		return nil
	})

	if v := len(DefaultMux.middlewares); v != 1 {
		t.Fatalf("unexpected: %v", v)
	}
}

type TargetOfHandlersScannerPlugin struct {
}

func (obj *TargetOfHandlersScannerPlugin) HandlersScannerProcess(m *ServeMux, rds []*RouteDefinition) error {
	m.HandleFunc("GET", "/api/test/{id}", func() {})

	return nil
}

func TestPluginWithPluginContainer(t *testing.T) {
	DefaultMux = NewServeMux()

	if v := len(DefaultMux.plugins); v != 0 {
		t.Fatalf("unexpected: %v", v)
	}

	Plugin(&pluginContainer{
		base: &TargetOfHandlersScannerPlugin{},
	})

	if v := len(DefaultMux.plugins); v != 1 {
		t.Fatalf("unexpected: %v", v)
	}
}

func TestPluginWithoutPluginContainer(t *testing.T) {
	DefaultMux = NewServeMux()

	if v := len(DefaultMux.plugins); v != 0 {
		t.Fatalf("unexpected: %v", v)
	}

	Plugin(&TargetOfHandlersScannerPlugin{})

	if v := len(DefaultMux.plugins); v != 1 {
		t.Fatalf("unexpected: %v", v)
	}
}

func TestPrepare(t *testing.T) {
	DefaultMux = NewServeMux()

	Plugin(&TargetOfHandlersScannerPlugin{})

	if v := len(DefaultMux.router.handlers); v != 0 {
		t.Fatalf("unexpected: %v", v)
	}

	DefaultMux.Prepare()

	if v := len(DefaultMux.router.handlers); v != 1 {
		t.Fatalf("unexpected: %v", v)
	}
}

func TestHandle(t *testing.T) {
	DefaultMux = NewServeMux()

	if v := len(DefaultMux.router.handlers); v != 0 {
		t.Fatalf("unexpected: %v", v)
	}

	Handle("GET", "/api/test", &handlerContainerImpl{
		handler: func() {},
		Context: background,
	})

	HandleFunc("GET", "/api/test/{id}", func() {})
	HandleFunc("PUT", "/api/test/{id}", func() {})

	if v := len(DefaultMux.router.handlers); v != 3 {
		t.Fatalf("unexpected: %v", v)
	}
}

func TestUconContextWithValue(t *testing.T) {
	var ctx Context = background
	if v := ctx.Value("a"); v != nil {
		t.Fatalf("unexpected: %v", v)
	}

	ctx = WithValue(ctx, "a", "b")

	if v := ctx.Value("a"); v.(string) != "b" {
		t.Fatalf("unexpected: %v", v)
	}

	ctx = WithValue(ctx, 1, 2)

	if v := ctx.Value("a"); v.(string) != "b" {
		t.Fatalf("unexpected: %v", v)
	}
	if v := ctx.Value(1); v.(int) != 2 {
		t.Fatalf("unexpected: %v", v)
	}
}
