package api

import "net/http"

type Endpoint struct {
    Path    string
    Method  string
    Auth    bool
    Handler http.HandlerFunc
}
