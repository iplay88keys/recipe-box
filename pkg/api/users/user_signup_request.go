package users

import (
    "fmt"
    "net/mail"
    "strconv"
    "unicode"

    "github.com/iplay88keys/recipe-box/pkg/api"
)

type UserSignupRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (u *UserSignupRequest) Validate() []api.Error {
    var errors []api.Error

    if u.Email == "" || u.Password == "" || u.Username == "" {
        errs := []string{"username, password, and email are all required"}
        errors = append(errors, api.Error{ErrorType: "basic", Errors: errs})
    }

    parser := mail.AddressParser{}
    _, err := parser.Parse(u.Email)
    if err != nil {
        errs := []string{"invalid email address"}
        errors = append(errors, api.Error{ErrorType: "email", Errors: errs})
    }

    passwordErrors := u.validatePassword()
    if len(passwordErrors) > 0 {
        errors = append(errors, api.Error{ErrorType: "password", Errors: passwordErrors})
    }

    return errors
}

func (u *UserSignupRequest) validatePassword() []string {
    var errors []string

    var uppercasePresent, lowercasePresent, numberPresent, specialCharPresent bool
    const minLength = 8
    const maxLength = 64
    var passLen int

    for _, ch := range u.Password {
        switch {
        case unicode.IsNumber(ch):
            numberPresent = true
            passLen++
        case unicode.IsUpper(ch):
            uppercasePresent = true
            passLen++
        case unicode.IsLower(ch):
            lowercasePresent = true
            passLen++
        case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
            specialCharPresent = true
            passLen++
        default:
            errors = append(errors, fmt.Sprintf("invalid character: %s", strconv.QuoteRune(ch)))
            passLen++
        }
    }

    if !lowercasePresent {
        errors = append(errors, "lowercase letter missing")
    }
    if !uppercasePresent {
        errors = append(errors, "uppercase letter missing")
    }
    if !numberPresent {
        errors = append(errors, "numeric character missing")
    }
    if !specialCharPresent {
        errors = append(errors, "special character missing")
    }
    if passLen < minLength || passLen > maxLength {
        errors = append(errors, fmt.Sprintf("must be between %d to %d characters long", minLength, maxLength))
    }

    return errors
}
