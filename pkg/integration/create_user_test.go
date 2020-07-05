package integration_test

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os/exec"
    "time"

    "github.com/onsi/gomega/gexec"

    "github.com/iplay88keys/recipe-box/pkg/api/users"
    . "github.com/iplay88keys/recipe-box/pkg/helpers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Signup", func() {
    It("creates a new user", func() {
        if !databaseVarsAvailable {
            Skip("Missing some database information. Run the tests from 'scripts/test.sh' to start up the database.")
        }

        client := &http.Client{
            Timeout: 10 * time.Second,
        }

        port, err := GetRandomPort()
        Expect(err).ToNot(HaveOccurred())

        cmd := exec.Command(pathToExecutable,
            "-port", port,
            "-databaseURL", fmt.Sprintf(`"%s"`, databaseURL),
        )

        session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
        Expect(err).ToNot(HaveOccurred())
        time.Sleep(1 * time.Second)

        body := []byte(`{
            "username": "some-user",
            "email": "someone@example.com",
            "password": "Pa3$word123"
        }`)
        req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s/api/v1/signup", port), bytes.NewBuffer(body))
        Expect(err).ToNot(HaveOccurred())

        resp, err := client.Do(req)
        Expect(err).ToNot(HaveOccurred())

        session.Kill()
        Eventually(session).Should(gexec.Exit())

        defer resp.Body.Close()
        bytes, err := ioutil.ReadAll(resp.Body)
        Expect(err).ToNot(HaveOccurred())

        var userSignupResponse users.UserSignupResponse
        err = json.Unmarshal(bytes, &userSignupResponse)
        Expect(err).ToNot(HaveOccurred())

        Expect(userSignupResponse).To(Equal(users.UserSignupResponse{
            EmailExisted: false,
        }))
    })

    It("returns an error if the json data is invalid", func() {
        if !databaseVarsAvailable {
            Skip("Missing some database information. Run the tests from 'scripts/test.sh' to start up the database.")
        }

        client := &http.Client{
            Timeout: 10 * time.Second,
        }

        port, err := GetRandomPort()
        Expect(err).ToNot(HaveOccurred())

        cmd := exec.Command(pathToExecutable,
            "-port", port,
            "-databaseURL", fmt.Sprintf(`"%s"`, databaseURL),
        )

        session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
        Expect(err).ToNot(HaveOccurred())
        time.Sleep(1 * time.Second)

        body := []byte(`{
            "username": "",
            "email": "",
            "password": ""
        }`)
        req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s/api/v1/signup", port), bytes.NewBuffer(body))
        Expect(err).ToNot(HaveOccurred())

        resp, err := client.Do(req)
        Expect(err).ToNot(HaveOccurred())

        session.Kill()
        Eventually(session).Should(gexec.Exit())

        Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

        bytes, err := ioutil.ReadAll(resp.Body)
        Expect(err).ToNot(HaveOccurred())



        Expect(string(bytes)).To(MatchJSON(` {
            "email_existed": false,
            "username_existed": false,
            "errors": [{
                "error_type": "basic",
                "errors": [
                    "username, password, and email are all required"
                ]
            }, {
                "error_type": "email",
                "errors": [
                    "invalid email address"
                ]
            }, {
                "error_type": "password",
                "errors": [
                    "lowercase letter missing",
                    "uppercase letter missing",
                    "numeric character missing",
                    "special character missing",
                    "must be between 8 to 64 characters long"
                ]
            }]
        }`))
    })
})
