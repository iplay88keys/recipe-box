package users_test

import (
    "bytes"
    "errors"
    "net/http"
    "net/http/httptest"

    "github.com/iplay88keys/my-recipe-library/pkg/token"

    "github.com/iplay88keys/my-recipe-library/pkg/api/users"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("login", func() {
    It("logs a user in", func() {
        verify := func(loginName, password string) (bool, int64, error) {
            return true, 0, nil
        }

        createToken := func(userid int64) (*token.Details, error) {
            return &token.Details{
                AccessToken:  "access token",
                RefreshToken: "refresh token",
            }, nil
        }

        storeToken := func(userid int64, details *token.Details) error {
            return nil
        }

        body := []byte(`{
            "login_name": "username",
            "password": "Pa3$12345"
        }`)

        req := httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Login(verify, createToken, storeToken).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusOK))
        Expect(rr.Body.String()).To(MatchJSON(`{
            "access_token": "access token",
            "refresh_token": "refresh token"
        }`))
    })

    It("returns unauthorized if the login fails due to bad credentials", func() {
        verify := func(loginName, password string) (bool, int64, error) {
            return false, 0, nil
        }

        createToken := func(userid int64) (*token.Details, error) {
            return &token.Details{

            }, nil
        }

        storeToken := func(userid int64, details *token.Details) error {
            return nil
        }

        body := []byte(`{
            "login_name": "username",
            "password": "bad-password"
        }`)

        req := httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Login(verify, createToken, storeToken).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusUnauthorized))
        Expect(rr.Body.String()).To(MatchJSON(`{
            "errors": {
                "alert": "Invalid login credentials"
            }
        }`))
    })

    It("returns a bad request if there is no body", func() {
        verify := func(loginName, password string) (bool, int64, error) {
            return false, 0, nil
        }

        createToken := func(userid int64) (*token.Details, error) {
            return &token.Details{}, nil
        }

        storeToken := func(userid int64, details *token.Details) error {
            return nil
        }

        req := httptest.NewRequest("POST", "/users/login", nil)
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Login(verify, createToken, storeToken).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusBadRequest))
    })

    It("returns an error if the login check fails", func() {
        verify := func(loginName, password string) (bool, int64, error) {
            return false, 0, errors.New("some error")
        }

        createToken := func(userid int64) (*token.Details, error) {
            return &token.Details{}, nil
        }

        storeToken := func(userid int64, details *token.Details) error {
            return nil
        }

        body := []byte(`{
            "login_name": "username",
            "password": "bad-password"
        }`)

        req := httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Login(verify, createToken, storeToken).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusInternalServerError))
    })

    It("returns an error if the token creation fails", func() {
        verify := func(loginName, password string) (bool, int64, error) {
            return true, 0, nil
        }

        createToken := func(userid int64) (*token.Details, error) {
            return nil, errors.New("some error")
        }

        storeToken := func(userid int64, details *token.Details) error {
            return nil
        }

        body := []byte(`{
            "login_name": "username",
            "password": "bad-password"
        }`)

        req := httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Login(verify, createToken, storeToken).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusInternalServerError))
    })

    It("returns an error if the token storing fails", func() {
        verify := func(loginName, password string) (bool, int64, error) {
            return true, 0, nil
        }

        createToken := func(userid int64) (*token.Details, error) {
            return &token.Details{}, nil
        }

        storeToken := func(userid int64, details *token.Details) error {
            return errors.New("some error")
        }

        body := []byte(`{
            "login_name": "username",
            "password": "bad-password"
        }`)

        req := httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Login(verify, createToken, storeToken).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusInternalServerError))
    })
})
