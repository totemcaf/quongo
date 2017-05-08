package main

import (
  "net/http"

  "github.com/gorilla/mux"
  "log"
)

type Route struct {
  Name        string
  Method      string
  Pattern     string
  HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(routes *Routes) *mux.Router {

  router := mux.NewRouter().StrictSlash(true)
  for _, route := range *routes {
    var handler http.Handler

    handler = route.HandlerFunc
    handler = Logger(handler, route.Name)

    log.Printf("Adding handler for %s as %s: %s", route.Name, route.Method, route.Pattern)
    router.
    Methods(route.Method).
      Path(route.Pattern).
      Name(route.Name).
      Handler(handler)
  }

  return router
}
