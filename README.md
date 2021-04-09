# uiserver
A simple Go HTTP server that combines UI content and a REST API, based on
[gorilla/mux](https://github.com/gorilla/mux), for the content to consume.

The intended purpose of this package is to allow creation of Go desktop apps that have a browser-based user interface, with the app providing a REST API to support the UI.

The idea is to build something similar to an Electron app, but with the underlying business logic written in Go and the UI displayed in a browser tab, rather than a full-blown windowed application.

## Usage

We construct a server by calling **uiserver.New()**, which returns a **UIServer** instance. The single parameter that is passed to uiserver.New() is a file system that contains the files that provide the user interface for the application.

```go
//go:embed ui/public/*
var ui embed.FS

func main() {
  server := uiserver.New(ui)
  
  //...
  
  server.Run(":1234")
}
```

This is sufficient to make the user interface available at **http://localhost:1234/**, but if we want the UI to persist data across sessions, we need to provide an API it can use to achieve that.

So, before calling **Run()** we can register routes along with handlers for them.

```go
server.Get("/todos", GetAllToDos)
server.Get("/todos/{id}", GetToDo)

server.Post("/todos", CreateNewToDo)
server.Patch("/todos/{id}", UpdateToDo)
server.Put("/todos/{id}", ReplaceToDo)
server.Delete("/todos/{id}", DeleteToDo)
```

When a request comes in that matches a provided route and method, the corresponding handler is called.  To avoid having a surfeit of parameters, it is passed a **uiserver.Context** struct that looks like

```go
type Context struct {
	Route   string
	Method  string
	Path    string
	Vars    map[string]string
	Query   url.Values
	Writer  http.ResponseWriter
	Request *http.Request
}
```

We recommend that each handler is specific to one route/method combination, so in most cases, the **Route**, **Method** and **Path** won't be of particular interest.  Route is what was passed when registering the handler; Path is the same, but with the route variables set to the actual values based on the request.

The  **Writer** and **Request** are the same as a normal handler function would be called with.

**Vars** is a map that provides the value of each route variable in the specific request and **Query** gives  access to the contents of the query string, if any.

Vars is a string to string map, whereas Query is a **[url.Values](https://golang.org/pkg/net/url/#Values)** type that maps a name to a slice of strings, as provided by the underlying **[gorilla/mux](gorilla/mux)** package that uiserver uses to provide the actual server.

### An Example Handler

An implementation of **GetToDo()** might look like

```go
func GetToDo(c Context) {
  idString, ok := c.Vars["id"]

  if !ok {
    app.Respond(c, nil, 400)
    return
  }

  id, err := strconv.Atoi(idString)

  if err != nil {
    app.Respond(c, nil, 400)
    return
  }

  item, ok := ... // Code to retrieve the todo with the given id

  if ok {
    app.Respond(c, item, 200)
  } else {
    app.Respond(c, nil, 404)
  }
}
```

Unless this function was being used to handle more than one route, we would expect the **id** variable to always be present; we include a check in the example code purely for completeness.

## A Complete Example

An example app that uses this package to implement a simple todo list whose contents are persisted into a JSON file can be found at [https://github.com/steatopygous/uiserver-example](https://github.com/steatopygous/uiserver-example).

The UI for the app is written in [Svelte](https://svelte.dev).

The API portion follows a style similar to that proposed by Mat Ryer in various talks like [this one](https://www.youtube.com/watch?v=rWBSMsLG8po), where each handler is wrapped in another function that forms a closure over any context it needs to perform its job, like the actual retrieval of the ToDo that's just a comment in the code above.

