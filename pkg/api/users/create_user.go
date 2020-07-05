package users

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/iplay88keys/recipe-box/pkg/api"
)

type UserSignupResponse struct {
    EmailExisted    bool        `json:"email_existed"`
    UsernameExisted bool        `json:"username_existed"`
    Errors          []api.Error `json:"errors,omitempty"`
}
type existsByUsername func(username string) (bool, error)
type existsByEmail func(email string) (bool, error)
type insertUser func(username, email, password string) (int64, error)

func Signup(existsByUsername existsByUsername, existsByEmail existsByEmail, insertUser insertUser) api.Endpoint {
    return api.Endpoint{
        Path:   "signup",
        Method: http.MethodPost,
        Handler: func(w http.ResponseWriter, r *http.Request) {
            defer r.Body.Close()
            var user UserSignupRequest
            err := json.NewDecoder(r.Body).Decode(&user)
            if err != nil {
                fmt.Println("Error decoding json body for user signup")
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            validationErrors := user.Validate()

            if len(validationErrors) > 0 {
                resp := &UserSignupResponse{
                    Errors: validationErrors,
                }
                respBytes, err := json.Marshal(resp)
                if err != nil {
                    fmt.Println("Error creating response for validation errors")
                    w.WriteHeader(http.StatusInternalServerError)
                    return
                }

                w.WriteHeader(http.StatusBadRequest)
                api.LogWriteErr(w.Write(respBytes))
                return
            }

            usernameExists, err := existsByUsername(user.Username)
            if err != nil {
                fmt.Println("Error checking if user exists by username")
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            emailExists, err := existsByEmail(user.Email)
            if err != nil {
                fmt.Println("Error checking if user exists by email")
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            resp := &UserSignupResponse{
                EmailExisted:    false,
                UsernameExisted: false,
            }

            if usernameExists || emailExists {
                resp.EmailExisted = emailExists
                resp.UsernameExisted = usernameExists
                respBytes, err := json.Marshal(resp)
                if err != nil {
                    fmt.Println("Error creating response for user previously existed")
                    w.WriteHeader(http.StatusInternalServerError)
                    return
                }

                api.LogWriteErr(w.Write(respBytes))
                return
            }

            _, err = insertUser(user.Username, user.Email, user.Password)
            if err != nil {
                fmt.Println("Failed to create user")
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            respBytes, err := json.Marshal(resp)
            if err != nil {
                fmt.Println("Error creating response for user created")
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            api.LogWriteErr(w.Write(respBytes))
            return
        },
    }
}
