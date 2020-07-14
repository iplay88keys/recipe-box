package integration_test

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"

    "github.com/iplay88keys/recipe-box/pkg/api/recipes"
    . "github.com/iplay88keys/recipe-box/pkg/helpers"
    "github.com/iplay88keys/recipe-box/pkg/repositories"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("GetRecipe", func() {
    var (
        username string
        password string
        recipeID int64
    )

    BeforeEach(func() {
        _, err := db.Exec("DELETE FROM users WHERE id IS NOT NULL")
        Expect(err).ToNot(HaveOccurred())

        username = "get_recipe_user"
        password = "Pa3$word123"

        usersRepo := repositories.NewUsersRepository(db)
        userID, err := usersRepo.Insert(username, username + "@example.com", password)
        Expect(err).ToNot(HaveOccurred())

        res, err := db.Exec(`INSERT INTO recipes (
            creator, name, description, servings, prep_time, total_time
            ) VALUES (?, ?, ?, ?, ?, ?)`,
            userID, "Ice Cream", "Yum.", 1, "1 m", "1 m")
        Expect(err).ToNot(HaveOccurred())

        recipeID, err = res.LastInsertId()
        Expect(err).ToNot(HaveOccurred())

        res, err = db.Exec("INSERT INTO ingredients (name) VALUES (?)",
            "Vanilla Ice Cream")
        Expect(err).ToNot(HaveOccurred())

        ingredientID, err := res.LastInsertId()
        Expect(err).ToNot(HaveOccurred())

        res, err = db.Exec("INSERT INTO measurements (name) VALUES (?)",
            "Scoop")
        Expect(err).ToNot(HaveOccurred())

        measurementID, err := res.LastInsertId()
        Expect(err).ToNot(HaveOccurred())

        _, err = db.Exec(`INSERT INTO recipe_ingredients (
            recipe_id, ingredient_id, ingredient_no, amount, measurement_id
            ) VALUES (?, ?, ?, ?, ?)`,
            recipeID, ingredientID, 1, 1, measurementID)
        Expect(err).ToNot(HaveOccurred())

        _, err = db.Exec(`INSERT INTO recipe_steps (
            recipe_id, step_no, instructions
            ) VALUES (?, ?, ?)`,
            recipeID, 1, "Place ice cream in bowl.")
        Expect(err).ToNot(HaveOccurred())
    })

    It("returns a single recipe if authenticated", func() {
        req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/api/v1/recipes/%d", port, recipeID), nil)
        Expect(err).ToNot(HaveOccurred())

        // Needs auth header

        resp, err := client.Do(req)
        Expect(err).ToNot(HaveOccurred())

        defer resp.Body.Close()
        bytes, err := ioutil.ReadAll(resp.Body)
        Expect(err).ToNot(HaveOccurred())

        var recipeList recipes.RecipeResponse
        err = json.Unmarshal(bytes, &recipeList)
        Expect(err).ToNot(HaveOccurred())

        Expect(recipeList).To(Equal(recipes.RecipeResponse{
            Recipe: repositories.Recipe{
                ID:          Int64Pointer(recipeID),
                Name:        StringPointer("Ice Cream"),
                Description: StringPointer("Yum."),
                Creator:     StringPointer(username),
                Servings:    IntPointer(1),
                PrepTime:    StringPointer("1 m"),
                CookTime:    nil,
                CoolTime:    nil,
                TotalTime:   StringPointer("1 m"),
                Source:      nil,
            },
            Ingredients: []*repositories.Ingredient{{
                Ingredient:       StringPointer("Vanilla Ice Cream"),
                IngredientNumber: IntPointer(1),
                Amount:           StringPointer("1"),
                Measurement:      StringPointer("Scoop"),
                Preparation:      nil,
            }},
            Steps: []*repositories.Step{{
                StepNumber:   IntPointer(1),
                Instructions: StringPointer("Place ice cream in bowl."),
            }},
        }))
    })
})
