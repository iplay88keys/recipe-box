package token_test

import (
    "net/http"
    "os"
    "time"

    "github.com/dgrijalva/jwt-go"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "github.com/iplay88keys/recipe-box/pkg/token"
)

var _ = Describe("token", func() {
    Context("CreateToken", func() {
        It("creates a token", func() {
            err := os.Setenv("ACCESS_SECRET", "secret value")
            Expect(err).ToNot(HaveOccurred())

            token, err := token.CreateToken(0)
            Expect(err).ToNot(HaveOccurred())
            Expect(len(token)).To(BeNumerically(">", 0))
        })
    })

    Context("ExtractToken", func() {
        It("returns the token string from a request", func() {
            req, err := http.NewRequest(http.MethodPost, "example.com", nil)
            Expect(err).ToNot(HaveOccurred())

            req.Header.Set("Authorization", "bearer some-token")

            token := token.ExtractToken(req)
            Expect(err).ToNot(HaveOccurred())
            Expect(token).To(Equal("some-token"))
        })

        It("returns an empty string if the token does not exist", func() {
            req, err := http.NewRequest(http.MethodPost, "example.com", nil)
            Expect(err).ToNot(HaveOccurred())

            req.Header.Set("Authorization", "bearer")

            token := token.ExtractToken(req)
            Expect(err).ToNot(HaveOccurred())
            Expect(token).To(Equal(""))
        })

        It("returns an empty string if the Authorization header does not exist", func() {
            req, err := http.NewRequest(http.MethodPost, "example.com", nil)
            Expect(err).ToNot(HaveOccurred())

            token := token.ExtractToken(req)
            Expect(err).ToNot(HaveOccurred())
            Expect(token).To(Equal(""))
        })
    })

    Context("ParseToken", func() {
        It("parses a JWT token", func() {
            err := os.Setenv("ACCESS_SECRET", "secret value")
            Expect(err).ToNot(HaveOccurred())

            // token from previous run of token.CreateToken(0)
            jwtToken, err := token.ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJleHAiOjE1OTQ2MTY0MzAsImlhdCI6MTU5NDYxMjgzMH0.9Bhbk7G25LG-fU9w_HOxDUw3u16cG0QULwIIUENXdJ4")
            Expect(err).ToNot(HaveOccurred())

            Expect(jwtToken).To(Equal(&jwt.Token{
                Raw:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJleHAiOjE1OTQ2MTY0MzAsImlhdCI6MTU5NDYxMjgzMH0.9Bhbk7G25LG-fU9w_HOxDUw3u16cG0QULwIIUENXdJ4",
                Method: &jwt.SigningMethodHMAC{Name: "HS256", Hash: 5},
                Header: map[string]interface{}{"alg": "HS256", "typ": "JWT"},
                Claims: &token.UserClaims{
                    UserID: 0,
                    StandardClaims: jwt.StandardClaims{
                        ExpiresAt: 1594616430,
                        IssuedAt:  1594612830,
                    },
                },
                Signature: "9Bhbk7G25LG-fU9w_HOxDUw3u16cG0QULwIIUENXdJ4",
                Valid:     true,
            }))
        })
    })

    Context("VerifyToken", func() {
        It("returns true if JWT Token is valid", func() {
            err := os.Setenv("ACCESS_SECRET", "secret value")
            Expect(err).ToNot(HaveOccurred())

            // token from previous run of token.CreateToken(0) passed through token.ParseToken
            valid, err := token.VerifyToken(&jwt.Token{
                Raw:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJleHAiOjE1OTQ2MTY0MzAsImlhdCI6MTU5NDYxMjgzMH0.9Bhbk7G25LG-fU9w_HOxDUw3u16cG0QULwIIUENXdJ4",
                Method: &jwt.SigningMethodHMAC{Name: "HS256", Hash: 5},
                Header: map[string]interface{}{"alg": "HS256", "typ": "JWT"},
                Claims: &token.UserClaims{
                    UserID: 0,
                    StandardClaims: jwt.StandardClaims{
                        IssuedAt:  time.Now().Unix(),
                        ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
                    },
                },
                Signature: "9Bhbk7G25LG-fU9w_HOxDUw3u16cG0QULwIIUENXdJ4",
                Valid:     true,
            })
            Expect(err).ToNot(HaveOccurred())
            Expect(valid).To(BeTrue())
        })

        It("returns false if JWT Token is invalid", func() {
            err := os.Setenv("ACCESS_SECRET", "secret value")
            Expect(err).ToNot(HaveOccurred())

            // token from previous run of token.CreateToken(0) passed through token.ParseToken
            valid, err := token.VerifyToken(&jwt.Token{
                Raw:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJleHAiOjE1OTQ2MTY0MzAsImlhdCI6MTU5NDYxMjgzMH0.9Bhbk7G25LG-fU9w_HOxDUw3u16cG0QULwIIUENXdJ4",
                Method: &jwt.SigningMethodHMAC{Name: "HS256", Hash: 5},
                Header: map[string]interface{}{"alg": "HS256", "typ": "JWT"},
                Claims: &token.UserClaims{
                    UserID: 0,
                    StandardClaims: jwt.StandardClaims{
                        IssuedAt:  time.Now().Unix(),
                        ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
                    },
                },
                Signature: "9Bhbk7G25LG-fU9w_HOxDUw3u16cG0QULwIIUENXdJ4",
                Valid:     false,
            })
            Expect(err).To(HaveOccurred())
            Expect(valid).To(BeFalse())
        })
    })
})
