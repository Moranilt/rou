package rou

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDynamicRoutes(t *testing.T) {
	t.Run("route with multiple params and equal length of chunks", func(t *testing.T) {
		router := NewRouter()

		expectedParams := map[string]string{"id": "10", "name": "melony"}
		spyHandler := func(ctx ContextParams) {
			t.Helper()

			for name, value := range expectedParams {
				if ctx.RouterParams().Get(name) != value {
					t.Errorf("Pram %q is not equal to expected value. Got - %s, want - %s",
						name,
						ctx.RouterParams().Get(name),
						value,
					)
					return
				}
			}

		}

		router.Get("/users/:id/friends/:name", spyHandler)
		newServer := httptest.NewServer(router)
		res, _ := http.Get(newServer.URL + "/users/10/friends/melony?name=Joe&age=48")
		bytesRes, _ := io.ReadAll(res.Body)

		if len(bytesRes) != 0 {
			t.Errorf("Response is not empty. Got - %s", bytesRes)
		}
	})

	t.Run("wrong length of called route", func(t *testing.T) {
		router := NewRouter()

		spyHandler := func(ctx ContextParams) {
			t.Helper()

			io.WriteString(ctx.ResponseWriter(), "Should not run this function")
		}

		router.Get("/users/:id/name/:name", spyHandler)
		newServer := httptest.NewServer(router)
		res, _ := http.Get(newServer.URL + "/users/10/name?name=Joe&age=48")
		bytesRes, _ := io.ReadAll(res.Body)

		if string(bytesRes) != "Should not run this function" && res.StatusCode == http.StatusMethodNotAllowed {
			t.Error("Response is empty")
		}
	})
}
