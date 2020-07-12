package integration_test

import (
    "database/sql"
    "fmt"
    "os"
    "testing"

    "github.com/onsi/gomega/gexec"

    _ "github.com/go-sql-driver/mysql"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestIntegration(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Integration Suite")
}

var (
    pathToExecutable      string
    databaseURL           string
    databaseVarsAvailable bool
    db                    *sql.DB
)

var _ = BeforeSuite(func() {
    user := os.Getenv("DATABASE_USERNAME")
    password := os.Getenv("DATABASE_PASSWORD")
    host := os.Getenv("DATABASE_HOST")
    port := os.Getenv("DATABASE_PORT")
    databaseName := os.Getenv("DATABASE_NAME")

    if user != "" && password != "" && host != "" && port != "" && databaseName != "" {
        databaseURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, databaseName)
        databaseVarsAvailable = true

        var err error
        db, err = sql.Open("mysql", databaseURL)
        if err != nil {
            panic(err)
        }

        pathToExecutable, err = gexec.Build("github.com/iplay88keys/recipe-box")
        Expect(err).ToNot(HaveOccurred())
    }
})

var _ = AfterSuite(func() {
    gexec.CleanupBuildArtifacts()
})
