package api

import "net/http"

type Endpoint struct {
    Path    string
    Method  string
    Handler http.HandlerFunc
}
