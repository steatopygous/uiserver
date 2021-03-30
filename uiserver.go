package uiserver

import (
	"io/fs"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)


// UIServer is the type returned by New(). All of its fields are private and purely
// for use by uiserver's implementation.

type UIServer struct {
	handlers []pathMethodHandler
	root     fs.FS
	mux      *mux.Router
}


// Context is the type passed to handler functions.  It provides the Route and REST Method
// that are being handled, the Path (which is the same as the route, but with specific values
// provided for the route variables).
//
// It also surfaces data from the underlying gorilla/mux router, including the Request, Writer
// along with a map of the route variables and the Query parameters.

type Context struct {
	Route   string
	Method  string
	Path    string
	Vars    map[string]string
	Query   url.Values
	Writer  http.ResponseWriter
	Request *http.Request
}


// New() constructs a new UIServer that uses the provided file system as the source for
// the UI's content.  Using this, the developer will make calls to the various REST verb
// functions to register handlers for the various resource paths the UI code can access.

func New(ui fs.FS) UIServer {
	root := getUIRoot(ui)

	var handlers []pathMethodHandler

	muxRouter := mux.NewRouter()

	return UIServer{handlers, root, muxRouter}
}


// Run() starts the server running on the specified port.

func (server UIServer) Run(port string) {
	server.mux.PathPrefix("/").Handler(http.FileServer(http.FS(server.root)))

	_ = http.ListenAndServe(port, server.mux)
}


// getUIRoot() extracts the top-level folder from the UI file system.  This allows
// all URLs referenced in the content to be relative to "/", rather than that folder.

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

