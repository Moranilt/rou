package rou

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var mockedRoutes = []Route{
	{
		Path: "/test-1",
	},
	{
		Path: "/test-2",
	},
}

var routesWithMethod = []struct {
	method string
	routes []Route
}{
	{
		method: "GET",
		routes: mockedRoutes,
	},
	{
		method: "POST",
		routes: mockedRoutes,
	},
	{
		method: "PUT",
		routes: mockedRoutes,
	},
	{
		method: "PATCH",
		routes: mockedRoutes,
	},
	{
		method: "DELETE",
		routes: mockedRoutes,
	},
	{
		method: "HEAD",
		routes: mockedRoutes,
	},
	{
		method: "OPTIONS",
		routes: mockedRoutes,
	},
}

func TestStoreRoutes(t *testing.T) {
	fakeHandler := func(ctx ContextParams) {}
	for _, test := range routesWithMethod {
		t.Run(fmt.Sprintf("add routes with method %s", test.method), func(t *testing.T) {
			router := NewRouter()

			for _, route := range test.routes {
				switch test.method {
				case "GET":
					router.Get(route.Path, fakeHandler)
				case "POST":
					router.Post(route.Path, fakeHandler)
				case "PUT":
					router.Put(route.Path, fakeHandler)
				case "PATCH":
					router.Patch(route.Path, fakeHandler)
				case "DELETE":
					router.Delete(route.Path, fakeHandler)
				case "HEAD":
					router.Head(route.Path, fakeHandler)
				case "OPTIONS":
					router.Options(route.Path, fakeHandler)
				}
			}

			got := router.GetRoutes(test.method)
			expected := test.routes

			assertRoutes(t, got, expected)
		})
	}
}

func assertRoutes(t testing.TB, got, want []Route) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("different length of array. Got - %d, want - %d", len(got), len(want))
		return
	}

	for i := 0; i < len(got); i++ {
		if got[i].Path != want[i].Path {
			t.Errorf("got is not equal expected routes. Got - %v, want - %v", got, want)
			return
		}
	}
}

func TestServeHTTP(t *testing.T) {
	t.Run("POST request", func(t *testing.T) {
		router := NewRouter()
		requestBody := []byte("Request body")
		responseBody := "Response body"
		spyHandler := func(ctx ContextParams) {
			t.Helper()
			bodyBytes, _ := io.ReadAll(ctx.Request().Body)
			if !reflect.DeepEqual(bodyBytes, requestBody) {
				t.Errorf("Body is not the same. Got - %v, want - %v", bodyBytes, requestBody)
				return
			}

			io.WriteString(ctx.ResponseWriter(), responseBody)
		}
		router.Post("/test-route", spyHandler)
		newServer := httptest.NewServer(router)
		res, _ := http.Post(newServer.URL+"/test-route", "application/json", bytes.NewBuffer(requestBody))

		resBytes, _ := io.ReadAll(res.Body)
		if !reflect.DeepEqual(resBytes, []byte(responseBody)) {
			t.Errorf("Response is not the same. Got - %s, want %s", responseBody, resBytes)
		}
	})

	t.Run("GET request with params", func(t *testing.T) {
		router := NewRouter()
		responseBody := "Response body"
		expectedParams := map[string]string{"name": "Joe", "age": "48"}
		spyHandler := func(ctx ContextParams) {
			t.Helper()

			for name, value := range expectedParams {
				if ctx.Params().Get(name) != value {
					t.Errorf("Pram %s is not equal to expected value. Got - %s, want - %s", name, ctx.Params().Get(name), value)
					return
				}
			}

			io.WriteString(ctx.ResponseWriter(), responseBody)
		}

		router.Get("/test-route", spyHandler)
		newServer := httptest.NewServer(router)
		res, _ := http.Get(newServer.URL + "/test-route?name=Joe&age=48")

		resBytes, _ := io.ReadAll(res.Body)
		if !reflect.DeepEqual(resBytes, []byte(responseBody)) {
			t.Errorf("Response is not the same. Got - %s, want %s", resBytes, responseBody)
		}
	})

	t.Run("Wrong route name in request", func(t *testing.T) {
		router := NewRouter()
		spyHandler := func(ctx ContextParams) {}

		router.Get("/test-route", spyHandler)
		newServer := httptest.NewServer(router)
		res, _ := http.Get(newServer.URL + "/bad-route")

		expectedResponse := `{"error":{"message":"Page not found","code":404},"body":null}`
		resBytes, _ := io.ReadAll(res.Body)
		if !reflect.DeepEqual(resBytes, []byte(expectedResponse)) {
			t.Errorf("Response is not the same. Got - %s, want %s", resBytes, expectedResponse)
		}
	})

	t.Run("Wrong route method", func(t *testing.T) {
		router := NewRouter()
		spyHandler := func(ctx ContextParams) {}

		router.Post("/test-route", spyHandler)
		newServer := httptest.NewServer(router)
		res, _ := http.Get(newServer.URL + "/test-route")

		expectedResponse := `{"error":{"message":"Method not allowed","code":405},"body":null}`
		resBytes, _ := io.ReadAll(res.Body)
		if !reflect.DeepEqual(resBytes, []byte(expectedResponse)) {
			t.Errorf("Response is not the same. Got - %s, want %s", resBytes, expectedResponse)
		}
	})

	t.Run("Wrong route method with dynamic params", func(t *testing.T) {
		router := NewRouter()
		spyHandler := func(ctx ContextParams) {}

		router.Post("/test-route/:id", spyHandler)
		newServer := httptest.NewServer(router)
		res, _ := http.Get(newServer.URL + "/test-route/10")

		expectedResponse := `{"error":{"message":"Method not allowed","code":405},"body":null}`
		resBytes, _ := io.ReadAll(res.Body)
		if !reflect.DeepEqual(resBytes, []byte(expectedResponse)) {
			t.Errorf("Response is not the same. Got - %s, want %s", resBytes, expectedResponse)
		}
	})

	t.Run("Success response", func(t *testing.T) {
		router := NewRouter()
		responseBody := "Response body"
		spyHandler := func(ctx ContextParams) {
			t.Helper()

			ctx.SuccessJSONResponse(responseBody)
		}

		router.Get("/test-route", spyHandler)
		newServer := httptest.NewServer(router)
		res, _ := http.Get(newServer.URL + "/test-route")

		expectedResponse := `{"error":null,"body":"Response body"}`
		resBytes, _ := io.ReadAll(res.Body)
		if !reflect.DeepEqual(resBytes, []byte(expectedResponse)) {
			t.Errorf("Response is not the same. Got - %s, want %s", expectedResponse, resBytes)
		}
	})
}
