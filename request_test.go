package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEmptyRequest(t *testing.T) {
	req := NewRequest("GET", "http://example.com")
	if req.Method != "GET" {
		t.Errorf("expected method GET, got %s", req.Method)
	}

	if req.URL != "http://example.com" {
		t.Errorf("expected URL http://example.com, got %s", req.URL)
	}
}

func TestRequestHeaders(t *testing.T) {
	req := NewRequest("GET", "http://example.com")
	req.AddHeader("Content-Type", "application/json")
	req.AddHeader("Accept", "application/json")

	if req.GetHeader("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", req.GetHeader("Content-Type"))
	}

	if req.GetHeader("Accept") != "application/json" {
		t.Errorf("expected Accept application/json, got %s", req.GetHeader("Accept"))
	}
}

func TestRequestQueries(t *testing.T) {
	req := NewRequest("GET", "http://example.com")
	req.AddQuery("foo", "bar")
	req.AddQuery("baz", "qux")

	if req.GetQuery("foo") != "bar" {
		t.Errorf("expected foo=bar, got %s", req.GetQuery("foo"))
	}

	if req.GetQuery("baz") != "qux" {
		t.Errorf("expected baz=qux, got %s", req.GetQuery("baz"))
	}
}

func TestGetUnsetHeader(t *testing.T) {
	req := NewRequest("GET", "http://example.com")
	if req.GetHeader("Content-Type") != "" {
		t.Errorf("expected empty string, got %s", req.GetHeader("Content-Type"))
	}
}

func TestGetUnsetQuery(t *testing.T) {
	req := NewRequest("GET", "http://example.com")
	if req.GetQuery("foo") != "" {
		t.Errorf("expected empty string, got %s", req.GetQuery("foo"))
	}
}

func TestMakeWithEmptyRequest(t *testing.T) {
	req := NewRequest("GET", "http://example.com")
	httpReq := req.make()
	if httpReq.Method != "GET" {
		t.Errorf("expected method GET, got %s", httpReq.Method)
	}

	if httpReq.URL.String() != "http://example.com" {
		t.Errorf("expected URL http://example.com, got %s", httpReq.URL.String())
	}
}

func TestMake(t *testing.T) {
	req := NewRequest("GET", "http://example.com")
	req.AddHeader("Content-Type", "application/json")
	req.AddHeader("Accept", "application/json")
	req.AddQuery("foo", "bar")
	req.AddQuery("baz", "qux")

	httpReq := req.make()
	if httpReq.Method != "GET" {
		t.Errorf("expected method GET, got %s", httpReq.Method)
	}

	if httpReq.URL.String() != "http://example.com?baz=qux&foo=bar" {
		t.Errorf("expected URL http://example.com?baz=qux&foo=bar, got %s", httpReq.URL.String())
	}

	if httpReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", httpReq.Header.Get("Content-Type"))
	}

	if httpReq.Header.Get("Accept") != "application/json" {
		t.Errorf("expected Accept application/json, got %s", httpReq.Header.Get("Accept"))
	}
}

func TestMakeWithBody(t *testing.T) {
	req := NewRequest("GET", "http://example.com")
	req.AddHeader("Content-Type", "application/json")
	req.AddHeader("Accept", "application/json")
	req.AddQuery("foo", "bar")
	req.AddQuery("baz", "qux")
	req.AddBody(strings.NewReader("hello world"))

	httpReq := req.make()
	if httpReq.Method != "GET" {
		t.Errorf("expected method GET, got %s", httpReq.Method)
	}

	if httpReq.URL.String() != "http://example.com?baz=qux&foo=bar" {
		t.Errorf("expected URL http://example.com?baz=qux&foo=bar, got %s", httpReq.URL.String())
	}

	if httpReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", httpReq.Header.Get("Content-Type"))
	}

	if httpReq.Header.Get("Accept") != "application/json" {
		t.Errorf("expected Accept application/json, got %s", httpReq.Header.Get("Accept"))
	}

	if httpReq.Body == nil {
		t.Errorf("expected non-nil body, got nil")
	}
}

func TestAddBodyString(t *testing.T) {
	req := NewRequest("GET", "http://example.com")
	req.AddBodyString("hello world")

	if req.Body == nil {
		t.Errorf("expected non-nil body, got nil")
	}
}

func TestAddBodyJSON(t *testing.T) {
	req := NewRequest("GET", "http://example.com")
	req.AddBodyJSON(map[string]string{"foo": "bar"})

	if req.Body == nil {
		t.Errorf("expected non-nil body, got nil")
	}
}

func TestAddBodyJSONError(t *testing.T) {
	req := NewRequest("GET", "http://example.com")
	req.AddBodyJSON(make(chan int))

	if req.Body != nil {
		t.Errorf("expected nil body, got %v", req.Body)
	}
}

func TestSend(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	req := NewRequest("GET", ts.URL)
	resp, err := req.Send()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	if resp.Duration == 0 {
		t.Errorf("expected non-zero duration, got %d", resp.Duration)
	}
}