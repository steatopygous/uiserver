package uiserver

import (
	"regexp"
)

// Get() registers a handler for GET calls on given route.
func (server *UIServer) Get(route string, handler Handler) {
	server.addHandler(route, "GET", handler)
}

// Post() registers a handler for POST calls on given route.
func (server *UIServer) Post(route string, handler Handler) {
	server.addHandler(route, "POST", handler)
}

// Put() registers a handler for PUT calls on given route.
func (server *UIServer) Put(route string, handler Handler) {
	server.addHandler(route, "PUT", handler)
}

// Patch() registers a handler for PATCH calls on the given route
func (server *UIServer) Patch(route string, handler Handler) {
	server.addHandler(route, "PATCH", handler)
}

// Patch() registers a handler for PATCH calls on given route.
func (server *UIServer) Delete(route string, handler Handler) {
	server.addHandler(route, "DELETE", handler)
}



// Handler is similar to http.Handler, but provides for functions that
// accept a uiserver.Context rather than http.ResponseWriter and http.Request.
type Handler func(context Context)


// addHandler() registers a handler for a route and method.
func (server *UIServer) addHandler(route string, method string, handler Handler) {
	path := pathForRoute(route)

	server.handlers = append(server.handlers, pathMethodHandler{route, path, method, handler})

	server.mux.Handle(route, server)
}


// pathForRoute() replaces any route variables with just a placeholder.  We convert the
// registered route into this form to make comparison of incoming paths to routes independent
// of the specific route variable values.
func pathForRoute(route string) string {
	re := regexp.MustCompile(`{.*?}`)

	return re.ReplaceAllString(route, `{}`)
}

