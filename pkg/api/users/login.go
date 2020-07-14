package users

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/iplay88keys/recipe-box/pkg/api"
)

type UserLoginRequest struct {
    LoginName string `json:"login_name"`
    Password  string `json:"password"`
}

type UserLoginResponse struct {
    Token  string            `json:"token,omitempty"`
    Errors map[string]string `json:"errors,omitempty"`
}

type verify func(loginName, password string) (bool, int64, error)
type createToken func(userid int64) (string, error)

func Login(verify verify, createToken createToken) api.Endpoint {
    return api.Endpoint{
        Path:   "users/login",
        Method: http.MethodPost,
        Handler: func(w http.ResponseWriter, r *http.Request) {
            defer r.Body.Close()
            var user UserLoginRequest
            err := json.NewDecoder(r.Body).Decode(&user)
            if err != nil {
                fmt.Println("Error decoding json body for login")
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            valid, userID, err := verify(user.LoginName, user.Password)
            if err != nil {
                fmt.Println("Error logging user in")
                fmt.Println(err)
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            if !valid {
                errors := make(map[string]string)
                errors["alert"] = "Invalid login credentials"

                resp := &UserLoginResponse{
                    Errors: errors,
                }

                respBytes, err := json.Marshal(resp)
                if err != nil {
                    fmt.Println("Error creating response for login errors")
                    w.WriteHeader(http.StatusInternalServerError)
                    return
                }

                w.WriteHeader(http.StatusUnauthorized)
                api.LogWriteErr(w.Write(respBytes))
                return
            }

            token, err := createToken(userID)
            if err != nil {
                fmt.Printf("Error creating token for user login: %s\n", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            resp := &UserLoginResponse{
                Token: token,
            }

            respBytes, err := json.Marshal(resp)
            if err != nil {
                fmt.Printf("Error creating response for login response with jwt: %s\n", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            w.WriteHeader(http.StatusOK)
            api.LogWriteErr(w.Write(respBytes))
            return
        },
    }
}
