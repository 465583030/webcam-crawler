package router

import (
	"strings"
)

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