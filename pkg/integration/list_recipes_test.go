package integration_test

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os/exec"
    "time"

    "github.com/onsi/gomega/gexec"

    "github.com/iplay88keys/recipe-box/pkg/api/recipes"
    . "github.com/iplay88keys/recipe-box/pkg/helpers"
    "github.com/iplay88keys/recipe-box/pkg/repositories"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("ListRecipes", func() {
    It("returns a list of recipes", func() {
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

        resp, err := client.Get(fmt.Sprintf("http://localhost:%s/api/v1/recipes", port))
        Expect(err).ToNot(HaveOccurred())

        session.Kill()
        Eventually(session).Should(gexec.Exit())

        defer resp.Body.Close()
        bytes, err := ioutil.ReadAll(resp.Body)
        Expect(err).ToNot(HaveOccurred())

        var recipeList recipes.RecipeListResponse
        err = json.Unmarshal(bytes, &recipeList)
        Expect(err).ToNot(HaveOccurred())

        Expect(recipeList).To(Equal(recipes.RecipeListResponse{
            Recipes: []*repositories.Recipe{{
                ID:          Int64Pointer(1),
                Name:        StringPointer("Root Beer Float"),
                Description: StringPointer("Delicious drink for a hot summer day."),
                Creator:     nil,
                PrepTime:    nil,
                Source:      nil,
            }, {
                ID:          Int64Pointer(2),
                Name:        StringPointer("Nana's Beans"),
                Description: StringPointer("Spruced up baked beans."),
                Creator:     nil,
                PrepTime:    nil,
                Source:      nil,
            }},
        }))
    })
})
