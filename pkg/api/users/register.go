package users

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/iplay88keys/my-recipe-library/pkg/api"
)

type UserSignupResponse struct {
    Errors map[string]string `json:"errors,omitempty"`
}

type existsByUsername func(username string) (bool, error)
type existsByEmail func(email string) (bool, error)
type insertUser func(username, email, password string) (int64, error)

func Register(existsByUsername existsByUsername, existsByEmail existsByEmail, insertUser insertUser) api.Endpoint {
    return api.Endpoint{
        Path:   "users/register",
        Method: http.MethodPost,
        Handler: func(w http.ResponseWriter, r *http.Request) {
            defer r.Body.Close()
            var user UserSignupRequest
            err := json.NewDecoder(r.Body).Decode(&user)
            if err != nil {
                fmt.Println("Error decoding json body for registration")
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            usernameExists, err := existsByUsername(user.Username)
            if err != nil {
                fmt.Println("Error checking if user exists by username for registration")
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            emailExists, err := existsByEmail(user.Email)
            if err != nil {
                fmt.Println("Error checking if user exists by email for registration")
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            validationErrors := user.Validate(usernameExists, emailExists)
            if len(validationErrors) > 0 {
                resp := &UserSignupResponse{
                    Errors: validationErrors,
                }
                respBytes, err := json.Marshal(resp)
                if err != nil {
                    fmt.Println("Error creating response for registration validation errors")
                    w.WriteHeader(http.StatusInternalServerError)
                    return
                }

                w.WriteHeader(http.StatusBadRequest)
                api.LogWriteErr(w.Write(respBytes))
                return
            }

            _, err = insertUser(user.Username, user.Email, user.Password)
            if err != nil {
                fmt.Println("Failed to register user")
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            w.WriteHeader(http.StatusOK)
            return
        },
    }
}
