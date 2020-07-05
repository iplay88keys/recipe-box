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

var _ = Describe("GetRecipe", func() {
    It("returns a single recipe", func() {
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

        req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/api/v1/recipes/1", port), nil)
        Expect(err).ToNot(HaveOccurred())

        resp, err := client.Do(req)
        Expect(err).ToNot(HaveOccurred())

        session.Kill()
        Eventually(session).Should(gexec.Exit())

        defer resp.Body.Close()
        bytes, err := ioutil.ReadAll(resp.Body)
        Expect(err).ToNot(HaveOccurred())

        var recipeList recipes.RecipeResponse
        err = json.Unmarshal(bytes, &recipeList)
        Expect(err).ToNot(HaveOccurred())

        Expect(recipeList).To(Equal(recipes.RecipeResponse{
            Recipe: repositories.Recipe{
                ID:          Int64Pointer(1),
                Name:        StringPointer("Root Beer Float"),
                Description: StringPointer("Delicious drink for a hot summer day."),
                Creator:     StringPointer("user"),
                Servings:    IntPointer(1),
                PrepTime:    StringPointer("5 m"),
                CookTime:    nil,
                CoolTime:    nil,
                TotalTime:   StringPointer("5 m"),
                Source:      nil,
            },
            Ingredients: []*repositories.Ingredient{{
                Ingredient:       StringPointer("Vanilla Ice Cream"),
                IngredientNumber: IntPointer(1),
                Amount:           StringPointer("1"),
                Measurement:      StringPointer("Scoop"),
                Preparation:      nil,
            }, {
                Ingredient:       StringPointer("Root Beer"),
                IngredientNumber: IntPointer(2),
                Amount:           nil,
                Measurement:      nil,
                Preparation:      nil,
            }},
            Steps: []*repositories.Step{{
                StepNumber:   IntPointer(1),
                Instructions: StringPointer("Place ice cream in glass."),
            }, {
                StepNumber:   IntPointer(2),
                Instructions: StringPointer("Top with Root Beer."),
            }},
        }))
    })
})
