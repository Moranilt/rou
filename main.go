package rou

import (
	"net/http"
	"strings"
)

const (
	MethodGet     = "GET"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodDelete  = "DELETE"
	MethodOptions = "OPTIONS"
	MethodHead    = "HEAD"
)

const (
	MessageBodyIsNotValid   = "Request body is not valid"
	MessageMethodNotAllowed = "Method not allowed"
	MessagePageNotFound     = "Page not found"
)

type routerParams struct {
	value map[string]string
}

func (r *routerParams) Delete(name string) {
	delete(r.value, name)
}

func (r routerParams) Has(name string) bool {
	_, has := r.value[name]
	return has
}

func (r *routerParams) Set(name, value string) {
	r.value[name] = value
}

func (r routerParams) Get(name string) string {
	return r.value[name]
}

type Route struct {
	Path    string
	Handler func(ContextParams)
}

type existingRoute struct {
	Method string
	Path   string
}

type routes struct {
	existingRoutesWithMethod map[existingRoute]bool
	routes                   map[string][]Route
}

// Stores route to if it is not exists
func (r *routes) storeRoute(method string, route string, handler func(ContextParams)) {
	newRoute := existingRoute{Method: method, Path: route}
	if !r.existingRoutesWithMethod[newRoute] {
		r.existingRoutesWithMethod[newRoute] = true
		r.routes[method] = append(r.routes[method], Route{Path: route, Handler: handler})
	}
}

// Check for route exists in Routes with given method and path
func (r routes) Exists(requestPath string) bool {
	for route := range r.existingRoutesWithMethod {
		_, equal := isEqualPaths(route.Path, requestPath)
		if equal {
			return true
		}
	}
	return false
}

func (r routes) GetRoutes(method string) []Route {
	return r.routes[method]
}

// Initial struct to create HTTP server provide this structure to http.ListenAndServe function
// It has a list of routes  which is stored to serve
type simpleRouter struct {
	Routes *routes
}

type SimpleRouter interface {
	// Get a list of routes with handlers by method
	GetRoutes(method string) []Route
	// Add route by method GET
	Get(route string, handler func(ContextParams))
	// Add route by method POST
	Post(route string, handler func(ContextParams))
	// Add route by method PUT
	Put(route string, handler func(ContextParams))
	// Add route by method PATCH
	Patch(route string, handler func(ContextParams))
	// Add route by method DELETE
	Delete(route string, handler func(ContextParams))
	// Add route by method OPTIONS
	Options(route string, handler func(ContextParams))
	// Add route by method HEAD
	Head(route string, handler func(ContextParams))
	// Implements an http.Handler interface to use it like server handler in http.ListenAndServe
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// Create as new SimpleRouter instance
func NewRouter() SimpleRouter {
	routes := routes{
		existingRoutesWithMethod: make(map[existingRoute]bool),
		routes:                   make(map[string][]Route),
	}
	return &simpleRouter{Routes: &routes}
}

func (sr simpleRouter) GetRoutes(method string) []Route {
	return sr.Routes.GetRoutes(method)
}

func (sr *simpleRouter) storeRoute(method string, route string, handler func(ContextParams)) {
	sr.Routes.storeRoute(method, route, handler)
}

// Add route by method GET
func (sr simpleRouter) Get(route string, handler func(ContextParams)) {
	sr.storeRoute(MethodGet, route, handler)
}

// Add route by method POST
func (sr simpleRouter) Post(route string, handler func(ContextParams)) {
	sr.storeRoute(MethodPost, route, handler)
}

// Add route by method PUT
func (sr simpleRouter) Put(route string, handler func(ContextParams)) {
	sr.storeRoute(MethodPut, route, handler)
}

// Add route by method PATCH
func (sr simpleRouter) Patch(route string, handler func(ContextParams)) {
	sr.storeRoute(MethodPatch, route, handler)
}

// Add route by method DELETE
func (sr simpleRouter) Delete(route string, handler func(ContextParams)) {
	sr.storeRoute(MethodDelete, route, handler)
}

// Add route by method OPTIONS
func (sr simpleRouter) Options(route string, handler func(ContextParams)) {
	sr.storeRoute(MethodOptions, route, handler)
}

// Add route by method HEAD
func (sr simpleRouter) Head(route string, handler func(ContextParams)) {
	sr.storeRoute(MethodHead, route, handler)
}

func (sr simpleRouter) createContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		responseWriter: w,
		request:        r,
		routeParams:    &routerParams{value: make(map[string]string)},
	}
}

func prepareURLChunks(url string) []string {
	return strings.Split(strings.Trim(url, "/"), "/")
}

func isEqualPaths(route string, requestPath string) (*map[string]string, bool) {
	params := make(map[string]string)

	if route == requestPath {
		return &params, true
	}

	clearRoutePath := prepareURLChunks(route)
	clearRequestPath := prepareURLChunks(requestPath)

	if len(clearRoutePath) != len(clearRequestPath) {
		return nil, false
	}

	for i := 0; i < len(clearRoutePath); i++ {
		routeChunk := clearRoutePath[i]
		if clearRoutePath[i][0] == ':' {
			params[routeChunk[1:]] = clearRequestPath[i]
		} else {
			if clearRequestPath[i] != routeChunk {
				return nil, false
			}
		}
	}
	return &params, true
}

func (sr *simpleRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := sr.createContext(w, r)

	routesByMethod := sr.GetRoutes(r.Method)

ROUTES_BY_METHOD:
	for _, route := range routesByMethod {
		if r.URL.Path == route.Path {
			route.Handler(ctx)
			return
		}

		params, equal := isEqualPaths(route.Path, r.URL.Path)
		if !equal {
			continue ROUTES_BY_METHOD
		}

		for name, value := range *params {
			ctx.RouterParams().Set(name, value)
		}
		route.Handler(ctx)
		return
	}

	if sr.Routes.Exists(r.URL.Path) {
		ctx.ErrorJSONResponse(http.StatusMethodNotAllowed, MessageMethodNotAllowed)
		return
	}
	ctx.ErrorJSONResponse(http.StatusNotFound, MessagePageNotFound)
}
