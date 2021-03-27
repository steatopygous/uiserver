package uiserver

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (server UIServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	method := request.Method

	item, err := server.handlerFor(path, method)

	if err != nil {
		fmt.Println("Sorry, I don't know how to handle", method, "for path", path)

		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		path := request.URL.Path
		method := request.Method
		vars := mux.Vars(request)
		query := request.URL.Query()

		context := Context{item.route, method, path, vars, query, writer, request}

		item.handler(context)
	}
}

