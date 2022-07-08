package rou

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Context struct {
	responseWriter http.ResponseWriter
	request        *http.Request
	routeParams    Storage
}

type Storage interface {
	Delete(name string)
	Has(name string) bool
	Set(name, value string)
	Get(name string) string
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

// Set response content type
func (ctx *Context) SetContentType(contentType string) {
	ctx.responseWriter.Header().Add("Content-Type", contentType)
}

// Returns RouterParams structure with params of route as key value
//
// If you have route "/user/:id/posts/:postId" and request URL path "/user/1/posts/45" RouterParams will have map with keys
// took from route and values which it will take from request URL path `["id": "1", "postId": "45"]`
func (c Context) RouterParams() Storage {
	return c.routeParams
}

func (c Context) ErrorJSONResponse(status int, message string) {
	c.ResponseWriter().Header().Add("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(status)
	responseMsg := ResponseObject[any]{Error: &ErrorObject{Message: message, Code: status}}
	jsonContent, _ := json.Marshal(responseMsg)
	io.WriteString(c.ResponseWriter(), string(jsonContent))
}

func (c Context) SuccessJSONResponse(body any) {
	c.ResponseWriter().Header().Add("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)
	responseMsg := ResponseObject[any]{Body: body}
	jsonContent, _ := json.Marshal(responseMsg)
	io.WriteString(c.ResponseWriter(), string(jsonContent))
}
