package auth

import (
    "net/http"

    "golang.org/x/net/context"

    "github.com/iplay88keys/my-recipe-library/pkg/token"
)

type ContextKey string

const ContextUserKey ContextKey = "user"

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . TokenService
type TokenService interface {
    ValidateToken(r *http.Request) (*token.AccessDetails, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RedisRepo
type RedisRepo interface {
    RetrieveTokenDetails(details *token.AccessDetails) (int64, error)
}

type Middleware struct {
    tokenService TokenService
    redis        RedisRepo
}

func NewMiddleware(tokenService TokenService, redis RedisRepo) *Middleware {
    return &Middleware{
        tokenService: tokenService,
        redis:        redis,
    }
}

func (a *Middleware) Handler(next http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        details, err := a.tokenService.ValidateToken(r)
        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        userID, err := a.redis.RetrieveTokenDetails(details)
        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), ContextUserKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}
