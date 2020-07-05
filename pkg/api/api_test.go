package api_test

import (
    "io/ioutil"
    "net/http"
    "time"

    "github.com/iplay88keys/recipe-box/pkg/api"
    "github.com/iplay88keys/recipe-box/pkg/helpers"

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

        server = api.New(&api.Config{
            Port:      port,
            StaticDir: "fixtures",
            Endpoints: []api.Endpoint{{
                Path:   "test-api-endpoint",
                Method: http.MethodGet,
                Handler: func(w http.ResponseWriter, r *http.Request) {
                    w.Write([]byte("Exists"))
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

    It("serves the api pages", func() {
        stop := server.Start()
        defer stop()

        client := &http.Client{
            Timeout: 15 * time.Second,
        }

        resp, err := client.Get("http://localhost:" + port + "/api/v1/test-api-endpoint")
        Expect(err).ToNot(HaveOccurred())
        Expect(resp.StatusCode).To(Equal(http.StatusOK))

        body, err := ioutil.ReadAll(resp.Body)
        Expect(err).ToNot(HaveOccurred())
        Expect(string(body)).To(ContainSubstring("Exists"))
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
})
