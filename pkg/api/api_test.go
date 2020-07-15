package api_test

import (
    "errors"
    "io/ioutil"
    "net/http"
    "time"

    "github.com/iplay88keys/my-recipe-library/pkg/token"

    "github.com/iplay88keys/my-recipe-library/pkg/api/auth/authfakes"

    "github.com/iplay88keys/my-recipe-library/pkg/api/auth"

    "github.com/iplay88keys/my-recipe-library/pkg/api"
    "github.com/iplay88keys/my-recipe-library/pkg/helpers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("API", func() {
    var (
        server *api.API
        port   string
    )

    BeforeEach(func() {
        var err error

        port, err = helpers.GetRandomPort()
        Expect(err).ToNot(HaveOccurred())

        tokenService := &authfakes.FakeTokenService{}
        redisRepo := &authfakes.FakeRedisRepo{}

        tokenService.ValidateTokenStub = func(r *http.Request) (*token.AccessDetails, error) {
            if r.Header.Get("Authorization") == "bearer token" {
                return &token.AccessDetails{
                    AccessUuid: "some-uuid",
                    UserId:     10,
                }, nil
            } else {
                return nil, errors.New("auth error")
            }
        }
        redisRepo.RetrieveTokenDetailsStub = func(details *token.AccessDetails) (int64, error) {
            return 10, nil
        }

        authMiddleware := auth.NewMiddleware(tokenService, redisRepo)

        server = api.New(&api.Config{
            Port:           port,
            StaticDir:      "fixtures",
            AuthMiddleware: authMiddleware,
            Endpoints: []api.Endpoint{{
                Path:   "test-unauthenticated-endpoint",
                Method: http.MethodGet,
                Handler: func(w http.ResponseWriter, r *http.Request) {
                    w.Write([]byte("Exists"))
                },
            }, {
                Path:   "test-authenticated-endpoint",
                Method: http.MethodGet,
                Auth:   true,
                Handler: func(w http.ResponseWriter, r *http.Request) {
                    w.Write([]byte("Authenticated"))
                },
            }},
        })
    })

    It("serves the index page for the react app", func() {
        stop := server.Start()
        defer stop()

        client := &http.Client{
            Timeout: 15 * time.Second,
        }

        resp, err := client.Get("http://localhost:" + port)
        Expect(err).ToNot(HaveOccurred())
        Expect(resp.StatusCode).To(Equal(http.StatusOK))

        body, err := ioutil.ReadAll(resp.Body)
        Expect(err).ToNot(HaveOccurred())
        Expect(string(body)).To(ContainSubstring("Test HTML"))
    })

    It("serves unauthenticated api pages", func() {
        stop := server.Start()
        defer stop()

        client := &http.Client{
            Timeout: 15 * time.Second,
        }

        resp, err := client.Get("http://localhost:" + port + "/api/v1/test-unauthenticated-endpoint")
        Expect(err).ToNot(HaveOccurred())
        Expect(resp.StatusCode).To(Equal(http.StatusOK))

        body, err := ioutil.ReadAll(resp.Body)
        Expect(err).ToNot(HaveOccurred())
        Expect(string(body)).To(ContainSubstring("Exists"))
    })

    It("serves authenticated api pages if the user is authenticated", func() {
        stop := server.Start()
        defer stop()

        client := &http.Client{
            Timeout: 15 * time.Second,
        }

        req, err := http.NewRequest(http.MethodGet, "http://localhost:" + port + "/api/v1/test-authenticated-endpoint", nil)
        Expect(err).ToNot(HaveOccurred())
        req.Header.Set("Authorization", "bearer token")

        resp, err := client.Do(req)
        Expect(err).ToNot(HaveOccurred())
        Expect(resp.StatusCode).To(Equal(http.StatusOK))

        body, err := ioutil.ReadAll(resp.Body)
        Expect(err).ToNot(HaveOccurred())
        Expect(string(body)).To(ContainSubstring("Authenticated"))
    })

    It("serves the static files directly", func() {
        stop := server.Start()
        defer stop()

        client := &http.Client{
            Timeout: 15 * time.Second,
        }

        resp, err := client.Get("http://localhost:" + port + "/static.html")
        Expect(err).ToNot(HaveOccurred())
        Expect(resp.StatusCode).To(Equal(http.StatusOK))

        body, err := ioutil.ReadAll(resp.Body)
        Expect(err).ToNot(HaveOccurred())
        Expect(string(body)).To(ContainSubstring("Static HTML"))
    })

    It("returns 404 if the api page does not exist", func() {
        stop := server.Start()
        defer stop()

        client := &http.Client{
            Timeout: 15 * time.Second,
        }

        resp, err := client.Get("http://localhost:" + port + "/api/v1/non-existent")
        Expect(err).ToNot(HaveOccurred())
        Expect(resp.StatusCode).To(Equal(http.StatusNotFound))

        body, err := ioutil.ReadAll(resp.Body)
        Expect(err).ToNot(HaveOccurred())
        Expect(string(body)).To(Equal("page not found"))
    })

    It("returns unauthorized if the user is not authenticated for an endpoint that requires auth", func() {
        stop := server.Start()
        defer stop()

        client := &http.Client{
            Timeout: 15 * time.Second,
        }

        resp, err := client.Get("http://localhost:" + port + "/api/v1/test-authenticated-endpoint")
        Expect(err).ToNot(HaveOccurred())
        Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
    })
})
