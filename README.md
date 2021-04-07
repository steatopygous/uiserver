# uiserver
A simple Go HTTP server that combines UI content and a REST API (based on gorilla/mux) for the content to consume.

The intended purpose of this package is to allow creation of Go desktop apps that have a browser-based user
interface, with the app providing a REST API to support the UI.

The idea is to build something similar to an Electron app, but with the underlying business logic written in Go
and the UI being displayed in a browser tab, rather than a full-blown windowed application.

An example app that uses this package to implement a simple
todo list app can be found at
[http://github/steatopygous/uiserver-example](https://github.com/steatopygous/uiserver-example).
