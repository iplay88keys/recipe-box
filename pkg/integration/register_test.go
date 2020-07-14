package integration_test

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "os/exec"
    "time"

    "github.com/onsi/gomega/gexec"

    . "github.com/iplay88keys/recipe-box/pkg/helpers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Register", func() {
    BeforeEach(func() {
        _, err := db.Exec("DELETE FROM users WHERE id IS NOT NULL")
        Expect(err).ToNot(HaveOccurred())
    })

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

        username := "some_user"
        body := []byte(fmt.Sprintf(`{
            "username": "%s",
            "email": "someone@example.com",
            "password": "Pa3$word123"
        }`, username))
        req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s/api/v1/users/register", port), bytes.NewBuffer(body))
        Expect(err).ToNot(HaveOccurred())

        resp, err := client.Do(req)
        Expect(err).ToNot(HaveOccurred())

        session.Kill()
        Eventually(session).Should(gexec.Exit())

        Expect(resp.StatusCode).To(Equal(http.StatusOK))

        var count int
        row := db.QueryRow("SELECT COUNT(*) FROM users WHERE username=?", username)
        err = row.Scan(&count)
        Expect(err).ToNot(HaveOccurred())
        Expect(count).To(Equal(1))
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
            "username": "a",
            "email": "a",
            "password": "a"
        }`)
        req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s/api/v1/users/register", port), bytes.NewBuffer(body))
        Expect(err).ToNot(HaveOccurred())

        resp, err := client.Do(req)
        Expect(err).ToNot(HaveOccurred())

        session.Kill()
        Eventually(session).Should(gexec.Exit())

        Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

        bytes, err := ioutil.ReadAll(resp.Body)
        Expect(err).ToNot(HaveOccurred())

        Expect(string(bytes)).To(MatchJSON(` {
            "errors": {
                "email": "Invalid email address",
                "password": "Uppercase letter missing, Numeric character missing, Special character missing, Must be between 6 and 64 characters long",
                "username": "Must be between 6 and 30 characters long"
            }
        }`))
    })
})
