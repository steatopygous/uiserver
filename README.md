# uiserver
A simple Go HTTP server that combines UI content and a REST API (based on gorilla/mux) for the content to consume.

The intended purpose of this package is to allow creation of Go desktop apps that have a browser-based user
interface, with the app providing a REST API to support the UI.

The idea is to build something similar to an Electron app, but with the underlying business logic written in Go
and the UI being displayed in a browser tab, rather than a full-blown windowed application.

In the near future, an example app that uses this package (along with the new Go 1.16 embed functionality to allow the
UI and business logic to be built into a single executable) will be made available at
[http://github/steatopygous/uiserver-example]().
