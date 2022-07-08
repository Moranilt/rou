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

type Storage interface {
	Delete(name string)
	Has(name string) bool
	Set(name, value string)
	Get(name string) string
}

type Route struct {
	Path    string
	Handler func(ContextParams)
}

type ExistingRoute struct {
	Method string
	Path   string
}

type Routes struct {
	existingRoutesWithMethod map[ExistingRoute]bool
	existingRoutesName       map[string]bool
	routes                   map[string][]Route
}

// Stores route to if it is not exists
func (r *Routes) storeRoute(method string, route string, handler func(ContextParams)) {
	newRoute := ExistingRoute{Method: method, Path: route}
	if !r.existingRoutesWithMethod[newRoute] {
		r.existingRoutesWithMethod[newRoute] = true
		r.existingRoutesName[route] = true
		r.routes[method] = append(r.routes[method], Route{Path: route, Handler: handler})
	}
}

// Check for route exists in Routes with given method and path
func (r Routes) Exists(requestPath string) bool {
	for route := range r.existingRoutesName {
		_, equal := isEqualPaths(route, requestPath)
		if equal {
			return true
		}
	}
	return false
}

func (r Routes) GetRoutes(method string) []Route {
	return r.routes[method]
}

// Initial struct to create HTTP server provide this structure to http.ListenAndServe function
// It has a list of routes  which is stored to serve
type SimpleRouter struct {
	Routes *Routes
}

// Create new SimpleRouter
func NewRouter() *SimpleRouter {
	routes := Routes{
		existingRoutesName:       make(map[string]bool),
		existingRoutesWithMethod: make(map[ExistingRoute]bool),
		routes:                   make(map[string][]Route),
	}
	return &SimpleRouter{Routes: &routes}
}

func (sr SimpleRouter) GetRoutes(method string) []Route {
	return sr.Routes.GetRoutes(method)
}

func (sr *SimpleRouter) storeRoute(method string, route string, handler func(ContextParams)) {
	sr.Routes.storeRoute(method, route, handler)
}

// Add route by method GET
func (sr SimpleRouter) Get(route string, handler func(ContextParams)) {
	sr.storeRoute(MethodGet, route, handler)
}

// Add route by method POST
func (sr SimpleRouter) Post(route string, handler func(ContextParams)) {
	sr.storeRoute(MethodPost, route, handler)
}

// Add route by method PUT
func (sr SimpleRouter) Put(route string, handler func(ContextParams)) {
	sr.storeRoute(MethodPut, route, handler)
}

// Add route by method PATCH
func (sr SimpleRouter) Patch(route string, handler func(ContextParams)) {
	sr.storeRoute(MethodPatch, route, handler)
}

// Add route by method DELETE
func (sr SimpleRouter) Delete(route string, handler func(ContextParams)) {
	sr.storeRoute(MethodDelete, route, handler)
}

// Add route by method OPTIONS
func (sr SimpleRouter) Options(route string, handler func(ContextParams)) {
	sr.storeRoute(MethodOptions, route, handler)
}

// Add route by method HEAD
func (sr SimpleRouter) Head(route string, handler func(ContextParams)) {
	sr.storeRoute(MethodHead, route, handler)
}

func (sr SimpleRouter) createContext(w http.ResponseWriter, r *http.Request) *Context {
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

func (sr *SimpleRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		ctx.errorJSONResponse(http.StatusMethodNotAllowed, MessageMethodNotAllowed)
		return
	}
	ctx.errorJSONResponse(http.StatusNotFound, MessagePageNotFound)
}
