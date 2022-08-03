# Simple HTTP router
 A simple HTTP router built on top of the Golang http library
 
 ## Install
 
 `rou` is a standard Go module which can be installed with:

```sh
go get github.com/Moranilt/rou
```
 

## Usage

```go
package main

import (
  "github.com/Moranilt/rou"
)

const (
  authError = "Authorization header should be provided"
  requestDataError = "X-Request-Data header should be provided"
)

func GET_UserHandler(ctx *rou.Context) {
  userId := ctx.RouterParams().Get("userId") // Extract params from router
  // etc
  io.WriteString(ctx.ResponseWriter(), "Response message")
}

func GET_UsersPostHandler(ctx *rou.Context) {
  userId := ctx.RouterParams().Get("userId") // Extract params from router
  postId := ctx.RouterParams().Get("postId") // Extract params from router
  // etc
  io.WriteString(ctx.ResponseWriter(), "Response message")
}

func POST_UserHandler(ctx *rou.Context) {
  // your actions
}

func AuthMiddleware(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Authorization") == "" {
		io.WriteString(w, authError)
		return false
	}
	return true
}

func RequestDataMiddleware(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("X-Request-Data") == "" {
		io.WriteString(w, requestDataError)
		return false
	}
	return true
}

func main() {
  router := rou.NewRouter()
  router.Middleware(AuthMiddleware)
  router.Get("/users/:userId", GET_UserHandler).Middleware(RequestDataMiddleware)
  router.Get("/users/:userId/posts/:postId", GET_UsersPostHandler)
  router.Post("/users/create", POST_UserHandler)
  router.Put(...)
  router.Patch(...)
  router.Delete(...)
  router.Head(...)
  router.Options(...)

  log.Fatal(router.RunServer(":8080")) // Runs server on http://localhost:8080
}

```
