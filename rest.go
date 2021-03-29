package uiserver

import (
	"regexp"
)

func (server UIServer) Get(path string, handler Handler) {
	server.addHandler(path, "GET", handler)
}

func (server UIServer) Post(path string, handler Handler) {
	server.addHandler(path, "POST", handler)
}

func (server UIServer) Put(path string, handler Handler) {
	server.addHandler(path, "PUT", handler)
}

func (server UIServer) Patch(path string, handler Handler) {
	server.addHandler(path, "PATCH", handler)
}

func (server UIServer) Delete(path string, handler Handler) {
	server.addHandler(path, "DELETE", handler)
}


// Implementation details

type Handler func(context Context)

func (server UIServer) addHandler(route string, method string, handler Handler) {
	path := pathForRoute(route)

	server.handlers = append(server.handlers, pathMethodHandler{route, path, method, handler})

	server.mux.Handle(route, server)
}

func pathForRoute(route string) string {
	re := regexp.MustCompile(`{.*?}`)

	return re.ReplaceAllString(route, `{}`)
}

