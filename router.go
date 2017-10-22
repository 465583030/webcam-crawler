package main

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

// PathParams contains the parameters contained in a route.
type PathParams map[string]string

// A Handler is a function handling requests.
type Handler func(http.ResponseWriter, *http.Request, PathParams) error

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

// An HTTPError exposes an error and an HTTP status code.
type HTTPError interface {
	error
	Status() int
}

// StatusError implements HTTPError interface.
type StatusError struct {
	Code int
	Err  error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

// Status returns the HTTP status code associated with the error.
func (se StatusError) Status() int {
	return se.Code
}

// NewRouter creates a new Router.
func NewRouter(defaultHandler Handler) *Router {
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
	handler, params, err := r.getHandler(req)
	if err != nil {
		err = r.callDefaultHandler(w, req)
	} else {
		err = handler(w, req, params)
	}

	if err != nil {
		switch e := err.(type) {
		case HTTPError:
			log.Printf("HTTP %d : %s\n", e.Status(), e)
			http.Error(w, http.StatusText(e.Status()), e.Status())
		default:
			log.Printf("ERROR: %s\n", e)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func (r *Router) createRootIfNeeded() {
	if r.root == nil {
		r.root = newTree()
	}
}

func (r *Router) getHandler(req *http.Request) (Handler, PathParams, error) {
	r.createRootIfNeeded()

	if node, params := r.root.findNode(req.URL.Path); node != nil {
		if handler, ok := node.handlers[req.Method]; ok {
			return handler, params, nil
		}

		return nil, nil, errors.New("No handler for method " + req.Method + " for path " + req.URL.Path)
	}

	return nil, nil, errors.New("No handler for path " + req.URL.Path)
}

func (r *Router) callDefaultHandler(w http.ResponseWriter, req *http.Request) error {
	if r.defaultHandler != nil {
		return r.defaultHandler(w, req, make(PathParams))
	}

	// If no default handler handler is set, return a successful empty response.
	return nil
}

type node struct {
	children []*node
	value    string
	isParam  bool
	handlers map[string]Handler
}

func (n *node) addNode(method string, path []string, handler Handler) {
	if len(path) == 0 {
		// This is the destination node, set the handler
		n.handlers[method] = handler
		return
	}

	value := path[0]

	// Look for next path node in children
	for _, c := range n.children {
		if c.value == value {
			c.addNode(method, path[1:], handler)
			return
		}
	}

	// If not found create a new node
	newValue := value
	isParam := len(value) != 0 && value[0] == ':'
	if isParam {
		newValue = value[1:]
	}

	newNode := node{
		value:    newValue,
		isParam:  isParam,
		handlers: make(map[string]Handler),
	}

	n.children = append(n.children, &newNode)
	newNode.addNode(method, path[1:], handler)
}

func (n *node) traverse(path []string, params PathParams) (*node, PathParams) {
	// Stop recursion if the path is empty
	if len(path) == 0 {
		return n, params
	}

	// Look for next path node in children
	for _, c := range n.children {
		if c.value == path[0] || c.isParam {
			// Verify that the rest of the path matches a route
			if result, _ := c.traverse(path[1:], params); result != nil {
				// If the node is a parameter, add it in the list
				if c.isParam {
					params[c.value] = path[0]
				}

				return result, params
			}
		}
	}

	return nil, params
}

func (n *node) findNode(path string) (*node, PathParams) {
	params := make(PathParams)
	return n.traverse(splitPath(path), params)
}

func newTree() *node {
	return &node{
		value:    "",
		isParam:  false,
		handlers: make(map[string]Handler),
		children: make([]*node, 0),
	}
}

func (n *node) addRoute(r Route) {
	n.addNode(r.Method, splitPath(r.Path), r.Handler)
}

func splitPath(path string) []string {
	result := make([]string, 0)

	values := strings.Split(path, "/")
	for _, v := range values {
		if v != "" {
			result = append(result, v)
		}
	}

	return result
}
