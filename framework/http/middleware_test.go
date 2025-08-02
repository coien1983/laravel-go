package http

import (
	"net/http"
	"testing"
)

// mockRequest/Response 用于测试

type mockRequest struct {
	path    string
	headers map[string]string
	params  map[string]string
}

func (r *mockRequest) Method() string                                             { return "GET" }
func (r *mockRequest) URL() string                                                { return r.path }
func (r *mockRequest) Path() string                                               { return r.path }
func (r *mockRequest) Header(key string) string                                   { return r.headers[key] }
func (r *mockRequest) Headers() map[string]string                                 { return r.headers }
func (r *mockRequest) IP() string                                                 { return "127.0.0.1" }
func (r *mockRequest) UserAgent() string                                          { return "test-agent" }
func (r *mockRequest) Param(key string) string                                    { return r.params[key] }
func (r *mockRequest) Params() map[string]string                                  { return r.params }
func (r *mockRequest) Query(key string) string                                    { return "" }
func (r *mockRequest) QueryInt(key string, _ ...int) int                          { return 0 }
func (r *mockRequest) QueryBool(key string, _ ...bool) bool                       { return false }
func (r *mockRequest) QueryArray(key string) []string                             { return nil }
func (r *mockRequest) Body() []byte                                               { return nil }
func (r *mockRequest) Json(data interface{}) error                                { return nil }
func (r *mockRequest) Form(key string) string                                     { return "" }
func (r *mockRequest) File(key string) ([]byte, string, error)                    { return nil, "", nil }
func (r *mockRequest) Cookie(key string) string                                   { return "" }
func (r *mockRequest) Raw() *http.Request                                         { return nil }

// mockResponse 用于测试

type mockResponse struct {
	status  int
	data    interface{}
	headers map[string]string
}

func (r *mockResponse) Status() int                { return r.status }
func (r *mockResponse) Data() interface{}          { return r.data }
func (r *mockResponse) Headers() map[string]string { return r.headers }
func (r *mockResponse) SetHeader(key, value string) Response {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	r.headers[key] = value
	return r
}
func (r *mockResponse) Send(_ http.ResponseWriter) {}

func newMockResponse(status int, data interface{}) Response {
	return &mockResponse{status: status, data: data, headers: make(map[string]string)}
}

// mockHandler 用于测试

type mockHandler struct {
	called *bool
}

func (h *mockHandler) Handle(request Request) Response {
	if h.called != nil {
		*h.called = true
	}
	return newMockResponse(200, "handler")
}

func TestPipelineChain(t *testing.T) {
	pipeline := NewPipeline()
	order := []string{}

	mw1 := MiddlewareFunc(func(req Request, next Next) Response {
		order = append(order, "mw1")
		return next(req)
	})
	mw2 := MiddlewareFunc(func(req Request, next Next) Response {
		order = append(order, "mw2")
		return next(req)
	})
	pipeline.Use(mw1, mw2)

	called := false
	handler := &mockHandler{called: &called}
	resp := pipeline.Process(&mockRequest{path: "/"}, handler)

	if !called {
		t.Error("Handler should be called")
	}
	if resp.Status() != 200 || resp.Data() != "handler" {
		t.Error("Handler response incorrect")
	}
	if len(order) != 2 || order[0] != "mw1" || order[1] != "mw2" {
		t.Errorf("Middleware order incorrect: %v", order)
	}
}

func TestConditionalMiddleware(t *testing.T) {
	pipeline := NewPipeline()
	flag := false
	condMw := NewConditionalMiddleware(
		func(req Request) bool { return req.Path() == "/pass" },
		MiddlewareFunc(func(req Request, next Next) Response {
			flag = true
			return next(req)
		}),
	)
	pipeline.Use(condMw)

	handler := &mockHandler{}
	pipeline.Process(&mockRequest{path: "/block"}, handler)
	if flag {
		t.Error("Conditional middleware should not run for /block")
	}
	pipeline.Process(&mockRequest{path: "/pass"}, handler)
	if !flag {
		t.Error("Conditional middleware should run for /pass")
	}
}

func TestPriorityPipeline(t *testing.T) {
	pp := NewPriorityPipeline()
	order := []string{}
	pp.Use(10, MiddlewareFunc(func(req Request, next Next) Response {
		order = append(order, "high")
		return next(req)
	}))
	pp.Use(1, MiddlewareFunc(func(req Request, next Next) Response {
		order = append(order, "low")
		return next(req)
	}))

	handler := &mockHandler{}
	pp.Process(&mockRequest{path: "/"}, handler)

	if len(order) != 2 || order[0] != "high" || order[1] != "low" {
		t.Errorf("Priority order incorrect: %v", order)
	}
}

func TestMiddlewareSetHeader(t *testing.T) {
	mw := MiddlewareFunc(func(req Request, next Next) Response {
		resp := next(req)
		resp.SetHeader("X-Test", "ok")
		return resp
	})
	pipeline := NewPipeline().Use(mw)
	handler := &mockHandler{}
	resp := pipeline.Process(&mockRequest{path: "/"}, handler)
	if resp.Headers()["X-Test"] != "ok" {
		t.Error("Middleware should set header")
	}
}
