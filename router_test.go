package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func helloHandler(w http.ResponseWriter, r *http.Request, p PathParams) error {
	fmt.Fprintf(w, "Hello from %s", r.URL.Path)
	return nil
}

func httpErrorHandler(w http.ResponseWriter, r *http.Request, p PathParams) error {
	return StatusError{404, errors.New("test HTTPError")}
}

func errorHandler(w http.ResponseWriter, r *http.Request, p PathParams) error {
	return errors.New("test error")
}

func handlerMustNotBeCalled(t *testing.T) Handler {
	return func(http.ResponseWriter, *http.Request, PathParams) error {
		t.Fail()
		return nil
	}
}

type testController struct {
	routes []Route
}

func (c *testController) GetRoutes() []Route {
	return c.routes
}

func getResponseBody(r *http.Response) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	return buf.String()
}

func TestRouterRootPath(t *testing.T) {
	router := NewRouter(handlerMustNotBeCalled(t))
	router.Mount("/", &testController{
		[]Route{
			Route{"GET", "/", helloHandler},
		},
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	res := rec.Result()

	if res.Status != "200 OK" {
		t.Errorf("Expected 200 status code, got %s\n", res.Status)
	}

	if getResponseBody(res) != "Hello from /" {
		t.Errorf("Unexpected request body: %s\n", getResponseBody(res))
	}
}

func TestRouterMountPath(t *testing.T) {
	router := NewRouter(handlerMustNotBeCalled(t))
	router.Mount("/there", &testController{
		[]Route{
			Route{"GET", "/", helloHandler},
		},
	})

	req := httptest.NewRequest("GET", "/there", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	res := rec.Result()

	if res.Status != "200 OK" {
		t.Errorf("Expected 200 status code, got %s\n", res.Status)
	}

	if getResponseBody(res) != "Hello from /there" {
		t.Errorf("Unexpected request body: %s\n", getResponseBody(res))
	}
}

func TestRouterMountPathAndControllerPath(t *testing.T) {
	router := NewRouter(handlerMustNotBeCalled(t))
	router.Mount("/there", &testController{
		[]Route{
			Route{"GET", "/here", helloHandler},
			Route{"GET", "/", handlerMustNotBeCalled(t)},
		},
	})

	req := httptest.NewRequest("GET", "/there/here", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	res := rec.Result()

	if res.Status != "200 OK" {
		t.Errorf("Expected 200 status code, got %s\n", res.Status)
	}

	if getResponseBody(res) != "Hello from /there/here" {
		t.Errorf("Unexpected request body: %s\n", getResponseBody(res))
	}
}

func TestRouterPathParams(t *testing.T) {
	router := NewRouter(handlerMustNotBeCalled(t))
	router.Mount("/", &testController{
		[]Route{
			Route{"GET", "/:name", helloHandler},
			Route{"GET", "/", handlerMustNotBeCalled(t)},
		},
	})

	req := httptest.NewRequest("GET", "/toto", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	res := rec.Result()

	if res.Status != "200 OK" {
		t.Errorf("Expected 200 status code, got %s\n", res.Status)
	}

	if getResponseBody(res) != "Hello from /toto" {
		t.Errorf("Unexpected request body: %s\n", getResponseBody(res))
	}
}

func TestRouterAmbiguousPathParams(t *testing.T) {
	router := NewRouter(handlerMustNotBeCalled(t))
	router.Mount("/", &testController{
		[]Route{
			Route{"GET", "/:name", handlerMustNotBeCalled(t)},
			Route{"GET", "/toto/:special", helloHandler},
		},
	})

	req := httptest.NewRequest("GET", "/toto/titi", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	res := rec.Result()

	if res.Status != "200 OK" {
		t.Errorf("Expected 200 status code, got %s\n", res.Status)
	}

	if getResponseBody(res) != "Hello from /toto/titi" {
		t.Errorf("Unexpected request body: %s\n", getResponseBody(res))
	}
}

func TestRouterDefaultHandler(t *testing.T) {
	router := NewRouter(helloHandler)
	router.Mount("/", &testController{
		[]Route{
			Route{"GET", "/", handlerMustNotBeCalled(t)},
		},
	})

	req := httptest.NewRequest("GET", "/notfound", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	res := rec.Result()

	if res.Status != "200 OK" {
		t.Errorf("Expected 200 status code, got %s\n", res.Status)
	}

	if getResponseBody(res) != "Hello from /notfound" {
		t.Errorf("Unexpected request body: %s\n", getResponseBody(res))
	}
}

func TestRouterNotFoundMethod(t *testing.T) {
	router := NewRouter(helloHandler)
	router.Mount("/", &testController{
		[]Route{
			Route{"GET", "/", handlerMustNotBeCalled(t)},
		},
	})

	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	res := rec.Result()

	if res.Status != "200 OK" {
		t.Errorf("Expected 200 status code, got %s\n", res.Status)
	}

	if getResponseBody(res) != "Hello from /" {
		t.Errorf("Unexpected request body: %s\n", getResponseBody(res))
	}
}

func TestRouterNotInitializedDoesNotPanic(t *testing.T) {
	router := &Router{}
	router.Mount("/", &testController{
		[]Route{
			Route{"GET", "/a", handlerMustNotBeCalled(t)},
		},
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	res := rec.Result()

	if res.Status != "200 OK" {
		t.Errorf("Expected 200 status code, got %s\n", res.Status)
	}
}

func TestRouterHandlerHTTPError(t *testing.T) {
	router := NewRouter(handlerMustNotBeCalled(t))
	router.Mount("/", &testController{
		[]Route{
			Route{"GET", "/", httpErrorHandler},
		},
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	res := rec.Result()

	if res.StatusCode != 404 {
		t.Errorf("Expected 404 status code, got %d\n", res.StatusCode)
	}
}

func TestRouterHandlerError(t *testing.T) {
	router := NewRouter(handlerMustNotBeCalled(t))
	router.Mount("/", &testController{
		[]Route{
			Route{"GET", "/", errorHandler},
		},
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	res := rec.Result()

	if res.StatusCode != 500 {
		t.Errorf("Expected 500 status code, got %d\n", res.StatusCode)
	}
}
