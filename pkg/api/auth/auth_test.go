package auth_test

import (
    "errors"
    "net/http"
    "net/http/httptest"

    . "github.com/onsi/gomega"

    "github.com/iplay88keys/recipe-box/pkg/api/auth"
    "github.com/iplay88keys/recipe-box/pkg/api/auth/authfakes"
    "github.com/iplay88keys/recipe-box/pkg/token"

    . "github.com/onsi/ginkgo"
)

var _ = Describe("Auth", func() {
    var (
        handle http.HandlerFunc
        req    *http.Request
    )

    BeforeEach(func() {
        var err error
        handle = func(w http.ResponseWriter, r *http.Request) {
            w.Write([]byte(""))
        }

        req, err = http.NewRequest("GET", "/", nil)
        Expect(err).ToNot(HaveOccurred())
    })

    It("calls the next handler if the token is valid", func() {
        tokenService := &authfakes.FakeTokenService{}
        redisRepo := &authfakes.FakeRedisRepo{}

        tokenService.ValidateTokenReturns(&token.AccessDetails{
            AccessUuid: "some-uuid",
            UserId:     10,
        }, nil)
        redisRepo.RetrieveTokenDetailsReturns(10, nil)

        rr := httptest.NewRecorder()

        validator := auth.NewMiddleware(tokenService, redisRepo)
        validator.Handler(http.HandlerFunc(handle)).ServeHTTP(rr, req)

        Expect(rr.Code).To(Equal(http.StatusOK))
    })

    It("returns unauthorized if the token is invalid", func() {
        tokenService := &authfakes.FakeTokenService{}
        redisRepo := &authfakes.FakeRedisRepo{}

        tokenService.ValidateTokenReturns(nil, errors.New("invalid token"))

        rr := httptest.NewRecorder()

        validator := auth.NewMiddleware(tokenService, redisRepo)
        validator.Handler(http.HandlerFunc(handle)).ServeHTTP(rr, req)

        Expect(rr.Code).To(Equal(http.StatusUnauthorized))
    })

    It("returns unauthorized if redis returns an error", func() {
        tokenService := &authfakes.FakeTokenService{}
        redisRepo := &authfakes.FakeRedisRepo{}

        tokenService.ValidateTokenReturns(&token.AccessDetails{
            AccessUuid: "some-uuid",
            UserId:     10,
        }, nil)
        redisRepo.RetrieveTokenDetailsReturns(-1, errors.New("some redis error"))

        rr := httptest.NewRecorder()

        validator := auth.NewMiddleware(tokenService, redisRepo)
        validator.Handler(http.HandlerFunc(handle)).ServeHTTP(rr, req)

        Expect(rr.Code).To(Equal(http.StatusUnauthorized))
    })
})
