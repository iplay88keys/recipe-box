package users

import (
    "net/http"

    "github.com/iplay88keys/my-recipe-library/pkg/token"

    "github.com/iplay88keys/my-recipe-library/pkg/api"
)

type validateToken func(r *http.Request) (*token.AccessDetails, error)
type deleteTokenDetails func(uuid string) (int64, error)

func Logout(validateToken validateToken, deleteTokenDetails deleteTokenDetails) api.Endpoint {
    return api.Endpoint{
        Path:   "users/logout",
        Method: http.MethodPost,
        Auth:   true,
        Handler: func(w http.ResponseWriter, r *http.Request) {
            details, err := validateToken(r)
            if err != nil {
                w.WriteHeader(http.StatusUnauthorized)
                return
            }

            userID, err := deleteTokenDetails(details.AccessUuid)
            if err != nil || userID == 0 {
                w.WriteHeader(http.StatusUnauthorized)
                return
            }

            api.LogWriteErr(w.Write([]byte("")))
            return
        },
    }
}
