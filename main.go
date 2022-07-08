package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
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
func (r Routes) Exists(route string) bool {
	return r.existingRoutesName[route]
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
	}
}

func (sr *SimpleRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := sr.createContext(w, r)

	routesByMethod := sr.GetRoutes(r.Method)

	for _, route := range routesByMethod {

		if r.URL.Path == route.Path {
			route.Handler(ctx)
			return
		}
	}

	if sr.Routes.Exists(r.URL.Path) {
		ctx.errorJSONResponse(http.StatusMethodNotAllowed, MessageMethodNotAllowed)
		return
	}
	ctx.errorJSONResponse(http.StatusNotFound, MessagePageNotFound)
}

type Context struct {
	responseWriter http.ResponseWriter
	request        *http.Request
}

type ContextParams interface {
	ResponseWriter() http.ResponseWriter
	Request() *http.Request
	Params() url.Values
}

type ErrorObject struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type ResponseObject[T any] struct {
	Error *ErrorObject `json:"error"`
	Body  T            `json:"body"`
}

// Returns basic HTTP response writer
func (c Context) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

// Returns basic HTTP request object
func (c Context) Request() *http.Request {
	return c.request
}

// Returns query params of request
func (c Context) Params() url.Values {
	return c.request.URL.Query()
}

func (c Context) errorJSONResponse(status int, message string) {
	c.ResponseWriter().Header().Add("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(status)
	responseMsg := ResponseObject[any]{Error: &ErrorObject{Message: message, Code: status}}
	jsonContent, _ := json.Marshal(responseMsg)
	io.WriteString(c.ResponseWriter(), string(jsonContent))
}

func (c Context) successJSONResponse(body any) {
	c.ResponseWriter().Header().Add("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	responseMsg := ResponseObject[any]{Body: body}
	jsonContent, _ := json.Marshal(responseMsg)
	io.WriteString(c.ResponseWriter(), string(jsonContent))
}
