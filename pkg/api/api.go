package api

import (
    "fmt"
    "net"
    "net/http"
    "time"

    "github.com/gorilla/mux"

    "github.com/iplay88keys/recipe-box/pkg/api/auth"
)

type Config struct {
    Port           string
    StaticDir      string
    AuthMiddleware *auth.Middleware
    Endpoints      []Endpoint
}

type API struct {
    Config *Config
    Server *http.Server
}

func New(config *Config) *API {
    r := mux.NewRouter()

    api := r.PathPrefix("/api/v1").Subrouter()

    for _, endpoint := range config.Endpoints {
        handler := endpoint.Handler
        if endpoint.Auth {
            handler = config.AuthMiddleware.Handler(endpoint.Handler)
        }

        api.HandleFunc(fmt.Sprintf("/%s", endpoint.Path), handler).Methods(endpoint.Method)
    }
    api.NotFoundHandler = http.HandlerFunc(notFound)

    spa := spaHandler{
        staticPath: config.StaticDir,
        indexPath:  "index.html",
    }
    r.PathPrefix("/").Handler(spa)

    return &API{
        Config: config,
        Server: &http.Server{
            Addr:         net.JoinHostPort("", config.Port),
            Handler:      r,
            ReadTimeout:  15 * time.Second,
            WriteTimeout: 15 * time.Second,
        },
    }
}

func (a *API) Start() (shutdown func()) {
    go a.Server.ListenAndServe()

    return func() {
        a.Server.Close()
    }
}

func notFound(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    LogWriteErr(w.Write([]byte("page not found")))
}

func LogWriteErr(_ int, err error) {
    if err != nil {
        fmt.Printf("failed to write response: %s", err.Error())
    }
}
