package uiserver

import (
	"io/fs"
	"net/http"
	"net/url"

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

