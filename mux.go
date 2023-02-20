package mux

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

type (
	// Mux is a simple HTTP request multiplexer
	Mux struct {
		routes          map[string][]*route
		middlewares     []func(http.HandlerFunc) http.HandlerFunc
		notFoundHandler http.HandlerFunc
	}

	route struct {
		path    string
		handler http.HandlerFunc
	}
)

// New creates a new instance of Mux
func New() *Mux {
	return &Mux{routes: make(map[string][]*route)}
}

// Connect adds a new route for the CONNECT HTTP method
func (m *Mux) Connect(path string, handler http.HandlerFunc) {
	m.append("CONNECT", path, handler)
}

// Delete adds a new route for the DELETE HTTP method
func (m *Mux) Delete(path string, handler http.HandlerFunc) {
	m.append("DELETE", path, handler)
}

// Get adds a new route for the GET HTTP method
func (m *Mux) Get(path string, handler http.HandlerFunc) {
	m.append("GET", path, handler)
}

// Head adds a new route for the HEAD HTTP method
func (m *Mux) Head(path string, handler http.HandlerFunc) {
	m.append("HEAD", path, handler)
}

// Options adds a new route for the OPTIONS HTTP method
func (m *Mux) Options(path string, handler http.HandlerFunc) {
	m.append("OPTIONS", path, handler)
}

// Patch adds a new route for the PATCH HTTP method
func (m *Mux) Patch(path string, handler http.HandlerFunc) {
	m.append("PATCH", path, handler)
}

// Post adds a new route for the POST HTTP method
func (m *Mux) Post(path string, handler http.HandlerFunc) {
	m.append("POST", path, handler)
}

// Put adds a new route for the PUT HTTP method
func (m *Mux) Put(path string, handler http.HandlerFunc) {
	m.append("PUT", path, handler)
}

// Trace adds a new route for the TRACE HTTP method
func (m *Mux) Trace(path string, handler http.HandlerFunc) {
	m.append("TRACE", path, handler)
}

// Use adds a new middleware function to the Mux
func (m *Mux) Use(middleware func(http.Handler) http.Handler) {
	m.middlewares = append(m.middlewares, func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			middleware(handler).ServeHTTP(w, r)
		}
	})
}

// NotFound sets the handler to be called when no matching route is found
func (m *Mux) NotFound(handler http.HandlerFunc) {
	m.notFoundHandler = handler
}

// NotFoundHandler returns the handler to be called when no matching route is found
func (m *Mux) NotFoundHandler() http.HandlerFunc {
	if m.notFoundHandler != nil {
		return m.notFoundHandler
	}
	return http.NotFound
}

// ServeHTTP implements the http.Handler interface and handles incoming requests
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.routes == nil {
		m.NotFoundHandler().ServeHTTP(w, r)
		return
	}

	routes := m.routes[r.Method]
	segments := strings.Split(r.URL.Path, "/")

	for _, route := range routes {
		routeSegments := strings.Split(route.path, "/")
		if len(routeSegments) != len(segments) {
			continue
		}

		match := true
		params := make(map[string]string)

		for i, routeSegment := range routeSegments {
			if strings.HasPrefix(routeSegment, ":") {
				params[routeSegment[1:]] = segments[i]
			} else if routeSegment != segments[i] {
				match = false
				break
			}
		}

		if match {
			ctx := context.WithValue(r.Context(), "params", params)
			r = r.WithContext(ctx)

			handler := route.handler
			for i := len(m.middlewares) - 1; i >= 0; i-- {
				handler = m.middlewares[i](handler)
			}

			handler(w, r)
			return
		}
	}

	m.NotFoundHandler().ServeHTTP(w, r)
}

func (m *Mux) append(method string, path string, handler http.HandlerFunc) {
	routes := m.routes[method]

	for _, route := range routes {
		if route.path == path {
			return
		}
	}

	m.routes[method] = append(routes, &route{path, handler})
}

// Param retrieves the value of a named path parameter from the request context
func Param(r *http.Request, param string) string {
	return r.Context().Value("params").(map[string]string)[param]
}

// ParamInt retrieves the value of a named path parameter from the request context as an integer
// If the value cannot be converted to an integer, 0 is returned
func ParamInt(r *http.Request, param string) int {
	v, err := strconv.Atoi(Param(r, param))
	if err != nil {
		return 0
	}
	return v
}
