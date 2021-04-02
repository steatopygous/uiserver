package uiserver

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)


// ServeHTTP() implements the http.Handler interface.  It routes the request to the
// appropriate handler, based on the path and method.
func (server UIServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	method := request.Method

	item, err := server.handlerFor(path, method)

	if err != nil {
		fmt.Println("Sorry, I don't know how to handle", method, "for path", path)

		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(request)
	query := request.URL.Query()

	context := Context{item.route, method, path, vars, query, writer, request}

	item.handler(context)
}


// handlerFor() finds the registered handler function, if any, for the given path and method.
func (server UIServer) handlerFor(path string, method string) (pathMethodHandler, error) {
	for _, item := range server.handlers {
		if item.matchesPathAndMethod(path, method) {
			fmt.Println(path, method, "is handled by", item)
			return item, nil
		}
	}

	return pathMethodHandler{"", "", "", nil}, errors.New(fmt.Sprintf("No Handler for %s on %s", method, path))
}

