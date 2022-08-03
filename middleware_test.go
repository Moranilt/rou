package rou

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var authError = "Authorization header should be provided"
var requestDataError = "X-Request-Data header should be provided"
var middlewares = []MiddlewareFunction{
	func(w http.ResponseWriter, r *http.Request) bool {
		if r.Header.Get("Authorization") == "" {
			io.WriteString(w, authError)
			return false
		}
		return true
	},
	func(w http.ResponseWriter, r *http.Request) bool {
		if r.Header.Get("X-Request-Data") == "" {
			io.WriteString(w, requestDataError)
			return false
		}
		return true
	},
}

func TestServerMiddleware(t *testing.T) {
	t.Run("not pass first middleware", func(t *testing.T) {
		router := NewRouter()

		spyHandler := func(ctx *Context) {
			t.Helper()

			io.WriteString(ctx.ResponseWriter(), "Should not run this function")
		}

		router.Use(middlewares...)
		router.Get("/users", spyHandler)

		newServer := httptest.NewServer(router)
		res, _ := http.Get(newServer.URL + "/users")
		bytesRes, _ := io.ReadAll(res.Body)

		if string(bytesRes) != authError {
			t.Errorf("expected %q, got %q", authError, string(bytesRes))
		}
	})

	t.Run("all middlewares successfully completed", func(t *testing.T) {
		router := NewRouter()

		succeedMsg := "Should run this function"
		spyHandler := func(ctx *Context) {
			t.Helper()

			io.WriteString(ctx.ResponseWriter(), succeedMsg)
		}

		router.Use(middlewares...)
		router.Get("/users", spyHandler)

		newServer := httptest.NewServer(router)
		request, _ := http.NewRequest("GET", newServer.URL+"/users", bytes.NewBuffer([]byte{}))
		request.Header.Add("Authorization", "secret key")
		request.Header.Add("X-Request-Data", "request data")
		client := &http.Client{}
		res, _ := client.Do(request)
		bytesRes, _ := io.ReadAll(res.Body)

		if string(bytesRes) != succeedMsg {
			t.Errorf("expected %q, got %q", succeedMsg, string(bytesRes))
		}
	})

}

func TestRouteMiddleware(t *testing.T) {
	t.Run("not pass middleware", func(t *testing.T) {
		router := NewRouter()

		succeedMsg := "Should run this function"
		spyHandler := func(ctx *Context) {
			t.Helper()

			io.WriteString(ctx.ResponseWriter(), succeedMsg)
		}
		router.Get("/users", spyHandler).Middleware(middlewares...)

		newServer := httptest.NewServer(router)
		res, _ := http.Get(newServer.URL + "/users")
		bytesRes, _ := io.ReadAll(res.Body)

		if string(bytesRes) != authError {
			t.Errorf("expected %q, got %q", authError, string(bytesRes))
		}
	})

	t.Run("all middlewares successfully completed", func(t *testing.T) {
		router := NewRouter()

		succeedMsg := "Should run this function"
		spyHandler := func(ctx *Context) {
			t.Helper()

			io.WriteString(ctx.ResponseWriter(), succeedMsg)
		}
		router.Get("/users", spyHandler).Middleware(middlewares...)

		newServer := httptest.NewServer(router)
		request, _ := http.NewRequest("GET", newServer.URL+"/users", bytes.NewBuffer([]byte{}))
		request.Header.Add("Authorization", "secret key")
		request.Header.Add("X-Request-Data", "request data")
		client := &http.Client{}
		res, _ := client.Do(request)
		bytesRes, _ := io.ReadAll(res.Body)

		if string(bytesRes) != succeedMsg {
			t.Errorf("expected %q, got %q", succeedMsg, string(bytesRes))
		}
	})

}
