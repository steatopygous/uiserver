package uiserver

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"regexp"

	"github.com/gorilla/mux"
)

type UIServer struct {
	handlers []pathMethodHandler
	root     fs.FS
	mux      *mux.Router
}

type Context struct {
	Route    string
	Method   string
	Path     string
	Vars     map[string]string
	Query    url.Values
	Response http.ResponseWriter
	Request  *http.Request
}

func New(ui fs.FS) UIServer {
	root := getUIRoot(ui)

	var handlers []pathMethodHandler

	muxRouter := mux.NewRouter()

	return UIServer{handlers, root, muxRouter}
}

func (server UIServer) Run(port string) {
	server.mux.PathPrefix("/").Handler(http.FileServer(http.FS(server.root)))

	_ = http.ListenAndServe(port, server.mux)
}

func (server UIServer) Get(path string, handler handler) {
	server.addHandler(path, "GET", handler)
}

func (server UIServer) Post(path string, handler handler) {
	server.addHandler(path, "POST", handler)
}

func (server UIServer) Put(path string, handler handler) {
	server.addHandler(path, "PUT", handler)
}

func (server UIServer) Patch(path string, handler handler) {
	server.addHandler(path, "PATCH", handler)
}

func (server UIServer) Delete(path string, handler handler) {
	server.addHandler(path, "DELETE", handler)
}

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

// Implementation details

func getUIRoot(ui fs.FS) fs.FS {
	// We expect the top level of the file system to contain a single folder
	// that contains all of the UI content, potentially in sub-folders.

	items, _ := fs.ReadDir(ui, ".")

	if len(items) != 1 {
		panic("the top level of the UI content file system can only contain a single folder")
	}

	contentFolder := items[0].Name()

	// Since we want the UI content to be accessible as "/", we need to use
	// the sub-folder as the root of our file system.

	root, _ := fs.Sub(ui, contentFolder)

	return root
}

type handler func(context Context)

func (server UIServer) addHandler(route string, method string, handler handler) {
	path := pathForRoute(route)

	server.handlers = append(server.handlers, pathMethodHandler{route, path, method, handler})

	server.mux.Handle(route, server)
}

func pathForRoute(route string) string {
	re := regexp.MustCompile(`{.*?}`)

	return re.ReplaceAllString(route, `{}`)
}

func (server UIServer) handlerFor(path string, method string) (pathMethodHandler, error) {
	for _, item := range server.handlers {
		if item.matchesPathAndMethod(path, method) {
			return item, nil
		}
	}

	return pathMethodHandler{"", "", "", nil}, errors.New(fmt.Sprintf("No handler for %s on %s", method, path))
}

