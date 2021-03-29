package uiserver

import "strings"

type pathMethodHandler struct {
	route   string
	path    string
	method  string
	handler Handler
}

func(pmh pathMethodHandler) matchesPathAndMethod(path string, method string) bool {
	pathParts := strings.Split(path, "/")
	parts := strings.Split(pmh.path, "/")

	if len(parts) == len(pathParts) {
		for i := 0; i < len(parts); i++ {
			if parts[i] != "{}" && parts[i] != pathParts[i] {
				break
			}
		}

		if pmh.method == method {
			return true
		}
	}

	return false
}

