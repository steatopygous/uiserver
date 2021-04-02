package uiserver

import (
	"strings"
)

// pathMethodHandler collects together a route, path and method, along with the
// handler function for the given route + method.
type pathMethodHandler struct {
	route   string
	path    string
	method  string
	handler Handler
}


// matchesPathAndMethod() checks whether a given pathMethodHandler instance corresponds
// to the provided path and method parameters.
func(pmh pathMethodHandler) matchesPathAndMethod(path string, method string) bool {
	if pmh.method != method {
		return false
	}

	pathParts := strings.Split(path, "/")
	parts := strings.Split(pmh.path, "/")

	if len(parts) != len(pathParts) {
		return false
	}

	for i := 0; i < len(parts); i++ {
		if parts[i] != "{}" && parts[i] != pathParts[i] {
			return false
		}
	}

	return true
}



// ByPrecedence represents a slice of request handlers for purposes
// of sorting by route precedence.
type ByPrecedence []pathMethodHandler

// Len() returns the length of a slice of handlers
func (handlers ByPrecedence) Len() int {
	return len(handlers)
}

// Less() determines the precedence of two handlers' routes
func (handlers ByPrecedence) Less(i, j int) bool {
	routeI := handlers[i].route
	routeJ := handlers[j].route

	// Routes with more parts are "less than" those with longer paths.

	partsI := strings.Split(routeI, "/")
	partsJ := strings.Split(routeJ, "/")

	if len(partsI) > len(partsJ) {
		return true
	}

	if len(partsI) < len(partsJ) {
		return false
	}

	// Routes that have no variables are "less than" those that do.

	hasVariablesI := hasVariables(routeI)
	hasVariablesJ := hasVariables(routeJ)

	if hasVariablesI && !hasVariablesJ {
		return false
	}

	if hasVariablesJ && !hasVariablesI {
		return true
	}

	// In other cases, we don't care

	return false
}

// Swap() swaps two handlers in a slice, given their indices
func (handlers ByPrecedence) Swap(i, j int) {
	handlers[i], handlers[j] = handlers[j], handlers[i]
}

// hasVariables() determines whether a given route contains variables
func hasVariables(route string) bool {
	return strings.Contains(route, "{")
}

