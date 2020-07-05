package users_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "github.com/iplay88keys/recipe-box/pkg/api"
    "github.com/iplay88keys/recipe-box/pkg/api/users"
)

var _ = Describe("Validate", func() {
    It("verifies all fields exist", func() {
        req := users.UserSignupRequest{
            Username: "",
            Email:    "",
            Password: "",
        }

        errors := req.Validate()
        Expect(errors).To(ContainElement(api.Error{
            ErrorType: "basic",
            Errors:    []string{"username, password, and email are all required"},
        }))
    })

    It("verifies the email address", func() {
        req := users.UserSignupRequest{
            Username: "test",
            Email:    "bad-email",
            Password: "Pa3$",
        }

        errors := req.Validate()
        Expect(errors).To(ContainElement(api.Error{
            ErrorType: "email",
            Errors:    []string{"invalid email address"},
        }))
    })

    It("verifies the password", func() {
        req := users.UserSignupRequest{
            Username: "test",
            Email:    "email@example.com",
            Password: " ",
        }

        errors := req.Validate()
        Expect(errors).To(Equal([]api.Error{{
            ErrorType: "password",
            Errors: []string{
                "invalid character: ' '",
                "lowercase letter missing",
                "uppercase letter missing",
                "numeric character missing",
                "special character missing",
                "must be between 8 to 64 characters long",
            },
        }}))
    })
})
