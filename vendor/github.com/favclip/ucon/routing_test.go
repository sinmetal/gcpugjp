package ucon

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type RequestOfRoutingInfoAddHandlers struct {
	ID     int    `json:"id"`
	Offset int    `json:"offset"`
	Text   string `json:"text"`
}

type ResponseOfRoutingInfoAddHandlers struct {
	Text string `json:"text"`
}

func TestRouterServeHTTP1(t *testing.T) {
	DefaultMux = NewServeMux()
	Orthodox()

	HandleFunc("PUT", "/api/test/{id}", func(req *RequestOfRoutingInfoAddHandlers) (*ResponseOfRoutingInfoAddHandlers, error) {
		if v := req.ID; v != 1 {
			t.Errorf("unexpected: %v", v)
		}
		if v := req.Offset; v != 100 {
			t.Errorf("unexpected: %v", v)
		}
		if v := req.Text; v != "Hi!" {
			t.Errorf("unexpected: %v", v)
		}
		return &ResponseOfRoutingInfoAddHandlers{Text: req.Text + "!"}, nil
	})

	DefaultMux.Prepare()

	resp := MakeHandlerTestBed(t, "PUT", "/api/test/1?offset=100", strings.NewReader("{\"text\":\"Hi!\"}"))

	if v := resp.StatusCode; v != 200 {
		t.Errorf("unexpected: %v", v)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if v := string(body); v != "{\"text\":\"Hi!!\"}" {
		t.Errorf("unexpected: %v", v)
	}
}

func TestRouterPickupBestRouteDefinition(t *testing.T) {
	DefaultMux = NewServeMux()

	h := func() {}

	HandleFunc("OPTIONS", "/", h)
	HandleFunc("*", "/foobar/", h)
	HandleFunc("GET", "/api/todo/{id}", h)
	HandleFunc("GET", "/api/todo", h)
	HandleFunc("POST", "/api/todo", h)
	HandleFunc("POST", "/api/todo/nested/too/long/", h)
	HandleFunc("POST", "/api/todo/nested/", h)

	{
		// OPTIONS /
		req, err := http.NewRequest("OPTIONS", "/api/todo", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.Method != "OPTIONS" || rd.PathTemplate.PathTemplate != "/" {
			t.Errorf("unexpected: %#v", rd)
		}
	}
	{
		// * /foobar/
		req, err := http.NewRequest("HEAD", "/foobar/buzz", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.Method != "*" || rd.PathTemplate.PathTemplate != "/foobar/" {
			t.Errorf("unexpected: %#v", rd)
		}
	}
	{
		// GET /api/todo/{id}
		req, err := http.NewRequest("GET", "/api/todo/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.Method != "GET" || rd.PathTemplate.PathTemplate != "/api/todo/{id}" {
			t.Errorf("unexpected: %#v", rd)
		}
	}
	{
		// handler is not exists
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd != nil {
			t.Fatalf("unexpected")
		}
	}
	{
		// POST /api/todo matches this request
		req, err := http.NewRequest("POST", "/api/todo/", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.Method != "POST" || rd.PathTemplate.PathTemplate != "/api/todo" {
			t.Errorf("unexpected: %#v", rd)
		}
	}
	{
		// POST /api/todo/nested/too/long/
		req, err := http.NewRequest("POST", "/api/todo/nested/too/long/", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.Method != "POST" || rd.PathTemplate.PathTemplate != "/api/todo/nested/too/long/" {
			t.Errorf("unexpected: %#v", rd)
		}
	}
	{
		// POST /api/todo/nested/
		req, err := http.NewRequest("POST", "/api/todo/nested/too", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.Method != "POST" || rd.PathTemplate.PathTemplate != "/api/todo/nested/" {
			t.Errorf("unexpected: %#v", rd)
		}
	}
}

func TestRouterPickupBestRouteDefinitionFocusOnMatchLength(t *testing.T) {
	DefaultMux = NewServeMux()

	h := func() {}

	HandleFunc("GET", "/", h)
	HandleFunc("GET", "/a", h)
	HandleFunc("GET", "/a/", h)
	HandleFunc("GET", "/a/b", h)

	HandleFunc("GET", "/c", h)
	HandleFunc("GET", "/c/d", h)

	{
		// GET / -> /
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.PathTemplate.PathTemplate != "/" {
			t.Fatalf("unexpected: %v <- %#v", req.URL.Path, rd.PathTemplate)
		}
	}
	{
		// GET /a -> /a
		req, err := http.NewRequest("GET", "/a", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.PathTemplate.PathTemplate != "/a" {
			t.Fatalf("unexpected: %v <- %#v", req.URL.Path, rd.PathTemplate)
		}
	}
	{
		// GET /a/b -> /a/b
		req, err := http.NewRequest("GET", "/a/b", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.PathTemplate.PathTemplate != "/a/b" {
			t.Fatalf("unexpected: %v <- %#v", req.URL.Path, rd.PathTemplate)
		}
	}
	{
		// GET /a/c -> /a/
		req, err := http.NewRequest("GET", "/a/c", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.PathTemplate.PathTemplate != "/a/" {
			t.Fatalf("unexpected: %v <- %#v", req.URL.Path, rd.PathTemplate)
		}
	}
	{
		// GET /x -> /
		req, err := http.NewRequest("GET", "/x", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.PathTemplate.PathTemplate != "/" {
			t.Fatalf("unexpected: %v <- %#v", req.URL.Path, rd.PathTemplate)
		}
	}
	{
		// GET /c/ -> /c
		req, err := http.NewRequest("GET", "/c/", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.PathTemplate.PathTemplate != "/c" {
			t.Fatalf("unexpected: %v <- %#v", req.URL.Path, rd.PathTemplate)
		}
	}
}

func TestRouterPickupBestRouteDefinitionFocusOnEarliestOne(t *testing.T) {
	DefaultMux = NewServeMux()

	h := func() {}

	HandleFunc("GET", "/api/todo/special", h)
	HandleFunc("GET", "/api/todo/{id}", h)

	{
		// GET /api/todo/special
		req, err := http.NewRequest("GET", "/api/todo/special", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.Method != "GET" || rd.PathTemplate.PathTemplate != "/api/todo/special" {
			t.Errorf("unexpected: %#v", rd)
		}
	}
	{
		// GET /api/todo/111 -> /api/todo/{id}
		req, err := http.NewRequest("GET", "/api/todo/111", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.Method != "GET" || rd.PathTemplate.PathTemplate != "/api/todo/{id}" {
			t.Errorf("unexpected: %#v", rd)
		}
	}
}

func TestRouterPickupBestRouteDefinitionFocusOnMultipleMethod(t *testing.T) {
	DefaultMux = NewServeMux()

	h := func() {}

	HandleFunc("OPTIONS", "/", h)
	HandleFunc("*", "/foo/bar", h)
	HandleFunc("GET,POST,PUT", "/foo", h)
	HandleFunc("PUT", "/foo/bar", h)

	{
		// OPTIONS /foo/bar -> /
		req, err := http.NewRequest("OPTIONS", "/foo/bar", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.Method != "OPTIONS" {
			t.Fatalf("unexpected: %#v", rd.Method)
		}
		if rd.PathTemplate.PathTemplate != "/" {
			t.Fatalf("unexpected: %v <- %#v", req.URL.Path, rd.PathTemplate)
		}
	}
	{
		// GET /foo
		req, err := http.NewRequest("GET", "/foo", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.Method != "GET" {
			t.Fatalf("unexpected: %#v", rd.Method)
		}
		if rd.PathTemplate.PathTemplate != "/foo" {
			t.Fatalf("unexpected: %v <- %#v", req.URL.Path, rd.PathTemplate)
		}
	}
	{
		// POST /foo/bar -> * /foo/bar
		// earlier is better
		req, err := http.NewRequest("POST", "/foo/bar", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.Method != "*" {
			t.Fatalf("unexpected: %#v", rd.Method)
		}
		if rd.PathTemplate.PathTemplate != "/foo/bar" {
			t.Fatalf("unexpected: %v <- %#v", req.URL.Path, rd.PathTemplate)
		}
	}
	{
		// PUT /foo/bar
		// exact match is better
		req, err := http.NewRequest("PUT", "/foo/bar", nil)
		if err != nil {
			t.Fatal(err)
		}
		rd := DefaultMux.router.pickupBestRouteDefinition(req)
		if rd == nil {
			t.Fatalf("unexpected")
		}
		if rd.Method != "PUT" {
			t.Fatalf("unexpected: %#v", rd.Method)
		}
		if rd.PathTemplate.PathTemplate != "/foo/bar" {
			t.Fatalf("unexpected: %v <- %#v", req.URL.Path, rd.PathTemplate)
		}
	}
}

func TestPathTemplateMatch(t *testing.T) {
	pt := ParsePathTemplate("/page/{id}")
	{
		reqURL, err := url.Parse("http://example.com/page/foo")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(reqURL.EscapedPath())
		match, params := pt.Match(reqURL.EscapedPath())
		if !match {
			t.Fatalf("unexpected")
		}
		if v := params["id"]; v != "foo" {
			t.Fatalf("unexpected: %v", v)
		}
	}
	{
		reqURL, err := url.Parse("http://example.com/page/foo%2Fbar")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(reqURL.EscapedPath())
		match, params := pt.Match(reqURL.EscapedPath())
		if !match {
			t.Fatalf("unexpected")
		}
		if v := params["id"]; v != "foo/bar" {
			t.Fatalf("unexpected: %v", v)
		}
	}
}

func TestPathTemplateMatch_fallback(t *testing.T) {
	pt := ParsePathTemplate("/")
	reqURL, err := url.Parse("http://example.com/index.js")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reqURL.EscapedPath())
	match, _ := pt.Match(reqURL.EscapedPath())
	if !match {
		t.Fatalf("unexpected")
	}
}

func TestPathTemplateMatch_aLotOfParameter(t *testing.T) {
	pt := ParsePathTemplate("/api/{id}")
	reqURL, err := url.Parse("http://example.com/api/1/2")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reqURL.EscapedPath())
	match, params := pt.Match(reqURL.EscapedPath())
	if !match {
		t.Fatalf("unexpected")
	}
	if v := params["id"]; v != "1" {
		t.Fatalf("unexpected: %v", v)
	}
}

func TestPathTemplateMatch_tooFewParameter(t *testing.T) {
	pt := ParsePathTemplate("/api/{id}")
	reqURL, err := url.Parse("http://example.com/api/")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reqURL.EscapedPath())
	match, params := pt.Match(reqURL.EscapedPath())
	if match {
		t.Fatalf("unexpected")
	}
	if _, ok := params["id"]; ok {
		t.Fatalf("unexpected: %v", ok)
	}
}

func TestPathTemplateMatch_failure(t *testing.T) {
	pt := ParsePathTemplate("/api/")
	reqURL, err := url.Parse("http://example.com/static/")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reqURL.EscapedPath())
	match, _ := pt.Match(reqURL.EscapedPath())
	if match {
		t.Fatalf("unexpected")
	}
}
