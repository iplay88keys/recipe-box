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

        req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Signup(existsByUsername, existsByEmail, insertUser).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusOK))
        Expect(rr.Body.String()).To(MatchJSON(`{
            "email_existed": false,
            "username_existed": false
        }`))
    })

    It("returns validation info", func() {
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
            "username": "",
            "email":    "",
            "password": ""
        }`)

        req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Signup(existsByUsername, existsByEmail, insertUser).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusBadRequest))
        Expect(rr.Body.String()).To(MatchJSON(`{
            "email_existed": false,
            "username_existed": false,
            "errors": [{
                "error_type": "basic",
                "errors": [
                    "username, password, and email are all required"
                ]
            }, {
                "error_type": "email",
                "errors": [
                    "invalid email address"
                ]
            }, {
                "error_type": "password",
                "errors": [
                    "lowercase letter missing",
                    "uppercase letter missing",
                    "numeric character missing",
                    "special character missing",
                    "must be between 8 to 64 characters long"
                ]
            }]
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

        req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Signup(existsByUsername, existsByEmail, insertUser).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusOK))
        Expect(rr.Body.String()).To(MatchJSON(`{
            "email_existed": false,
            "username_existed": true
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

        req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Signup(existsByUsername, existsByEmail, insertUser).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusOK))
        Expect(rr.Body.String()).To(MatchJSON(`{
            "email_existed": true,
            "username_existed": false
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

        req := httptest.NewRequest("POST", "/signup", nil)
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Signup(existsByUsername, existsByEmail, insertUser).Handler)

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

        req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Signup(existsByUsername, existsByEmail, insertUser).Handler)

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

        req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Signup(existsByUsername, existsByEmail, insertUser).Handler)

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

        req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(users.Signup(existsByUsername, existsByEmail, insertUser).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusInternalServerError))
    })
})
