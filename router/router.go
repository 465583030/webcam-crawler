package router

import (
	"net/http"
)

// PathParams contains the parameters contained in a route.
type PathParams map[string]string

// A Handler is a function handling requests.
type Handler func(http.ResponseWriter, *http.Request, PathParams)

// A Route matches a method and a path to a Handler.
type Route struct {
	Method  string
	Path    string
	Handler Handler
}

// A Router serves requests with its registered controllers.
type Router struct {
	root           *node
	defaultHandler Handler
}

// A Controller defines a slice of routes.
type Controller interface {
	GetRoutes() []Route
}

// New creates a new Router.
func New(defaultHandler Handler) *Router {
	return &Router{
		newTree(),
		defaultHandler,
	}
}

// Mount a controller on a path.
func (r *Router) Mount(path string, controller Controller) {
	r.createRootIfNeeded()

	for _, route := range controller.GetRoutes() {
		route.Path = path + route.Path
		r.root.addRoute(route)
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.createRootIfNeeded()

	if node, params := r.root.findNode(req.URL.Path); node != nil {
		if handler := node.handlers[req.Method]; handler != nil {
			handler(w, req, params)
		} else {
			r.callDefaultHandler(w, req)
		}
	} else {
		r.callDefaultHandler(w, req)
	}
}

func (r *Router) createRootIfNeeded() {
	if r.root == nil {
		r.root = newTree()
	}
}

func (r *Router) callDefaultHandler(w http.ResponseWriter, req *http.Request) {
	if r.defaultHandler != nil {
		r.defaultHandler(w, req, make(PathParams))
	}
}
