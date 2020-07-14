package token

import (
    "errors"
    "fmt"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
    UserID int64
    jwt.StandardClaims
}

func (u UserClaims) Valid() error {
    return nil
}

func CreateToken(userID int64) (string, error) {
    claims := UserClaims{
        userID,
        jwt.StandardClaims{
            IssuedAt:  time.Now().Unix(),
            ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
        },
    }

    unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    token, err := unsignedToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
    if err != nil {
        return "", err
    }

    return token, nil
}

func ExtractToken(r *http.Request) string {
    auth := strings.Split(r.Header.Get("Authorization"), " ")
    if len(auth) == 2 {
        return auth[1]
    }

    return ""
}

func ParseToken(tokenString string) (*jwt.Token, error) {
    token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("ACCESS_SECRET")), nil
    })
    if err != nil {
        return nil, err
    }

    return token, nil
}

func VerifyToken(token *jwt.Token) (bool, error) {
    if _, ok := token.Claims.(*UserClaims); !ok || !token.Valid {
        return false, errors.New("invalid token")
    }

    return true, nil
}
