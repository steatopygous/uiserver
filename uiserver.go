package uiserver

import (
	"fmt"
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
func New(fs fs.FS) UIServer {
	root := getUIRoot(fs)

	var handlers []pathMethodHandler

	muxRouter := mux.NewRouter()

	return UIServer{handlers, root, muxRouter}
}


// Run() starts the server running on the specified port.
func (server UIServer) Run(port string) error {
	server.mux.PathPrefix("/").Handler(http.FileServer(http.FS(server.root)))

	err := http.ListenAndServe(port, server.mux)

	if err != nil {
		fmt.Println("uiserver.Run() - ListenAndServe() returned an error...", err)
		return err
	}

	return nil
}


// getUIRoot() traverses the provided file system to find the UI content.
func getUIRoot(ui fs.FS) fs.FS {
	// The content folder is the first folder that contains more
	// than one file.

	for {
		items, _ := fs.ReadDir(ui, ".")

		if len(items) > 1 {
			break
		}

		contentFolder := items[0].Name()

		ui, _ = fs.Sub(ui, contentFolder)
	}

	return ui
}

