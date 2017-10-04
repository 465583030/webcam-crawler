package router

import (
	"net/http"
	"net/url"
	"testing"
)

func handlerMustBeCalled(called *bool) Handler {
	return func(http.ResponseWriter, *http.Request, PathParams) {
		*called = true
	}
}

func handlerMustBeCalledWithParam(called *bool, paramName string, paramValue string) Handler {
	return func(w http.ResponseWriter, r *http.Request, p PathParams) {
		if p[paramName] == paramValue {
			*called = true
		}
	}
}

func handlerMustNotBeCalled(t *testing.T) Handler {
	return func(http.ResponseWriter, *http.Request, PathParams) {
		t.Fail()
	}
}

type testController struct {
	routes []Route
}

func (c *testController) GetRoutes() []Route {
	return c.routes
}

func TestRouterRootPath(t *testing.T) {
	called := false

	router := New(handlerMustNotBeCalled(t))
	router.Mount("/", &testController{
		[]Route{
			Route{"GET", "/", handlerMustBeCalled(&called)},
		},
	})

	router.ServeHTTP(nil, &http.Request{
		Method: "GET",
		URL: &url.URL{
			Path: "/",
		},
	})

	if !called {
		t.Fail()
	}
}

func TestRouterMountPath(t *testing.T) {
	called := false

	router := New(handlerMustNotBeCalled(t))
	router.Mount("/there", &testController{
		[]Route{
			Route{"GET", "/", handlerMustBeCalled(&called)},
		},
	})

	router.ServeHTTP(nil, &http.Request{
		Method: "GET",
		URL: &url.URL{
			Path: "/there",
		},
	})

	if !called {
		t.Fail()
	}
}

func TestRouterMountPathAndControllerPath(t *testing.T) {
	called := false

	router := New(handlerMustNotBeCalled(t))
	router.Mount("/there", &testController{
		[]Route{
			Route{"GET", "/here", handlerMustBeCalled(&called)},
			Route{"GET", "/", handlerMustNotBeCalled(t)},
		},
	})

	router.ServeHTTP(nil, &http.Request{
		Method: "GET",
		URL: &url.URL{
			Path: "/there/here",
		},
	})

	if !called {
		t.Fail()
	}
}

func TestRouterPathParams(t *testing.T) {
	called := false

	router := New(handlerMustNotBeCalled(t))
	router.Mount("/", &testController{
		[]Route{
			Route{"GET", "/:name", handlerMustBeCalledWithParam(&called, "name", "toto")},
			Route{"GET", "/", handlerMustNotBeCalled(t)},
		},
	})

	router.ServeHTTP(nil, &http.Request{
		Method: "GET",
		URL: &url.URL{
			Path: "/toto",
		},
	})

	if !called {
		t.Fail()
	}
}

func TestRouterAmbiguousPathParams(t *testing.T) {
	called := false

	router := New(handlerMustNotBeCalled(t))
	router.Mount("/", &testController{
		[]Route{
			Route{"GET", "/:name", handlerMustNotBeCalled(t)},
			Route{"GET", "/toto/:special", handlerMustBeCalledWithParam(&called, "special", "titi")},
		},
	})

	router.ServeHTTP(nil, &http.Request{
		Method: "GET",
		URL: &url.URL{
			Path: "/toto/titi",
		},
	})

	if !called {
		t.Fail()
	}
}

func TestRouterDefaultHandler(t *testing.T) {
	called := false

	router := New(handlerMustBeCalled(&called))
	router.Mount("/", &testController{
		[]Route{
			Route{"GET", "/", handlerMustNotBeCalled(t)},
		},
	})

	router.ServeHTTP(nil, &http.Request{
		Method: "GET",
		URL: &url.URL{
			Path: "/notfound",
		},
	})

	if !called {
		t.Fail()
	}
}

func TestRouterNotFoundMethod(t *testing.T) {
	called := false

	router := New(handlerMustBeCalled(&called))
	router.Mount("/", &testController{
		[]Route{
			Route{"GET", "/", handlerMustNotBeCalled(t)},
		},
	})

	router.ServeHTTP(nil, &http.Request{
		Method: "POST",
		URL: &url.URL{
			Path: "/",
		},
	})

	if !called {
		t.Fail()
	}
}

func TestRouterNotInitializedDoesNotPanic(t *testing.T) {
	router := &Router{}
	router.Mount("/", &testController{
		[]Route{
			Route{"GET", "/a", handlerMustNotBeCalled(t)},
		},
	})

	router.ServeHTTP(nil, &http.Request{
		Method: "GET",
		URL: &url.URL{
			Path: "/",
		},
	})
}
