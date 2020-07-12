package users_test

import (
    "bytes"
    "errors"
    "net/http"
    "net/http/httptest"

    "github.com/iplay88keys/recipe-box/pkg/api/users"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("createUser", func() {
    It("creates a user", func() {
        existsByUsername := func(username string) (bool, error) {
            return false, nil
        }

        existsByEmail := func(email string) (bool, error) {
            return false, nil
        }

        insertUser := func(username, email, password string) (int64, error) {
            return 1, nil
        }

        body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

        req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Register(existsByUsername, existsByEmail, insertUser).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusOK))
    })

    It("returns validation info", func() {
        existsByUsername := func(username string) (bool, error) {
            return false, nil
        }

        existsByEmail := func(email string) (bool, error) {
            return false, nil
        }

        insertUser := func(username, email, password string) (int64, error) {
            return 1, nil
        }

        body := []byte(`{
            "username": "",
            "email":    "",
            "password": ""
        }`)

        req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Register(existsByUsername, existsByEmail, insertUser).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusBadRequest))
        Expect(rr.Body.String()).To(MatchJSON(`{
            "errors": {
                "email": "Required",
                "password": "Required",
                "username": "Required"
            }
        }`))
    })

    It("returns info if the username already exists", func() {
        existsByUsername := func(username string) (bool, error) {
            return true, nil
        }

        existsByEmail := func(email string) (bool, error) {
            return false, nil
        }

        insertUser := func(username, email, password string) (int64, error) {
            return 1, nil
        }

        body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

        req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Register(existsByUsername, existsByEmail, insertUser).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusBadRequest))
        Expect(rr.Body.String()).To(MatchJSON(`{
            "errors": {
                "username": "Username already in use"
            }
        }`))
    })

    It("returns info if the email already exists", func() {
        existsByUsername := func(username string) (bool, error) {
            return false, nil
        }

        existsByEmail := func(email string) (bool, error) {
            return true, nil
        }

        insertUser := func(username, email, password string) (int64, error) {
            return 1, nil
        }

        body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

        req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Register(existsByUsername, existsByEmail, insertUser).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusBadRequest))
        Expect(rr.Body.String()).To(MatchJSON(`{
            "errors": {
                "email": "Email already in use"
            }
        }`))
    })

    It("returns bad request if there is no body", func() {
        existsByUsername := func(username string) (bool, error) {
            return false, nil
        }

        existsByEmail := func(email string) (bool, error) {
            return false, nil
        }

        insertUser := func(username, email, password string) (int64, error) {
            return -1, errors.New("some error")
        }

        req := httptest.NewRequest("POST", "/users/register", nil)
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Register(existsByUsername, existsByEmail, insertUser).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusBadRequest))
    })

    It("returns an error if the username check fails", func() {
        existsByUsername := func(username string) (bool, error) {
            return false, errors.New("error username")
        }

        existsByEmail := func(email string) (bool, error) {
            return false, nil
        }

        insertUser := func(username, email, password string) (int64, error) {
            return 1, nil
        }

        body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

        req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Register(existsByUsername, existsByEmail, insertUser).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusInternalServerError))
    })

    It("returns an error if the email check fails", func() {
        existsByUsername := func(username string) (bool, error) {
            return false, nil
        }

        existsByEmail := func(email string) (bool, error) {
            return false, errors.New("error email")
        }

        insertUser := func(username, email, password string) (int64, error) {
            return 1, nil
        }

        body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

        req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Register(existsByUsername, existsByEmail, insertUser).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusInternalServerError))
    })

    It("returns an error if the user insert fails", func() {
        existsByUsername := func(username string) (bool, error) {
            return false, errors.New("error username")
        }

        existsByEmail := func(email string) (bool, error) {
            return false, nil
        }

        insertUser := func(username, email, password string) (int64, error) {
            return -1, errors.New("some error")
        }

        body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

        req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Register(existsByUsername, existsByEmail, insertUser).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusInternalServerError))
    })
})
